package nats

import (
	"context"
	"crypto/tls"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/trigger"

	nats "github.com/nats-io/nats.go"
)

var triggerMd = trigger.NewMetadata(&Settings{}, &HandlerSettings{}, &Output{})

func init() {
	_ = trigger.Register(&Trigger{}, &Factory{})
}

// Handler is a NATS subject handler
type Handler struct {
	handlerSettings	*HandlerSettings
	logger         log.Logger
	messageChannel chan *nats.Msg
	natsConn       *nats.Conn
	natsSubscription *nats.Subscription
	stopChannel    chan bool
	triggerHandler trigger.Handler
}

// Factory struct
type Factory struct {
}

// Trigger struct
type Trigger struct {
	id       string
	triggerSettings *Settings
	natsHandlers []*Handler
}

// New trigger method of Factory
func (*Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	s := &Settings{}
	err := metadata.MapToStruct(config.Settings, s, true)
	if err != nil {
		return nil, err
	}

	return &Trigger{id: config.Id, triggerSettings: s}, nil

}

// Metadata method of Factory
func (f *Factory) Metadata() *trigger.Metadata {
	return triggerMd
}

// Metadata implements trigger.Trigger.Metadata
func (t *Trigger) Metadata() *trigger.Metadata {
	return triggerMd
}

// Initialize method of trigger
func (t *Trigger) Initialize(ctx trigger.InitContext) error {

	logger := ctx.Logger()

	for _, handler := range ctx.GetHandlers() {

		handlerSettings := &HandlerSettings{}

		if err := metadata.MapToStruct(handler.Settings(), handlerSettings, true); err != nil {
			return err
		}

		nc, err := getNatsConnection(logger, t.triggerSettings)
		if err != nil {
			return err
		}

		messageChannel := make(chan *nats.Msg)
		stopChannel := make(chan bool)
		natsHandler := &Handler{
			handlerSettings: handlerSettings,
			logger: logger,
			messageChannel: messageChannel,
			natsConn: nc,
			stopChannel: stopChannel,
			triggerHandler: handler,
		}
		t.natsHandlers = append(t.natsHandlers, natsHandler)

	}

	return nil
}

// Start implements util.Managed.Start
func (t *Trigger) Start() error {
	for _, handler := range t.natsHandlers {
		_ = handler.Start()
	}
	return nil
}

// Stop implements util.Managed.Stop
func (t *Trigger) Stop() error {
	for _, handler := range t.natsHandlers {
		_ = handler.Stop()
	}
	return nil
}

func (h *Handler) handleMessage() {
	for {
		select {
		case done := <-h.stopChannel:
			if done {
				return
			}
		case msg := <-h.messageChannel:
			var (
				err error
				// results map[string]interface{}
			)
			out := &Output{}
			out.Message, err = coerce.ToString(msg.Data)
			if err != nil {
				h.logger.Errorf("Run action for handler [%v] failed for reason [%v] message lost", h.triggerHandler.Name(), err)
			}
			_, err = h.triggerHandler.Handle(context.Background(), out.ToMap)
			if err != nil {
				h.logger.Errorf("Run action for handler [%v] failed for reason [%v] message lost", h.triggerHandler.Name(), err)
			}

			// reply := &Reply{}
			// _ = reply.FromMap(results)

			//do something with the reply in future

		}
	}
}

// Start starts the handler
func (h *Handler) Start() error {
	var err error
	go h.handleMessage()
	if len(h.handlerSettings.Queue) > 0 {
		h.natsSubscription, err =  h.natsConn.QueueSubscribe(h.handlerSettings.Subject, h.handlerSettings.Queue, func(m *nats.Msg) {
			h.messageChannel <- m
		})
		if err != nil {
			return err
		}
	} else {
		h.natsSubscription, err = h.natsConn.Subscribe(h.handlerSettings.Subject, func(m *nats.Msg) {
			h.messageChannel <- m
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Stop stops the handler
func (h *Handler) Stop() error {
	_ = h.natsConn.Drain()
	h.natsConn.Close()
	h.stopChannel <- true
	close(h.messageChannel)
	close(h.stopChannel)
	return nil
}

func getNatsConnection(logger log.Logger, settings *Settings) (*nats.Conn, error) {
	var (
		err error 
		authOpts []nats.Option
		reconnectOpts []nats.Option
		sslConfigOpts []nats.Option
		urlString string
	)

	// Check ClusterUrls
	if err := checkClusterUrls(settings); err != nil {
		return nil, err
	}

	urlString = settings.ClusterUrls

	authOpts, err = getNatsConnAuthOpts(settings)
	if err != nil {
		return nil, err
	}

	reconnectOpts, err = getNatsConnReconnectOpts(settings)
	if err != nil {
		return nil, err
	}

	sslConfigOpts, err = getNatsConnSslConfigOpts(settings)
	if err != nil {
		return nil, err
	}
	
	natsOptions := append(authOpts, reconnectOpts...)
	natsOptions = append(natsOptions, sslConfigOpts...)

	// Check ConnName
	if len(settings.ConnName) > 0 {
		natsOptions = append(natsOptions, nats.Name(settings.ConnName))
	}
	
	return nats.Connect(urlString, natsOptions...)

}

// checkClusterUrls is function to all valid NATS cluster urls
func checkClusterUrls(settings *Settings) error {
	// Check ClusterUrls
	clusterUrls := strings.Split(settings.ClusterUrls, ",")
	if len(clusterUrls) < 1 {
		return fmt.Errorf("ClusterUrl [%v] is invalid, require at least one url", settings.ClusterUrls)
	}
	for _, v := range clusterUrls {
		if err := validateClusterURL(v); err != nil {
			return err
		}
	}
	return nil
}

// validateClusterUrl is function to check NATS cluster url specificaiton
func validateClusterURL(url string) error {
	hostPort := strings.Split(url, ":")
	if len(hostPort) < 2 || len(hostPort) > 3 {
		return fmt.Errorf("ClusterUrl must be composed of sections like \"{nats|tls}://host[:port]\"")
	}
	if len(hostPort) == 3 {
		i, err := strconv.Atoi(hostPort[2])
		if err != nil || i < 0 || i > 32767 {
			return fmt.Errorf("port specification [%v] is not numeric and between 0 and 32767", hostPort[2])
		}
	}
	if (hostPort[0] != "nats") && (hostPort[0] != "tls") {
		return fmt.Errorf("protocol schema [%v] is not nats or tls", hostPort[0])
	}
	
	return nil
}

// getNatsConnAuthOps return slice of nats.Option specific for NATS authentication
func getNatsConnAuthOpts(settings *Settings) ([]nats.Option, error) {
	opts := make([]nats.Option, 0)
	// Check auth setting
	if settings.Auth != nil {
		if username, ok := settings.Auth["username"]; ok { // Check if usename is defined
			password, ok := settings.Auth["password"] // check if password is defined
			if !ok {
				return nil, fmt.Errorf("Missing password")
			} else {
				// Create UserInfo NATS option
				opts = append(opts, nats.UserInfo(username.(string), password.(string)))
			}
		} else if token, ok := settings.Auth["token"]; ok { // Check if token is defined
			opts = append(opts, nats.Token(token.(string)))
		} else if nkeySeedfile, ok := settings.Auth["nkeySeedfile"]; ok { // Check if nkey seed file is defined
			nkey, err := nats.NkeyOptionFromSeed(nkeySeedfile.(string))
			if err != nil {
				return nil, err
			}
			opts = append(opts, nkey)
		} else if credfile, ok := settings.Auth["credfile"]; ok { // Check if credential file is defined
			opts = append(opts, nats.UserCredentials(credfile.(string)))
		}
	}
	return opts, nil
}

func getNatsConnReconnectOpts(settings *Settings) ([]nats.Option, error) {
	opts := make([]nats.Option, 0)
	// Check reconnect setting
	if settings.Reconnect != nil {

		// Enable autoReconnect
		if autoReconnect, ok := settings.Reconnect["autoReconnect"]; ok {
			if !autoReconnect.(bool) {
				opts = append(opts, nats.NoReconnect())
			}
		}

		// Max reconnect attempts
		if maxReconnects, ok := settings.Reconnect["maxReconnects"]; ok {
			opts = append(opts, nats.MaxReconnects(maxReconnects.(int)))
		}

		// Don't randomize 
		if dontRandomize, ok := settings.Reconnect["dontRandomize"]; ok {
			if dontRandomize.(bool) {
				opts = append(opts, nats.DontRandomize())
			}
		}

		// Reconnect wait in seconds
		if reconnectWait, ok := settings.Reconnect["reconnectWait"]; ok {
			duration, err := time.ParseDuration(fmt.Sprintf("%vs", reconnectWait))
			if err != nil {
				return nil, err
			}
			opts = append(opts, nats.ReconnectWait(duration))
		}

		// Reconnect buffer size in bytes
		if reconnectBufSize, ok := settings.Reconnect["reconnectBufSize"]; ok {
			opts = append(opts, nats.ReconnectBufSize(reconnectBufSize.(int)))
		}
	}
	return opts, nil
}

func getNatsConnSslConfigOpts(settings *Settings) ([]nats.Option, error) {
	opts := make([]nats.Option, 0)

	// Check sslConfig setting
	if settings.SslConfig != nil {
		
		// Skip verify
		if skipVerify, ok := settings.SslConfig["skipVerify"]; ok {
			opts = append(opts, nats.Secure(&tls.Config{
				InsecureSkipVerify: skipVerify.(bool),
			}))
		}

		// CA Root
		if caFile, ok := settings.SslConfig["caFile"]; ok {
			opts = append(opts, nats.RootCAs(caFile.(string)))
			// Cert file
			if certFile, ok := settings.SslConfig["certFile"]; ok {
				if keyFile, ok := settings.SslConfig["keyFile"]; ok {
					opts = append(opts, nats.ClientCert(certFile.(string), keyFile.(string)))
				} else {
					return nil, fmt.Errorf("Missing keyFile setting")
				}
			} else {
				return nil, fmt.Errorf("Missing certFile setting")
			}
		} else {
			return nil, fmt.Errorf("Missing caFile setting")
		}

	}
	return opts, nil
}