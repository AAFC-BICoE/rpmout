package main

import (
	"os/exec"
//	"fmt"
	"log"
	"bufio"
)

func runExec(commandAndArgs []string) (chan string, chan bool){
	command := commandAndArgs[0]
	args := commandAndArgs[1:]

	lines := make(chan string, 500)
	done := make(chan bool)

	cmd := exec.Command(command, args[0:]...)
	//cmd := exec.Command("rpm", "-q")
	stdout, err := cmd.StdoutPipe()
	
	
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	go func(){
		r := bufio.NewReader(stdout)
		s, e := Readln(r)
		for e == nil {
			//fmt.Println(s)
			lines <- s
			s,e = Readln(r)
		}
		stdout.Close()
		if err := cmd.Wait(); err != nil {
			log.Fatal(err)
		}

		done <- true
	}()

	return lines, done
}


func Readln(r *bufio.Reader) (string, error) {

  var (isPrefix bool = true
       err error = nil
       line, ln []byte
      )
  for isPrefix && err == nil {
      line, isPrefix, err = r.ReadLine()
      ln = append(ln, line...)
  }
  return string(ln),err
}
