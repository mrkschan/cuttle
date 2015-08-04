package main

import (
	log "github.com/Sirupsen/logrus"
	"time"
)

type LimitController interface {
	Acquire()
	Release()
}

var (
	// TODO(mrkschan): Should support different contollers per downstream.
	controller LimitController

	// TODO(mrkschan): Should support different pending queue per downstream.
	pendingChan = make(chan uint)

	// TODO(mrkschan): Should support different ready queue per downstream.
	readyChan = make(chan uint)
)

func setLimitController(c LimitController) {
	controller = c
}

func startLimitControl() {
	go func() {
		for {
			controller.Release()
		}
	}()
}

type RPSController struct {
	limit           uint
	lastRequestNano int64
}

func (c *RPSController) Acquire() {
	pendingChan <- 1
	<-readyChan
}

func (c *RPSController) Release() {
	<-pendingChan

	elapsed := (time.Now().UnixNano() - c.lastRequestNano)
	elapsedMilli := elapsed / int64(time.Millisecond)
	log.Debug("RPS Limit Control: ", "elapsed=", elapsedMilli)
	if elapsedMilli < 1000 {
		time.Sleep(time.Duration(elapsedMilli) * time.Millisecond)
	}
	c.lastRequestNano = time.Now().UnixNano()

	readyChan <- 1
}
