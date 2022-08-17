FROM golang:1.18-alpine as build

WORKDIR /go/src/github.com/realestate-com-au/shush
COPY . /go/src/github.com/realestate-com-au/shush
RUN go mod vendor && go install

FROM alpine:3.16.2@sha256:bc41182d7ef5ffc53a40b044e725193bc10142a1243f395ee852a8d9730fc2ad

RUN mkdir -p /go/bin

USER nobody
ENV PATH /go/bin:$PATH
COPY --from=build /go/bin/shush /go/bin/shush

ENTRYPOINT ["/go/bin/shush"]
