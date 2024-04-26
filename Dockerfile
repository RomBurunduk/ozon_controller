FROM golang:alpine

WORKDIR /opt/app

COPY go.sum go.sum
COPY go.mod go.mod
RUN go mod download

ENV DBNAME=test \
    HOST=localhost \
    PASSWORD=test \
    PORT=5432

COPY . .
RUN go install ./internal/grpc
ENTRYPOINT /go/bin/grpc


EXPOSE 9094
EXPOSE 9095
