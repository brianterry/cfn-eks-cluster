package resource

import (
	"errors"
	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrorEvent(t *testing.T) {
	cases := []struct {
		Err          error
		ExpectedType string
	}{
		{errors.New("arbitrary error"), cloudformation.HandlerErrorCodeGeneralServiceException},
		{makeAwsError(eks.ErrCodeResourceLimitExceededException), cloudformation.HandlerErrorCodeServiceLimitExceeded},
		{makeAwsError(eks.ErrCodeInvalidParameterException), cloudformation.HandlerErrorCodeInvalidRequest},
		{makeAwsError(eks.ErrCodeUnsupportedAvailabilityZoneException), cloudformation.HandlerErrorCodeInvalidRequest},
		{makeAwsError(eks.ErrCodeNotFoundException), cloudformation.HandlerErrorCodeNotFound},
		{makeAwsError(eks.ErrCodeResourceInUseException), cloudformation.HandlerErrorCodeAlreadyExists},
		{makeAwsError("arbitrary aws error"), cloudformation.HandlerErrorCodeGeneralServiceException},
	}
	for _, tc := range cases {
		t.Run(tc.Err.Error(), func(t *testing.T) {
			progressEvent := errorEvent(&Model{}, tc.Err)
			assert.Equal(t, handler.Failed, progressEvent.OperationStatus)
			assert.Equal(t, tc.ExpectedType, progressEvent.HandlerErrorCode)
			assert.Equal(t, tc.Err.Error(), progressEvent.Message)
		})
	}
}

func TestSuccessEvent(t *testing.T) {
	progressEvent := successEvent(&Model{})
	assert.Equal(t, handler.Success, progressEvent.OperationStatus)
	assert.Equal(t, int64(0), progressEvent.CallbackDelaySeconds)
	assert.Nil(t, progressEvent.CallbackContext)
}

func TestInProgressEvent(t *testing.T) {
	clusterName := "clusterName"
	progressEvent := inProgressEvent(&Model{Name: &clusterName}, "message", true)
	assert.Equal(t, handler.InProgress, progressEvent.OperationStatus)
	assert.Equal(t, callbackDelay, progressEvent.CallbackDelaySeconds)
	assert.Equal(t, clusterName, *progressEvent.CallbackContext["ClusterName"].(*string))
	assert.Equal(t, "message", progressEvent.Message)
}
