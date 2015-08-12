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
		zone := Zone{
			Host:    config["host"].(string),
			Shared:  config["shared"].(bool),
			Control: config["control"].(string),
			Limit:   config["limit"].(int),
		}
		zone.Activate()

		zones[i] = zone
	}

	// // Config limit controller.
	// var controller LimitController
	// control := viper.GetString("limitcontrol.controller")
	// if control == "rps" {
	// 	limit := viper.GetInt("limitcontrol.limit")
	// 	controller = &RPSControl{
	// 		Limit: limit,
	// 	}
	// } else {
	// 	log.Fatal("Unknown limit control: ", control)
	// }
	//
	// Config proxy.
	addr := viper.GetString("addr")
	verbose := viper.GetBool("verbose")
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = verbose

	// Starts now.
	// controller.Start()
	//
	// proxy.OnRequest().DoFunc(
	// 	func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	// 		// Acquire permission to forward request to downstream.
	// 		controller.Acquire()
	//
	// 		return r, nil // Forward request.
	// 	})

	log.Fatal(http.ListenAndServe(addr, proxy))
}
