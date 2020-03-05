package resource

import (
	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/eks"
)

const (
	callbackDelay int64 = 120
)

func errorEvent(model *Model, err error) handler.ProgressEvent {
	errorType := cloudformation.HandlerErrorCodeGeneralServiceException
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case eks.ErrCodeResourceLimitExceededException:
			errorType = cloudformation.HandlerErrorCodeServiceLimitExceeded
		case eks.ErrCodeInvalidParameterException:
			errorType = cloudformation.HandlerErrorCodeInvalidRequest
		case eks.ErrCodeUnsupportedAvailabilityZoneException:
			errorType = cloudformation.HandlerErrorCodeInvalidRequest
		case eks.ErrCodeNotFoundException:
			errorType = cloudformation.HandlerErrorCodeNotFound
		case eks.ErrCodeResourceInUseException:
			errorType = cloudformation.HandlerErrorCodeAlreadyExists
		}
	}
	return handler.ProgressEvent{
		OperationStatus:  handler.Failed,
		HandlerErrorCode: errorType,
		Message:          err.Error(),
		ResourceModel:    model,
	}
}

func successEvent(model *Model) handler.ProgressEvent {
	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		ResourceModel:   model,
	}
}

func inProgressEvent(model *Model, message string, opComplete bool) handler.ProgressEvent {
	return handler.ProgressEvent{
		OperationStatus:      handler.InProgress,
		ResourceModel:        model,
		Message:              message,
		CallbackContext:      map[string]interface{}{"ClusterName": model.Name, "OpComplete": opComplete},
		CallbackDelaySeconds: callbackDelay,
	}
}
