package handlers

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/vladimiroff/logmon/clf"
)

const (
	noSectionsMsg  = "No sections have been hit"
	mostUsedMsg    = "Most used section(s):\n"
	sectionsFormat = "  - %s: total hits: %d, average bytes served: %d\n"
)

var sectionRe = regexp.MustCompile(`^(/[A-Za-z0-9_\-%\.]+).*$`)

// Section is a container for hits and total content length of a visited
// section.
type Section struct {
	hits uint64
	size uint64
}

// SectionCounter counts all visited sections and reports most visited of them.
type SectionCounter struct {
	*Handler

	all  map[string]Section
	best []string
}

// NewSections creates new handler for rating sections and starts its loop.
func NewSections() *SectionCounter {
	s := SectionCounter{
		Handler: NewHandler(),
		all:     make(map[string]Section),
		best:    make([]string, 0),
	}
	go s.loop()
	return &s
}

func (s *SectionCounter) loop() {
	var (
		line clf.Line
		out  chan string
	)

	for {
		select {
		case line = <-s.input:
			s.process(line)
		case out = <-s.cycle:
			s.sum(out)
		case <-s.stop:
			return
		}
	}
}

func (s *SectionCounter) process(l clf.Line) {
	name := l.Resource
	section := s.all[name]
	section.hits++
	section.size += l.ContentLength
	s.all[name] = section

	switch {
	case len(s.best) == 0:
		fallthrough
	case s.all[name].hits > s.all[s.best[0]].hits:
		s.best = []string{name}
	case s.all[name].hits == s.all[s.best[0]].hits:
		if dups(s.best, name) == 0 {
			s.best = append(s.best, name)
		}
	}
}

func (s *SectionCounter) sum(out chan string) {
	var buf bytes.Buffer

	defer func() {
		s.all = make(map[string]Section)
		s.best = make([]string, 0)
	}()

	if len(s.best) == 0 {
		out <- noSectionsMsg
		return
	}

	buf.WriteString(mostUsedMsg)
	for _, name := range s.best {
		buf.WriteString(fmt.Sprintf(sectionsFormat, name,
			s.all[name].hits, s.all[name].size/s.all[name].hits))
	}
	out <- buf.String()
}

func dups(x []string, element string) int {
	var result int
	for _, v := range x {
		if v == element {
			result++
		}
	}
	return result
}
func extractSection(resource string) string {
	match := sectionRe.FindStringSubmatch(resource)
	if len(match) < 2 {
		return resource
	}
	return match[1]
}
