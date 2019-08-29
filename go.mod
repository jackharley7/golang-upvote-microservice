module github.com/jackharley7/golang-upvote-microservice

go 1.12

require (
	cloud.google.com/go v0.38.0
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/googleapis/gax-go v2.0.2+incompatible // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/jackharley7/discussproto v0.0.24
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/stretchr/testify v1.3.0
	github.com/urfave/cli v1.20.0
	github.com/xdg/scram v0.0.0-20180814205039-7eeb5667e42c // indirect
	github.com/xdg/stringprep v1.0.0 // indirect
	go.mongodb.org/mongo-driver v1.0.3
	go.opencensus.io v0.22.0 // indirect
	google.golang.org/api v0.7.0 // indirect
	google.golang.org/grpc v1.21.1
)

// replace github.com/jackharley7/discussproto => ../discuss-proto-pkg
