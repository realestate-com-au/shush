FROM alpine:3.3

RUN apk --no-cache add git go

ENV GOPATH /go
ENV GOBIN $GOPATH/bin
ENV PATH $GOBIN:$PATH

WORKDIR /go/src/github.com/realestate-com-au/shush
COPY . /go/src/github.com/realestate-com-au/shush
ENV GOPATH /go/src/github.com/realestate-com-au/shush/Godeps/_workspace:$GOPATH
RUN go get .

ENTRYPOINT ["/go/bin/shush"]
