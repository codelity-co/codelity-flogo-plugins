package minio

import (
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/app/resolve"
)

type Settings struct {
	Endpoint string `md:"endpoint,required"`
	AccessKey string `md:"accessKey,required"`
	SecretKey string `md:"secretKey,required"`
	EnableSsl bool `md:"enableSsl"`
	BucketName string `md:"bucketName,required"`
	Region string `md:"region"`
	MethodName string `md:"methodName,required"` 
	MethodOptions map[string]interface{} `md:"methodOptions"`
	DataType string `md:"dataType,required"`
}

func (s *Settings) FromMap(values map[string]interface{}) error {

	var (
		err error
	)

	s.Endpoint, err = coerce.ToString(values["endpoint"])
	if err != nil {
		return err
	}

	s.AccessKey, err = coerce.ToString(values["accessKey"])
	if err != nil {
		return err
	}

	s.SecretKey, err = coerce.ToString(values["secretKey"])
	if err != nil {
		return err
	}

	s.EnableSsl, err = coerce.ToBool(values["enableSsl"])
	if err != nil {
		return err
	}

	s.BucketName, err = coerce.ToString(values["bucketName"])
	if err != nil {
		return err
	}

	s.Region, err = coerce.ToString(values["region"])
	if err != nil {
		return err
	}

	s.MethodName, err = coerce.ToString(values["methodName"])
	if err != nil {
		return err
	}

	s.DataType, err = coerce.ToString(values["dataType"])
	if err != nil {
		return err
	}

	if values["methodOptions"] != nil {
		s.MethodOptions = make(map[string]interface{})
		for k, v := range values["methodOptions"].(map[string]interface{}) {
			s.MethodOptions[k], err = s.MapValue(v)
			if err != nil {
				return err
			}
		}
	}

	return nil

}

func (s *Settings) ToMap() map[string]interface{} {

	return map[string]interface{}{
		"endpoint": s.Endpoint,
		"accessKey":    s.AccessKey,
		"secretKey": s.SecretKey,
		"enableSsl": s.EnableSsl,
		"bucketName": s.BucketName,
		"region": s.Region,
		"methodName": s.MethodName,
		"methodOptions": s.MethodOptions,
		"dataType":    s.DataType,
	}

}

func (s *Settings) MapValue(value interface{}) (interface{}, error) {
	var (
		err      error
		anyValue interface{}
	)

	switch val := value.(type) {
	case string:
		strVal := val
		if len(strVal) > 0 && strVal[0] == '=' {
			anyValue, err = resolve.Resolve(strVal[1:], nil)
			if err != nil {
				return nil, err
			}
		} else {
			anyValue, err = coerce.ToAny(val)
			if err != nil {
				return nil, err
			}
		}
	default:
		anyValue, err = coerce.ToAny(val)
		if err != nil {
			return nil, err
		}
	}

	return anyValue, nil
}

type Input struct {
	ObjectName string `md:"objectName,required"`
	Data interface{} `md:"data,required"`
}

func (r *Input) FromMap(values map[string]interface{}) error {
	var err error

	r.ObjectName, err = coerce.ToString(values["objectName"])
	if err != nil {
		return err
	}

	r.Data, err = coerce.ToAny(values["data"])
	if err != nil {
		return err
	}
	
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
	var err error

	o.Status, err = coerce.ToString(values["status"])
	if err != nil {
		return err
	}

	o.Result, err = coerce.ToObject(values["result"])
	if err != nil {
		return err
	}

	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"status": o.Status,
		"result": o.Result,
	}
}
