package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"sync"
)

var (
	colorized = flag.Bool("c", true, "Print colorized output")
	timestamp = flag.Bool("t", false, "Print timestamp to output")
	indicator = flag.Bool("i", true, "Printa command index indicator")
)

func usage() {
	fmt.Printf("Usage: %s [OPTIONS] COMMANDS\n\n", os.Args[0])
	fmt.Print("Background command runner and combine output into stdout\n\n")
	fmt.Println("Options:")
	flag.PrintDefaults()
}

func formatCmdOutput(cmd *command, index int) {
	cl := (index % 6) + 31 // Limit number of color due to available constants in color lib
	cmd.outputFormatter = newNoopFormatter()
	if *indicator {
		cmd.outputFormatter = newPrefixFormatter(fmt.Sprintf("[%d] ", index+1))
	}
	if *timestamp {
		cmd.wrapFormatter(newTimestampFormatter())
	}
	if *colorized {
		cmd.wrapFormatter(newColorFormatter(cl))
	}
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
	for i, cmd := range cmds {
		wg.Add(1)
		go func(index int, s string) {
			defer wg.Done()
			c, err := newCmd(s, strconv.Itoa(index))
			if err != nil {
				fmt.Println(err)
			}
			c.setOutput(os.Stdout)
			formatCmdOutput(c, index)

			if err := c.Run(); err != nil {
				fmt.Println(err)
			}
		}(i, cmd)
	}

	wg.Wait()
}
