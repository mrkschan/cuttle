package main

import (
	"crypto/tls"
	"flag"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/elazarl/goproxy"
	"gopkg.in/yaml.v2"

    "github.com/mrkschan/cuttle/cuttle"
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

	cfg := Config{Addr: ":3128", CACert: "", CAKey: "", TLSVerify: true}
	if err := yaml.Unmarshal(bytes, &cfg); err != nil {
		log.Errorf("Malformed YAML in %s.", filename)
		log.Fatal(err)
	}

	zones := make([]cuttle.Zone, len(cfg.Zones))
	for i, c := range cfg.Zones {
		if c.Path == "" {
			c.Path = "/"
		}

		if c.LimitBy == "" {
			c.LimitBy = "host"
		}

		log.Debugf("ZoneConfig: host - %s, path - %s, limitby - %s, shared - %t, control - %s, rate - %d",
			c.Host, c.Path, c.LimitBy, c.Shared, c.Control, c.Rate)

		zones[i] = *cuttle.NewZone(c.Host, c.Path, c.LimitBy, c.Shared, c.Control, c.Rate)
	}

	// Config CA Cert.
	cert := goproxy.GoproxyCa // Use goproxy CA as default
	if cfg.CACert != "" && cfg.CAKey != "" {
		cert, err = tls.LoadX509KeyPair(cfg.CACert, cfg.CAKey)
		if err != nil {
			log.Warnf("Cannot load CA certificate from %s and %s.", cfg.CACert, cfg.CAKey)
		}
	}

	// Config proxy.
	proxy := goproxy.NewProxyHttpServer()
	proxy.Tr = &http.Transport{
		// Config TLS cert verification.
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !cfg.TLSVerify},
		Proxy:           http.ProxyFromEnvironment,
	}

	var httpsHandler goproxy.FuncHttpsHandler = func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		action := &goproxy.ConnectAction{
			Action:    goproxy.ConnectMitm,
			TLSConfig: goproxy.TLSConfigFromCA(&cert),
		}
		return action, host
	}
	proxy.OnRequest().HandleConnect(httpsHandler)
	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			var zone *cuttle.Zone
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
	Addr      string // Optional, default ":3128"
	CACert    string // Optional, default ""
	CAKey     string // Optional, default ""
	TLSVerify bool   // Optional, default "true"

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
