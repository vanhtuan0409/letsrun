package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

var (
	delimiter = flag.String("F", ";;", "Commands delimiter")
)

type command struct {
	c *exec.Cmd
}

func newCmd(s string, out io.Writer) (*command, error) {
	parts, err := splitArgs(s)
	if err != nil {
		return nil, err
	}
	osCmd := exec.Command(parts[0], parts[1:]...)
	osCmd.Stdout = out
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
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		panic("Invalid arguments. Missing commands")
	}

	cmds := splitCommand(args[0], *delimiter)

	var wg sync.WaitGroup
	cmdList := []*command{}
	for i, cmd := range cmds {
		wg.Add(1)
		go func(index int, s string) {
			defer wg.Done()
			w := newWriter(index)
			c, err := newCmd(s, w)
			if err != nil {
				fmt.Println(err)
			}
			cmdList = append(cmdList, c)
			if err := c.Run(); err != nil {
				fmt.Println(err)
			}
		}(i, cmd)
	}

	sigs := make(chan os.Signal)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		for _, cmd := range cmdList {
			cmd.c.Process.Signal(sig)
		}
	}()

	wg.Wait()
}
