package main

import (
	"container/list"
	log "github.com/Sirupsen/logrus"
	"time"
)

// LimitController defines behaviors of all rate limit controls.
type LimitController interface {
	Start()
	Acquire()
}

// NoopControl does not do any rate limit.
type NoopControl struct {
	// Label of this control.
	Label string
}

func NewNoopControl(label string) *NoopControl {
	return &NoopControl{label}
}

// Start running NoopControl.
func (c *NoopControl) Start() {
	log.Debugf("NoopControl[%s]: Activated.", c.Label)
}

// Acquire permission to forward request from NoopControl.
func (c *NoopControl) Acquire() {
	log.Debugf("NoopControl[%s]: Seeking permission.", c.Label)
	log.Debugf("NoopControl[%s]: Granted permission.", c.Label)
}

// RPSControl provides requests per second rate limit control.
type RPSControl struct {
	// Label of this control.
	Label string
	// Limit holds the number of requests per second.
	Limit int

	pendingChan chan uint
	readyChan   chan uint
	seen        *list.List
}

func NewRPSControl(label string, limit int) *RPSControl {
	return &RPSControl{label, limit, make(chan uint), make(chan uint), list.New()}
}

// Start running RPSControl.
func (c *RPSControl) Start() {
	go func() {
		log.Debugf("RPSControl[%s]: Activated.", c.Label)

		for {
			<-c.pendingChan

			log.Debugf("RPSControl[%s]: Limited at %dreq/s.", c.Label, c.Limit)
			if c.seen.Len() == c.Limit {
				front := c.seen.Front()
				nanoElapsed := time.Now().UnixNano() - front.Value.(int64)
				milliElapsed := nanoElapsed / int64(time.Millisecond)
				log.Debugf("RPSControl[%s]: Elapsed %dms since first request.", c.Label, milliElapsed)

				if waitTime := 1000 - milliElapsed; waitTime > 0 {
					log.Debugf("RPSControl[%s]: Waiting for %dms.", c.Label, waitTime)
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

// Acquire permission to forward request from RPSControl.
func (c *RPSControl) Acquire() {
	log.Debugf("RPSControl[%s]: Seeking permission.", c.Label)
	c.pendingChan <- 1
	<-c.readyChan
	log.Debugf("RPSControl[%s]: Granted permission.", c.Label)
}
