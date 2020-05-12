package datetimex

import (
	"time"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	_ = function.Register(&fnUtcNano{})
}

type fnUtcNano struct {
}

func init() {
	function.Register(&fnUtcNano{})
}

func (s *fnUtcNano) Name() string {
	return "utcNano"
}

func (s *fnUtcNano) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{}, false
}

func (s *fnUtcNano) Eval(in ...interface{}) (interface{}, error) {
	return time.Now().UTC().UnixNano(), nil
}
