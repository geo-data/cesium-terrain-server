export GOPATH=$(CURDIR)
export GOBIN=$(GOPATH)/bin

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
	wget --directory-prefix=data https://raw.github.com/geo-data/cesium-terrain-builder/master/data/smallterrain-blank.terrain
