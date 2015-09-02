package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/elazarl/goproxy"
	"gopkg.in/yaml.v2"
)

func main() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	switch os.Getenv("LOGLEVEL") {
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "WARNING", "WARN":
		log.SetLevel(log.WarnLevel)
	case "ERROR", "ERR", "FATAL":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	filename := "./cuttle.yml"
	flag.StringVar(&filename, "f", filename, "Configuration file to be loaded.")
	flag.Parse()

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Errorf("Failed to load configuration from %s.", filename)
		log.Fatal(err)
	}

	cfg := Config{Addr: ":3128"}
	if err := yaml.Unmarshal(bytes, &cfg); err != nil {
		log.Errorf("Malformed YAML in %s.", filename)
		log.Fatal(err)
	}

	zones := make([]Zone, len(cfg.Zones))
	for i, c := range cfg.Zones {
		if c.Path == "" {
			c.Path = "/"
		}

		if c.LimitBy == "" {
			c.LimitBy = "host"
		}

		log.Debugf("ZoneConfig: host - %s, path - %s, limitby - %s, shared - %t, control - %s, rate - %d",
			c.Host, c.Path, c.LimitBy, c.Shared, c.Control, c.Rate)

		zones[i] = *NewZone(c.Host, c.Path, c.LimitBy, c.Shared, c.Control, c.Rate)
	}

	// Config proxy.
	proxy := goproxy.NewProxyHttpServer()

	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			var zone *Zone
			for _, z := range zones {
				if z.MatchHost(r.URL.Host) && z.MatchPath(r.URL.Path) {
					zone = &z
					break
				}
			}

			if zone != nil {
				// Acquire permission to forward request to upstream server.
				zone.GetController(r.URL.Host, r.URL.Path).Acquire()
			} else {
				// No rate limit applied.
				log.Warnf("Main: No zone is applied to %s", r.URL)
			}

			// Forward request.
			log.Infof("Main: Forwarding request to %s", r.URL)
			return r, nil
		})

	log.Infof("Listening on %s", cfg.Addr)
	log.Fatalln(http.ListenAndServe(cfg.Addr, proxy))
}

type Config struct {
	Addr string // Optional, default ":3128"

	Zones []ZoneConfig
}

type ZoneConfig struct {
	Host    string
	Path    string // Optional, default "/"
	LimitBy string // Optional, default "host"
	Shared  bool   // Optional, default "false"
	Control string
	Rate    int
}
