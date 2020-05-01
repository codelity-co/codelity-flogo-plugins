package sample

import (
	"testing"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/support/test"
	"github.com/stretchr/testify/assert"
)

type MinioActivityTestSuite struct {
	suite.Suite
}

func TestMinioActivityTestSuite(t *testing.T) {
	suite.Run(t, new(MinioActivityTestSuite))
}

func (suite *MinioActivityTestSuite) SetupTest() {}

func (suite *NatsTriggerTestSuite) TestMinioActivity_Register() {

	ref := activity.GetRef(&Activity{})
	act := activity.Get(ref)

	assert.NotNil(suite.T(), act)
}

func (suite *NatsTriggerTestSuite) TestMinioActivity_Eval() {

}
