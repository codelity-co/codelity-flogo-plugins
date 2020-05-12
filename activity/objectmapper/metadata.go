package objectmapper

import (
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct{}

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

type Output struct {
	OutVar interface{} `md:"out,required"`
}

func (o *Output) FromMap(values map[string]interface{}) error {
	var err error

	o.OutVar, err = coerce.ToAny(values["out"])
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
