SHELL := /bin/bash

# Makefile
APPNAME=riotpot
DOCKER=build/docker/
PLUGINS_DIR=pkg/plugin
EXCLUDE_PLUGINS=

# docker cmd below
.PHONY:  first build-docker-first docker-build-doc docker-doc-up up down up-all build build-plugins build-all build-ui statik
first:
	sudo npm install n -g
	sudo n stable
#安装npm
build-docker-first:
	sudo apt-get remove docker docker-engine docker.io containerd runc
	sudo apt update
	sudo apt-get install ca-certificates curl gnupg lsb-release
	curl -fsSL http://mirrors.aliyun.com/docker-ce/linux/ubuntu/gpg | sudo apt-key add –
	sudo add-apt-repository "deb [arch=amd64] http://mirrors.aliyun.com/docker-ce/linux/ubuntu $(lsb_release -cs) stable"
	sudo apt-get install docker-ce docker-ce-cli containerd.io
	systemctl start docker
	sudo apt-get -y install apt-transport-https ca-certificates curl software-properties-common
	service docker restart
#安装docker
docker-build-doc:
	docker build -f $(DOCKER)Dockerfile.documentation . -t $(APPNAME)/v1
docker-doc-up: docker-build-doc
	docker run -p 6060:6060 -it $(APPNAME)/v1
up:
	docker-compose -p riotpot -f ${DOCKER}docker-compose.yaml up -d --build
down:
	docker-compose -p riotpot -f ${DOCKER}docker-compose.yaml down -v
up-all:
	make docker-build-doc
	make docker-doc-up
	make up
#docker系统的创建
build:
	@go build -gcflags='all=-N -l' -o ./bin/ ./cmd/riotpot/.
	@echo "Finished building Binary"
build-plugins: $(PLUGINS_DIR)/*
	@IFS=' ' read -r -a exclude <<< "${EXCLUDE_PLUGINS}"; \
	for folder in $^ ; do \
		result=$${folder%%+(/)}; \
		result=$${result##*/}; \
		result=$${result:-/}; \
		if ! [[ $${exclude[*]} =~ "$${result}" ]]; then \
			go build -buildmode=plugin --mod=mod -gcflags='all=-N -l' -o bin/plugins/$${result}.so $${folder}/*.go; \
		fi \
	done
	@echo "Finished building plugins"
build-ui:
	@npm --prefix=./ui run build
	@echo "Finished building UI"
build-all: \
	build \
	build-plugins
statik:
	@statik -src=/api/swagger
