package cuttle

import (
	"fmt"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// A Zone holds the settings and states of the rate limited location(s).
type Zone struct {
	// Host specifies the URL host of the location(s).
	// It supports using wildcard '*' for matching part of the host name.
	// E.g. "*.github.com" matches both "www.github.com" and "api.github.com".
	Host string
	// Path specifies the URL path of the location(s).
	// It supports using wildcard '*' for matching part of the path.
	// E.g. "/*" matches both "/atom" and "/github".
	Path string

	// LimitBy specifies the rate limit subject of the location(s).
	// Rate limit can be performed by "host" or "path".
	LimitBy string
	// Shared specifies whether the rate limit is shared among all location(s) in the Zone.
	Shared bool

	// Control specifies which rate limit controller is used.
	Control string
	// Rate specifies the rate of the rate limit controller.
	Rate int
    Nseconds int
	controllers map[string]LimitController
}

// NewZone returns a new Zone given the configurations.
func NewZone(host string, path string, limitby string, shared bool, control string, rate int, nseconds int) *Zone {
	return &Zone{
		host, path, limitby, shared, control, rate, nseconds,
		make(map[string]LimitController),
	}
}

// MatchHost determines whether the host name of a location belongs to the Zone.
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

// MatchPath determines whether the URL path of a location belongs to the Zone.
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

// GetController returns the rate limit controller of a location.
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
		case "rpm":
			controller = NewRPMControl(key, z.Rate)
        case "rpns":
			controller = NewRPNSControl(key, z.Rate, z.Nseconds)
		case "noop":
			controller = NewNoopControl(key)
		case "ban":
			controller = NewBanControl(key)
		}

		z.controllers[key] = controller
		controller.Start()
	}
	log.Debugf("Zone.GetController: zone - %s%s, key - %s, control - %s", z.Host, z.Path, key, z.Control)

	return z.controllers[key]
}
