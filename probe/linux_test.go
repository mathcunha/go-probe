package probe

import (
	"testing"
	"time"
)

func TestStats(t *testing.T) {
	stats := new(Stats)
	GetAllStats(stats)
	//to load CPU data
	time.Sleep(2000 * time.Millisecond)

	GetAllStats(stats)
	t.Log(stats)
}
