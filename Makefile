GOROOT := $(GOROOT)
GO := $(GOROOT)/bin/go
PWD := $(PWD)
TAG := `git describe --tags | cut -d '-' -f 1 `.`git rev-parse --short HEAD`
DATETIME := `TZ=$(TZ) date +%Y%m%d.%H%M%S`

build:
	$(GO) build -ldflags "-X main.Version=$(TAG) -X main.CompileDate=$(DATETIME)" -a -o ./secretary
