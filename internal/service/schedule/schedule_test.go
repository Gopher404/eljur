package schedules

import (
	"context"
	"eljur/internal/config"
	"eljur/internal/storage"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetActualSchedule(t *testing.T) {
	cnf, err := config.GetConfig("C:\\Users\\79212\\GolandProjects\\eljur\\config\\config.yaml")
	fmt.Printf("%+v", cnf)
	require.NoError(t, err)
	s, err := storage.New(&cnf.DB)
	require.NoError(t, err)
	schedule := New(s, &cnf.Schedule)
	ctx := context.Background()
	sc, err := schedule.GetActualSchedule(ctx)
	fmt.Printf("%#v\n", sc)
	for _, day := range sc.Days {
		fmt.Printf("%#v\n", day)
	}
	require.NoError(t, err)

}
