package resource

import (
	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/service/eks"
)

func Create(req handler.Request, _ *Model, model *Model) (handler.ProgressEvent, error) {
	return createCluster(eks.New(req.Session), model, req.CallbackContext), nil
}

func Read(req handler.Request, _ *Model, model *Model) (handler.ProgressEvent, error) {
	return describeCluster(eks.New(req.Session), model), nil
}

func Update(req handler.Request, _ *Model, model *Model) (handler.ProgressEvent, error) {
	return updateCluster(eks.New(req.Session), model, req.CallbackContext), nil
}

func Delete(req handler.Request, _ *Model, model *Model) (handler.ProgressEvent, error) {
	return deleteCluster(eks.New(req.Session), model, req.CallbackContext), nil
}

func List(req handler.Request, _ *Model, _ *Model) (handler.ProgressEvent, error) {
	return listClusters(eks.New(req.Session)), nil
}
