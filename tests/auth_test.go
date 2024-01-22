package tests

import (
	"context"
	"eljur/internal/domain/models"
	"eljur/tests/suite"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRegLog(t *testing.T) {
	auth, err := suite.GetAuthClient()
	require.NoError(t, err)
	ctx := context.Background()

	login := "test"
	pass := "pass"

	err = auth.Register(ctx, login, pass)
	require.NoError(t, err)

	token, err := auth.Login(ctx, login, pass)
	require.NoError(t, err)

	err = auth.SetPermission(ctx, login, models.PermAdmin)
	require.NoError(t, err)

	fmt.Println(token)
}
