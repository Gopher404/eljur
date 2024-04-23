package metric

import (
	"bytes"
	"eljur/internal/pkg/logger"
	"eljur/pkg/tr"
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
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

func GetXLSXLogs() (io.ReadSeeker, error) {
	logs, err := logger.GetLogs()
	if err != nil {
		return nil, tr.Trace(err)
	}

	file := excelize.NewFile()

	sheetIndex, err := file.NewSheet("logs")
	if err != nil {
		return nil, tr.Trace(err)
	}
	file.SetActiveSheet(sheetIndex)

	headers := []string{"time", "level", "message", "URL", "Method", "Remote Addr", "Body"}
	for i, header := range headers {
		if err := file.SetCellValue("logs", fmt.Sprintf("%s%d", string(rune(65+i)), 1), header); err != nil {
			return nil, tr.Trace(err)
		}
	}

	logsM := make([]map[string]string, 0)

	if err := json.Unmarshal(logs, &logsM); err != nil {
		return nil, tr.Trace(err)
	}
	data := make([][]string, len(logsM))

	for i, logM := range logsM {
		logS := []string{logM["time"], logM["level"], logM["msg"], logM["URL"], logM["Method"], logM["Remote"], logM["body"]}

		for _, k := range []string{"time", "level", "msg", "URL", "Method", "Remote", "body"} {
			delete(logM, k)
		}
		for _, v := range logM {
			logS = append(logS, v)
		}
		data[i] = logS
	}

	for i, row := range data {
		dataRow := i + 2
		for j, col := range row {
			if err := file.SetCellValue("logs", fmt.Sprintf("%s%d", string(rune(65+j)), dataRow), col); err != nil {
				return nil, tr.Trace(err)
			}
		}
	}

	buf, err := file.WriteToBuffer()
	if err != nil {
		return nil, tr.Trace(err)
	}
	b, err := io.ReadAll(buf)
	if err != nil {
		return nil, tr.Trace(err)
	}
	return bytes.NewReader(b), nil
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
