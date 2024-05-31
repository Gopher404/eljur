package schedules

import (
	"bytes"
	"eljur/internal/config"
	"eljur/pkg/tr"
	"fmt"
	"github.com/ledongthuc/pdf"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type docsGetter interface {
	DocsGet(token string, groupId string) (*documentsResp, error)
}

func newDocsWithCache(api docsGetter, cacheTTL time.Duration) docsGetter {
	docsGetrWithCache := &docsGetterWithCache{
		api:   api,
		cache: make(map[string]*documentsResp),
	}
	go func() {
		time.Sleep(cacheTTL)
		for k := range docsGetrWithCache.cache {
			delete(docsGetrWithCache.cache, k)
		}
	}()
	return docsGetrWithCache
}

type docsGetterWithCache struct {
	api   docsGetter
	cache map[string]*documentsResp
}

func (c *docsGetterWithCache) DocsGet(token string, groupId string) (*documentsResp, error) {
	resp, ok := c.cache[token+"&"+groupId]
	if ok {
		return resp, nil
	}

	resp, err := c.api.DocsGet(token, groupId)
	if err != nil {
		return nil, tr.Trace(err)
	}

	c.cache[token+"&"+groupId] = resp
	return resp, nil
}

func newParser(api docsGetter, vkServerConf config.VKSeverConfig, cacheTTL time.Duration) *Parser {
	pdf.DebugOn = true

	startServer()
	time.Sleep(time.Second)

	parser := &Parser{
		vkAPI:         newDocsWithCache(api, cacheTTL),
		vkServerConf:  vkServerConf,
		documentCache: make(map[string]string),
	}
	go func() {
		for {
			time.Sleep(cacheTTL)
			for _, k := range parser.documentCache {
				delete(parser.documentCache, k)
			}
		}
	}()
	token, _ := getVkToken(vkServerConf)
	parser.token = token
	return parser
}

type Parser struct {
	vkAPI         docsGetter
	token         string
	vkServerConf  config.VKSeverConfig
	documentCache map[string]string
}

func newVkErr(status int, msg string) error {
	return fmt.Errorf("vk error: code %d msg %s", status, msg)
}

func (p *Parser) getListDocuments(groupId string) ([]*documentInfo, error) {
	resp, err := p.vkAPI.DocsGet(p.token, groupId)
	if err != nil {
		return nil, tr.Trace(err)
	}

	if resp.Error != nil {
		if resp.Error.Code == 5 {
			token, err := getVkToken(p.vkServerConf)
			if err != nil {
				return nil, tr.Trace(fmt.Errorf("error get token: %e, > %e",
					err, newVkErr(resp.Error.Code, resp.Error.Msg)))
			}
			p.token = token

		} else {
			return nil, tr.Trace(newVkErr(resp.Error.Code, resp.Error.Msg))
		}
	}
	resp, err = p.vkAPI.DocsGet(p.token, groupId)
	if err != nil {
		return nil, tr.Trace(err)
	}
	if resp.Error != nil {
		return nil, tr.Trace(newVkErr(resp.Error.Code, resp.Error.Msg))
	}
	return resp.Response.Items, nil
}

func (*Parser) getDateFromDocInfo(docInfo *documentInfo) (time.Time, bool) {
	lensDigit := []int{0, 1, 3, 4, 6, 7, 8, 9}
	lensDot := []int{2, 5}
	dateS := ""
	for _, r := range docInfo.Title {
		a := string(r)
		if unicode.IsDigit(r) && slices.Index(lensDigit, len(dateS)) > -1 {
			dateS += a
			continue
		}
		if a == "." && slices.Index(lensDot, len(dateS)) > -1 {
			dateS += a
			continue
		}
		if len(dateS) == 10 {
			break
		}
		dateS = ""
	}

	if len(dateS) != 10 {
		return time.Time{}, false
	}

	date, err := time.Parse(dateLayout, dateS)
	if err != nil {
		return time.Time{}, false
	}
	return date, true
}

func (*Parser) getWeekFromDocument(doc string) string {
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
		s2 := sl[len(sl)-2]
		if strings.Index(s2, "/") > -1 || (len(s2) < 4 && !isInits(s2)) {
			sl = append(sl[:len(sl)-2], sl[len(sl)-1])
			ch.Auditorium = s2 + ch.Auditorium
		}
		name2 := ""
		fi := 0
		if strings.ToLower(sl[len(sl)-1]) == "ничего" {
			ch.Auditorium = ""
			ch.Name = "ничего"

		} else if strings.ToLower(sl[len(sl)-3]) == "замена" {
			ch.Name = "з/а"

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
				if !isInits(sl[i]) && strings.ToLower(sl[i]) != "ничего" {
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

func (p *Parser) getDocument(docInfo *documentInfo) (string, error) {
	docS, ok := p.documentCache[docInfo.Url]
	if ok {
		return docS, nil
	}
	resp, err := http.Get(docInfo.Url)
	if err != nil {
		return "", tr.Trace(err)
	}
	if resp.StatusCode != 200 {
		return "", tr.Trace(fmt.Errorf("error get file, status: %d", resp.StatusCode))
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
	docS = string(b)
	p.documentCache[docInfo.Url] = docS
	return docS, nil
}

func readerToReaderAt(r io.Reader) (io.ReaderAt, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(buf), nil
}

func filterSlice(sl []string) (res []string) {
	for _, s := range sl {
		if len(s) == 0 {
			continue
		}
		f := false
		for _, a := range s {
			if string(a) != " " {
				f = true
				break
			}
		}
		if f {
			res = append(res, s)
		}
	}
	return
}
