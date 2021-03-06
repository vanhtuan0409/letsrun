package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type formatter interface {
	format(s string) string
	wrap(f formatter) formatter
}

type noopFormatter struct {
}

func newNoopFormatter() *noopFormatter {
	return new(noopFormatter)
}

func (n *noopFormatter) wrap(f formatter) formatter {
	return n
}

func (n *noopFormatter) format(s string) string {
	return s
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

type colorFormatter struct {
	c      *color.Color
	parent formatter
}

func newColorFormatter(colorCode int) *colorFormatter {
	f := new(colorFormatter)
	f.c = color.New(color.Attribute(colorCode))
	return f
}

func (c *colorFormatter) wrap(f formatter) formatter {
	c.parent = f
	return c
}

func (c *colorFormatter) format(s string) string {
	preformatted := s
	if c.parent != nil {
		preformatted = c.parent.format(s)
	}
	return c.c.Sprint(preformatted)
}
