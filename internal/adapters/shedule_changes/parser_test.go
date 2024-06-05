package schedule_changes

import (
	"eljur/internal/config"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetWeekFromDocument(t *testing.T) {
	cnf, err := config.GetConfig("C:\\Users\\79212\\GolandProjects\\eljur\\config\\config.yaml")
	fmt.Printf("%+v", cnf)

	p := NewParser(cnf.ScheduleChanges)

	docsInf, err := p.GetListDocuments()
	require.NoError(t, err)

	doc, err := p.GetDocument(docsInf[0])
	require.NoError(t, err)

	week := p.GetWeekFromDocument(doc)
	fmt.Println(week)
}

func TestGetChangesFromDocument(t *testing.T) {
	cnf, err := config.GetConfig("C:\\Users\\79212\\GolandProjects\\eljur\\config\\config.yaml")
	fmt.Printf("%+v", cnf)

	p := NewParser(cnf.ScheduleChanges)

	docsInf, err := p.GetListDocuments()
	require.NoError(t, err)

	doc, err := p.GetDocument(docsInf[0])
	require.NoError(t, err)

	ch := p.GetChangesFromDocument(doc, cnf.Schedule.GroupName)
	fmt.Println(ch)
}
