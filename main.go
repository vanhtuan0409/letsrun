package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"sync"
)

func usage() {
	fmt.Printf("Usage: %s [OPTIONS] COMMANDS\n\n", os.Args[0])
	fmt.Print("Background command runner and combine output into stdout\n\n")
	flag.PrintDefaults()
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
			c, err := newCmd(s, strconv.Itoa(index))
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
