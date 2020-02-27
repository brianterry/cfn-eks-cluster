package resource

import (
	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/service/eks"
)

func Create(req handler.Request, _ *Model, model *Model) (handler.ProgressEvent, error) {
	svc := eks.New(req.Session)
	if req.CallbackContext == nil {
		return createCluster(svc, model), nil
	}
	return stabilize(svc, model, "ACTIVE"), nil
}

func Read(req handler.Request, _ *Model, model *Model) (handler.ProgressEvent, error) {
	svc := eks.New(req.Session)
	input := &eks.DescribeClusterInput{
		Name: model.Name,
	}
	response, err := svc.DescribeCluster(input)
	if err != nil {
		return errorEvent(model, err), nil
	}
	describeClusterToModel(*response.Cluster, model)
	return successEvent(model), nil
}

func Update(req handler.Request, _ *Model, model *Model) (handler.ProgressEvent, error) {
	svc := eks.New(req.Session)
	if req.CallbackContext == nil {
		return updateCluster(svc, model), nil
	}
	return stabilize(svc, model, "ACTIVE"), nil
}

func Delete(req handler.Request, _ *Model, model *Model) (handler.ProgressEvent, error) {
	svc := eks.New(req.Session)
	if req.CallbackContext == nil {
		return deleteCluster(svc, model), nil
	}
	return stabilize(svc, model, "DELETED"), nil
}

func List(req handler.Request, _ *Model, _ *Model) (handler.ProgressEvent, error) {
	svc := eks.New(req.Session)
	response, err := svc.ListClusters(&eks.ListClustersInput{})
	if err != nil {
		return errorEvent(nil, err), err
	}
	models := make([]interface{}, 1)
	for _, m := range response.Clusters {
		input := &eks.DescribeClusterInput{
			Name: m,
		}
		describeResponse, err := svc.DescribeCluster(input)
		if err != nil {
			return errorEvent(nil, err), nil
		}
		model := &Model{}
		describeClusterToModel(*describeResponse.Cluster, model)
		models = append(models, model)
	}
	return handler.ProgressEvent{
		ResourceModels:  models,
		OperationStatus: handler.Success,
	}, nil
}
