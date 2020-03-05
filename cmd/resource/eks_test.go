package resource

import (
	"errors"
	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var anErr error = errors.New("upstream error")

type mockEKSClient struct {
	eksiface.EKSAPI
	MockCluster       *eks.Cluster
	MockCreateError   error
	MockDescribeError error
	MockUpdateError   error
	MockDeleteError   error
	MockListError     error
}

func (m *mockEKSClient) CreateCluster(input *eks.CreateClusterInput) (*eks.CreateClusterOutput, error) {
	return &eks.CreateClusterOutput{
		Cluster: &eks.Cluster{
			Arn: m.MockCluster.Arn,
			CertificateAuthority: &eks.Certificate{
				Data: m.MockCluster.CertificateAuthority.Data,
			},
			ClientRequestToken: m.MockCluster.ClientRequestToken,
			CreatedAt:          &time.Time{},
			Endpoint:           m.MockCluster.Endpoint,
			Identity: &eks.Identity{
				Oidc: &eks.OIDC{
					Issuer: m.MockCluster.Identity.Oidc.Issuer,
				},
			},
			Logging: &eks.Logging{
				ClusterLogging: []*eks.LogSetup{{
					Enabled: aws.Bool(false),
					Types:   []*string{},
				}},
			},
			Name:            input.Name,
			PlatformVersion: m.MockCluster.PlatformVersion,
			ResourcesVpcConfig: &eks.VpcConfigResponse{
				ClusterSecurityGroupId: m.MockCluster.ResourcesVpcConfig.ClusterSecurityGroupId,
				EndpointPrivateAccess:  input.ResourcesVpcConfig.EndpointPrivateAccess,
				EndpointPublicAccess:   input.ResourcesVpcConfig.EndpointPublicAccess,
				PublicAccessCidrs:      input.ResourcesVpcConfig.PublicAccessCidrs,
				SecurityGroupIds:       input.ResourcesVpcConfig.SecurityGroupIds,
				SubnetIds:              input.ResourcesVpcConfig.SubnetIds,
				VpcId:                  m.MockCluster.ResourcesVpcConfig.VpcId,
			},
			RoleArn: input.RoleArn,
			Status:  m.MockCluster.Status,
			Tags:    input.Tags,
			Version: input.Version,
		},
	}, m.MockCreateError
}

func (m *mockEKSClient) UpdateClusterConfig(input *eks.UpdateClusterConfigInput) (*eks.UpdateClusterConfigOutput, error) {
	return &eks.UpdateClusterConfigOutput{
		Update: &eks.Update{
			CreatedAt: &time.Time{},
			Errors:    nil,
			Id:        aws.String("Id"),
			Params:    []*eks.UpdateParam{},
			Status:    aws.String(eks.ClusterStatusUpdating),
			Type:      aws.String(""),
		},
	}, m.MockUpdateError
}

func (m *mockEKSClient) DeleteCluster(input *eks.DeleteClusterInput) (*eks.DeleteClusterOutput, error) {
	return &eks.DeleteClusterOutput{
		Cluster: &eks.Cluster{
			Arn: m.MockCluster.Arn,
			CertificateAuthority: &eks.Certificate{
				Data: m.MockCluster.CertificateAuthority.Data,
			},
			ClientRequestToken: m.MockCluster.ClientRequestToken,
			CreatedAt:          &time.Time{},
			Endpoint:           m.MockCluster.Endpoint,
			Identity: &eks.Identity{
				Oidc: &eks.OIDC{
					Issuer: m.MockCluster.Identity.Oidc.Issuer,
				},
			},
			Logging: &eks.Logging{
				ClusterLogging: []*eks.LogSetup{{
					Enabled: aws.Bool(false),
					Types:   []*string{},
				}},
			},
			Name:            input.Name,
			PlatformVersion: m.MockCluster.PlatformVersion,
			ResourcesVpcConfig: &eks.VpcConfigResponse{
				ClusterSecurityGroupId: m.MockCluster.ResourcesVpcConfig.ClusterSecurityGroupId,
				EndpointPrivateAccess:  m.MockCluster.ResourcesVpcConfig.EndpointPrivateAccess,
				EndpointPublicAccess:   m.MockCluster.ResourcesVpcConfig.EndpointPublicAccess,
				PublicAccessCidrs:      m.MockCluster.ResourcesVpcConfig.PublicAccessCidrs,
				SecurityGroupIds:       m.MockCluster.ResourcesVpcConfig.SecurityGroupIds,
				SubnetIds:              m.MockCluster.ResourcesVpcConfig.SubnetIds,
				VpcId:                  m.MockCluster.ResourcesVpcConfig.VpcId,
			},
			RoleArn: m.MockCluster.RoleArn,
			Status:  m.MockCluster.Status,
			Tags:    m.MockCluster.Tags,
			Version: m.MockCluster.Version,
		},
	}, m.MockDeleteError
}

func (m *mockEKSClient) DescribeCluster(input *eks.DescribeClusterInput) (*eks.DescribeClusterOutput, error) {
	return &eks.DescribeClusterOutput{
		Cluster: &eks.Cluster{
			Arn: m.MockCluster.Arn,
			CertificateAuthority: &eks.Certificate{
				Data: m.MockCluster.CertificateAuthority.Data,
			},
			ClientRequestToken: m.MockCluster.ClientRequestToken,
			CreatedAt:          &time.Time{},
			Endpoint:           m.MockCluster.Endpoint,
			Identity: &eks.Identity{
				Oidc: &eks.OIDC{
					Issuer: m.MockCluster.Identity.Oidc.Issuer,
				},
			},
			Logging: &eks.Logging{
				ClusterLogging: []*eks.LogSetup{{
					Enabled: aws.Bool(false),
					Types:   []*string{},
				}},
			},
			Name:            input.Name,
			PlatformVersion: m.MockCluster.PlatformVersion,
			ResourcesVpcConfig: &eks.VpcConfigResponse{
				ClusterSecurityGroupId: m.MockCluster.ResourcesVpcConfig.ClusterSecurityGroupId,
				EndpointPrivateAccess:  m.MockCluster.ResourcesVpcConfig.EndpointPrivateAccess,
				EndpointPublicAccess:   m.MockCluster.ResourcesVpcConfig.EndpointPublicAccess,
				PublicAccessCidrs:      m.MockCluster.ResourcesVpcConfig.PublicAccessCidrs,
				SecurityGroupIds:       m.MockCluster.ResourcesVpcConfig.SecurityGroupIds,
				SubnetIds:              m.MockCluster.ResourcesVpcConfig.SubnetIds,
				VpcId:                  m.MockCluster.ResourcesVpcConfig.VpcId,
			},
			RoleArn: m.MockCluster.RoleArn,
			Status:  m.MockCluster.Status,
			Tags:    m.MockCluster.Tags,
			Version: m.MockCluster.Version,
		},
	}, m.MockDescribeError
}

func makeCluster() *eks.Cluster {
	return &eks.Cluster{
		Arn:                  aws.String("MockArn"),
		CertificateAuthority: &eks.Certificate{Data: aws.String("CertificateAuthority")},
		Endpoint:             aws.String("MockEndpoint"),
		Status:               aws.String(eks.ClusterStatusCreating),
		Identity:             &eks.Identity{Oidc: &eks.OIDC{Issuer: aws.String("issuer")}},
		PlatformVersion:      aws.String("1.14.7"),
		ResourcesVpcConfig: &eks.VpcConfigResponse{
			ClusterSecurityGroupId: aws.String("ClusterSecurityGroupId"),
			VpcId:                  aws.String("VpcId"),
			SecurityGroupIds:       []*string{aws.String("sg-1")},
			SubnetIds:              []*string{aws.String("subnet-1")},
			PublicAccessCidrs:      []*string{aws.String("10.0.0.0/16")},
			EndpointPublicAccess:   aws.Bool(true),
			EndpointPrivateAccess:  aws.Bool(false),
		},
		RoleArn: aws.String("RoleArn"),
		Version: aws.String("Version"),
		Tags:    map[string]*string{},
	}
}

func makeModel() *Model {
	return &Model{
		Name:    aws.String("test"),
		RoleArn: aws.String("role"),
		Version: aws.String("1.14"),
		ResourcesVpcConfig: &ResourcesVpcConfig{
			SecurityGroupIds: []string{"sg-1"},
			SubnetIds:        []string{"subnet-1"},
		},
	}
}

func TestCreateCluster(t *testing.T) {
	mockSvc := &mockEKSClient{
		MockCluster: makeCluster(),
	}
	mockSvc.MockCluster.Status = aws.String(eks.ClusterStatusCreating)

	model := makeModel()
	var callbackContext map[string]interface{}
	t.Run("in progress", func(t *testing.T) {
		progress := createCluster(mockSvc, model, callbackContext)
		assert.Equal(t, handler.InProgress, progress.OperationStatus)
	})
	t.Run("aws api error", func(t *testing.T) {
		mockSvc.MockCreateError = awserr.New(eks.ErrCodeResourceInUseException, "mock aws error", anErr)
		progress := createCluster(mockSvc, model, callbackContext)
		assert.Equal(t, handler.Failed, progress.OperationStatus)
	})
	t.Run("success", func(t *testing.T) {
		mockSvc.MockCreateError = nil
		callbackContext = map[string]interface{}{"ClusterName": "test", "OpComplete": true}
		mockSvc.MockCluster.Status = aws.String(eks.ClusterStatusActive)
		progress := createCluster(mockSvc, model, callbackContext)
		assert.Equal(t, handler.Success, progress.OperationStatus)
	})
}

func TestUpdateCluster(t *testing.T) {
	mockSvc := &mockEKSClient{
		MockCluster: makeCluster(),
	}
	mockSvc.MockCluster.Status = aws.String(eks.ClusterStatusUpdating)

	model := makeModel()
	var callbackContext map[string]interface{}
	t.Run("in progress", func(t *testing.T) {
		progress := updateCluster(mockSvc, model, callbackContext)
		assert.Equal(t, handler.InProgress, progress.OperationStatus)
	})
	t.Run("aws api error", func(t *testing.T) {
		mockSvc.MockUpdateError = awserr.New(eks.ErrCodeClientException, "mock aws error", anErr)
		progress := updateCluster(mockSvc, model, callbackContext)
		assert.Equal(t, handler.Failed, progress.OperationStatus)
	})
	t.Run("success", func(t *testing.T) {
		mockSvc.MockDescribeError = nil
		mockSvc.MockCluster.Status = aws.String(eks.ClusterStatusActive)
		callbackContext = map[string]interface{}{"ClusterName": "test", "OpComplete": true}
		progress := updateCluster(mockSvc, model, callbackContext)
		assert.Equal(t, handler.Success, progress.OperationStatus)
	})
	t.Run("update already in progress", func(t *testing.T) {
		callbackContext = nil
		mockSvc.MockUpdateError = awserr.New(eks.ErrCodeResourceInUseException, "mock aws error", anErr)
		progress := updateCluster(mockSvc, model, callbackContext)
		assert.Equal(t, handler.InProgress, progress.OperationStatus)
		assert.Equal(t, false, progress.CallbackContext["OpComplete"].(bool))
	})
}

func TestDeleteCluster(t *testing.T) {
	mockSvc := &mockEKSClient{
		MockCluster: makeCluster(),
	}
	mockSvc.MockCluster.Status = aws.String(eks.ClusterStatusDeleting)

	model := makeModel()
	var callbackContext map[string]interface{}
	t.Run("in progress", func(t *testing.T) {
		progress := deleteCluster(mockSvc, model, callbackContext)
		assert.Equal(t, handler.InProgress, progress.OperationStatus)
	})
	t.Run("aws api error", func(t *testing.T) {
		mockSvc.MockDeleteError = awserr.New(eks.ErrCodeClientException, "mock aws error", anErr)
		progress := deleteCluster(mockSvc, model, callbackContext)
		assert.Equal(t, handler.Failed, progress.OperationStatus)
	})
	t.Run("success", func(t *testing.T) {
		mockSvc.MockDescribeError = awserr.New(eks.ErrCodeResourceNotFoundException, "mock aws error", anErr)
		callbackContext = map[string]interface{}{"ClusterName": "test", "OpComplete": true}
		progress := deleteCluster(mockSvc, model, callbackContext)
		assert.Equal(t, handler.Success, progress.OperationStatus)
	})
	t.Run("update in progress", func(t *testing.T) {
		callbackContext = nil
		mockSvc.MockDeleteError = awserr.New(eks.ErrCodeResourceInUseException, "mock aws error", anErr)
		progress := deleteCluster(mockSvc, model, callbackContext)
		assert.Equal(t, handler.InProgress, progress.OperationStatus)
		assert.Equal(t, false, progress.CallbackContext["OpComplete"].(bool))
	})
}

func TestStabilize(t *testing.T) {
	mockSvc := &mockEKSClient{
		MockCluster: makeCluster(),
	}
	mockSvc.MockCluster.Status = aws.String(eks.ClusterStatusActive)

	model := makeModel()
	t.Run("in progress", func(t *testing.T) {
		progress := stabilize(mockSvc, model, "DELETED", true)
		assert.Equal(t, handler.InProgress, progress.OperationStatus)
	})
	t.Run("active", func(t *testing.T) {
		progress := stabilize(mockSvc, model, eks.ClusterStatusActive, true)
		assert.Equal(t, handler.Success, progress.OperationStatus)
	})
	t.Run("cluster in failed state", func(t *testing.T) {
		mockSvc.MockCluster.Status = aws.String(eks.ClusterStatusFailed)
		progress := stabilize(mockSvc, model, eks.ClusterStatusActive, true)
		assert.Equal(t, handler.Failed, progress.OperationStatus)
	})
	t.Run("deleted", func(t *testing.T) {
		mockSvc.MockDescribeError = awserr.New(eks.ErrCodeResourceNotFoundException, "mock aws error", anErr)
		progress := stabilize(mockSvc, model, "DELETED", true)
		assert.Equal(t, handler.Success, progress.OperationStatus)
	})
	t.Run("aws api error", func(t *testing.T) {
		mockSvc.MockDescribeError = awserr.New(eks.ErrCodeResourceInUseException, "mock aws error", anErr)
		progress := stabilize(mockSvc, model, eks.ClusterStatusActive, true)
		assert.Equal(t, handler.Failed, progress.OperationStatus)
	})
}

func TestDescribeCluster(t *testing.T) {
	mockSvc := &mockEKSClient{
		MockCluster: makeCluster(),
	}
	model := makeModel()
	t.Run("success", func(t *testing.T) {
		progress := describeCluster(mockSvc, model)
		assert.Equal(t, handler.Success, progress.OperationStatus)
	})
	t.Run("aws api err", func(t *testing.T) {
		mockSvc.MockDescribeError = awserr.New(eks.ErrCodeResourceNotFoundException, "mock aws error", anErr)
		progress := describeCluster(mockSvc, model)
		assert.Equal(t, handler.Failed, progress.OperationStatus)
	})
}

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

func TestResourceNotFound(t *testing.T) {
	assert.True(t, resourceNotFound(makeAwsError(eks.ErrCodeResourceNotFoundException)))
	assert.False(t, resourceNotFound(makeAwsError(eks.ErrCodeResourceInUseException)))
	assert.False(t, resourceNotFound(errors.New("arbitrary error")))
}
