FROM golang:1.16-alpine as build

WORKDIR /go/src/github.com/realestate-com-au/shush
COPY . /go/src/github.com/realestate-com-au/shush
RUN go mod vendor && go install

FROM alpine:3.13@sha256:e1c082e3d3c45cccac829840a25941e679c25d438cc8412c2fa221cf1a824e6a

RUN mkdir -p /go/bin

USER nobody
ENV PATH /go/bin:$PATH
COPY --from=build /go/bin/shush /go/bin/shush

ENTRYPOINT ["/go/bin/shush"]
