// Package handlers defines handlers and monitors of CLF lines.
package handlers

import (
	"time"

	"github.com/vladimiroff/logmon/clf"
)

const dateFormat = "2/Jan/2006:15:04:05 -0700"

var now = time.Now

// Handler interface defines handler's behaviour.
type Handler interface {
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
