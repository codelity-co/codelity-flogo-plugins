package objectmapper

import (
	"encoding/json"
	"testing"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/support/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

)

type ObjectMapperActivityTestSuite struct {
	suite.Suite
}

func TestObjectMapperActivityTestSuite(t *testing.T) {
	suite.Run(t, new(ObjectMapperActivityTestSuite))
}

func (suite *ObjectMapperActivityTestSuite) TestObjectMapperActivity_Register() {
	t := suite.T()

	ref := activity.GetRef(&Activity{})
	act := activity.Get(ref)

	assert.NotNil(t, act)
}

func (suite *ObjectMapperActivityTestSuite) TestObjectMapperActivity_Eval() {
	var (
		err error
		done bool
	)

	t := suite.T()
	
	act := &Activity{}
	tc := test.NewActivityContext(act.Metadata())
	testMappingString := `{
		"mapping": {
			"abc": "test", 
			"cde": "=coerce.toInt32(\"1\")",
			"efg": {
				"fgh": "ijk",
			}
		}
	}`
	testMapping := make(map[string]interface{})

	err = json.Unmarshal([]byte(testMappingString), &testMapping)
	assert.Nil(t, err)

	input := &Input{
		Mapping: testMapping,
	}
	err = tc.SetInputObject(input)
	assert.Nil(t, err)

	done, err = act.Eval(tc)
	assert.True(t, done)
	assert.Nil(t, err)

	output := &Output{}
	err = tc.GetOutputObject(output)
	assert.Nil(t, err)
	var expectedValue interface{} = map[string]interface{}{
		"abc": "test", 
		"cde": 1,
		"efg": map[string]interface{}{
			"fgh": "ijk",
		},
	}
	assert.Equal(t, expectedValue, output.OutVar)
}
