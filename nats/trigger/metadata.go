package nats

import (
	"github.com/project-flogo/core/data/coerce"
)

// Settings is NATS connection setting construct
type Settings struct {
	ClusterUrls string `md:"clusterUrls,required"` 			// The NATS cluster to connect to
	ConnName    string `md:"connName"`
	Auth       	map[string]interface{} `md:"auth"` 			// Auth setting
	Reconnect 	map[string]interface{} `md:"reconnect"` // Reconnect setting
	SslConfig 	map[string]interface{} `md:"sslConfig"` // SSL config setting
	Streaming   map[string]interface{} `md:"streaming"` // NATS streaming config
}

type HandlerSettings struct {
	Subject string `md:"subject,required"`
	Queue string `md:"queue"`
	ChannelId string `md:"channelId"`
	DurableName string `md:"durableName"`
	MaxInFlight int `md:"maxInFlight"`
}

type Output struct {
	Message string `md:"message"`
}

func (o *Output) FromMap(values map[string]interface{}) error {

	var err error
	o.Message, err = coerce.ToString(values["message"])
	if err != nil {
		return err
	}

	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"message": o.Message,
	}
}

// type Reply struct {
// 	AReply interface{} `md:"aReply"`
// }

// func (r *Reply) FromMap(values map[string]interface{}) error {
// 	r.AReply = values["aReply"]
// 	return nil
// }

// func (r *Reply) ToMap() map[string]interface{} {
// 	return map[string]interface{}{
// 		"aReply": r.AReply,
// 	}
// }
