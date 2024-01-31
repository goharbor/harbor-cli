GITCOMMIT := $(shell git rev-parse --short=8 HEAD)
PROJECT_PKG = github.com/goharbor/harbor-cli
RELEASE_CHANNEL="edge"
LDFLAGS = "-w -s -X $(PROJECT_PKG)/version.GitCommit=$(GITCOMMIT) -X $(PROJECT_PKG)/version.ReleaseChannel=$(RELEASE_CHANNEL)"
ARCH := amd64
GO_EXE = go

make: 
	gofmt -l -s -w .
	go build -ldflags=${LDFLAGS} -o harbor cmd/harbor/main.go

windows:
	go build -ldflags=${LDFLAGS} -o harbor.exe cmd/harbor/main.go

.PHONY: build-win-amd64
build-win-amd64:  ## build for windows amd64
	CGO_ENABLED=0 GOARCH=$(ARCH) GOOS=windows $(GO_EXE) build -v --ldflags=$(LDFLAGS) \
		-o bin/harbor-windows-$(ARCH).exe ./cmd/harbor/main.go
.PHONY: build-linux-amd64
build-linux-amd64:  ## build for linux amd64
	CGO_ENABLED=0 GOARCH=$(ARCH) GOOS=linux $(GO_EXE) build -v --ldflags=$(LDFLAGS) \
		-o bin/harbor-linux-$(ARCH) ./cmd/harbor/main.go

.PHONY: build-darwin-amd64
build-darwin-amd64:  ## build for darwin amd64
	CGO_ENABLED=0 GOARCH=$(ARCH) GOOS=darwin $(GO_EXE) build -v --ldflags=$(LDFLAGS) \
		-o bin/harbor-darwin-$(ARCH) ./cmd/harbor/main.go

.PHONY: clean
clean:
	rm -rf bin

.PHONY: lint
lint:
	golangci-lint run --timeout 5m