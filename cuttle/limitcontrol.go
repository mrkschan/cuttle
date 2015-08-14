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
}

func NewNoopControl() *NoopControl {
	return &NoopControl{}
}

// Start running NoopControl.
func (c *NoopControl) Start() {
}

// Acquire permission to forward request from NoopControl.
func (c *NoopControl) Acquire() {
}

// RPSControl provides requests per second rate limit control.
type RPSControl struct {
	// Limit holds the number of requests per second.
	Limit int

	pendingChan chan uint
	readyChan   chan uint
	seen        *list.List
}

func NewRPSControl(limit int) *RPSControl {
	return &RPSControl{limit, make(chan uint), make(chan uint), list.New()}
}

// Start running RPSControl.
func (c *RPSControl) Start() {
	go func() {
		log.Debugf("RPSControl: activated.")

		for {
			<-c.pendingChan

			log.Debugf("RPSControl: limit - %d", c.Limit)
			if c.seen.Len() == c.Limit {
				front := c.seen.Front()
				nanoElapsed := time.Now().UnixNano() - front.Value.(int64)
				milliElapsed := nanoElapsed / int64(time.Millisecond)
				log.Debugf("RPSControl: elapsed - %dms", milliElapsed)

				if waitTime := 1000 - milliElapsed; waitTime > 0 {
					log.Debugf("RPSControl: waiting - %dms", waitTime)
					time.Sleep(time.Duration(waitTime) * time.Millisecond)
				}

				c.seen.Remove(front)
			}
			c.seen.PushBack(time.Now().UnixNano())

			c.readyChan <- 1
		}

		log.Debugf("RPSControl: deactivated.")
	}()
}

// Acquire permission to forward request from RPSControl.
func (c *RPSControl) Acquire() {
	log.Debugf("RPSControl: permission requested.")
	c.pendingChan <- 1
	<-c.readyChan
	log.Debugf("RPSControl: permission granted.")
}
