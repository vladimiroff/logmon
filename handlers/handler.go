// Package handlers defines handlers and monitors of CLF lines.
package handlers

import (
	"time"

	"github.com/vladimiroff/logmon/clf"
)

const dateFormat = "2/Jan/2006:15:04:05 -0700"

var now = time.Now

// Interface interface defines handler'h behaviour.
type Interface interface {
	// Input returns a channel for handling new lines.
	Input() chan<- clf.Line

	// Sum returns a channel for requesting summaries.
	Sum() chan<- chan string

	// Stop the alerter's loop so no more inputs are going to be handled.
	//
	// Once stopped a handler cannot be restarted and stopping a stopped
	// handler is not supposed to happen (i.e. undefined behaviour).
	Stop()
}

// Handler is a base ready for embedding implementation of a bare handler.
type Handler struct {
	input chan clf.Line
	cycle chan chan string
	stop  chan struct{}
}

// NewHandler creates a Handler instance.
func NewHandler() *Handler {
	return &Handler{
		input: make(chan clf.Line),
		cycle: make(chan chan string),
		stop:  make(chan struct{}),
	}
}

// Input returns a channel for handling new lines.
func (h *Handler) Input() chan<- clf.Line {
	return h.input
}

// Sum returns a channel for requesting summaries.
func (h *Handler) Sum() chan<- chan string {
	return h.cycle
}

// Stop the sections'h loop.
func (h *Handler) Stop() {
	close(h.stop)
}
