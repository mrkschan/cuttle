package main

import (
	"testing"
	"time"
)

func TestRPSControl(t *testing.T) {
	var control LimitController
	var startT, endT int64

	control = NewRPSControl(2)
	control.Start()

	startT = time.Now().UnixNano()
	control.Acquire() // Expect no wait time.
	control.Acquire() // Expect no wait time.
	endT = time.Now().UnixNano()

	// Expecting no delay in 2 consecutive Acquire() with Limit=2.
	if elapsed := (endT - startT) / int64(time.Millisecond); elapsed > 1000 {
		t.Errorf("2x RPSControl.Acquire() elapsed %dms, want %dms", elapsed, 0)
	}

	control = NewRPSControl(2)
	control.Start()

	startT = time.Now().UnixNano()
	control.Acquire() // Expect no wait time.
	time.Sleep(time.Duration(500) * time.Millisecond)
	control.Acquire() // Expect no wait time.
	time.Sleep(time.Duration(300) * time.Millisecond)
	control.Acquire() // Expect 200ms wait time.
	endT = time.Now().UnixNano()

	// Expecting delay in 3 consecutive Acquire() with Limit=2.
	if elapsed := (endT - startT) / int64(time.Millisecond); elapsed < 1000 {
		t.Errorf("3x RPSControl.Acquire() elapsed %dms, want > %dms", elapsed, 1000)
	}

	control = NewRPSControl(2)
	control.Start()

	startT = time.Now().UnixNano()
	control.Acquire() // Expect no wait time.
	time.Sleep(time.Duration(500) * time.Millisecond)
	control.Acquire() // Expect no wait time.
	time.Sleep(time.Duration(300) * time.Millisecond)
	control.Acquire() // Expect 200ms wait time.
	time.Sleep(time.Duration(400) * time.Millisecond)
	control.Acquire() // Expect 100ms wait time.
	endT = time.Now().UnixNano()

	// Expecting delay in 4 consecutive Acquire() with Limit=2.
	if elapsed := (endT - startT) / int64(time.Millisecond); elapsed < 1500 {
		t.Errorf("4x RPSControl.Acquire() elapsed %dms, want > %dms", elapsed, 1500)
	}
}
