package minio

import (
	"bytes"
	"encoding/json"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/support/log"

	"github.com/minio/minio-go/v6"
)

func init() {
	_ = activity.Register(&Activity{}, New) //activity.Register(&Activity{}, New) to create instances using factory method 'New'
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

//New optional factory method, should be used if one activity instance per configuration is desired
func New(ctx activity.InitContext) (activity.Activity, error) {

	logger := ctx.Logger()

	var (
		minioClient *minio.Client
	)

	s := &Settings{}
	err := s.FromMap(ctx.Settings())
	if err != nil {
		return nil, err
	}

	logger.Debugf("Setting: %v", s)

	minioClient, err = minio.New(s.Endpoint, s.AccessKey, s.SecretKey, s.EnableSsl)
	if err != nil {
		return nil, err
	}
	
	act := &Activity{
		activitySettings: s,
		logger: logger,
		minioClient: minioClient,
	} 

	return act, nil
}

// Activity is an sample Activity that can be used as a base to create a custom activity
type Activity struct {
	activitySettings *Settings
	logger log.Logger
	minioClient *minio.Client
}

// Metadata returns the activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements api.Activity.Eval - Logs the Message
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	input := Input{}
	err = ctx.GetInputObject(&input)
	if err != nil {
		return true, err
	}

	a.logger.Debugf("Input: %v", input)

	dataBytes := getDataBytes(input.Data)
	switch a.activitySettings.MethodName {
	case "PutObject":
		numberOfBytes, err := a.minioClient.PutObject(a.activitySettings.BucketName, input.ObjectName, bytes.NewReader(dataBytes), int64(len(dataBytes)), minio.PutObjectOptions{})
		if err != nil {
			_ = ctx.SetOutputObject(&Output{
				Status: "ERROR", 
				Result: map[string]interface{}{
					"error": err.Error(),
				},
			})
			return true, err
		}
		if err := ctx.SetOutputObject(&Output{
			Status: "SUCCESS",
			Result: map[string]interface{}{
				"bytes": numberOfBytes,
			},
		}); err != nil {
			return true, err
		}
	}

	return true, nil
}

func getDataBytes(data interface{}) []byte {
	switch value := data.(type) {
	case string:
		return []byte(value)
	case map[string]interface{}:
		dataBytes, _ := json.Marshal(value)
		return dataBytes
	case []byte:
		return value
	}
	return nil
}
