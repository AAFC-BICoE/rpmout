package main

import (
	"os/exec"
	//	"fmt"
	"bufio"
	"log"
)

type Line struct {
	val string
}

// run program and write each line from stdout into a channel
func runExec(commandAndArgs []string) chan *Line {
	command := commandAndArgs[0]
	args := commandAndArgs[1:]

	lines := make(chan *Line, 1024)

	cmd := exec.Command(command, args[0:]...)
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	go func() {
		var si *Line

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			si = new(Line)
			si.val = line
			lines <- si
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		close(lines)
		stdout.Close()
		if err := cmd.Wait(); err != nil {
			log.Fatal(err)
		}
	}()
	return lines
}

func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}
