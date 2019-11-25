package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os/exec"
	"strings"
)

type command struct {
	c               *exec.Cmd
	id              string
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

func newCmd(s string, id string) (*command, error) {
	var err error
	parts, err := splitArgs(s)
	if err != nil {
		return nil, err
	}

	c := new(command)
	c.id = id
	osCmd := exec.Command(parts[0], parts[1:]...)
	if c.outputScanner, err = getCmdOutputScanner(osCmd); err != nil {
		return nil, err
	}
	c.c = osCmd
	c.outputFormatter = newPrefixFormatter(fmt.Sprintf("[%s] ", c.id))
	return c, nil
}

func (c *command) Run() error {
	go func() {
		for c.outputScanner.Scan() {
			line := c.outputScanner.Text()
			if c.outputFormatter != nil {
				fmt.Println(c.outputFormatter.format(line))
			}
		}
	}()
	return c.c.Run()
}
