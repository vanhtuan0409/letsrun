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
	wrap(f formatter)
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

func (p *prefixFormatter) wrap(f formatter) {
	p.parent = f
}

func (p *prefixFormatter) format(s string) string {
	preformatted := s
	if p.parent != nil {
		preformatted = p.parent.format(s)
	}
	return fmt.Sprintf("%s%s", p.prefix, preformatted)
}
