install:
	go install

build:
	gox -osarch "linux/amd64 darwin/amd64"
