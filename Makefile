.PHONY: build test e2e clean

build:
	cfn generate
	env GOOS=linux go build -ldflags="-s -w" -tags="logging callback scheduler" -o bin/handler cmd/main.go

test:
	cfn generate
	env GOOS=linux go build -ldflags="-s -w" -o bin/handler cmd/main.go
	go test ./cmd/resource

e2e:
	aws cloudformation delete-stack --stack-name testeks --region us-west-2
	make build
	cfn submit --verbose --region us-west-2 --set-default
	aws cloudformation create-stack --stack-name testeks --template-body file://test.template.yaml --region us-west-2

clean:
	rm -rf bin
