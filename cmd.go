package main

import (
	"bufio"
	"encoding/csv"
	"os/exec"
	"strings"
)

type command struct {
	c             *exec.Cmd
	id            string
	out           chan<- string
	outputScanner *bufio.Scanner
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

func newCmd(s string, id string, out chan<- string) (*command, error) {
	var err error
	parts, err := splitArgs(s)
	if err != nil {
		return nil, err
	}

	c := new(command)
	c.out = out
	osCmd := exec.Command(parts[0], parts[1:]...)
	if c.outputScanner, err = getCmdOutputScanner(osCmd); err != nil {
		return nil, err
	}
	c.c = osCmd
	return c, nil
}

func (c *command) Run() error {
	go func() {
		for c.outputScanner.Scan() {
			c.out <- c.outputScanner.Text()
		}
	}()
	return c.c.Run()
}
