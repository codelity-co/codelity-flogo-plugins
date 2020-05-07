package minio

import (
	"bytes"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"

	"github.com/minio/minio-go/v6"
)

func init() {
	_ = activity.Register(&Activity{}) //activity.Register(&Activity{}, New) to create instances using factory method 'New'
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

//New optional factory method, should be used if one activity instance per configuration is desired
func New(ctx activity.InitContext) (activity.Activity, error) {

	logger := ctx.Logger()

	var (
		minioClient *minio.Client
	)

	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
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

	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return true, err
	}

	a.logger.Debugf("ObjectName: %v", input.ObjectName)
	a.logger.Debugf("Data: %v", input.Data)

	result := make(map[string]interface{})
	output := &Output{}
	switch a.activitySettings.MethodName {
	case "PutObject":
		dataBytes := getDataBytes(a.activitySettings.DataType, input.Data)
		numberOfBytes, err := a.minioClient.PutObject(a.activitySettings.BucketName, input.ObjectName, bytes.NewReader(dataBytes), int64(len(dataBytes)), minio.PutObjectOptions{})
		if err != nil {
			return true, err
		}
		result["bytes"] = numberOfBytes
		output.Status = "SUCCESS"
		output.Result = result
	}
	
	err = ctx.SetOutputObject(output)
	if err != nil {
		return true, err
	}

	return true, nil
}

func getDataBytes(dataType string, data interface{}) []uint8 {
	switch dataType {
	case "string":
		return []uint8(data.(string))
	}
	return nil
}