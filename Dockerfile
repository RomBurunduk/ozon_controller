FROM golang:alpine

WORKDIR /opt/app

COPY go.sum go.sum
COPY go.mod go.mod
RUN go mod download



COPY . .
RUN go install ./internal/grpc
ENTRYPOINT /go/bin/grpc


EXPOSE 50051
