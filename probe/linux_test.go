package probe

import (
	"testing"
	"time"
)

func TestStats(t *testing.T) {
	stats := new(Stats)
	getAllStats(stats)
	//to load CPU data
	time.Sleep(2000 * time.Millisecond)

	getAllStats(stats)
	t.Log(stats)
}
