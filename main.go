package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var (
	delimiter = flag.String("F", ";;", "Commands delimiter")
)

func usage() {
	fmt.Printf("Usage: %s [OPTIONS] COMMANDS\n\n", os.Args[0])
	fmt.Print("Background command runner and combine output into stdout\n\n")
	fmt.Println("Options:")
	flag.PrintDefaults()
}

type command struct {
	c *exec.Cmd
}

func newCmd(s string, out io.Writer, outErr io.Writer) (*command, error) {
	parts, err := splitArgs(s)
	if err != nil {
		return nil, err
	}
	osCmd := exec.Command(parts[0], parts[1:]...)
	osCmd.Stdout = out
	osCmd.Stderr = outErr
	c := new(command)
	c.c = osCmd
	return c, nil
}

func (c *command) Run() error {
	return c.c.Run()
}

func splitCommand(s string, d string) []string {
	d = strings.TrimSpace(d)
	cmds := strings.Split(s, d)
	for i := range cmds {
		cmds[i] = strings.TrimSpace(cmds[i])
	}
	return cmds
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

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Invalid arguments. Missing commands")
		flag.Usage()
		os.Exit(1)
	}

	cmds := splitCommand(args[0], *delimiter)

	var wg sync.WaitGroup
	cmdList := []*command{}
	for i, cmd := range cmds {
		wg.Add(1)
		go func(index int, s string) {
			defer wg.Done()
			out := newWriter(os.Stdout, index)
			outErr := newWriter(os.Stderr, index)
			c, err := newCmd(s, out, outErr)
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
