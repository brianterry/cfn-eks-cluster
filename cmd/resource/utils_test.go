package resource

import (
	"errors"
	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/eks"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestGenerateClusterName(t *testing.T) {
	name := *generateClusterName(nil)
	expectedLength := len(generatedClusterNamePrefix) + generatedClusterNameSuffixLength
	t.Run("cluster name length", func(t *testing.T) {
		assert.Equal(t, len(name), expectedLength)
	})
	anotherName := *generateClusterName(nil)
	t.Run("cluster name uniqueness", func(t *testing.T) {
		assert.NotEqual(t, name, anotherName)
	})
	existingName := "existing-name"
	name = *generateClusterName(&existingName)
	t.Run("dont generate if name is not nil", func(t *testing.T) {
		assert.Equal(t, name, existingName)
	})
	emptyName := ""
	name = *generateClusterName(&emptyName)
	t.Run("generate if name is empty string", func(t *testing.T) {
		assert.NotEqual(t, name, emptyName)
	})
}

func TestDescribeClusterToModel(t *testing.T) {
	cluster := eks.Cluster{
		Arn:                  aws.String("Arn"),
		CertificateAuthority: &eks.Certificate{Data: aws.String("CertificateAuthority")},
		Endpoint:             aws.String("Endpoint"),
		Name:                 aws.String("Name"),
		ResourcesVpcConfig: &eks.VpcConfigResponse{
			SecurityGroupIds: []*string{aws.String("sg-1"), aws.String("sg-2")},
			SubnetIds:        []*string{aws.String("subnet-1"), aws.String("subnet-2")},
		},
		RoleArn: aws.String("RoleArn"),
		Version: aws.String("Version"),
	}
	model := Model{}
	describeClusterToModel(cluster, &model)
	cases := []struct{ Name, Actual, Expected string }{
		{"Arn", *model.Arn, *cluster.Arn},
		{"CertificateAuthorityData", *model.CertificateAuthorityData, *cluster.CertificateAuthority.Data},
		{"Endpoint", *model.Endpoint, *cluster.Endpoint},
		{"Name", *model.Name, *cluster.Name},
		{"security group 1", model.ResourcesVpcConfig.SecurityGroupIds[0], *cluster.ResourcesVpcConfig.SecurityGroupIds[0]},
		{"security group 2", model.ResourcesVpcConfig.SecurityGroupIds[1], *cluster.ResourcesVpcConfig.SecurityGroupIds[1]},
		{"subnet 1", model.ResourcesVpcConfig.SubnetIds[0], *cluster.ResourcesVpcConfig.SubnetIds[0]},
		{"subnet 2", model.ResourcesVpcConfig.SubnetIds[1], *cluster.ResourcesVpcConfig.SubnetIds[1]},
		{"RoleArn", *model.RoleArn, *cluster.RoleArn},
		{"Version", *model.Version, *cluster.Version},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			assert.Equal(t, tc.Expected, tc.Actual)
		})
	}
}

func makeAwsError(code string) error {
	return awserr.New(code, "generated error", errors.New("original error"))
}

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
	progressEvent := inProgressEvent(&Model{Name: &clusterName}, "message")
	assert.Equal(t, handler.InProgress, progressEvent.OperationStatus)
	assert.Equal(t, callbackDelay, progressEvent.CallbackDelaySeconds)
	assert.Equal(t, clusterName, *progressEvent.CallbackContext["ClusterName"].(*string))
	assert.Equal(t, "message", progressEvent.Message)
}

func TestResourceNotFound(t *testing.T) {
	assert.True(t, resourceNotFound(makeAwsError(eks.ErrCodeResourceNotFoundException)))
	assert.False(t, resourceNotFound(makeAwsError(eks.ErrCodeResourceInUseException)))
	assert.False(t, resourceNotFound(errors.New("arbitrary error")))
}
