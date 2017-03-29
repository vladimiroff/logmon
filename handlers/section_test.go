package handlers

import (
	"fmt"
	"strings"
	"testing"

	"github.com/vladimiroff/logmon/clf"
)

var (
	scCases = []struct {
		all  []string
		best []string
	}{
		{[]string{}, []string{}},
		{[]string{"foo", "bar", "baz", "baz"}, []string{"baz"}},
		{[]string{"foo", "bar", "baz"}, []string{"foo", "bar", "baz"}},
		{[]string{"foo", "foo", "bar", "baz", "baz"}, []string{"foo", "baz"}},
	}
	sceCases = []struct {
		in  string
		out string
	}{
		{"foo", "foo"},
		{"/foo", "/foo"},
		{"/foo.html", "/foo.html"},
		{"/foo?a=1", "/foo"},
		{"/foo/", "/foo"},
		{"/foo/bar", "/foo"},
		{"/foo/bar/", "/foo"},
	}
	scCountError    = "Expected %d lines of output, got %d: %s"
	scNoSectionsErr = "Didn't expect any sections, got %s"
	scWrongBestErr  = "Expected report line \n'%s', got \n'%s'"
	scExtractError  = "SectionCounter.extractSection(%q) should return %q, got %q"
)

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

func TestSections(t *testing.T) {

	for _, c := range scCases {
		t.Run(fmt.Sprintf("best %d of %d", len(c.best), len(c.all)), func(t *testing.T) {
			s := NewSections()
			defer s.Stop()
			out := make(chan string)

			for _, name := range c.all {
				s.Input() <- clf.Line{Resource: name}
			}
			s.Sum() <- out
			output := <-out

			if len(c.best) == 0 && output != noSectionsMsg {
				t.Errorf(scNoSectionsErr, output)
			}

			lines := strings.Split(output, "\n")
			if len(c.best) > 0 && len(lines)-2 != len(c.best) {
				t.Errorf(scCountError, len(c.best), len(lines)-2, output)
			}

			for i, line := range lines {
				// Skip first informative or latest empty line (after newline)
				if i == 0 || len(line) == 0 {
					continue
				}

				if i > len(c.best) {
					fmt.Printf("c.best = %+v\n", c.best)
					fmt.Printf("i = %d\n", i)
					fmt.Printf("lines = %s\n", lines)
				}
				expected := fmt.Sprintf(sectionsFormat, c.best[i-1],
					dups(c.all, c.best[i-1]), 0)
				if line+"\n" != expected {
					t.Errorf(scWrongBestErr, expected, line+"\n")
				}
			}
		})
	}
}

func TestExtractSection(t *testing.T) {
	for _, c := range sceCases {
		t.Run(fmt.Sprintf("%s from %s", c.out, c.in), func(t *testing.T) {
			section := extractSection(c.in)
			if section != c.out {
				t.Errorf(scExtractError, c.in, c.out, section)
			}
		})
	}
}
