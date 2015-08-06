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

// RPSControl provides requests per second rate limit control.
type RPSControl struct {
	// Limit holds the number of requests per second.
	Limit int

	pendingChan chan uint
	readyChan   chan uint
	seen        *list.List
}

// Start running RPSControl.
func (c *RPSControl) Start() {
	c.pendingChan = make(chan uint)
	c.readyChan = make(chan uint)
	c.seen = list.New()

	go func() {
		for {
			<-c.pendingChan

			if c.seen.Len() == c.Limit {
				front := c.seen.Front()
				nanoElapsed := time.Now().UnixNano() - front.Value.(int64)
				milliElapsed := nanoElapsed / int64(time.Millisecond)
				log.Debug("RPS control: ", "elapsed=", milliElapsed)

				if waitTime := 1000 - milliElapsed; waitTime > 0 {
					log.Debug("RPS control: ", "wait=", waitTime)
					time.Sleep(time.Duration(waitTime) * time.Millisecond)
				}

				c.seen.Remove(front)
			}
			c.seen.PushBack(time.Now().UnixNano())

			c.readyChan <- 1
		}
	}()
}

// Acquire permission to forward request from RPSControl.
func (c *RPSControl) Acquire() {
	c.pendingChan <- 1
	<-c.readyChan
}
