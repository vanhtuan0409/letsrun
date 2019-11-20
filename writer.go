package main

import (
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/segmentio/textio"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func newWriter(w io.Writer, cmdIndex int) io.Writer {
	prefix := fmt.Sprintf("[%d] ", cmdIndex+1)
	prefixer := textio.NewPrefixWriter(w, prefix)
	return prefixer
}
