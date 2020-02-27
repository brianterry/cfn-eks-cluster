package resource

import (
	"errors"
	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"log"
)

func stabilize(svc *eks.EKS, model *Model, desiredState string) handler.ProgressEvent {
	input := &eks.DescribeClusterInput{Name: model.Name}
	response, err := svc.DescribeCluster(input)
	if err != nil {
		log.Printf("describe error: %s", err.Error())
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
	return inProgressEvent(model, "cluster "+*response.Cluster.Status)
}

func createCluster(svc *eks.EKS, model *Model) handler.ProgressEvent {
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
	return inProgressEvent(model, "Cluster creation initiated")
}

func updateCluster(svc *eks.EKS, model *Model) handler.ProgressEvent {
	input := &eks.UpdateClusterConfigInput{
		Name: model.Name,
		ResourcesVpcConfig: &eks.VpcConfigRequest{
			SecurityGroupIds: aws.StringSlice(model.ResourcesVpcConfig.SecurityGroupIds),
			SubnetIds:        aws.StringSlice(model.ResourcesVpcConfig.SubnetIds),
		},
	}
	_, err := svc.UpdateClusterConfig(input)
	if err != nil {
		return errorEvent(model, err)
	}
	return inProgressEvent(model, "Cluster update initiated")
}

func deleteCluster(svc *eks.EKS, model *Model) handler.ProgressEvent {
	input := &eks.DeleteClusterInput{
		Name: model.Name,
	}
	_, err := svc.DeleteCluster(input)
	if err != nil {
		return errorEvent(model, err)
	}
	return inProgressEvent(model, "Cluster deletion initiated")
}
