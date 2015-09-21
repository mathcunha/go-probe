package probe

import (
	"fmt"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	tail := Load("/tmp/jmvmetrics.go")
	time.AfterFunc(time.Second, func() { tail.QuitChannel() <- true })

loop:
	for {
		select {
		case s, ok := <-tail.OutChannel():
			if !ok {
				break loop
			}
			fmt.Print(s)
		case <-tail.QuitChannel():
			tail.Close()
		}
	}
	fmt.Println("all done!")
}
