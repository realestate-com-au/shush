install:
	go install

test:
	go test ./... -timeout 60s -v -race -covermode=atomic

build:
	gox -osarch "linux/amd64 darwin/amd64 windows/amd64"
