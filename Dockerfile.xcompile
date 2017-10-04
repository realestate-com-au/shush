FROM golang:1.8@sha256:32c769bf92205580d6579d5b93c3c705f787f6c648105f00bb88a35024c7f8e4

RUN go get github.com/mitchellh/gox

ENV CGO_ENABLED=0

COPY . /go/src/github.com/realestate-com-au/shush
WORKDIR /go/src/github.com/realestate-com-au/shush
RUN gox -osarch "linux/amd64 darwin/amd64 windows/amd64"
