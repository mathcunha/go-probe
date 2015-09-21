package probe

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

func ReadJVMMetrics() {
	cmd := exec.Command("tail", "-n", "0", "-f", "/tmp/foo.txt")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	ch := make(chan string)
	quit := make(chan bool)
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if n != 0 {
				ch <- string(buf[:n])
			}
			if err != nil {
				break
			}
		}
		fmt.Println("Goroutine finished")
		close(ch)
	}()

	time.AfterFunc(time.Second, func() { quit <- true })

loop:
	for {
		select {
		case s, ok := <-ch:
			if !ok {
				break loop
			}
			fmt.Print(s)
		case <-quit:
			cmd.Process.Kill()
		}
	}
	fmt.Println("all done!")
}
