package probe

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type TailCommand interface {
	OutChannel(line bool) (out chan string)
	Close() error
}

type Tail struct {
	ch   chan string
	line chan string
	cmd  *exec.Cmd
}

func Load(file string) TailCommand {
	cmd := exec.Command("tail", "-n", "0", "-f", file)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	tail := &Tail{make(chan string), make(chan string), cmd}
	go func() {
		buf := make([]byte, 8)
		for {
			n, err := stdout.Read(buf)
			if n != 0 {
				tail.ch <- string(buf[:n])
			}
			if err != nil {
				break
			}
		}
		fmt.Println("Goroutine finished")
		close(tail.ch)
		close(tail.line)
	}()
	return tail
}

func (t *Tail) OutChannel(line bool) (out chan string) {
	if line {
		go func() {
			strLine := ""
		loop:
			for {
				select {
				case s, ok := <-t.ch:
					if !ok {
						break loop
					}
					if index := strings.Index(s, "\n"); -1 != index {
						//log.Printf("s=[%v], index=%v, len=%v\n", s, index, len(s))
						strLine += s[:index]
						t.line <- strLine
						strLine = ""
					} else {
						strLine += s
					}
				}
			}
		}()
		return t.line
	}
	return t.ch
}

func (t *Tail) Close() error {
	return t.cmd.Process.Kill()
}
