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
	out     chan string
	outLine chan string
	cmd     *exec.Cmd
}

//Loads a file using the command tail and sends the output to a channel
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
		buf := make([]byte, 2048)
		for {
			n, err := stdout.Read(buf)
			if n != 0 {
				tail.out <- string(buf[:n])
			}
			if err != nil {
				break
			}
		}
		fmt.Println("Goroutine finished")
		close(tail.out)
		close(tail.outLine)
	}()
	return tail
}

//if line==true this method returns a channel of lines finished by \n
func (t *Tail) OutChannel(line bool) (out chan string) {
	if line {
		out = t.outLine
		go func() {
			strLine := ""
		loop:
			for {
				select {
				case s, ok := <-t.out:
					if !ok {
						break loop
					}
				hasNext:
					for {
						if index := strings.Index(s, "\n"); -1 != index {
							strLine += s[:index]
							out <- strLine
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
	} else {
		out = t.out
	}
	return
}

func (t *Tail) Close() error {
	return t.cmd.Process.Kill()
}
