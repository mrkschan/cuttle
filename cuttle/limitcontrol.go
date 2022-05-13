package cuttle

import (
	"container/list"
	"time"

	log "github.com/sirupsen/logrus"
)

// LimitController defines behaviors of a rate limit control.
type LimitController interface {
	// Start the rate limit controller.
	Start()
	// Acquire permission to perform certain things.
	// The permission is granted according to the rate limit rule.
	Acquire() bool
}

// NoopControl does not perform any rate limit.
type NoopControl struct {
	// Label of this control.
	Label string
}

// NewNoopControl return a new NoopControl with the given label.
func NewNoopControl(label string) *NoopControl {
	return &NoopControl{label}
}

// Start running NoopControl.
func (c *NoopControl) Start() {
	log.Debugf("NoopControl[%s]: Activated.", c.Label)
}

// Acquire permission from NoopControl.
// Permission is granted immediately since it does not perform any rate limit.
func (c *NoopControl) Acquire() bool {
	log.Debugf("NoopControl[%s]: Seeking permission.", c.Label)
	log.Debugf("NoopControl[%s]: Granted permission.", c.Label)

	return true
}

// BanControl bans all the request.
type BanControl struct {
	// Label of this control.
	Label string
}

// NewBanControl return a new BanControl with the given label.
func NewBanControl(label string) *BanControl {
	return &BanControl{label}
}

// Start running BanControl.
func (c *BanControl) Start() {
	log.Debugf("BanControl[%s]: Activated.", c.Label)
}

// Acquire permission from BanControl.
// Permission is never granted.
func (c *BanControl) Acquire() bool {
	log.Debugf("BanControl[%s]: Seeking permission.", c.Label)
	log.Debugf("BanControl[%s]: No permission granted.", c.Label)

	return false
}

// RPSControl provides requests per second rate limit control.
type RPSControl struct {
	// Label of this control.
	Label string
	// Rate holds the number of requests per second.
	Rate int

	pendingChan chan uint
	readyChan   chan uint
	seen        *list.List
}

// NewRPSControl return a new RPSControl with the given label and rate.
func NewRPSControl(label string, rate int) *RPSControl {
	return &RPSControl{label, rate, make(chan uint), make(chan uint), list.New()}
}

// Start running RPSControl.
// A goroutine is launched to govern the rate limit of Acquire().
func (c *RPSControl) Start() {
	go func() {
		log.Debugf("RPSControl[%s]: Activated.", c.Label)

		for {
			<-c.pendingChan

			log.Debugf("RPSControl[%s]: Limited at %dreq/s.", c.Label, c.Rate)
			if c.seen.Len() == c.Rate {
				front := c.seen.Front()
				nanoElapsed := time.Now().UnixNano() - front.Value.(int64)
				milliElapsed := nanoElapsed / int64(time.Millisecond)
				log.Debugf("RPSControl[%s]: Elapsed %dms since first request.", c.Label, milliElapsed)

				if waitTime := 1000 - milliElapsed; waitTime > 0 {
					log.Infof("RPSControl[%s]: Waiting for %dms.", c.Label, waitTime)
					time.Sleep(time.Duration(waitTime) * time.Millisecond)
				}

				c.seen.Remove(front)
			}
			c.seen.PushBack(time.Now().UnixNano())

			c.readyChan <- 1
		}

		log.Debugf("RPSControl[%s]: Deactivated.", c.Label)
	}()
}

// Acquire permission from RPSControl.
// Permission is granted at a rate of N requests per second.
func (c *RPSControl) Acquire() bool {
	log.Debugf("RPSControl[%s]: Seeking permission.", c.Label)
	c.pendingChan <- 1
	<-c.readyChan
	log.Debugf("RPSControl[%s]: Granted permission.", c.Label)

	return true
}

// RPMControl provides requests per minute rate limit control.
type RPMControl struct {
	// Label of this control.
	Label string
	// Rate holds the number of requests per minute.
	Rate int

	pendingChan chan uint
	readyChan   chan uint
	seen        *list.List
}

// NewRPMControl return a new RPMControl with the given label and rate.
func NewRPMControl(label string, rate int) *RPMControl {
	return &RPMControl{label, rate, make(chan uint), make(chan uint), list.New()}
}

// Start running RPMControl.
// A goroutine is launched to govern the rate limit of Acquire().
func (c *RPMControl) Start() {
	go func() {
		log.Debugf("RPMControl[%s]: Activated.", c.Label)

		for {
			<-c.pendingChan

			log.Debugf("RPMControl[%s]: Limited at %dreq/m.", c.Label, c.Rate)
			if c.seen.Len() == c.Rate {
				front := c.seen.Front()
				nanoElapsed := time.Now().UnixNano() - front.Value.(int64)
				milliElapsed := nanoElapsed / int64(time.Millisecond)
				secondElapsed := milliElapsed / 1000
				log.Debugf("RPMControl[%s]: Elapsed %ds since first request.", c.Label, secondElapsed)

				if waitTime := 60 - secondElapsed; waitTime > 0 {
					log.Infof("RPMControl[%s]: Waiting for %ds.", c.Label, waitTime)
					time.Sleep(time.Duration(waitTime) * time.Second)
				}

				c.seen.Remove(front)
			}
			c.seen.PushBack(time.Now().UnixNano())

			c.readyChan <- 1
		}

		log.Debugf("RPMControl[%s]: Deactivated.", c.Label)
	}()
}

// Acquire permission from RPMControl.
// Permission is granted at a rate of N requests per minute.
func (c *RPMControl) Acquire() bool {
	log.Debugf("RPMControl[%s]: Seeking permission.", c.Label)
	c.pendingChan <- 1
	<-c.readyChan
	log.Debugf("RPMControl[%s]: Granted permission.", c.Label)

	return true
}
