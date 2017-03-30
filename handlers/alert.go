package handlers

import "fmt"

const (
	triggerAlertFormat = "High traffic generated an alert - hits = %d, triggered at %s"
	recoverAlertFormat = "Recovered from high traffic alert at %s"
)

// TrafficAlerter counts handled requests and fires an alert when going above
// given threshold.
type TrafficAlerter struct {
	*Handler

	count     uint64
	threshold uint64
	alert     bool
}

// NewTrafficAlerter creates new traffic alerter and starts its loop.
func NewTrafficAlerter(threshold uint64) *TrafficAlerter {
	ta := TrafficAlerter{
		Handler:   NewHandler(),
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
