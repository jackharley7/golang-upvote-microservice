DOCKER_IMAGE ?=	jackharley7/golang-upvote-microservice

codegen: proto/notification.proto
	GO111MODULE=off go get github.com/mwitkow/go-proto-validators/protoc-gen-govalidators
	GO111MODULE=off go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go install $(shell go list -f '{{ .Dir }}' -m github.com/golang/protobuf)/protoc-gen-go
	protoc \
		--proto_path=${GOPATH}/src \
		--proto_path=${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=plugins=grpc:. -I . proto/notification.proto proto/counts.proto\
		--govalidators_out=.
	gofmt -w proto

.PHONY: mockgen
mockgen: 
	GO111MODULE=off go get github.com/vektra/mockery/.../
	GO111MODULE=on mockery -dir=internal/upvote/ -all -output=internal/upvote/mock/