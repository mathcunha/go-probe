package main

import (
	"flag"
	"github.com/mathcunha/amon"
	"github.com/mathcunha/amon/scheduler"
	"github.com/mathcunha/go-probe/probe"
	"regexp"
	"time"
)

type JVMMetric struct {
	Date                         time.Time `json:"@timestamp"`
	Milliseconds                 int64
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
	GcCountCopy                  int32
	GcTimeMillisCopy             int32
	GcCountMarkSweepCompact      int32
	GcTimeMillisMarkSweepCompact int32
	GcCount                      int32
	GcTimeMillis                 int32
	ThreadsNew                   int32
	ThreadsRunnable              int32
	ThreadsBlocked               int32
	ThreadsWaiting               int32
	ThreadsTimedWaiting          int32
	ThreadsTerminated            int32
	LogFatal                     int32
	LogError                     int32
	LogWarn                      int32
	LogInfo                      int32
}

var fileName string
var interval string

type task struct{}

func init() {
	flag.StringVar(&interval, "interval", "5s", "interval between metric collection")
	flag.StringVar(&fileName, "jvmmetrics", "/usr/local/hadoop/nodemanager-jvm-metrics.out", "hadoop jvm metrics file")
}

func (t *task) Interval() string {
	return interval
}

func (t *task) Run() {
	stats := new(probe.Stats)
	probe.GetAllStats(stats)
	probe.PostStats(stats)
}

func main() {
	flag.Parse()
	scheduler.Schedule([]scheduler.Task{scheduler.Task(&task{})})
	readMetrics()
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
			jvm := new(JVMMetric)
			amon.LoadAttributes(jvm, []byte(s))
			probe.PostStats(jvm)
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
	e.Milliseconds = amon.ParseInt64((*mp)["Milliseconds"])
	e.Date = time.Unix(0, e.Milliseconds*int64(time.Millisecond))
	e.ProcessName = (*mp)["ProcessName"]
	e.SessionId = (*mp)["e.SessionId"]
	e.ProcessName = (*mp)["ProcessName"]
	e.SessionId = (*mp)["SessionId"]
	e.Hostname = (*mp)["Hostname"]
	e.MemNonHeapUsedM = amon.ParseFloat((*mp)["MemNonHeapUsedM"])
	e.MemNonHeapCommittedM = amon.ParseFloat((*mp)["MemNonHeapCommittedM"])
	e.MemNonHeapMaxM = amon.ParseFloat((*mp)["MemNonHeapMaxM"])
	e.MemHeapUsedM = amon.ParseFloat((*mp)["MemHeapUsedM"])
	e.MemHeapCommittedM = amon.ParseFloat((*mp)["MemHeapCommittedM"])
	e.MemHeapMaxM = amon.ParseFloat((*mp)["MemHeapMaxM"])
	e.MemMaxM = amon.ParseFloat((*mp)["MemMaxM"])
	e.GcCountCopy = amon.ParseInt((*mp)["GcCountCopy"])
	e.GcTimeMillisCopy = amon.ParseInt((*mp)["GcTimeMillisCopy"])
	e.GcCountMarkSweepCompact = amon.ParseInt((*mp)["GcCountMarkSweepCompact"])
	e.GcTimeMillisMarkSweepCompact = amon.ParseInt((*mp)["GcTimeMillisMarkSweepCompact"])
	e.GcCount = amon.ParseInt((*mp)["GcCount"])
	e.GcTimeMillis = amon.ParseInt((*mp)["GcTimeMillis"])
	e.ThreadsNew = amon.ParseInt((*mp)["ThreadsNew"])
	e.ThreadsRunnable = amon.ParseInt((*mp)["ThreadsRunnable"])
	e.ThreadsBlocked = amon.ParseInt((*mp)["ThreadsBlocked"])
	e.ThreadsWaiting = amon.ParseInt((*mp)["ThreadsWaiting"])
	e.ThreadsTimedWaiting = amon.ParseInt((*mp)["ThreadsTimedWaiting"])
	e.ThreadsTerminated = amon.ParseInt((*mp)["ThreadsTerminated"])
	e.LogFatal = amon.ParseInt((*mp)["LogFatal"])
	e.LogError = amon.ParseInt((*mp)["LogError"])
	e.LogWarn = amon.ParseInt((*mp)["LogWarn"])
	e.LogInfo = amon.ParseInt((*mp)["LogInfo"])
}
