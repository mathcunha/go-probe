package probe

import (
	"fmt"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	tail := Load("/tmp/jmvmetrics.go")
	quit := make(chan bool)
	time.AfterFunc(4*time.Second, func() { quit <- true })

loop:
	for {
		select {
		case s, ok := <-tail.OutChannel(true):
			if !ok {
				break loop
			}
			fmt.Printf("%v\n", s)
		case <-quit:
			tail.Close()
		}
	}
	fmt.Println("all done!")
}
