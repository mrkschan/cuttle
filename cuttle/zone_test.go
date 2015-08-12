package main

import (
	"testing"
)

func TestZone(t *testing.T) {
	var zone Zone
	var c1, c2 LimitController

	zone = *NewZone("*.github.com", true, "rps", 2)

	c1 = zone.GetController("www.github.com")
	c2 = zone.GetController("api.github.com")
	if c1 != c2 {
		t.Errorf("Shared zone should return shared controller.")
	}

	zone = *NewZone("*.github.com", false, "rps", 2)

	c1 = zone.GetController("www.github.com")
	c2 = zone.GetController("api.github.com")
	if c1 == c2 {
		t.Errorf("Non-shared zone should return individual controller.")
	}
}
