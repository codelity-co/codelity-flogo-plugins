package sample

import "github.com/project-flogo/core/data/coerce"

type Settings struct {
	Endpoint string `md:"endpoint,required"`
	AccessKey string `md:"accessKey,required"`
	SecretKey string `md:"secretKey,required"`
	EnableSsl bool `md:"enableSsl"`
}

type Input struct {
	Method string `md:"method,required"`
	Params map[string]interface{} `md:"params,required"`
}

func (r *Input) FromMap(values map[string]interface{}) error {
	r.Method, _ := values["method"].(string)
	r.Params, _ := values["params"].(map[string]interface{})
	return nil
}

func (r *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"method": r.Method,
		"params": r.Params,
	}
}

type Output struct {
	Status string `md:"status,required"`,
	Result map[string]interface{} `md:"result"`
}

func (o *Output) FromMap(values map[string]interface{}) error {
	o.Status, _ := values["status"].(string)
	o.Result, _ := values["result"].(map[string]interface{})
	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"status": o.Status,
		"result": o.Result,
	}
}
