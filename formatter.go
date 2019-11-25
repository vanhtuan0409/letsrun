package main

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type formatter interface {
	format(s string) string
	wrap(f formatter) formatter
}

type prefixFormatter struct {
	prefix string
	parent formatter
}

func newPrefixFormatter(prefix string) *prefixFormatter {
	p := new(prefixFormatter)
	p.prefix = prefix
	return p
}

func (p *prefixFormatter) wrap(f formatter) formatter {
	p.parent = f
	return p
}

func (p *prefixFormatter) format(s string) string {
	preformatted := s
	if p.parent != nil {
		preformatted = p.parent.format(s)
	}
	return fmt.Sprintf("%s%s", p.prefix, preformatted)
}

type timestampFormatter struct {
	parent formatter
}

func newTimestampFormatter() *timestampFormatter {
	return new(timestampFormatter)
}

func (t *timestampFormatter) wrap(f formatter) formatter {
	t.parent = f
	return t
}

func (t *timestampFormatter) format(s string) string {
	preformatted := s
	if t.parent != nil {
		preformatted = t.parent.format(s)
	}
	timeStr := time.Now().Format(time.RFC822)
	return fmt.Sprintf("%s %s", timeStr, preformatted)
}
