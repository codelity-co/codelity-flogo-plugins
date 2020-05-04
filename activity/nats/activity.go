package nats

import (
	"crypto/tls"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"

	nats "github.com/nats-io/nats.go"
	stan "github.com/nats-io/stan.go"
)

func init() {
	_ = activity.Register(&Activity{}) //activity.Register(&Activity{}, New) to create instances using factory method 'New'
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

//New optional factory method, should be used if one activity instance per configuration is desired
func New(ctx activity.InitContext) (activity.Activity, error) {

	logger := ctx.Logger()

	logger.Debugf("Running New method of activity...")

	s := &Settings{}
	logger.Debugf("Mapping Settings struct...")
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		logger.Errorf("Map settings error: %v", err)
		return nil, err
	}
	logger.Debugf("Mapped Settings struct successfully")

	logger.Debugf("Setting: %v", s)

	logger.Debugf("Getting NATS connection...")
	nc, err := getNatsConnection(logger, s)
	if err != nil {
		logger.Errorf("NATS connection error: %v", err)
		return nil, err
	}
	logger.Debugf("Got NATS connection")

	logger.Debugf("Creating Activity struct...")
	act := &Activity{
		activitySettings: s,
		logger:           logger,
		natsConn:         nc,
		natsStreaming:    false,
	}
	logger.Debugf("Created Activity struct successfully")

	if enableStreaming, ok := s.Streaming["enableStreaming"]; ok {
		logger.Debugf("Enabling NATS streaming...")
		act.natsStreaming = enableStreaming.(bool)
		if act.natsStreaming {
			logger.Debugf("Getting STAN connection...")
			act.stanConn, err = getStanConnection(s, nc)
			if err != nil {
				logger.Errorf("STAN connection error: %v", err)
				return nil, err
			}
			logger.Debugf("Got STAN connection")
		}
		logger.Debugf("Enabled NATS streaming successfully")
	}

	logger.Debugf("Finished New method of activity")
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
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	a.logger.Debugf("Running Eval method of activity...")
	input := &Input{}
	a.logger.Debugf("Getting Input object from context...")
	err = ctx.GetInputObject(input)
	if err != nil {
		return true, err
	}
	a.logger.Debugf("Got Input object successfully")
	a.logger.Debugf("Input: %v", input)

	a.logger.Debugf("Converting input.Data to []uint8...")
	dataBytes := []uint8(input.Data)
	a.logger.Debugf("Converted input.Data to []uint8")

	if !a.natsStreaming {
		a.logger.Debugf("Publishing data to NATS subject...")
		if err := a.natsConn.Publish(input.Subject, dataBytes); err != nil {
			return true, err
		}
		a.logger.Debugf("Published data to NATS subject")
	} else {
		a.logger.Debugf("Publishing data to STAN Channel...")
		if err := a.stanConn.Publish(input.ChannelId, dataBytes); err != nil {
			return true, err
		}
		a.logger.Debugf("Published data to STAN channel")
	}

	a.logger.Debugf("Createing Ouptut struct...")
	output := &Output{Status: "SUCCESS"}
	a.logger.Debugf("Setting output object in context...")
	err = ctx.SetOutputObject(output)
	if err != nil {
		return true, err
	}
	a.logger.Debugf("Successfully set output object in context")

	return true, nil
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
	logger.Debugf("Checking clusterUrls...")
	if err := checkClusterUrls(settings); err != nil {
		return nil, err
	}
	logger.Debugf("Checked")

	urlString = settings.ClusterUrls

	logger.Debugf("Getting NATS connection auth settings...")
	authOpts, err = getNatsConnAuthOpts(settings)
	if err != nil {
		return nil, err
	}
	logger.Debugf("Got NATS connection auth settings")

	logger.Debugf("Getting NATS connection reconnect settings...")
	reconnectOpts, err = getNatsConnReconnectOpts(settings)
	if err != nil {
		return nil, err
	}
	logger.Debugf("Got NATS connection reconnect settings")

	logger.Debugf("Getting NATS connection sslConfig settings...")
	sslConfigOpts, err = getNatsConnSslConfigOpts(settings)
	if err != nil {
		return nil, err
	}
	logger.Debugf("Got NATS connection sslConfig settings")

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

func getStanConnection(ts *Settings, conn *nats.Conn) (stan.Conn, error) {

	var (
		err       error
		clusterId interface{}
		ok        bool
		hostname  string
		sc        stan.Conn
	)

	clusterId, ok = ts.Streaming["clusterId"]
	if !ok {
		return nil, fmt.Errorf("clusterId not found")
	}

	hostname, err = os.Hostname()
	hostname = strings.Split(hostname, ".")[0]
	hostname = strings.Split(hostname, ":")[0]

	fmt.Println(hostname)

	if err != nil {
		return nil, err
	}

	sc, err = stan.Connect(clusterId.(string), hostname, stan.NatsConn(conn))
	if err != nil {
		return nil, err
	}

	return sc, nil
}
