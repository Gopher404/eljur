package schedule_changes

import (
	"eljur/internal/config"
	"eljur/internal/domain/models"
	"eljur/pkg/tr"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func newVKAPI(cnf config.VKAPIConfig) *vkAPI {
	return &vkAPI{
		cnf: cnf,
	}
}

type vkAPI struct {
	cnf config.VKAPIConfig
}

func (api *vkAPI) DocsGet(token string) (*models.DocumentsResp, error) {
	url := fmt.Sprintf("https://api.vk.com/method/docs.get?owner_id=%s&v=%s&access_token=%s",
		api.cnf.GroupId, api.cnf.Version, token)
	resp, err := http.Get(url)

	if err != nil {
		return nil, tr.Trace(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, tr.Trace(err)
	}

	var docsResp models.DocumentsResp

	if err := json.Unmarshal(body, &docsResp); err != nil {
		return nil, tr.Trace(err)
	}
	return &docsResp, nil
}
