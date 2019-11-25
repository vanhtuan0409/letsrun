package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type command struct {
	c               *exec.Cmd
	id              string
	out             io.Writer
	outputFormatter formatter
	outputScanner   *bufio.Scanner
}

func splitArgs(cmd string) ([]string, error) {
	r := csv.NewReader(strings.NewReader(cmd))
	r.Comma = ' '
	fields, err := r.Read()
	if err != nil {
		return nil, err
	}
	return fields, nil
}

func newCmd(s string, id string, out io.Writer) (*command, error) {
	var err error
	parts, err := splitArgs(s)
	if err != nil {
		return nil, err
	}

	c := new(command)
	c.id = id
	c.out = out
	osCmd := exec.Command(parts[0], parts[1:]...)
	if c.outputScanner, err = getCmdOutputScanner(osCmd); err != nil {
		return nil, err
	}
	c.c = osCmd
	c.outputFormatter = newPrefixFormatter(fmt.Sprintf("[%s] ", c.id))
	return c, nil
}

func (c *command) captureOutput() {
	for c.outputScanner.Scan() {
		line := c.outputScanner.Text() + "\n"
		if c.outputFormatter != nil {
			line = c.outputFormatter.format((line))
		}
		fmt.Fprintf(c.out, line)
	}

}

func (c *command) Run() error {
	go c.captureOutput()
	return c.c.Run()
}
