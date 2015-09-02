package main

import (
	"fmt"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type Zone struct {
	Host    string
	Path    string
	LimitBy string
	Shared  bool
	Control string
	Rate    int

	controllers map[string]LimitController
}

func NewZone(host string, path string, limitby string, shared bool, control string, rate int) *Zone {
	return &Zone{
		host, path, limitby, shared, control, rate,
		make(map[string]LimitController),
	}
}

func (z *Zone) MatchHost(host string) bool {
	log.Debugf("Zone.MatchHost: zone - %s%s, host - %s", z.Host, z.Path, host)

	pattern := z.Host
	pattern = strings.Replace(pattern, ".", "\\.", -1)
	pattern = strings.Replace(pattern, "*", "[^\\.]+", -1)

	matched, err := regexp.MatchString(fmt.Sprintf("^%s", pattern), host)
	if err != nil {
		log.Warn(err)
		return false
	}
	return matched
}

func (z *Zone) MatchPath(path string) bool {
	log.Debugf("Zone.MatchHost: zone - %s%s, path - %s", z.Host, z.Path, path)

	pattern := z.Path
	pattern = strings.Replace(pattern, "*", "([^/]*)", -1)

	matched, err := regexp.MatchString(fmt.Sprintf("^%s", pattern), path)
	if err != nil {
		log.Warn(err)
		return false
	}
	return matched
}

func (z *Zone) GetController(host string, path string) LimitController {
	var key string
	switch z.LimitBy {
	case "host":
		if z.Shared {
			key = "host:*"
		} else {
			key = "host:" + host
		}
	case "path":
		if z.Shared {
			key = "path:*"
		} else {
			pattern := z.Path
			pattern = strings.Replace(pattern, "*", "([^/]*)", -1)
			matches := regexp.MustCompile(pattern).FindStringSubmatch(path)
			key = "path:" + strings.Join(matches[1:], "-")
		}
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
	log.Debugf("Zone.GetController: zone - %s%s, key - %s, control - %s", z.Host, z.Path, key, z.Control)

	return z.controllers[key]
}
