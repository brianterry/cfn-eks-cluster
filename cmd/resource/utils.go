package resource

import (
	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/eks"
	"math/rand"
	"time"
)

func generateClusterName(name *string) *string {
	if name != nil {
		if *name != "" {
			return name
		}
	}
	rand.Seed(time.Now().UnixNano())
	letters := []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, generatedClusterNameSuffixLength)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	generated := generatedClusterNamePrefix + string(b)
	return &generated
}

func describeClusterToModel(cluster eks.Cluster, model *Model) {
	model.Name = cluster.Name
	model.RoleArn = cluster.RoleArn
	model.Version = cluster.Version
	model.ResourcesVpcConfig = &ResourcesVpcConfig{
		SecurityGroupIds: aws.StringValueSlice(cluster.ResourcesVpcConfig.SecurityGroupIds),
		SubnetIds:        aws.StringValueSlice(cluster.ResourcesVpcConfig.SubnetIds),
	}
	model.Arn = cluster.Arn
	model.CertificateAuthorityData = cluster.CertificateAuthority.Data
	model.ClusterSecurityGroupId = cluster.ResourcesVpcConfig.ClusterSecurityGroupId
	model.Endpoint = cluster.Endpoint
}

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

func inProgressEvent(model *Model, message string) handler.ProgressEvent {
	return handler.ProgressEvent{
		OperationStatus:      handler.InProgress,
		ResourceModel:        model,
		Message:              message,
		CallbackContext:      map[string]interface{}{"ClusterName": model.Name},
		CallbackDelaySeconds: callbackDelay,
	}
}

func resourceNotFound(err error) bool {
	if aerr, ok := err.(awserr.Error); ok {
		if aerr.Code() == eks.ErrCodeResourceNotFoundException {
			return true
		}
	}
	return false
}
