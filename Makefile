export GOPATH=$(CURDIR)
export GOBIN=$(GOPATH)/bin
cesium_version=$(shell cat $(CURDIR)/docker/cesium-version.txt)

server: src/github.com/gorilla/handlers src/github.com/gorilla/mux src/cesium-terrain-server/server.go bin/go-bindata src/cesium-terrain-server/assets/assets.go
	go build src/cesium-terrain-server/server.go

src/github.com/gorilla/handlers:
	go get github.com/gorilla/handlers

src/github.com/gorilla/mux:
	go get github.com/gorilla/mux

src/cesium-terrain-server/assets.go: bin/go-bindata data
	bin/go-bindata -nocompress -pkg="assets" -o src/cesium-terrain-server/assets/assets.go data

bin/go-bindata: data/smallterrain-blank.terrain
	go get github.com/jteeuwen/go-bindata/... && touch bin/go-bindata

data/smallterrain-blank.terrain:
	wget --no-verbose --directory-prefix=data https://raw.github.com/geo-data/cesium-terrain-builder/master/data/smallterrain-blank.terrain

docker-local: docker/local/cesium-terrain-server.tar.gz docker/local/Cesium-$(cesium_version).zip
	docker build -t geodata/cesium-terrain-server:local docker

docker/local/Cesium-$(cesium_version).zip: docker/cesium-version.txt
	wget --no-verbose --directory-prefix=docker/local https://cesiumjs.org/releases/Cesium-$(cesium_version).zip

docker/local/cesium-terrain-server.tar.gz:
	tar -czvf docker/local/cesium-terrain-server.tar.gz docker/cesium-version.txt src/cesium-terrain-server Makefile

.PHONY: docker-local
