package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func usage() {
	fmt.Printf("Usage: %s [OPTIONS] COMMANDS\n\n", os.Args[0])
	fmt.Print("Background command runner and combine output into stdout\n\n")
	flag.PrintDefaults()
}

type command struct {
	c      *exec.Cmd
	stdout *bufio.Scanner
	stderr *bufio.Scanner
}

func newCmd(s string) (*command, error) {
	parts, err := splitArgs(s)
	if err != nil {
		return nil, err
	}

	c := new(command)

	osCmd := exec.Command(parts[0], parts[1:]...)
	stdout, err := osCmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	c.stdout = bufio.NewScanner(stdout)
	stderr, err := osCmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	c.stderr = bufio.NewScanner(stderr)

	c.c = osCmd
	return c, nil
}

func (c *command) Run() error {
	go func() {
		for c.stdout.Scan() {
			fmt.Println(c.stdout.Text())
		}
	}()
	go func() {
		for c.stderr.Scan() {
			fmt.Println(c.stderr.Text())
		}
	}()
	return c.c.Run()
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

func main() {
	flag.Usage = usage
	flag.Parse()

	cmds := flag.Args()
	if len(cmds) < 1 {
		fmt.Println("Invalid arguments. Missing commands")
		flag.Usage()
		os.Exit(1)
	}

	var wg sync.WaitGroup
	cmdList := []*command{}
	for i, cmd := range cmds {
		wg.Add(1)
		go func(index int, s string) {
			defer wg.Done()
			c, err := newCmd(s)
			if err != nil {
				fmt.Println(err)
			}
			cmdList = append(cmdList, c)
			if err := c.Run(); err != nil {
				fmt.Println(err)
			}
		}(i, cmd)
	}

	wg.Wait()
}
