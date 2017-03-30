// Package monitor provides functionality for monitoring CLF formatted log files.
package monitor

import (
	"log"
	"time"

	"github.com/vladimiroff/logmon/clf"
	"github.com/vladimiroff/logmon/handlers"
)

type handlerRecord struct {
	inst handlers.Interface
	tick *time.Ticker
}

// Monitor holds a register of handlers, sends all lines to them and makes sure
// to gather summary over provided period of time.
type Monitor struct {
	handlers []handlerRecord
	Input    chan<- clf.Line
}

// New instanciates new monitor.
func New() *Monitor {
	input := make(chan clf.Line)
	m := &Monitor{
		handlers: make([]handlerRecord, 0),
		Input:    input,
	}

	go func(input <-chan clf.Line) {
		for line := range input {
			for _, handle := range m.handlers {
				handle.inst.Input() <- line
			}
		}
	}(input)

	return m
}

// AddHandler creates new handler and registers it into the monitor.
func (m *Monitor) AddHandler(h handlers.Interface, d time.Duration) {
	handler := handlerRecord{
		inst: h,
		tick: time.NewTicker(d),
	}
	m.handlers = append(m.handlers, handler)

	go func(h handlerRecord, d time.Duration) {
		var out = make(chan string)

		for range h.tick.C {
			h.inst.Sum() <- out
			summary := <-out
			if len(summary) > 0 {
				log.Printf("%s", summary)
			}
		}
	}(handler, d)
}
