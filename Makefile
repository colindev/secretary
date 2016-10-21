GOROOT := $(GOROOT)
GO := $(GOROOT)/bin/go
PWD := $(PWD)
TAG := `git describe --tags | cut -d '-' -f 1 `.`git rev-parse --short HEAD`
DATETIME := `TZ=$(TZ) date +%Y%m%d.%H%M%S`

APP := secretary
DESCRIPTION := 'secretary - 仿 cron 排程工具'
SERVICE_FILE := /lib/systemd/system/$(APP).service

test:
	$(GO) test -v ./process

build: test
	$(GO) build -ldflags "-X main.Version=$(TAG) -X main.CompileDate=$(DATETIME)" -a -o ./$(APP)

deploy:
	@read -r -p "deoloy dist: ex. user@127.0.0.1:/dir: " DIST; if [ -n "$$DIST" ] ; then scp -r ./$(APP) ./Makefile ./custom/{.env.sample,schedule.sample} ./util $$DIST; fi

service:
	mkdir -p ./custom/util/systemd
	cat util/systemd/service.temp |\
		sed 's@{CMD}@'"$(APP)"'@g' |\
		sed 's/{DESCRIPTION}/'"$(DESCRIPTION)"'/' |\
		sed "s/{USER}/$(USER)/g" |\
		sed 's@{PWD}@'"$(PWD)"'@g' > custom/util/systemd/$(APP).service

install-service:
	cp custom/util/systemd/$(APP).service $(SERVICE_FILE)
	systemctl daemon-reload
	systemctl enable $(APP)

uninstall-service:
	systemctl stop $(APP)
	systemctl disable $(APP)
	rm -i $(SERVICE_FILE)
	systemctl daemon-reload
