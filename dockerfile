FROM golang:latest AS build_base
COPY ./pkg /go/src/CRUD/blogserver/pkg
COPY ./internal /go/src/CRUD/blogserver/internal

WORKDIR /go/src/CRUD/blogserver

# Force the go compiler to use modules
ENV GO111MODULE=on

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

# This image builds the server
FROM build_base AS server_builder
# Here we copy the rest of the source code
COPY . .
# And compile the project
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -a -tags netgo -ldflags '-w -extldflags "-static"' ./cmd/upvote

FROM alpine:latest
RUN apk --no-cache add ca-certificates

# FOR DEV ONLY

COPY --from=server_builder /go/bin/upvote /bin/upvote
COPY ./gcp-creds* /bin/gcp-creds

RUN GRPC_HEALTH_PROBE_VERSION=v0.2.0 && \
  wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
  chmod +x /bin/grpc_health_probe

CMD ["./bin/upvote"]