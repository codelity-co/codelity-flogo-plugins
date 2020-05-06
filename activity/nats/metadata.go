package nats

import (
	"reflect"
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct {
	ClusterUrls string                 `md:"clusterUrls,required"` // The NATS cluster to connect to
	ConnName    string                 `md:"connName"`
	Auth        map[string]interface{} `md:"auth"`      // Auth setting
	Reconnect   map[string]interface{} `md:"reconnect"` // Reconnect setting
	SslConfig   map[string]interface{} `md:"sslConfig"` // SSL config setting
	Streaming   map[string]interface{} `md:"streaming"` // NATS streaming config
	DataType    string                 `md:"dataType"`  // Data type
}

func (s *Settings) FromMap(values map[string]interface{}) error {
	var err error
	s.ClusterUrls, err = coerce.ToString(values["clusterUrls"])
	if err != nil {
		return err
	}

	s.ConnName, err = coerce.ToString(values["connName"])
	if err != nil {
		return err
	}

	s.DataType, err = coerce.ToString(values["dataType"])
	if err != nil {
		return err
	}

	s.Auth = make(map[string]interface{})
	for k, v := range coerce.ToObject(values["auth"]) {
		s.Auth[k.(string)] := s.MapValue(v)
	}

	s.Reconnect = make(map[string]interface{})
	for k, v := range coerce.ToObject(values["reconnect"]) {
		s.Auth[k.(string)] := s.MapValue(v)
	}

	s.SslConfig = make(map[string]interface{})
	for k, v := range coerce.ToObject(values["sslConfig"]) {
		s.Auth[k.(string)] := s.MapValue(v)
	}

	s.Streaming = make(map[string]interface{})
	for k, v := range coerce.ToObject(values["streaming"]) {
		s.Auth[k.(string)] := s.MapValue(v)
	}

}

func (s *Settings) ToMap() map[string]interface{} {

	return map[string]interface{}{
		"clusterUrls": s.ClusterUrls,
		"connName": s.ConnName,
		"dataType": s.DataType,
		"auth": s.Auth,
		"reconnect": s.Reconnect,
		"sslConfig": s.SslConfig,
		"streaming": s.Streaming,
	}

}

func (s *Settings) MapValue(value interface{}) interface{} {
	return coerce.ToAny(value)
}

type Input struct {
	Subject   string `md:"subject"`
	ChannelId string `md:"channelId"`
	Data      string `md:"data"`
}

func (i *Input) FromMap(values map[string]interface{}) error {
	var err error
	i.Subject, err = coerce.ToString(values["subject"])
	if err != nil {
		return err
	}
	i.ChannelId, err = coerce.ToString(values["channelId"])
	if err != nil {
		return err
	}

	i.Data, err = coerce.ToString(values["data"])
	if err != nil {
		return err
	}
	return nil
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"subject":   i.Subject,
		"channelId": i.ChannelId,
		"data":      i.Data,
	}
}

type Output struct {
	Status string                 `md:"status"`
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
