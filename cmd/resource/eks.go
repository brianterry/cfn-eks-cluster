package resource

import (
	"errors"
	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"math/rand"
	"time"
)

const (
	generatedClusterNameSuffixLength = 8
	generatedClusterNamePrefix       = "EKS-"
)

func stabilize(svc eksiface.EKSAPI, model *Model, desiredState string, opComplete bool) handler.ProgressEvent {
	input := &eks.DescribeClusterInput{Name: model.Name}
	response, err := svc.DescribeCluster(input)
	if err != nil {
		if resourceNotFound(err) == true && desiredState == "DELETED" {
			return successEvent(model)
		}
		return errorEvent(model, err)
	}
	describeClusterToModel(*response.Cluster, model)
	if *response.Cluster.Status == desiredState {
		return successEvent(model)
	}
	if *response.Cluster.Status == "FAILED" {
		return errorEvent(model, errors.New("cluster status is FAILED"))
	}
	return inProgressEvent(model, "cluster "+*response.Cluster.Status, opComplete)
}

func createCluster(svc eksiface.EKSAPI, model *Model, callbackContext map[string]interface{}) handler.ProgressEvent {
	if callbackContext != nil {
		return stabilize(svc, model, "ACTIVE", true)
	}
	model.Name = generateClusterName(model.Name)
	input := &eks.CreateClusterInput{
		Name: model.Name,
		ResourcesVpcConfig: &eks.VpcConfigRequest{
			SecurityGroupIds: aws.StringSlice(model.ResourcesVpcConfig.SecurityGroupIds),
			SubnetIds:        aws.StringSlice(model.ResourcesVpcConfig.SubnetIds),
		},
		RoleArn: model.RoleArn,
		Version: model.Version,
	}
	response, err := svc.CreateCluster(input)
	if err != nil {
		return errorEvent(model, err)
	}
	describeClusterToModel(*response.Cluster, model)
	return inProgressEvent(model, "Cluster creation initiated", true)
}

func describeCluster(svc eksiface.EKSAPI, model *Model) handler.ProgressEvent {
	input := &eks.DescribeClusterInput{Name: model.Name}
	response, err := svc.DescribeCluster(input)
	if err != nil {
		return errorEvent(model, err)
	}
	describeClusterToModel(*response.Cluster, model)
	return successEvent(model)
}

func updateCluster(svc eksiface.EKSAPI, model *Model, callbackContext map[string]interface{}) handler.ProgressEvent {
	if callbackContext != nil {
		opComplete := callbackContext["OpComplete"].(bool)
		progress := stabilize(svc, model, "ACTIVE", opComplete)
		if progress.OperationStatus == handler.Success && opComplete == true {
			return progress
		}
	}
	input := &eks.UpdateClusterConfigInput{
		Name: model.Name,
		ResourcesVpcConfig: &eks.VpcConfigRequest{
			SecurityGroupIds: aws.StringSlice(model.ResourcesVpcConfig.SecurityGroupIds),
			SubnetIds:        aws.StringSlice(model.ResourcesVpcConfig.SubnetIds),
		},
	}
	_, err := svc.UpdateClusterConfig(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == eks.ErrCodeResourceInUseException {
				return inProgressEvent(model, aerr.Error(), false)
			}
		}
		return errorEvent(model, err)
	}
	return inProgressEvent(model, "Cluster update initiated", true)
}

func deleteCluster(svc eksiface.EKSAPI, model *Model, callbackContext map[string]interface{}) handler.ProgressEvent {
	if callbackContext != nil {
		opComplete := callbackContext["OpComplete"].(bool)
		progress := stabilize(svc, model, "DELETED", opComplete)
		if progress.OperationStatus == handler.Success && opComplete == true {
			return progress
		}
	}
	input := &eks.DeleteClusterInput{
		Name: model.Name,
	}
	_, err := svc.DeleteCluster(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == eks.ErrCodeResourceInUseException {
				return inProgressEvent(model, aerr.Error(), false)
			}
		}
		return errorEvent(model, err)
	}
	return inProgressEvent(model, "Cluster deletion initiated", true)
}

func listClusters(svc eksiface.EKSAPI) handler.ProgressEvent {
	response, err := svc.ListClusters(&eks.ListClustersInput{})
	if err != nil {
		return errorEvent(nil, err)
	}
	models := make([]interface{}, 1)
	for _, m := range response.Clusters {
		model := &Model{Name: m}
		p := describeCluster(svc, model)
		if p.OperationStatus == handler.Failed {
			return p
		}
		models = append(models, model)
	}
	return handler.ProgressEvent{
		ResourceModels:  models,
		OperationStatus: handler.Success,
	}
}

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

func resourceNotFound(err error) bool {
	if aerr, ok := err.(awserr.Error); ok {
		if aerr.Code() == eks.ErrCodeResourceNotFoundException {
			return true
		}
	}
	return false
}
