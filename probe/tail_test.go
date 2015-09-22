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
	line := tail.OutChannel(true)

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
		}
	}
	fmt.Println("all done!")
}
