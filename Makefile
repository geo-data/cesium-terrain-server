cesium_version:=$(shell cat $(CURDIR)/docker/cesium-version.txt)
checkout:=$(shell cat $(CURDIR)/docker/cts-checkout.txt)
GOFILES:=$(shell find . -name '*.go')

server: server.go .go-bindata assets/assets.go
	go get -d ./... && go build server.go

assets/assets.go: .go-bindata data
	go-bindata -ignore \\.gitignore -nocompress -pkg="assets" -o assets/assets.go data

.go-bindata: data/smallterrain-blank.terrain
	go get github.com/jteeuwen/go-bindata/... && touch .go-bindata

data/smallterrain-blank.terrain:
	curl --location --progress-bar https://raw.github.com/geo-data/cesium-terrain-builder/master/data/smallterrain-blank.terrain > data/smallterrain-blank.terrain

docker-local: docker/local/cesium-terrain-server-$(checkout).tar.gz docker/local/Cesium-$(cesium_version).zip
	docker build -t geodata/cesium-terrain-server:local docker

docker/local/Cesium-$(cesium_version).zip: docker/cesium-version.txt
	curl --location --progress-bar https://cesiumjs.org/releases/Cesium-$(cesium_version).zip > docker/local/Cesium-$(cesium_version).zip

docker/local/cesium-terrain-server-$(checkout).tar.gz: $(GOFILES) docker/cts-checkout.txt docker/cesium-version.txt Makefile
	git archive HEAD --prefix=cesium-terrain-server-$(checkout)/ --format=tar.gz -o docker/local/cesium-terrain-server-$(checkout).tar.gz

.PHONY: docker-local
