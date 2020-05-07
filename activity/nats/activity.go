package nats

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/support/log"

	nats "github.com/nats-io/nats.go"
	stan "github.com/nats-io/stan.go"
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

//New optional factory method, should be used if one activity instance per configuration is desired
func New(ctx activity.InitContext) (activity.Activity, error) {

	var (
		err error 
		nc *nats.Conn
	)
	logger := ctx.Logger()

	logger.Debug("Running New method of activity...")

	s := &Settings{}

	logger.Debug("Mapping Settings struct...")
	err = s.FromMap(ctx.Settings())
	if err != nil {
		logger.Errorf("Map settings error: %v", err)
		return nil, err
	}
	logger.Debug("Mapped Settings struct successfully")

	logger.Debugf("From Map Setting: %v", s)

	logger.Debug("Getting NATS connection...")
	nc, err = getNatsConnection(logger, s)
	if err != nil {
		logger.Errorf("NATS connection error: %v", err)
		return nil, err
	}
	logger.Debug("Got NATS connection")

	logger.Debug("Creating Activity struct...")
	act := &Activity{
		activitySettings: s,
		logger:           logger,
		natsConn:         nc,
		natsStreaming:    false,
	}
	logger.Debug("Created Activity struct successfully")

	logger.Debugf("Streaming: %v", s.Streaming)

	if enableStreaming, ok := s.Streaming["enableStreaming"]; ok {
		logger.Debug("Enabling NATS streaming...")
		act.natsStreaming = enableStreaming.(bool)
		if err != nil {
			return nil, err
		}
		if act.natsStreaming {
			logger.Debug("Getting STAN connection...")
			act.stanConn, err = getStanConnection(s.Streaming, nc)
			if err != nil {
				logger.Errorf("STAN connection error: %v", err)
				return nil, err
			}
			logger.Debug("Got STAN connection")
		}
		logger.Debug("Enabled NATS streaming successfully")
	}


	logger.Debug("Finished New method of activity")
	return act, nil
}

// Activity is an sample Activity that can be used as a base to create a custom activity
type Activity struct {
	activitySettings *Settings
	logger           log.Logger
	natsConn         *nats.Conn
	natsStreaming    bool
	stanConn         stan.Conn
}

// Metadata returns the activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements api.Activity.Eval - Logs the Message
func (a *Activity) Eval(ctx activity.Context) (bool, error) {

	var (
		err    error
		result map[string]interface{}
	)

	result = make(map[string]interface{})

	a.logger.Debug("Running Eval method of activity...")
	input := &Input{}
	a.logger.Debug("Getting Input object from context...")
	err = ctx.GetInputObject(input)
	if err != nil {
		a.logger.Errorf("Error getting Input object: %v", err)
		_ = a.OutputToContext(ctx, nil, err)
		return true, err
	}
	a.logger.Debug("Got Input object successfully")
	a.logger.Debugf("Input: %v", input)

	if !a.natsStreaming {
		a.logger.Debug("Publishing data to NATS subject...")
		if err = a.natsConn.Publish(input.Subject, []byte(input.Data)); err != nil {
			a.logger.Errorf("Error publishing data to NATS subject: %v", err)
			_ = a.OutputToContext(ctx, nil, err)
			return true, err
		}
		a.logger.Debug("Published data to NATS subject")
	} else {
		
		message := map[string]interface{}{
			"subject": input.Subject,
			"message": input.Data,
		}

		var messageBytes []byte
		messageBytes, err = json.Marshal(message)
		if err != nil {
			a.logger.Errorf("Marshal error: %v", err)
			return true, err
		}
		a.logger.Debug("Publishing data to STAN Channel...")
		result["ackedNuid"], err = a.stanConn.PublishAsync(input.ChannelId, messageBytes, func(ackedNuid string, err error) {
			if err != nil {
				a.logger.Errorf("STAN acknowledgement error: %v", err)
			}
		})
		if err != nil {
			a.logger.Errorf("Error publishing data to STAN channel: %v", err)
			_ = a.OutputToContext(ctx, nil, err)
			return true, err
		}
		a.logger.Debugf("Published data to STAN channel: %v", result)
	}

	err = a.OutputToContext(ctx, result, nil)
	if err != nil {
		a.logger.Errorf("Error setting output object in context: %v", err)
		return true, err
	}
	a.logger.Debug("Successfully set output object in context")

	return true, nil
}

func (a *Activity) OutputToContext(ctx activity.Context, result map[string]interface{}, err error) error {
	a.logger.Debug("Createing Ouptut struct...")
	var output *Output
	if err != nil {
		output = &Output{Status: "ERROR", Result: map[string]interface{}{"errorMessage": fmt.Sprintf("%v", err)}}
	} else {
		output = &Output{Status: "SUCCESS", Result: result}
	}
	a.logger.Debug("Setting output object in context...")
	return ctx.SetOutputObject(output)
}

func getNatsConnection(logger log.Logger, settings *Settings) (*nats.Conn, error) {
	var (
		err           error
		authOpts      []nats.Option
		reconnectOpts []nats.Option
		sslConfigOpts []nats.Option
		urlString     string
	)

	// Check ClusterUrls
	logger.Debug("Checking clusterUrls...")
	if err := checkClusterUrls(logger, settings); err != nil {
		logger.Errorf("Error checking clusterUrls: %v", err)
		return nil, err
	}
	logger.Debug("Checked")

	urlString = settings.ClusterUrls

	logger.Debug("Getting NATS connection auth settings...")
	authOpts, err = getNatsConnAuthOpts(logger, settings)
	if err != nil {
		logger.Errorf("Error getting NATS connection auth settings:: %v", err)
		return nil, err
	}
	logger.Debug("Got NATS connection auth settings")

	logger.Debug("Getting NATS connection reconnect settings...")
	reconnectOpts, err = getNatsConnReconnectOpts(logger, settings)
	if err != nil {
		logger.Errorf("Error getting NATS connection reconnect settings:: %v", err)
		return nil, err
	}
	logger.Debug("Got NATS connection reconnect settings")

	logger.Debug("Getting NATS connection sslConfig settings...")
	sslConfigOpts, err = getNatsConnSslConfigOpts(logger, settings)
	if err != nil {
		logger.Errorf("Error getting NATS connection sslConfig settings:: %v", err)
		return nil, err
	}
	logger.Debug("Got NATS connection sslConfig settings")

	natsOptions := append(authOpts, reconnectOpts...)
	natsOptions = append(natsOptions, sslConfigOpts...)

	// Check ConnName
	if len(settings.ConnName) > 0 {
		natsOptions = append(natsOptions, nats.Name(settings.ConnName))
	}

	return nats.Connect(urlString, natsOptions...)

}

// checkClusterUrls is function to all valid NATS cluster urls
func checkClusterUrls(logger log.Logger, settings *Settings) error {
	// Check ClusterUrls
	clusterUrls := strings.Split(settings.ClusterUrls, ",")
	logger.Debugf("clusterUrls: %v", clusterUrls)
	if len(clusterUrls) < 1 {
		return fmt.Errorf("ClusterUrl [%v] is invalid, require at least one url", settings.ClusterUrls)
	}
	for _, v := range clusterUrls {
		logger.Debugf("v: %v", v)
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
func getNatsConnAuthOpts(logger log.Logger, settings *Settings) ([]nats.Option, error) {
	opts := make([]nats.Option, 0)
	// Check auth setting
	logger.Debugf("settings.Auth: %v", settings.Auth)
	if settings.Auth != nil && len(settings.Auth) > 0 {
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

func getNatsConnReconnectOpts(logger log.Logger, settings *Settings) ([]nats.Option, error) {
	opts := make([]nats.Option, 0)
	// Check reconnect setting
	logger.Debugf("settings.Reconnect: %v", settings.Reconnect)
	if settings.Reconnect != nil && len(settings.Reconnect) > 0 {

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

func getNatsConnSslConfigOpts(logger log.Logger, settings *Settings) ([]nats.Option, error) {
	opts := make([]nats.Option, 0)

	// Check sslConfig setting
	logger.Debugf("settings.SslConfig: %v", settings.SslConfig)
	if settings.SslConfig != nil && len(settings.SslConfig) > 0 {

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

func getStanConnection(mapping map[string]interface{}, conn *nats.Conn) (stan.Conn, error) {

	var (
		err       error
		clusterID string
		ok        bool
		hostname  string
		sc        stan.Conn
	)

	if _, ok = mapping["clusterId"]; !ok {
		return nil, fmt.Errorf("clusterId not found")
	}

	clusterID = mapping["clusterId"].(string)

	hostname, err = os.Hostname()
	hostname = strings.Split(hostname, ".")[0]
	hostname = strings.Split(hostname, ":")[0]

	fmt.Println(hostname)

	if err != nil {
		return nil, err
	}

	sc, err = stan.Connect(clusterID, hostname, stan.NatsConn(conn))
	if err != nil {
		return nil, err
	}

	return sc, nil
}
