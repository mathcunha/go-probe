package probe

import (
	"fmt"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	tail := Load("/tmp/jmvmetrics.go")
	quit := make(chan bool)
	line := make(chan string)
	time.AfterFunc(4*time.Second, func() { quit <- true })
	tail.ReadLine(line)

loop:
	for {
		select {
		case s, ok := <-line:
			if !ok {
				break loop
			}
			fmt.Printf("%v\n", s)
		case <-quit:
			tail.Close()
			close(line)
		}
	}
	fmt.Println("all done!")
}
