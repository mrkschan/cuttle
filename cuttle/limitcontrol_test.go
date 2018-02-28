package cuttle

import (
	"testing"
	"time"
)

func TestNoopControl(t *testing.T) {
	var control LimitController
	var startT, endT int64

	control = NewNoopControl("label")
	control.Start()

	acquired := true

	startT = time.Now().UnixNano()
	acquired = acquired && control.Acquire() // Expect no wait time.
	acquired = acquired && control.Acquire() // Expect no wait time.
	endT = time.Now().UnixNano()

	// Expecting acquired is true
	if acquired != true {
		t.Errorf("Permission cannot be acquired from NoopControl.Acquire()")
	}

	// Expecting no delay in 2 consecutive Acquire().
	if elapsed := (endT - startT) / int64(time.Millisecond); elapsed > 1 {
		t.Errorf("2x NoopControl.Acquire() elapsed %dms, want %dms", elapsed, 0)
	}
}

func TestBanControl(t *testing.T) {
	var control LimitController

	control = NewBanControl("label")
	control.Start()

	acquired := control.Acquire()

	// Expecting acquired is false
	if acquired != false {
		t.Errorf("Permission acquired from BanControl.Acquire()")
	}
}

func TestRPSControl(t *testing.T) {
	var control LimitController
	var startT, endT int64

	control = NewRPSControl("label", 2)
	control.Start()

	acquired := true

	startT = time.Now().UnixNano()
	acquired = acquired && control.Acquire() // Expect no wait time.
	acquired = acquired && control.Acquire() // Expect no wait time.
	endT = time.Now().UnixNano()

	// Expecting acquired is true
	if acquired != true {
		t.Errorf("Permission cannot be acquired from RPSControl.Acquire()")
	}

	// Expecting no delay in 2 consecutive Acquire() with Rate=2.
	if elapsed := (endT - startT) / int64(time.Millisecond); elapsed > 1000 {
		t.Errorf("2x RPSControl.Acquire() elapsed %dms, want %dms", elapsed, 0)
	}

	control = NewRPSControl("label", 2)
	control.Start()

	acquired = true

	startT = time.Now().UnixNano()
	acquired = acquired && control.Acquire() // Expect no wait time.
	time.Sleep(time.Duration(500) * time.Millisecond)
	acquired = acquired && control.Acquire() // Expect no wait time.
	time.Sleep(time.Duration(300) * time.Millisecond)
	acquired = acquired && control.Acquire() // Expect 200ms wait time.
	endT = time.Now().UnixNano()

	// Expecting acquired is true
	if acquired != true {
		t.Errorf("Permission cannot be acquired from RPSControl.Acquire()")
	}

	// Expecting delay in 3 consecutive Acquire() with Rate=2.
	if elapsed := (endT - startT) / int64(time.Millisecond); elapsed < 1000 {
		t.Errorf("3x RPSControl.Acquire() elapsed %dms, want > %dms", elapsed, 1000)
	}

	control = NewRPSControl("label", 2)
	control.Start()

	acquired = true

	startT = time.Now().UnixNano()
	acquired = acquired && control.Acquire() // Expect no wait time.
	time.Sleep(time.Duration(500) * time.Millisecond)
	acquired = acquired && control.Acquire() // Expect no wait time.
	time.Sleep(time.Duration(300) * time.Millisecond)
	acquired = acquired && control.Acquire() // Expect 200ms wait time.
	time.Sleep(time.Duration(400) * time.Millisecond)
	acquired = acquired && control.Acquire() // Expect 100ms wait time.
	endT = time.Now().UnixNano()

	// Expecting acquired is true
	if acquired != true {
		t.Errorf("Permission cannot be acquired from RPSControl.Acquire()")
	}

	// Expecting delay in 4 consecutive Acquire() with Rate=2.
	if elapsed := (endT - startT) / int64(time.Millisecond); elapsed < 1500 {
		t.Errorf("4x RPSControl.Acquire() elapsed %dms, want > %dms", elapsed, 1500)
	}
}

func TestRPMControl(t *testing.T) {
	var control LimitController
	var startT, endT int64

	control = NewRPMControl("label", 2)
	control.Start()

	acquired := true

	startT = time.Now().UnixNano()
	acquired = acquired && control.Acquire() // Expect no wait time.
	acquired = acquired && control.Acquire() // Expect no wait time.
	endT = time.Now().UnixNano()

	// Expecting acquired is true
	if acquired != true {
		t.Errorf("Permission cannot be acquired from RPMControl.Acquire()")
	}

	// Expecting no delay in 2 consecutive Acquire() with Rate=2.
	if elapsed := (endT - startT) / 1000; elapsed > 1000 {
		t.Errorf("2x RPMControl.Acquire() elapsed %dms, want %dms", elapsed, 0)
	}

	control = NewRPMControl("label", 30)
	control.Start()

	acquired = true

	startT = time.Now().UnixNano()
	acquired = acquired && control.Acquire() // Expect no wait time.
	time.Sleep(time.Duration(15) * time.Second)
	acquired = acquired && control.Acquire() // Expect no wait time.
	time.Sleep(time.Duration(30) * time.Second)
	acquired = acquired && control.Acquire() // Expect 30s wait time.
	endT = time.Now().UnixNano()

	// Expecting acquired is true
	if acquired != true {
		t.Errorf("Permission cannot be acquired from RPMControl.Acquire()")
	}

	// Expecting delay in 30 consecutive Acquire() with Rate=30.
	if elapsed := (endT - startT) / 1000; elapsed < 60 {
		t.Errorf("3x RPMControl.Acquire() elapsed %dms, want > %dms", elapsed, 60)
	}

	control = NewRPSControl("label", 30)
	control.Start()

	acquired = true

	startT = time.Now().UnixNano()
	acquired = acquired && control.Acquire() // Expect no wait time.
	time.Sleep(time.Duration(20) * time.Second)
	acquired = acquired && control.Acquire() // Expect no wait time.
	time.Sleep(time.Duration(30) * time.Second)
	acquired = acquired && control.Acquire() // Expect 30s wait time.
	time.Sleep(time.Duration(60) * time.Second)
	acquired = acquired && control.Acquire() // Expect 30s wait time.
	endT = time.Now().UnixNano()

	// Expecting acquired is true
	if acquired != true {
		t.Errorf("Permission cannot be acquired from RPSControl.Acquire()")
	}

	// Expecting delay in 4 consecutive Acquire() with Rate=2.
	if elapsed := (endT - startT) / 1000; elapsed < 120 {
		t.Errorf("4x RPSControl.Acquire() elapsed %dms, want > %dms", elapsed, 120)
	}
}
