package varmapping

import (
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/app/resolve"
)

type Settings struct {}

type Input struct {
	InVar map[string]interface{} `md:"in,required"`
}

func (i *Input) FromMap(values map[string]interface{}) error {
	var err error 

	i.InVar, err = coerce.ToObject(values["in"])
	if err != nil {
		return err
	}

	return nil
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"in": i.InVar,
	}
}

func (i *Input) MapValue(value interface{}) (interface{}, error) {

	var (
		err      error
		anyValue interface{}
	)

	switch val := value.(type) {
	case string:
		if len(val) > 0 && val[0] == '=' {
			anyValue, err = resolve.Resolve(val[1:], nil)
			if err != nil {
				return nil, err
			}
		} else {
			anyValue, err = coerce.ToAny(val)
			if err != nil {
				return nil, err
			}
		}
		
	case map[string]interface{}:
		dataMap := make(map[string]interface{})
		for k, v := range val {
			dataMap[k], err = i.MapValue(v)
			if err != nil {
				return nil, err
			}
		}
		anyValue = dataMap

	default:
		anyValue, err = coerce.ToAny(val)
		if err != nil {
			return nil, err
		}
	}

	return anyValue, nil
}

type Output struct {
	OutVar map[string]interface{} `md:"out,required"`
}

func (o *Output) FromMap(values map[string]interface{}) error {
	var err error

	o.OutVar, err = coerce.ToObject(values["out"])
	if err != nil {
		return err
	}

	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"out": o.OutVar,
	}
}
