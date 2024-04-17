package converter

import (
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"strings"
)

func LogsJSONToXLS(logs string) (io.Reader, error) {
	file := excelize.NewFile()

	sheetIndex, err := file.NewSheet("logs")
	if err != nil {
		return nil, err
	}
	file.SetActiveSheet(sheetIndex)

	headers := []string{"time", "level", "message", "URL", "Method", "Remote Addr", "Body"}
	for i, header := range headers {
		if err := file.SetCellValue("logs", fmt.Sprintf("%s%d", string(rune(65+i)), 1), header); err != nil {
			return nil, err
		}
	}

	logsSp := strings.Split(logs, "\n")
	data := make([][]string, len(logsSp))

	for i, log := range logsSp {
		logM := make(map[string]string)
		if err := json.Unmarshal([]byte(log), &logM); err != nil {
			return nil, err
		}
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
		fmt.Println(dataRow, row)
		for j, col := range row {
			if err := file.SetCellValue("logs", fmt.Sprintf("%s%d", string(rune(65+j)), dataRow), col); err != nil {
				return nil, err
			}
		}
	}

	return file.WriteToBuffer()
}
