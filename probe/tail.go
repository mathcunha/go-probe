package probe

import (
	"fmt"
	"log"
	"os/exec"
)

type TailCommand interface {
	OutChannel() (out chan string)
	QuitChannel() (quit chan bool)
	Close()
}

type Tail struct {
	ch   chan string
	quit chan bool
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
	tail := &Tail{make(chan string), make(chan bool), cmd}
	go func() {
		buf := make([]byte, 1024)
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

func (t *Tail) QuitChannel() (quit chan bool) {
	return t.quit
}

func (t *Tail) OutChannel() (out chan string) {
	return t.ch
}

func (t *Tail) Close() {
	t.cmd.Process.Kill()
}
