install:
	go install

test:
	go test ./... -v -race -covermode=atomic

build:
	gox -osarch "linux/amd64 darwin/amd64 windows/amd64"
