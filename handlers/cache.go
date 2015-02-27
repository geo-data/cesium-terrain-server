package handlers

import (
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/geo-data/cesium-terrain-server/log"
	"net/http"
	"net/url"
)

type Cache struct {
	mc      *memcache.Client
	handler http.Handler
}

func NewCache(connstr string, handler http.Handler) http.Handler {
	return &Cache{
		mc:      memcache.New(connstr),
		handler: handler,
	}
}

func (this *Cache) generateKey(r *http.Request) string {
	var u *url.URL
	if referer, ok := r.Header["Referer"]; ok {
		u, _ = url.Parse(referer[0])
	} else {
		// Copy the request URL
		u, _ = url.Parse(r.URL.String())
	}

	return "tiles" + u.RequestURI()
}

func (this *Cache) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rec := NewRecorder()

	// Write to both the recorder and original writer
	tee := MultiWriter(w, rec)
	this.handler.ServeHTTP(tee, r)

	key := this.generateKey(r)
	log.Debug(fmt.Sprintf("setting key: %s", key))
	if err := this.mc.Set(&memcache.Item{Key: key, Value: rec.Body.Bytes()}); err != nil {
		log.Err(err.Error())
	}

	return
}
