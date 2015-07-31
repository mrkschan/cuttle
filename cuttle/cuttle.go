package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/cloudflare/conf"
	"github.com/elazarl/goproxy"
)

func main() {
	log.SetFormatter(&log.TextFormatter{})

	config, err := conf.ReadConfigFile("cuttle.conf")
	if err != nil {
		log.Error("Failed to load config from 'cuttle.conf'.")
		log.Fatal(err)
	}
	addr := config.GetString("addr", ":8123")
	verbose := config.GetUint("verbose", 0)

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = verbose == 1
	log.Fatal(http.ListenAndServe(addr, proxy))
}
