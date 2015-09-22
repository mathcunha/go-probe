package probe

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type TailCommand interface {
	outChannel() (out chan string)
	ReadLine(chLine chan string)
	Close() error
}

type Tail struct {
	ch  chan string
	cmd *exec.Cmd
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
	tail := &Tail{make(chan string), cmd}
	go func() {
		buf := make([]byte, 1)
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
	}()
	return tail
}

func (t *Tail) outChannel() (out chan string) {
	return t.ch
}

func (t *Tail) ReadLine(chLine chan string) {
	go func() {
		strLine := ""
	loop:
		for {
			select {
			case s, ok := <-t.ch:
				if !ok {
					break loop
				}
			hasNext:
				for {
					if index := strings.Index(s, "\n"); -1 != index {
						strLine += s[:index]
						chLine <- strLine
						strLine = ""
						if len(s) > index+1 {
							s = s[index+1 : len(s)]
						} else {
							break hasNext
						}
					} else {
						strLine += s
						break hasNext
					}
				}
			}
		}
	}()
}

func (t *Tail) Close() error {
	return t.cmd.Process.Kill()
}
