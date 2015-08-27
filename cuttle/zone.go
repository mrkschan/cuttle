package main

import (
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type Zone struct {
	Host    string
	Shared  bool
	Control string
	Rate    int

	controllers map[string]LimitController
	re          string
}

func NewZone(host string, shared bool, control string, rate int) *Zone {
	re := strings.Replace(host, ".", "\\.", -1)
	re = strings.Replace(host, "*", "[^\\.]+", -1)

	return &Zone{
		host, shared, control, rate,
		make(map[string]LimitController), re,
	}
}

func (z *Zone) MatchHost(host string) bool {
	log.Debugf("Zone.MatchHost: zone - %s, host - %s", z.Host, host)

	matched, err := regexp.MatchString(z.re, host)
	if err != nil {
		log.Warn(err)
		return false
	}
	return matched
}

func (z *Zone) GetController(host string) LimitController {
	var key string
	if z.Shared {
		key = "*"
	} else {
		key = host
	}

	_, ok := z.controllers[key]
	if !ok {
		var controller LimitController

		switch z.Control {
		case "rps":
			controller = NewRPSControl(key, z.Rate)
		case "noop":
			controller = NewNoopControl(key)
		}

		z.controllers[key] = controller
		controller.Start()
	}
	log.Debugf("Zone.GetController: zone - %s, key - %s, control - %s", z.Host, key, z.Control)

	return z.controllers[key]
}
