package models

type Change struct {
	Number     int8
	Auditorium string
	Name       string
	Teacher    string
}

type DocumentInfo struct {
	Title string `json:"title"`
	Ext   string `json:"ext"`
	Url   string `json:"url"`
	Size  int64  `json:"size"`
}

type DocumentsResp struct {
	Error *struct {
		Code int    `json:"error_code"`
		Msg  string `json:"error_msg"`
	} `json:"error"`
	Response struct {
		Count int             `json:"count"`
		Items []*DocumentInfo `json:"items"`
	} `json:"response"`
}
