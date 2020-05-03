package sample

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/support/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type NatsActivityTestSuite struct {
	suite.Suite
}

func TestNatsActivityTestSuite(t *testing.T) {
	command := exec.Command("docker", "start", "nats-streaming")
	err := command.Run()
	if err != nil {
		fmt.Println(err.Error())
		command := exec.Command(
			"docker", "run", 
			"-p", "4222:4222", 
			"-p", "6222:6222", 
			"-p", "8222:8222", 
			"--name", "nats-streaming", 
			"-d", 
			"nats-streaming",
      "--addr=0.0.0.0",
      "--port=4222",
      "--http_port=8222",
      "--cluster_id=flogo",
			"--store=MEMORY",
		)
		err := command.Run()
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
	}
	suite.Run(t, new(NatsActivityTestSuite))
}

func (suite *NatsActivityTestSuite) SetupTest() {}

func  (suite *NatsActivityTestSuite) TestNatsActivity_Register() {
	t := suite.T()
	
	ref := activity.GetRef(&Activity{})
	act := activity.Get(ref)

	assert.NotNil(t, act)
}

func  (suite *NatsActivityTestSuite) TestNatsActivity_Publish() {
	t := suite.T()
	settings := &Settings{
		ClusterUrls: "nats://localhost:4222",
	}

	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	if err != nil {
		fmt.Println(err)
	}
	assert.Nil(t, err)



	ac := test.NewActivityContext(act.Metadata())
	ac.SetInput("subject", "flogo")
	ac.SetInput("data", []uint8("ABC"))
	
	_, err = act.Eval(ac)
	assert.Nil(t, err)

}
