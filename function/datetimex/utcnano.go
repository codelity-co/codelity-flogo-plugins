package datetimex

import (
	"time"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	_ = function.Register(&fnUtcNano{})
}

type fnUtcNano struct {}

func (s *fnUtcNano) Name() string {
	return "utcNano"
}

func (s *fnUtcNano) GetCategory() string {
	return "datetimex"
}

func (s *fnUtcNano) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{}, false
}

func (s *fnUtcNano) Eval(params ...interface{}) (interface{}, error) {
	return time.Now().UTC().UnixNano(), nil
}
