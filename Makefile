all: build install test

install:
	go install
build:
	go build
test:
	go test -v
deps:
	go get
rmq_test:
	RMQ_INTEGRATION_TEST=True go test -v
