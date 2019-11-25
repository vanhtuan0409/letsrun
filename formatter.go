package main

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type formatter interface {
	Format(s string) string
}

type prefixFormatter struct {
	prefix string
}

func newPrefixFormatter(prefix string) *prefixFormatter {
	p := new(prefixFormatter)
	p.prefix = prefix
	return p
}
