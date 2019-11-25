package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

func usage() {
	fmt.Printf("Usage: %s [OPTIONS] COMMANDS\n\n", os.Args[0])
	fmt.Print("Background command runner and combine output into stdout\n\n")
	flag.PrintDefaults()
}

type command struct {
	c             *exec.Cmd
	id            string
	out           chan<- string
	outputScanner *bufio.Scanner
}

func getCmdOutputScanner(cmd *exec.Cmd) (*bufio.Scanner, error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	return bufio.NewScanner(io.MultiReader(stdout, stderr)), nil
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

	outChan := make(chan string)
	defer close(outChan)

	var wg sync.WaitGroup
	cmdList := []*command{}
	for i, cmd := range cmds {
		wg.Add(1)
		go func(index int, s string) {
			defer wg.Done()
			c, err := newCmd(s, strconv.Itoa(index), outChan)
			if err != nil {
				fmt.Println(err)
			}
			cmdList = append(cmdList, c)
			if err := c.Run(); err != nil {
				fmt.Println(err)
			}
		}(i, cmd)
	}

	go func() {
		for line := range outChan {
			fmt.Println(line)
		}
	}()

	wg.Wait()
}
