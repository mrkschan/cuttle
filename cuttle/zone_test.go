package main

import (
	"testing"
)

func TestZone(t *testing.T) {
	var zone Zone
	var c1, c2 LimitController

	zone = Zone{
		Host:    "*.github.com",
		Shared:  true,
		Control: "rps",
		Limit:   2,
	}
	zone.Activate()

	c1 = zone.GetController("www.github.com")
	c2 = zone.GetController("api.github.com")
	if c1 != c2 {
		t.Errorf("Shared zone should return shared controller.")
	}

	zone = Zone{
		Host:    "*.github.com",
		Shared:  false,
		Control: "rps",
		Limit:   2,
	}
	zone.Activate()

	c1 = zone.GetController("www.github.com")
	c2 = zone.GetController("api.github.com")
	if c1 == c2 {
		t.Errorf("Non-shared zone should return individual controller.")
	}
}
