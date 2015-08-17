package main

import (
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/elazarl/goproxy"
	"gopkg.in/yaml.v2"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	filename := "cuttle.yml"
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Errorf("Failed to load configuration from %s.", filename)
		log.Fatal(err)
	}

	cfg := Config{Addr: ":8123"}
	if err := yaml.Unmarshal(bytes, &cfg); err != nil {
		log.Errorf("Malformed YAML in %s.", filename)
		log.Fatal(err)
	}

	zones := make([]Zone, len(cfg.Zones))
	for i, c := range cfg.Zones {
		zones[i] = *NewZone(c.Host, c.Shared, c.Control, c.Limit)
	}

	// Config proxy.
	proxy := goproxy.NewProxyHttpServer()

	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			for _, zone := range zones {
				if !zone.MatchHost(r.URL.Host) {
					continue
				}

				// Acquire permission to forward request to upstream server.
				zone.GetController(r.URL.Host).Acquire()

				// Forward request.
				log.Infof("Main: Forwarding request to %s", r.URL)
				return r, nil
			}

			// Forward request without rate limit.
			log.Warnf("Main: No zone is applied to %s", r.URL)
			log.Infof("Main: Forwarding request to %s", r.URL)
			return r, nil
		})

	log.Fatal(http.ListenAndServe(cfg.Addr, proxy))
}

type Config struct {
	Addr    string

	Zones []ZoneConfig
}

type ZoneConfig struct {
	Host    string
	Shared  bool
	Control string
	Limit   int
}
