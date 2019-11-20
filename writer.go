package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/segmentio/textio"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func newWriter(cmdIndex int) io.Writer {
	prefix := fmt.Sprintf("[%d] ", cmdIndex+1)
	prefixer := textio.NewPrefixWriter(os.Stdout, prefix)
	return prefixer
}
