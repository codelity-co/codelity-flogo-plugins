package nats

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"testing"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/support"
	"github.com/project-flogo/core/support/test"
	"github.com/project-flogo/core/trigger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type NatsTriggerTestSuite struct {
	suite.Suite
	testConfig string
}

func TestNatsTriggerTestSuite(t *testing.T) {
	command := exec.Command("docker", "start", "mqtt")
	err := command.Run()
	if err != nil {
		fmt.Println(err.Error())
		command := exec.Command("docker", "run", "-p", "1883:1883", "-p", "9001:9001", "--name", "mqtt", "-d", "eclipse-mosquitto")
		err := command.Run()
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
	}
	suite.Run(t, new(NatsTriggerTestSuite))
}

func (suite *NatsTriggerTestSuite) SetupTest() {
	suite.testConfig = `{
		"id": "nats-trigger",
		"ref": "github.com/codelity-co/codelity-flogo-plugins/nats/trigger",
		"settings": {
			"clusterUrls": "nats://localhost:4222"
		},
		"handlers": [
			{
				"settings": {
					"async": true,
					"subject": "flogo"
				},
				"action": {
					"id": "dummy"
				}
			}
		]
	}`
}

func (suite *NatsTriggerTestSuite) TestNatsTrigger_Register() {
	ref := support.GetRef(&Trigger{})
	f := trigger.GetFactory(ref)
	assert.NotNil(suite.T(), f)
}

func (suite *NatsTriggerTestSuite) TestNatsTrigger_Initialize() {
	t := suite.T()
	f := &Factory{}
	config := &trigger.Config{}

	err := json.Unmarshal([]byte(suite.testConfig), config)
	assert.Nil(t, err)

	actions := map[string]action.Action{"dummy": test.NewDummyAction(func() {
	})}

	trg, err := test.InitTrigger(f, config, actions)

	assert.Nil(t, err)
	assert.NotNil(t, trg)
	err = trg.Start()
	assert.Nil(t, err)
	err = trg.Stop()
	assert.Nil(t, err)
}
