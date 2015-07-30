package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/elazarl/goproxy"
)

func main() {
	log.SetFormatter(&log.TextFormatter{})

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	log.Fatal(http.ListenAndServe(":8123", proxy))
}
