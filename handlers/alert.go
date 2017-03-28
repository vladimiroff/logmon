package handlers

import (
	"fmt"

	"github.com/vladimiroff/logmon/clf"
)

const (
	triggerAlertFormat = "High traffic generated an alert - hits = %d, triggered at %s"
	recoverAlertFormat = "Recovered from high traffic alert at %s"
)

// TrafficAlerter counts handled requests and fires an alert when going above
// given threshold.
type TrafficAlerter struct {
	count     uint64
	input     chan clf.Line
	cycle     chan chan string
	stop      chan struct{}
	threshold uint64
	alert     bool
}

// NewTrafficAlerter creates new traffic alerter and starts its loop.
func NewTrafficAlerter(threshold uint64) *TrafficAlerter {
	ta := TrafficAlerter{
		input:     make(chan clf.Line),
		cycle:     make(chan chan string),
		stop:      make(chan struct{}),
		threshold: threshold,
	}
	go ta.loop()
	return &ta
}

func (ta *TrafficAlerter) loop() {
	for {
		select {
		case <-ta.input:
			ta.count++
		case out := <-ta.cycle:
			ta.sum(out)
		case <-ta.stop:
			return
		}
	}
}

// Input returns a channel for handling new lines.
func (ta *TrafficAlerter) Input() chan<- clf.Line {
	return ta.input
}

// Sum returns a channel for requesting summaries.
func (ta *TrafficAlerter) Sum() chan<- chan string {
	return ta.cycle
}

// Stop the alerter's loop.
//
// NOTE: Stopping a stopped alerter will panic.
func (ta *TrafficAlerter) Stop() {
	close(ta.stop)
}

func (ta *TrafficAlerter) sum(out chan string) {
	defer func() { ta.count = 0 }()

	if ta.count > ta.threshold && !ta.alert {
		ta.alert = true
		out <- fmt.Sprintf(triggerAlertFormat, ta.count, now().Format(dateFormat))
	} else if ta.count <= ta.threshold && ta.alert {
		ta.alert = false
		out <- fmt.Sprintf(recoverAlertFormat, now().Format(dateFormat))
	} else {
		out <- ""
	}
}
