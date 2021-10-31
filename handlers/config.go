package handlers

import (
	"net/http/httputil"
	"net/url"

	"github.com/xsoroton/autopilot-proxy/store"
	"github.com/xsoroton/autopilot-proxy/util"
)

type Handlers struct {
	Proxy      *httputil.ReverseProxy
	Remote     *url.URL
	CacheStore store.Store
}

func NewHandlersFromEnv() (handlers Handlers) {
	var err error
	remoteHost := util.GetEnv("AUTOPILOT_API", "https://private-1f378a-autopilot.apiary-proxy.com")
	handlers.Remote, err = url.Parse(remoteHost)
	if err != nil {
		panic(err)
	}
	// Note: store.NewMemStore() to use MEM store
	handlers.CacheStore = store.NewRedisStoreFromEnv()
	handlers.Proxy = httputil.NewSingleHostReverseProxy(handlers.Remote)
	return
}
