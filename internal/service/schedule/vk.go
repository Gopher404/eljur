package schedules

import (
	"eljur/pkg/tr"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func newVKAPI(apiVersion string) *vkAPI {
	return &vkAPI{
		apiVersion: apiVersion,
	}
}

type vkAPI struct {
	apiVersion string
}

func (api *vkAPI) DocsGet(token string, groupId string) (*documentsResp, error) {
	url := fmt.Sprintf("https://api.vk.com/method/docs.get?owner_id=%s&v=%s&access_token=%s",
		groupId, api.apiVersion, token)
	resp, err := http.Get(url)

	if err != nil {
		return nil, tr.Trace(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, tr.Trace(err)
	}

	var docsResp documentsResp

	if err := json.Unmarshal(body, &docsResp); err != nil {
		return nil, tr.Trace(err)
	}
	return &docsResp, nil
}

type documentInfo struct {
	Title string `json:"title"`
	Ext   string `json:"ext"`
	Url   string `json:"url"`
	Size  int64  `json:"size"`
}

type documentsResp struct {
	Error *struct {
		Code int    `json:"error_code"`
		Msg  string `json:"error_msg"`
	} `json:"error"`
	Response struct {
		Count int             `json:"count"`
		Items []*documentInfo `json:"items"`
	} `json:"response"`
}
