package tests

import (
	"eljur/tests/suite"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserStorage(t *testing.T) {
	cases := []struct {
		id       int
		fullName string
	}{
		{
			id:       1,
			fullName: "ыфсвфв свфысы",
		},
		{
			id:       2,
			fullName: "ввв",
		},
		{
			id:       3,
			fullName: "ббб",
		},
		{
			id:       4,
			fullName: "ааа",
		},
	}

	s, err := suite.GetStorage()
	require.NoError(t, err)

	for _, c := range cases {
		err = s.Users.NewUser(c.fullName)
		require.NoError(t, err, fmt.Sprintf(" case: %d", c.id))
	}

	users, err := s.Users.GetAll()
	require.NoError(t, err)

	assert.Equal(t, len(cases), len(users))

	for i := range users {
		assert.Equal(t, cases[len(cases)-i-1].fullName, users[i].FullName, fmt.Sprintf("user id: %d", users[i].Id))

		name, err := s.Users.GetById(users[i].Id)
		require.NoError(t, err, fmt.Sprintf("user id: %d", users[i].Id))
		assert.Equal(t, users[i].FullName, name, fmt.Sprintf("user id: %d", users[i].Id))

		err = s.Users.Delete(users[i].Id)
		require.NoError(t, err, fmt.Sprintf("user id: %d", users[i].Id))
	}
}

func TestSubjectStorage(t *testing.T) {
	cases := []struct {
		id   int
		name string
	}{
		{
			id:   1,
			name: "test1",
		},
		{
			id:   2,
			name: "test test",
		},
		{
			id:   3,
			name: "Математика",
		},
	}

	s, err := suite.GetStorage()
	require.NoError(t, err)

	for _, c := range cases {
		err = s.Users.NewUser(c.name)
		require.NoError(t, err, fmt.Sprintf("case: %d", c.id))
	}

	subjects, err := s.Subjects.GetAll()
	require.NoError(t, err)

	assert.Equal(t, len(cases), len(subjects))

	for i := range subjects {
		assert.Equal(t, cases[i].name, subjects[i].Name, fmt.Sprintf("case: %d", cases[i].id))

		name, err := s.Subjects.GetById(subjects[i].Id)
		require.NoError(t, err)
		assert.Equal(t, subjects[i].Name, name, fmt.Sprintf("case: %d", cases[i].id))

		err = s.Subjects.Delete(subjects[i].Id)
		require.NoError(t, err)
	}

}

func TestGrades(t *testing.T) {

}
