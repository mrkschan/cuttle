package cuttle

import (
	"testing"
)

func TestZone(t *testing.T) {
	var zone Zone
	var c1, c2, c3 LimitController

	zone = *NewZone("*.github.com", "/", "host", true, "rps", 2)
	c1 = zone.GetController("www.github.com", "/")
	c2 = zone.GetController("api.github.com", "/")
	if c1 != c2 {
		t.Errorf("Shared zone should return shared controller.")
	}

	zone = *NewZone("*.github.com", "/", "host", false, "rps", 2)
	c1 = zone.GetController("www.github.com", "/")
	c2 = zone.GetController("api.github.com", "/")
	if c1 == c2 {
		t.Errorf("Non-shared zone should return individual controller.")
	}

	zone = *NewZone("www.github.com", "/*", "path", true, "rps", 2)
	c1 = zone.GetController("www.github.com", "/")
	c2 = zone.GetController("www.github.com", "/github/")
	c3 = zone.GetController("www.github.com", "/atom/")
	if c1 != c2 {
		t.Errorf("Shared zone should return shared controller.")
	}
	if c2 != c3 {
		t.Errorf("Shared zone should return shared controller.")
	}

	zone = *NewZone("www.github.com", "/*", "path", false, "rps", 2)
	c1 = zone.GetController("www.github.com", "/")
	c2 = zone.GetController("www.github.com", "/github/")
	c3 = zone.GetController("www.github.com", "/atom/")
	if c1 == c2 {
		t.Errorf("Non-shared zone should return individual controller.")
	}
	if c2 == c3 {
		t.Errorf("Non-shared zone should return individual controller.")
	}

	zone = *NewZone("*", "/", "host", false, "rps", 2)
	if !zone.MatchHost("github.com") {
		t.Errorf("zone(%s).MatchHost(%s) should be %s", "*", "github.com", true)
	}

	zone = *NewZone("*.com", "/", "host", false, "rps", 2)
	if !zone.MatchHost("github.com") {
		t.Errorf("zone(%s).MatchHost(%s) should be %s", "*.com", "github.com", true)
	}
	if zone.MatchHost("github.org") {
		t.Errorf("zone(%s).MatchHost(%s) should be %s", "*.com", "github.org", false)
	}

	zone = *NewZone("github.com", "/", "host", false, "rps", 2)
	if !zone.MatchHost("github.com") {
		t.Errorf("zone(%s).MatchHost(%s) should be %s", "github.com", "github.com", true)
	}
	if zone.MatchHost("www.github.com") {
		t.Errorf("zone(%s).MatchHost(%s) should be %s", "github.com", "www.github.com", false)
	}

	zone = *NewZone("*.github.com", "/", "host", false, "rps", 2)
	if !zone.MatchHost("www.github.com") {
		t.Errorf("zone(%s).MatchHost(%s) should be %s", "*.github.com", "www.github.com", true)
	}
	if zone.MatchHost("github.com") {
		t.Errorf("zone(%s).MatchHost(%s) should be %s", "*.github.com", "github.com", false)
	}
	if zone.MatchHost("hubgit.com") {
		t.Errorf("zone(%s).MatchHost(%s) should be %s", "*.github.com", "hubgit.com", false)
	}

	zone = *NewZone("*.*.github.com", "/", "host", false, "rps", 2)
	if !zone.MatchHost("x.www.github.com") {
		t.Errorf("zone(%s).MatchHost(%s) should be %s", "*.*.github.com", "x.www.github.com", true)
	}
	if zone.MatchHost("www.github.com") {
		t.Errorf("zone(%s).MatchHost(%s) should be %s", "*.*.github.com", "www.github.com", false)
	}
	if zone.MatchHost("github.com") {
		t.Errorf("zone(%s).MatchHost(%s) should be %s", "*.*.github.com", "github.com", false)
	}

	zone = *NewZone("github.com", "/", "path", false, "rps", 2)
	if !zone.MatchPath("/") {
		t.Errorf("zone(%s).MatchPath(%s) should be %s", "/", "/", true)
	}
	if !zone.MatchPath("/github") {
		t.Errorf("zone(%s).MatchPath(%s) should be %s", "/", "/github", true)
	}

	zone = *NewZone("github.com", "/*", "path", false, "rps", 2)
	if !zone.MatchPath("/") {
		t.Errorf("zone(%s).MatchPath(%s) should be %s", "/*", "/", true)
	}
	if !zone.MatchPath("/github") {
		t.Errorf("zone(%s).MatchPath(%s) should be %s", "/*", "/github", true)
	}
	if !zone.MatchPath("/github/hub") {
		t.Errorf("zone(%s).MatchPath(%s) should be %s", "/*", "/github/hub", true)
	}

	zone = *NewZone("github.com", "/github/*", "path", false, "rps", 2)
	if zone.MatchPath("/") {
		t.Errorf("zone(%s).MatchPath(%s) should be %s", "/github/*", "/", false)
	}
	if zone.MatchPath("/atom") {
		t.Errorf("zone(%s).MatchPath(%s) should be %s", "/github/*", "/atom", false)
	}
	if zone.MatchPath("/github") {
		t.Errorf("zone(%s).MatchPath(%s) should be %s", "/github/*", "/github", false)
	}
	if !zone.MatchPath("/github/") {
		t.Errorf("zone(%s).MatchPath(%s) should be %s", "/github/*", "/github", true)
	}
	if !zone.MatchPath("/github/hub") {
		t.Errorf("zone(%s).MatchPath(%s) should be %s", "/github/*", "/github/hub", true)
	}
}
