package main

import (
	"time"
)

type LimitController interface {
	Unblock()
}

var (
	// TODO(mrkschan): Should support different contollers per downstream.
	controller LimitController

	// TODO(mrkschan): Should support different FIFO per downstream.
	RateLimitFIFO = make(chan uint)
)

func setLimitController(c LimitController) {
	controller = c
}

func startLimitControl() {
	go func() {
		for {
			// TODO(mrkschan): Implement the main loop
			controller.Unblock()
			time.Sleep(1)
		}
	}()
}

type RPSController struct {
	limit uint
}

func (c RPSController) Unblock() {
	// TODO(mrkschan): Implement request per second rate limit
	RateLimitFIFO <- 1
}
