package main

import (
	"flag"
	"github.com/mathcunha/amon"
	"github.com/mathcunha/go-probe/probe"
	"log"
	"regexp"
	"time"
)

type JVMMetric struct {
	Date                         time.Time
	Milliseconds                 int32
	ProcessName                  string
	SessionId                    string
	Hostname                     string
	MemNonHeapUsedM              float32
	MemNonHeapCommittedM         float32
	MemNonHeapMaxM               float32
	MemHeapUsedM                 float32
	MemHeapCommittedM            float32
	MemHeapMaxM                  float32
	MemMaxM                      float32
	GcCountCopy                  int
	GcTimeMillisCopy             int
	GcCountMarkSweepCompact      int
	GcTimeMillisMarkSweepCompact int
	GcCount                      int
	GcTimeMillis                 int
	ThreadsNew                   int
	ThreadsRunnable              int
	ThreadsBlocked               int
	ThreadsWaiting               int
	ThreadsTimedWaiting          int
	ThreadsTerminated            int
	LogFatal                     int
	LogError                     int
	LogWarn                      int
	LogInfo                      int
}

var fileName string

func init() {
	flag.StringVar(&fileName, "jvmmetrics", "/usr/local/hadoop/nodemanager-jvm-metrics.out", "hadoop jvm metrics file")
}

func main() {
	readMetrics()
	log.Println("hello world!")
}

func readMetrics() {

	tail := probe.Load(fileName)
	line := tail.OutChannel(true)

loop:
	for {
		select {
		case s, ok := <-line:
			if !ok {
				break loop
			}
			log.Printf("%v\n", s)
			jvm := new(JVMMetric)
			amon.LoadAttributes(jvm, []byte(s))
		}
	}

}

func (e *JVMMetric) BuildEvents(event *amon.Event) *[]amon.Event {
	event.Field = e
	events := []amon.Event{*event}
	return &events
}

func (e JVMMetric) Patterns() []*regexp.Regexp {
	return []*regexp.Regexp{regexp.MustCompile("(?P<Milliseconds>[0-9]+) jvm.JvmMetrics: Context=jvm, ProcessName=(?P<ProcessName>.*), SessionId=(?P<SessionId>.*), Hostname=(?P<Hostname>.*), MemNonHeapUsedM=(?P<MemNonHeapUsedM>.*), MemNonHeapCommittedM=(?P<MemNonHeapCommittedM>.*), MemNonHeapMaxM=(?P<MemNonHeapMaxM>.*), MemHeapUsedM=(?P<MemHeapUsedM>.*), MemHeapCommittedM=(?P<MemHeapCommittedM>.*), MemHeapMaxM=(?P<MemHeapMaxM>.*), MemMaxM=(?P<MemMaxM>.*), GcCountCopy=(?P<GcCountCopy>[0-9]+), GcTimeMillisCopy=(?P<GcTimeMillisCopy>[0-9]+), GcCountMarkSweepCompact=(?P<GcCountMarkSweepCompact>[0-9]+), GcTimeMillisMarkSweepCompact=(?P<GcTimeMillisMarkSweepCompact>[0-9]+), GcCount=(?P<GcCount>[0-9]+), GcTimeMillis=(?P<GcTimeMillis>[0-9]+), ThreadsNew=(?P<ThreadsNew>[0-9]+), ThreadsRunnable=(?P<ThreadsRunnable>[0-9]+), ThreadsBlocked=(?P<ThreadsBlocked>[0-9]+), ThreadsWaiting=(?P<ThreadsWaiting>[0-9]+), ThreadsTimedWaiting=(?P<ThreadsTimedWaiting>[0-9]+), ThreadsTerminated=(?P<ThreadsTerminated>[0-9]+), LogFatal=(?P<LogFatal>[0-9]+), LogError=(?P<LogError>[0-9]+), LogWarn=(?P<LogWarn>[0-9]+), LogInfo=(?P<LogInfo>[0-9]+)")}
}

func (e *JVMMetric) Load(mp *map[string]string) {
	e.Milliseconds = amon.ParseInt((*mp)["Milliseconds"])
	e.ProcessName = (*mp)["ProcessName"]
	log.Println((*mp))
}
