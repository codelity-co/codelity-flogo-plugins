package datetimex

import (
	"fmt"
	"reflect"
	"github.com/project-flogo/core/data/expression/function"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	function.ResolveAliases()
}

func TestFnUtcNano_Eval(t *testing.T) {
	var in = &fnUtcNano{}
	final, err := in.Eval()
	assert.Nil(t, err)
	assert.IsType(t, "int64", fmt.Sprintf("%v", reflect.TypeOf(final)))
}