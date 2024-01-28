make: 
	gofmt -l -s -w .
	go build -o harbor cmd/harbor/main.go