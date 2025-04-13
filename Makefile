
build::
	go build -o harbor-cli cmd/harbor/main.go

lint::
	gofmt -s -w .