package nats

import (
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct {
	ClusterUrls string `md:"clusterUrls,required"` 			// The NATS cluster to connect to
	ConnName    string `md:"connName"`
	Auth       	map[string]interface{} `md:"auth"` 			// Auth setting
	Reconnect 	map[string]interface{} `md:"reconnect"` // Reconnect setting
	SslConfig 	map[string]interface{} `md:"sslConfig"` // SSL config setting
	Streaming   map[string]interface{} `md:"streaming"` // NATS streaming config
	DataType    string                 `md:"dataType"`  // Data type
}

type Input struct {
	Subject string `md:"subject"`
	ChannelId string `md:"channelId"`
	Data string `md:"data"`
}

func (r *Input) FromMap(values map[string]interface{}) error {
	var err error
	r.Subject, err = coerce.ToString(values["subject"])
	if err != nil {
		return err
	}
	r.ChannelId, err = coerce.ToString(values["channelId"])
	if err != nil {
		return err
	}

	r.Data, err = coerce.ToString(values["data"])
	if err != nil {
		return err
	}
	return nil
}

func (r *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"subject": r.Subject,
		"channelId": r.ChannelId,
		"data": r.Data,
	}
}

type Output struct {
	Status string `md:"status"`
	Result map[string]interface{} `md:"result"`
}

func (o *Output) FromMap(values map[string]interface{}) error {
	o.Status, _ = values["status"].(string)
	o.Result, _ = values["result"].(map[string]interface{})
	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"status": o.Status,
		"result": o.Result,
	}
}
