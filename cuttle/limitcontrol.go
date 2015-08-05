package main

import (
	log "github.com/Sirupsen/logrus"
	"time"
)

// LimitController defines behaviors of all rate limit controls.
type LimitController interface {
	Start()
	Acquire()
}

// RPSController provides request per second rate limit.
type RPSController struct {
	Limit uint

	pendingChan chan uint
	readyChan   chan uint
	lastNano    int64
}

// Start running RPSController.
func (c *RPSController) Start() {
	c.pendingChan = make(chan uint)
	c.readyChan = make(chan uint)

	go func() {
		for {
			<-c.pendingChan

			nanoElapsed := time.Now().UnixNano() - c.lastNano
			milliElapsed := nanoElapsed / int64(time.Millisecond)
			log.Debug("RPS Limit Control: ", "elapsed=", milliElapsed)

			if milliElapsed < 1000 {
				time.Sleep(time.Duration(milliElapsed) * time.Millisecond)
			}
			c.lastNano = time.Now().UnixNano()

			c.readyChan <- 1
		}
	}()
}

// Acquire permission to forward request from RPSController.
func (c *RPSController) Acquire() {
	c.pendingChan <- 1
	<-c.readyChan
}
