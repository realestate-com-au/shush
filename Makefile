install:
	go install

build:
	gox -osarch "linux/amd64 darwin/amd64 windows/amd64"

test:
	go test ./...
