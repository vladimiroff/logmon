// Package handlers defines handlers and monitors of CLF lines.
package handlers

import "github.com/vladimiroff/logmon/clf"

// Handler interface defines handler's behaviour.
type Handler interface {
	// Input returns a channel for handling new lines.
	Input() chan<- clf.Line

	// Sum returns a channel for requesting summaries.
	Sum() chan<- chan string

	// Stop stops the alerter's loop and no more inputs and sum requests
	// are being handled.
	//
	// Once stopped a handler cannot be restarted and stopping a stopped
	// handler is not supposed to happen (i.e. undefined behaviour).
	Stop()
}
