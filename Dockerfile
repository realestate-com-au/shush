FROM alpine:3.2

RUN apk add --update git go && \
    rm -rf /var/cache/apk/*

ENV GOPATH /go
ENV GOBIN $GOPATH/bin
ENV PATH $GOBIN:$PATH

WORKDIR /go/src/github.com/realestate-com-au/shush
COPY . /go/src/github.com/realestate-com-au/shush
ENV GOPATH /go/src/github.com/realestate-com-au/shush/Godeps/_workspace:$GOPATH
RUN go get .

ENTRYPOINT ["/go/bin/shush"]
