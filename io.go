package main

import (
	"bufio"
	"io"
	"os/exec"
)

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
