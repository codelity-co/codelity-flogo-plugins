package sample

type Settings struct {
	Endpoint string `md:"endpoint,required"`
	AccessKey string `md:"accessKey,required"`
	SecretKey string `md:"secretKey,required"`
	EnableSsl bool `md:"enableSsl"`
	BucketName string `md:"bucketName,required"`
	Region string `md:"region"`
	MethodName string `md:"methodName,required"` 
	MethodOptions map[string]interface{} `md:"methodOpts"`
	DataType string `md:"dataType,required"`
}

type Input struct {
	ObjectName string `md:"objectName,required"`
	Data interface{} `md:"data,required"`
}

func (r *Input) FromMap(values map[string]interface{}) error {
	r.ObjectName, _ = values["objectName"].(string)
	r.Data = values["data"]
	return nil
}

func (r *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"objectName": r.ObjectName,
		"data": r.Data,
	}
}

type Output struct {
	Status string `md:"status,required"`
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
