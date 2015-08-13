package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/elazarl/goproxy"
	"github.com/spf13/viper"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	viper.SetConfigType("yaml")
	viper.SetConfigName("cuttle")
	viper.AddConfigPath("/etc/cuttle/")

	err := viper.ReadInConfig()
	if err != nil {
		log.Error("Failed to load config from 'cuttle.yml'.")
		log.Fatal(err)
	}

	viper.SetDefault("addr", ":8123")
	viper.SetDefault("verbose", false)

	defaults := []map[string]interface{}{{"host": "*", "shared": true, "control": "noop"}}
	viper.SetDefault("zones", defaults)

	configs := viper.Get("zones").([]interface{})
	zones := make([]Zone, len(configs))
	for i, v := range configs {
		config := v.(map[interface{}]interface{})
		zones[i] = *NewZone(
			config["host"].(string),
			config["shared"].(bool),
			config["control"].(string),
			config["limit"].(int),
		)
	}

	// Config proxy.
	addr := viper.GetString("addr")
	verbose := viper.GetBool("verbose")
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = verbose

	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			for _, zone := range zones {
				if !zone.MatchHost(r.URL.Host) {
					continue
				}

				// Acquire permission to forward request to upstream server.
				zone.GetController(r.URL.Host).Acquire()

				return r, nil // Forward request.
			}

			log.Warn("No zone is applied. - ", r.URL)
			return r, nil // Forward request without rate limit.
		})

	log.Fatal(http.ListenAndServe(addr, proxy))
}
