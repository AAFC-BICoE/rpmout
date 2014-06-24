package main

import (
	"os/exec"
	//	"fmt"
	"bufio"
	"log"
)

type StringInfo struct {
	val  string
	done bool
}

func runExec(commandAndArgs []string) chan *StringInfo {
	command := commandAndArgs[0]
	args := commandAndArgs[1:]

	lines := make(chan *StringInfo, 1024)

	cmd := exec.Command(command, args[0:]...)
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	go func() {
		var si *StringInfo
		r := bufio.NewReader(stdout)
		s, e := Readln(r)
		for e == nil {
			si = new(StringInfo)
			si.val = s
			si.done = false
			lines <- si
			s, e = Readln(r)
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
