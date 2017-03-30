# logmon

Simple console tool for monitoring CLF-formatted HTTP access log.

## Installation

    go get github.com/vladimiroff/logmon/...
    cd $GOPATH/src/github.com/vladimiroff/logmon
    dep ensure
    go install github.com/vladimiroff/logmon/cmd/logmon

## Usage

    logmon [options]

      -path string
            Path to log file
      -poll
            Poll for file changes instead of using inotify
      -threshold uint
            Raise alerts on traffic higher than provided threshold (default 1000)
