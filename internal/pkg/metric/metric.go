package metric

import (
	"eljur/internal/pkg/logger"
	"eljur/pkg/tr"
	"sync/atomic"
	"time"
)

type Metric struct {
	Logs            string
	Rps             int32
	RenderPerSecond int32
}

var (
	rps                 int32 = 0
	lastRps             int32 = 0
	renderPerSecond     int32 = 0
	lastRenderPerSecond int32 = 0
)

func GetMetric() (*Metric, error) {
	logs, err := getLogs()
	if err != nil {
		return nil, tr.Trace(err)
	}

	return &Metric{
		Logs:            logs,
		Rps:             (lastRps + rps + 1) / 2,
		RenderPerSecond: (lastRenderPerSecond + renderPerSecond + 1) / 2,
	}, nil
}

func getLogs() (string, error) {
	logs, err := logger.GetLogs()
	if err != nil {
		return "", tr.Trace(err)
	}
	return string(logs), nil
}

func CountRPS() {
	rpsT := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-rpsT.C:
				lastRps = rps
				atomic.AddInt32(&rps, -rps)
				lastRenderPerSecond = renderPerSecond
				atomic.AddInt32(&renderPerSecond, -renderPerSecond)
			}
		}
	}()
}

func HandleRequest() {
	atomic.AddInt32(&rps, 1)
}
func HandleRender() {
	atomic.AddInt32(&renderPerSecond, 1)
}
