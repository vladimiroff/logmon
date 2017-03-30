package handlers

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/vladimiroff/logmon/clf"
)

var (
	taCases = []struct {
		inputs uint64
		alert  bool
		change bool
		now    time.Time
	}{
		{4, false, false, time.Now()},
		{5, false, false, time.Now()},
		{6, true, true, time.Unix(1490721979, 0)},
		{7, true, false, time.Unix(1490722023, 0)},
		{2, false, true, time.Now()},
	}
	taErrFormat   = "Expected c.alert to be %t, got %t instead"
	taDateFormat  = "Expected to output %s at the end of new alert, got:\n%s"
	taUnexpFormat = "Output is not expected without a change, got: %s"
	taMsgError    = "Unexpected change message when alert is %t: %s"
)

func TestTrafficAlerts(t *testing.T) {

	ta := NewTrafficAlerter(5)
	for _, c := range taCases {
		t.Run(fmt.Sprintf("%t of %d inputs", c.alert, c.inputs), func(t *testing.T) {
			out := make(chan string)
			now = func() time.Time { return c.now }

			for j := uint64(0); j < c.inputs; j++ {
				ta.Input() <- clf.Line{}
			}

			ta.Sum() <- out
			output := <-out

			if ta.alert != c.alert {
				t.Errorf(taErrFormat, c.alert, ta.alert)
			}
			if c.change {
				var msg string
				if c.alert {
					msg = triggerAlertFormat
				} else {
					msg = recoverAlertFormat
				}

				if !strings.HasPrefix(output, strings.Split(msg, "%")[0]) {
					t.Errorf(taMsgError, c.alert, output)
				}

				if !strings.HasSuffix(output, now().Format(dateFormat)) {
					t.Errorf(taDateFormat, now().Format(dateFormat), output)
				}
			} else {
				if len(output) > 0 {
					t.Errorf(taUnexpFormat, output)
				}
			}

		})
	}
}
