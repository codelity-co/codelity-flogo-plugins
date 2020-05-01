package sample

import (
	"testing"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/support/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MinioActivityTestSuite struct {
	suite.Suite
}

func TestMinioActivityTestSuite(t *testing.T) {
	suite.Run(t, new(MinioActivityTestSuite))
}

func (suite *MinioActivityTestSuite) SetupTest() {}

func (suite *MinioActivityTestSuite) TestMinioActivity_Register() {

	ref := activity.GetRef(&Activity{})
	act := activity.Get(ref)

	assert.NotNil(suite.T(), act)
}

func (suite *MinioActivityTestSuite) TestMinioActivity_Settings() {
	t := suite.T()

	settings := &Settings{
		Endpoint: "miio:9000",
		AccessKey: "minio",
		SecretKey: "minio123", 
		EnableSsl: false,
	}
	
	iCtx := test.NewActivityInitContext(settings, nil)
	_, err := New(iCtx)
	assert.NotNil(t, err)

	settings = &Settings{
		Endpoint: "minio:9000",
		AccessKey: "minio123",
		SecretKey: "minio123", 
		EnableSsl: false,
	}

	iCtx = test.NewActivityInitContext(settings, nil)
	_, err = New(iCtx)
	assert.NotNil(t, err)

	settings = &Settings{
		Endpoint: "minio:9000",
		AccessKey: "minio",
		SecretKey: "minio123", 
		EnableSsl: false,
	}

	iCtx = test.NewActivityInitContext(settings, nil)
	_, err = New(iCtx)
	assert.Nil(t, err)
}

func (suite *MinioActivityTestSuite) TestMinioActivity_PutObject() {
	t := suite.T()
	
	settings := &Settings{
		Endpoint: "minio:9000",
		AccessKey: "minio",
		SecretKey: "minio123", 
		EnableSsl: false,
	}
	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(act.Metadata())
	tc.SetInput("method", "PutObject")
	params := make(map[string]interface{})
	params["bucketName"] = "flogo"
	params["objectName"] = "inbox/testing.json"
	params["stringData"] = "{\"abc\": \"123\"}"
	tc.SetInput("params", params)
	act.Eval(tc)
}
