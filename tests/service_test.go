package tests

import (
	"context"
	"eljur/tests/suite"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetGradesByMonthForUser(t *testing.T) {
	g, err := suite.GetGradesService()
	require.NoError(t, err)
	ctx := context.Background()
	grades, err := g.GetUserGradesByMonth(ctx, "test", 1, 1)
	require.NoError(t, err)

	fmt.Print("         ")

	fmt.Print("\n")
	for i, subject := range grades.SubjectsNames {
		fmt.Printf("%s: %+v \n", subject, grades.Grades[i])
	}

}
