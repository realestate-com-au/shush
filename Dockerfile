FROM golang:1.13-alpine

WORKDIR /go/src/github.com/realestate-com-au/shush
COPY . /go/src/github.com/realestate-com-au/shush
RUN go install

ENTRYPOINT ["/go/bin/shush"]
