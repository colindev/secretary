GOROOT := $(GOROOT)
GO := $(GOROOT)/bin/go
PWD := $(PWD)

docker-img:
	docker build -t centos6-gcc .

build:
	if test -n "$(OS)" ; then \
		$(GO) build -o ./release/secretary.$(OS); \
	else \
		$(GO) build -o ./release/secretary; \
	fi

centos6:
	docker run --rm -v $(PWD):/go-src -v $(GOROOT):/go-root -e "GOROOT=/go-root" -e "OS=centos6" centos6-gcc /go-src/compile
