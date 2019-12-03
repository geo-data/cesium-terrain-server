CESIUM_VERSION:=1.63.1
checkout:=$(or $(FRIENDLY_CHECKOUT),$(shell git branch --show-current))
FRIENDLY_CHECKOUT:=$(or $(FRIENDLY_CHECKOUT),$(shell echo $(checkout) | sed 's/\//-/g'))
GOFILES:=$(shell find . -name '*.go')
GOROOT:=$(or $(GOROOT),/usr/local/go)
GOBIN:=$(or $(GOBIN),/usr/local/go/bin)
GOBINDATA:=$(GOBIN)/go-bindata
DOCKER_REPO:=geo-data/cesium-terrain-server
DOCKER_LOCAL_NAME:=$(DOCKER_REPO):local
LATEST_TAG=$(shell git tag -l --sort=-v:refname | awk 'FNR == 1')
LATEST_TAG_STABLE=$(shell git tag -l --sort=-v:refname | grep -v alpha | awk 'FNR == 1')

install: $(GOFILES) assets/assets.go
	go get ./... && go install ./...

assets/assets.go: $(GOBINDATA) data
	$(GOBINDATA) -ignore \\.gitignore -nocompress -pkg="assets" -o assets/assets.go data

$(GOBINDATA):
	go get -u github.com/go-bindata/go-bindata/...

data/smallterrain-blank.terrain:
	curl --location --progress-bar https://raw.github.com/geo-data/cesium-terrain-builder/master/data/smallterrain-blank.terrain > data/smallterrain-blank.terrain

docker-local: docker/local/cesium-terrain-server-$(FRIENDLY_CHECKOUT).tar.gz docker/local/Cesium-$(CESIUM_VERSION).zip
	docker build --build-arg FRIENDLY_CHECKOUT=$(FRIENDLY_CHECKOUT) --build-arg CESIUM_VERSION=$(CESIUM_VERSION) -t $(DOCKER_LOCAL_NAME) docker

docker/local/Cesium-$(CESIUM_VERSION).zip:
	curl --location --progress-bar https://github.com/AnalyticalGraphicsInc/cesium/releases/download/$(CESIUM_VERSION)/Cesium-$(CESIUM_VERSION).zip > docker/local/Cesium-$(CESIUM_VERSION).zip

docker/local/cesium-terrain-server-$(FRIENDLY_CHECKOUT).tar.gz: $(GOFILES) Makefile
	git archive HEAD --prefix=cesium-terrain-server-$(FRIENDLY_CHECKOUT)/ --format=tar.gz -o docker/local/cesium-terrain-server-$(FRIENDLY_CHECKOUT).tar.gz

docker-tag:
	docker tag $(DOCKER_LOCAL_NAME) $(DOCKER_REPO):$(TO_VERSION)

docker-push:
	docker push $(DOCKER_REPO):$(TO_VERSION)

docker-tag-latest:
	make docker-tag TO_VERSION=latest

docker-push-latest:
	make docker-push TO_VERSION=latest

docker-tag-version:
	make docker-tag TO_VERSION=$(LATEST_TAG)

docker-push-version:
	make docker-push TO_VERSION=$(LATEST_TAG)

docker-tag-stable:
	make docker-tag TO_VERSION=$(LATEST_TAG_STABLE)

docker-push-stable:
	make docker-push TO_VERSION=$(LATEST_TAG_STABLE)


debug:
	echo CESIUM_VERSION: $(CESIUM_VERSION)
	echo checkout: $(checkout)
	echo FRIENDLY_CHECKOUT: $(FRIENDLY_CHECKOUT)
	echo GOFILES: $(GOFILES)
	echo GOBINDATA: $(GOBINDATA)

.PHONY: docker-local install docker-tag docker-push git-latest-tag debug
