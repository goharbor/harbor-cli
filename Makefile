BINARY_NAME := harbor-cli

build::
	go build -o ${BINARY_NAME} cmd/harbor/main.go

lint::
	gofmt -s -w .