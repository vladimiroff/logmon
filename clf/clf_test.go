package clf

import (
	"fmt"
	"testing"
	"time"
)

const example = `127.0.0.1 user-identifier frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326`

var cases = []struct {
	line  string
	valid bool
}{
	{example, true},
	{`example.com user-identifier frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326`, true},
	{`127.0.0.1 - - [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326`, true},
	{`127.0.0.1 user-identifier frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 a-lot`, false},
	{`127.0.0.1 user-identifier frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" "200 OK" 2326`, false},
	{`127.0.0.1 user-identifier frank [40/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326`, false},
	{`127.0.0.1 user-identifier frank [Mon, 02 Jan 2006 15:04:05 MST] "GET /apache_pb.gif HTTP/1.0" 200 2326`, false},
}

func Test(t *testing.T) {
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.line), func(t *testing.T) {
			var l Line
			err := l.UnmarshalText([]byte(c.line))

			if c.valid {
				if err != nil {
					t.Fatalf("Unexpected error: %s", err.Error())
				}
			} else {
				if err == nil {
					t.Fatalf("Expected error, got nil and unmarshalled: %#v", l)
				}
			}
		})

	}
}

func TestRFCExample(t *testing.T) {
	var l Line
	err := l.UnmarshalText([]byte(example))

	if err != nil {
		t.Fatalf("Unmarshal error: %s", err.Error())
	}

	expected := Line{
		RemoteHost:    "127.0.0.1",
		Identity:      "user-identifier",
		UserID:        "frank",
		Date:          time.Date(2000, 10, 10, 13, 55, 36, 0, l.Date.Location()),
		Method:        "GET",
		Resource:      "/apache_pb.gif",
		Protocol:      "HTTP/1.0",
		StatusCode:    200,
		ContentLength: 2326,
	}

	if l != expected {
		t.Fatalf("Unmarshal expect: \n%v\ngot: \n%v", expected, l)
	}
}
