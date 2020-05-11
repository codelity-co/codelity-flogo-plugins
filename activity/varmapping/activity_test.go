package varmapping

import (
	"testing"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/support/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type VarMappingActivityTestSuite struct {
	suite.Suite
}


func TestVarMappingActivityTestSuite(t *testing.T) {
	suite.Run(t, new(VarMappingActivityTestSuite))
}

func (suite *VarMappingActivityTestSuite) TestVarMappingActivity_Register() {
	t := suite.T()

	ref := activity.GetRef(&Activity{})
	act := activity.Get(ref)

	assert.NotNil(t, act)
}

func (suite *VarMappingActivityTestSuite) TestVarMappingActivity_Eval() {
	t := suite.T()

	act := &Activity{}
	tc := test.NewActivityContext(act.Metadata())
	var testValue interface{} = map[string]interface{}{
		"abc": "test", 
		"efg": map[string]interface{}{
			"fgh": "ijk",
		},
	}
	input := &Input{
		InVar: testValue,
	}
	err := tc.SetInputObject(input)
	assert.Nil(t, err)

	done, err := act.Eval(tc)
	assert.True(t, done)
	assert.Nil(t, err)

	output := &Output{}
	err = tc.GetOutputObject(output)
	assert.Nil(t, err)

	assert.Equal(t, testValue, output.OutVar)
}
