CESIUM_VERSION:=1.63.1
checkout:=$(or $(FRIENDLY_CHECKOUT),$(shell git branch --show-current))
FRIENDLY_CHECKOUT:=$(or $(FRIENDLY_CHECKOUT),$(shell echo $(checkout) | sed 's/\//-/g'))
GOFILES:=$(shell find . -name '*.go')
GOROOT:=$(or $(GOROOT),/usr/local/go)
GOBIN:=$(or $(GOBIN),/usr/local/go/bin)
GOBINDATA:=$(GOBIN)/go-bindata

install: $(GOFILES) assets/assets.go
	go get ./... && go install ./...

assets/assets.go: $(GOBINDATA) data
	$(GOBINDATA) -ignore \\.gitignore -nocompress -pkg="assets" -o assets/assets.go data

$(GOBINDATA):
	go get -u github.com/go-bindata/go-bindata/...

data/smallterrain-blank.terrain:
	curl --location --progress-bar https://raw.github.com/nmccready/cesium-terrain-builder/master/data/smallterrain-blank.terrain > data/smallterrain-blank.terrain

docker-local: docker/local/cesium-terrain-server-$(FRIENDLY_CHECKOUT).tar.gz docker/local/Cesium-$(CESIUM_VERSION).zip
	docker build --build-arg FRIENDLY_CHECKOUT=$(FRIENDLY_CHECKOUT) --build-arg CESIUM_VERSION=$(CESIUM_VERSION) -t nmccready/cesium-terrain-server:local docker

docker/local/Cesium-$(CESIUM_VERSION).zip:
	curl --location --progress-bar https://github.com/AnalyticalGraphicsInc/cesium/releases/download/$(CESIUM_VERSION)/Cesium-$(CESIUM_VERSION).zip > docker/local/Cesium-$(CESIUM_VERSION).zip

docker/local/cesium-terrain-server-$(FRIENDLY_CHECKOUT).tar.gz: $(GOFILES) Makefile
	git archive HEAD --prefix=cesium-terrain-server-$(FRIENDLY_CHECKOUT)/ --format=tar.gz -o docker/local/cesium-terrain-server-$(FRIENDLY_CHECKOUT).tar.gz

debug:
	echo CESIUM_VERSION: $(CESIUM_VERSION)
	echo checkout: $(checkout)
	echo FRIENDLY_CHECKOUT: $(FRIENDLY_CHECKOUT)
	echo GOFILES: $(GOFILES)
	echo GOBINDATA: $(GOBINDATA)

.PHONY: docker-local install debug
