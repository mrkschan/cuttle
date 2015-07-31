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

	// Config limit controller.
	control := config.GetString("limitcontrol", "rps")
	if control == "rps" {
		limit := config.GetUint("rps-limit", 2)
		setLimitController(RPSController{
			limit: limit,
		})
	} else {
		log.Fatal("Unknown limit control: ", control)
	}

	// Config proxy.
	addr := config.GetString("addr", ":8123")
	verbose := config.GetUint("verbose", 0)
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = verbose == 1

	// Starts now.
	startLimitControl()

	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			// Block until permission granted from contoller.
			<-RateLimitFIFO

			// Send request.
			return r, nil
		})

	log.Fatal(http.ListenAndServe(addr, proxy))
}
