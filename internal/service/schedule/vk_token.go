package schedules

import (
	"eljur/internal/config"
	"eljur/pkg/tr"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func getVkToken(config config.VKSeverConfig) (string, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s:%d/get_token", config.IP, config.Port))
	if err != nil {
		return "", tr.Trace(err)
	}
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", errors.New(string(body))
	}
	return string(body), nil
}
