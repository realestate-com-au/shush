FROM golang:1.21-alpine as build

WORKDIR /go/src/github.com/realestate-com-au/shush
COPY . /go/src/github.com/realestate-com-au/shush
RUN go mod vendor && go install

FROM alpine:3.18.2@sha256:82d1e9d7ed48a7523bdebc18cf6290bdb97b82302a8a9c27d4fe885949ea94d1

RUN mkdir -p /go/bin

USER nobody
ENV PATH /go/bin:$PATH
COPY --from=build /go/bin/shush /go/bin/shush

ENTRYPOINT ["/go/bin/shush"]
