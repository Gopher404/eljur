package tests

import (
	"context"
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
		_, err = s.Users.NewUser(context.Background(), c.fullName, "sasa")
		require.NoError(t, err, fmt.Sprintf(" case: %d", c.id))
	}

	users, err := s.Users.GetAll(context.Background())
	require.NoError(t, err)

	assert.Equal(t, len(cases), len(users))

	for i := range users {
		assert.Equal(t, cases[len(cases)-i-1].fullName, users[i].FullName, fmt.Sprintf("user id: %d", users[i].Id))

		name, err := s.Users.GetById(context.Background(), users[i].Id)
		require.NoError(t, err, fmt.Sprintf("user id: %d", users[i].Id))
		assert.Equal(t, users[i].FullName, name, fmt.Sprintf("user id: %d", users[i].Id))

		err = s.Users.Delete(context.Background(), users[i].Id)
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
		_, err = s.Subjects.NewSubject(context.Background(), models.Subject{Name: "Test", Semester: 1, Course: 1})
		require.NoError(t, err, fmt.Sprintf("case: %d", c.id))
	}

	subjects, err := s.Subjects.GetAll(context.Background())
	require.NoError(t, err)

	assert.Equal(t, len(cases), len(subjects))

	for i := range subjects {
		assert.Equal(t, cases[i].name, subjects[i].Name, fmt.Sprintf("case: %d", cases[i].id))

		name, err := s.Subjects.GetById(context.Background(), subjects[i].Id)
		require.NoError(t, err)
		assert.Equal(t, subjects[i].Name, name, fmt.Sprintf("case: %d", cases[i].id))

		err = s.Subjects.Delete(context.Background(), subjects[i].Id)
		require.NoError(t, err)
	}

}

func TestGradesStorage(t *testing.T) {
	s, err := suite.GetStorage()
	require.NoError(t, err)
	grade := models.Grade{
		UserId:    1,
		SubjectId: 2,
		Value:     "12345",
		Day:       5,
		Month:     1,
		Course:    3,
	}

	id, err := s.Grades.NewGrade(context.Background(), &grade)
	require.NoError(t, err)

	gradeR, err := s.Grades.Find(context.Background(), models.GradesFindOpts{Id: &id})
	require.NoError(t, err)
	grade.Id = id

	assert.Equal(t, grade, *gradeR[0])

	grades, err := s.Grades.Find(context.Background(), models.GradesFindOpts{
		SubjectId: &grade.SubjectId,
		Month:     &grade.Month,
	})
	assert.Equal(t, grade, *grades[0])

	grades, err = s.Grades.Find(context.Background(), models.GradesFindOpts{
		UserId: &grade.UserId,
		Month:  &grade.Month,
	})
	assert.Equal(t, grade, *grades[0])

	grade2 := models.Grade{
		Id:        id,
		UserId:    2,
		SubjectId: 2,
		Value:     "12345",
		Day:       5,
		Month:     1,
		Course:    3,
	}

	err = s.Grades.Update(context.Background(), models.MinGrade{
		Id:    grade2.Id,
		Value: grade2.Value,
	})
	require.NoError(t, err)

	gradeR, err = s.Grades.Find(context.Background(), models.GradesFindOpts{Id: &grade2.Id})
	require.NoError(t, err)
	assert.Equal(t, grade2, *gradeR[0])

	err = s.Grades.Delete(context.Background(), id)
	require.NoError(t, err)
}
