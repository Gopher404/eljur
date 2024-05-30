package schedules

import (
	"eljur/internal/config"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetWeekFromDocument(t *testing.T) {
	cnf, err := config.GetConfig("C:\\Users\\79212\\GolandProjects\\eljur\\config\\config.yaml")
	fmt.Printf("%+v", cnf)
	api := newVKAPI(cnf.Schedule.VKAPI.Version)
	p := newParser(api, cnf.Schedule.VKSever, cnf.Schedule.VKAPI.CacheTTL)

	docsInf, err := p.getListDocuments(cnf.Schedule.VKAPI.GroupId)
	require.NoError(t, err)

	doc, err := p.getDocument(docsInf[0])
	require.NoError(t, err)

	week := p.getWeekFromDocument(doc)
	fmt.Println(week)
}

func TestGetChangesFromDocument(t *testing.T) {
	cnf, err := config.GetConfig("C:\\Users\\79212\\GolandProjects\\eljur\\config\\config.yaml")
	fmt.Printf("%+v", cnf)
	api := newVKAPI(cnf.Schedule.VKAPI.Version)
	p := newParser(api, cnf.Schedule.VKSever, cnf.Schedule.VKAPI.CacheTTL)

	docsInf, err := p.getListDocuments(cnf.Schedule.VKAPI.GroupId)
	require.NoError(t, err)

	doc, err := p.getDocument(docsInf[0])
	require.NoError(t, err)

	ch := p.getChangesFromDocument(doc, cnf.Schedule.GroupName)
	fmt.Println(ch)
}
