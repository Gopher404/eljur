package tests

import (
	"eljur/internal/domain/models"
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

func TestGradesStorage(t *testing.T) {
	s, err := suite.GetStorage()
	require.NoError(t, err)
	grade := models.Grade{
		UserId:    1,
		SubjectId: 2,
		Value:     4,
		Day:       5,
		Month:     1,
		Course:    3,
	}

	id, err := s.Grades.NewGrade(grade)
	require.NoError(t, err)

	gradeR, err := s.Grades.Find(models.GradesFindOpts{Id: &id})
	require.NoError(t, err)
	grade.Id = id

	assert.Equal(t, grade, *gradeR[0])

	grades, err := s.Grades.Find(models.GradesFindOpts{
		SubjectId: &grade.SubjectId,
		Month:     &grade.Month,
	})
	assert.Equal(t, grade, *grades[0])

	grades, err = s.Grades.Find(models.GradesFindOpts{
		UserId: &grade.UserId,
		Month:  &grade.Month,
	})
	assert.Equal(t, grade, *grades[0])

	grade2 := models.Grade{
		Id:        id,
		UserId:    2,
		SubjectId: 3,
		Value:     5,
		Day:       6,
		Month:     2,
		Course:    4,
	}

	err = s.Grades.Update(grade2)
	require.NoError(t, err)

	gradeR, err = s.Grades.Find(models.GradesFindOpts{Id: &grade2.Id})
	require.NoError(t, err)
	assert.Equal(t, grade2, *gradeR[0])

	err = s.Grades.Delete(id)
	require.NoError(t, err)
}
