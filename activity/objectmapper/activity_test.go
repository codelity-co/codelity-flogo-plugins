package objectmapper

import (
	"encoding/json"
	"fmt"
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
		"abc": "test", 
		"cde": 1.0,
		"efg": {
			"fgh": "ijk"
		}
	}`
	testMapping := make(map[string]interface{})
	err = json.Unmarshal([]byte(testMappingString), &testMapping)
	assert.Nil(t, err)
	assert.NotNil(t, testMapping)
	fmt.Printf("testMapping: %v: \n", testMapping)

	input := &Input{
		InVar: testMapping,
	}
	fmt.Printf("Input: %v\n", input)

	err = tc.SetInputObject(input)
	assert.Nil(t, err)
	assert.NotNil(t, input)

	done, err = act.Eval(tc)
	assert.True(t, done)
	assert.Nil(t, err)
	output := &Output{}
	err = tc.GetOutputObject(output)
	assert.Nil(t, err)
	fmt.Printf("Output: %v\n", output)

	expectedValueString := `{
		"abc": "test", 
		"cde": 1.0,
		"efg": {
			"fgh": "ijk"
		}
	}`
	expectedValue := make(map[string]interface{})
	err = json.Unmarshal([]byte(expectedValueString), &expectedValue)
	assert.Nil(t, err)
	assert.Equal(t, expectedValue, output.OutVar.(map[string]interface{}))
}
