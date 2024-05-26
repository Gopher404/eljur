package schedules

import (
	"bytes"
	"eljur/pkg/tr"
	"encoding/json"
	"fmt"
	"github.com/ledongthuc/pdf"
	"io"
	"net/http"
	"strconv"
	"strings"
	"unicode"
)

func newParser(vkGroupID int64, vkAccessToken string) *Parser {
	pdf.DebugOn = true
	return &Parser{
		vkGroupID:     vkGroupID,
		vkAccessToken: vkAccessToken,
	}
}

type Parser struct {
	vkGroupID     int64
	vkAccessToken string
}

type document struct {
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
		Count int        `json:"count"`
		Items []document `json:"items"`
	} `json:"response"`
}

const vkGetDocsUrlTmpl = "https://api.vk.com/method/docs.get?owner_id=%d&v=5.236&access_token=%s"

func (p *Parser) getListDocuments() ([]document, error) {
	resp, err := http.Get(fmt.Sprintf(vkGetDocsUrlTmpl, p.vkGroupID, p.vkAccessToken))
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
	if docsResp.Error != nil {
		return nil, fmt.Errorf("vk error: code %d msg %s", docsResp.Error.Code, docsResp.Error.Msg)
	}
	fmt.Printf("%+v\n", docsResp)
	return docsResp.Response.Items, nil
}

func (p *Parser) getWeekFromDocument(doc string) string {
	const searchText = "("
	idx := strings.Index(doc, searchText)

	res := ""
	fReading := false
	for _, b := range doc[idx+len(searchText):] {
		a := string(b)
		if fReading {
			if a == ")" {
				break
			}
			res += a
		}

		if a == "-" || a == "–" {
			fReading = true
		}
	}

	return strings.ToLower(strings.ReplaceAll(res, " ", ""))
}

type change struct {
	Number     int8
	Auditorium string
	Name       string
	Teacher    string
}

func (p *Parser) getChangesFromDocument(doc string, groupName string) []change {
	var changes []change
	//var stringChanges []string
	c := strings.Count(doc, groupName)
	for i := 0; i < c; i++ {
		idx := strings.Index(doc, groupName)
		doc = doc[idx+len(groupName):]

		// find next group name
		countUpper := 0
		countDigit := 0
		stringChange := ""

		for _, b := range doc {
			a := string(b)
			stringChange += a
			if unicode.IsUpper(b) {
				countUpper++
			} else if unicode.IsDigit(b) && countUpper > 0 {
				countDigit++
			} else {
				countUpper, countDigit = 0, 0
			}
			if countUpper > 0 && countDigit == 2 {
				stringChange = stringChange[:len(stringChange)-countUpper-countDigit-2]
				break
			}
		}
		stringChange = strings.Trim(stringChange, " ")

		var ch change

		number, _ := strconv.Atoi(string(stringChange[0]))
		ch.Number = int8(number)

		stringChange = strings.Trim(stringChange[1:], " ")

		// find names

		isInits := func(s string) bool {
			findDot, findUpper := false, false
			for _, a := range s {
				if unicode.IsUpper(a) {
					findUpper = true
				} else if string(a) == "." {
					findDot = true
				} else {
					return false
				}
			}
			return findDot && findUpper
		}

		sl := strings.Split(stringChange, " ")
		sl = filterSlice(sl)
		ch.Auditorium = sl[len(sl)-1]
		name2 := ""
		fi := 0
		if sl[len(sl)-1] == "НИЧЕГО" {
			ch.Auditorium = ""
			ch.Name = "НИЧЕГО"

		} else if sl[len(sl)-3] == "Замена" {
			ch.Name = "замена"

		} else {
			if isInits(sl[len(sl)-2]) {
				name2 = sl[len(sl)-3] + " " + sl[len(sl)-2]
			} else {
				name2 = sl[len(sl)-2]
			}
			fi = len(sl) - 4
			if idx := strings.Index(ch.Auditorium, ","); idx > 0 {
				if isInits(sl[len(sl)-4]) {
					name2 = sl[len(sl)-5] + " " + sl[len(sl)-4] + " " + name2
					fi = len(sl) - 6
				} else {
					name2 = sl[len(sl)-4] + " " + name2
					fi = len(sl) - 5
				}
			}
			ch.Teacher = name2
			for i := fi; i >= 0; i-- {
				if !isInits(sl[i]) {
					ch.Name = sl[i] + " " + ch.Name
				} else {
					break
				}
			}
		}

		changes = append(changes, ch)

	}
	return changes
}

func (p *Parser) getDocument(docInfo document) (string, error) {
	resp, err := http.Get(docInfo.Url)
	if err != nil {
		return "", tr.Trace(err)
	}

	readerAt, err := readerToReaderAt(resp.Body)
	if err != nil {
		return "", tr.Trace(err)
	}
	doc, err := pdf.NewReader(readerAt, docInfo.Size)
	if err != nil {
		return "", tr.Trace(err)
	}
	r, err := doc.GetPlainText()
	if err != nil {
		return "", tr.Trace(err)
	}
	b, err := io.ReadAll(r)
	if err != nil {
		return "", tr.Trace(err)
	}
	return string(b), nil
}

func readerToReaderAt(r io.Reader) (io.ReaderAt, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(buf), nil
}

func filterSlice(sl []string) (res []string) {
MAINLOOP:
	for _, s := range sl {
		if len(s) == 0 {
			continue
		}
		for _, a := range s {
			if string(a) != " " {
				res = append(res, s)
				continue MAINLOOP
			}
		}

	}
	return
}
