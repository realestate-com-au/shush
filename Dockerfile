FROM golang:1.16-alpine as build

WORKDIR /go/src/github.com/realestate-com-au/shush
COPY . /go/src/github.com/realestate-com-au/shush
RUN go mod vendor && go install

FROM alpine:3.11@sha256:b276d875eeed9c7d3f1cfa7edb06b22ed22b14219a7d67c52c56612330348239

RUN mkdir -p /go/bin

USER nobody
ENV PATH /go/bin:$PATH
COPY --from=build /go/bin/shush /go/bin/shush

ENTRYPOINT ["/go/bin/shush"]
