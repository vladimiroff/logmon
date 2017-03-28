// Package clf unmarshals NCSA Common Log Format lines.
package clf

import (
	"fmt"
	"time"
)

const (
	lineFormat = "%s %s %s [%s %5s] %q %d %d"
	dateFormat = "2/Jan/2006:15:04:05 -0700"
	reqFormat  = "%s %s %s"
)

// Line represents one log line in NCSA Common Log Format.
type Line struct {
	RemoteHost    string
	Identity      string
	UserID        string
	Date          time.Time
	Method        string
	Resource      string
	Protocol      string
	StatusCode    uint16
	ContentLength uint64
}

// UnmarshalText decodes the receiver into UTF-8 encoded CLF line.
func (l *Line) UnmarshalText(text []byte) error {
	var d, t, r string

	_, err := fmt.Sscanf(string(text), lineFormat, &l.RemoteHost,
		&l.Identity, &l.UserID, &d, &t, &r, &l.StatusCode, &l.ContentLength)

	if err != nil {
		return unmarshalError("line", err)
	}

	fmt.Sscanf(r, reqFormat, &l.Method, &l.Resource, &l.Protocol)
	if err != nil {
		return unmarshalError("request", err)
	}

	l.Date, err = time.Parse(dateFormat, d+" "+t)
	if err != nil {
		return unmarshalError("date", err)
	}

	return nil
}

func unmarshalError(step string, err error) error {
	return fmt.Errorf("clf: cannot unmarshal %s: %s", step, err)
}
