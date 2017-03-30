package main

import (
	"flag"
	"log"
	"time"

	"github.com/hpcloud/tail"

	"github.com/vladimiroff/logmon/clf"
	"github.com/vladimiroff/logmon/handlers"
	"github.com/vladimiroff/logmon/monitor"
)

var (
	path      = flag.String("path", "", "Path to log file")
	poll      = flag.Bool("poll", false, "Poll for file changes instead of using inotify")
	threshold = flag.Uint64("threshold", 1000, "Raise alerts on traffic higher than provided threshold")

	alertPeriod    = 2 * time.Minute
	sectionsPeriod = 10 * time.Second
)

func main() {
	var (
		clfLine clf.Line
		err     error
	)

	flag.Parse()
	if len(*path) == 0 {
		log.Fatal("no path to log file provided")
	}

	config := tail.Config{Follow: true, ReOpen: true, Poll: *poll}
	f, err := tail.TailFile(*path, config)
	if err != nil {
		log.Fatal(err.Error())
	}

	monitor := monitor.New()
	monitor.AddHandler(handlers.NewTrafficAlerter(*threshold), alertPeriod)
	monitor.AddHandler(handlers.NewSections(), sectionsPeriod)

	for line := range f.Lines {
		if line.Err != nil {
			log.Printf("[ERROR] Can't fetch line: %s", line.Err)
		}

		err = clfLine.UnmarshalText([]byte(line.Text))
		if err != nil {
			log.Printf("[ERROR] Can't parse line: %s", err)
		}

		monitor.Input <- clfLine
	}
}
