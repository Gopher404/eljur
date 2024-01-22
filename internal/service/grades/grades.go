package grades

import (
	"eljur/internal/domain/models"
	"eljur/internal/storage"
	"fmt"
)

type GradeService struct {
	gradesStorage   storage.Grades
	subjectsStorage storage.Subjects
	userStorage     storage.Users
}

func New(gradesStorage storage.Grades, subjectsStorage storage.Subjects) *GradeService {
	return &GradeService{
		gradesStorage:   gradesStorage,
		subjectsStorage: subjectsStorage,
	}
}

type MinGradeStringWithDay struct {
	Id    int  `json:"id"`
	Value int8 `json:"value"`
	Day   int8 `json:"day"`
}

type UserGradesByMonth struct {
	SubjectsNames []string                  `json:"subject_names"`
	Grades        [][]MinGradeStringWithDay `json:"grades"`
}

type SubjectGradesByMonth struct {
	Days   []int8              `json:"days"`
	Users  []string            `json:"users"`
	Grades [][]models.MinGrade `json:"grades"`
}

func (g *GradeService) GetUserGradesByMonth(userId int, month int8, course int8) (*UserGradesByMonth, error) {
	const op = "GradeService.GetGradesByMonthForUser"
	grades, err := g.gradesStorage.Find(models.GradesFindOpts{
		UserId: &userId,
		Month:  &month,
		Course: &course,
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	subjects, err := g.subjectsStorage.GetAll()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var userGradesByMonth UserGradesByMonth

	for _, subject := range subjects {
		userGradesByMonth.SubjectsNames = append(userGradesByMonth.SubjectsNames, subject.Name)

		var newGradesSlice []MinGradeStringWithDay
		for _, grade := range grades {
			if subject.Id == grade.SubjectId {
				newGradesSlice = append(newGradesSlice, MinGradeStringWithDay{
					Id:    grade.Id,
					Value: grade.Value,
					Day:   grade.Day,
				})

			}

		}
		userGradesByMonth.Grades = append(userGradesByMonth.Grades, newGradesSlice)
	}

	return &userGradesByMonth, nil
}

func (g *GradeService) GetByMonthAndSubject(month int8, subjectId int, course int8) (*SubjectGradesByMonth, error) {
	const op = "GradeService.GetByMonth"
	grades, err := g.gradesStorage.Find(models.GradesFindOpts{
		SubjectId: &subjectId,
		Month:     &month,
		Course:    &course,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	users, err := g.userStorage.GetAll()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var subjectGradesByMonth SubjectGradesByMonth

	for i, user := range users {
		subjectGradesByMonth.Users = append(subjectGradesByMonth.Users, user.FullName)
		var newGradesSlice []models.MinGrade
		for _, grade := range grades {
			if grade.UserId == user.Id {
				if i == 0 {
					subjectGradesByMonth.Days = append(subjectGradesByMonth.Days, grade.Day)
				}
				newGradesSlice = append(newGradesSlice, models.MinGrade{
					Id:    grade.Id,
					Value: grade.Value,
				})
			}
		}
		subjectGradesByMonth.Grades = append(subjectGradesByMonth.Grades, newGradesSlice)
	}

	return &subjectGradesByMonth, nil
}

func (g *GradeService) SaveGrades(grades []*models.Grade) error {
	const op = "GradeService.UpdateGrades"
	for _, grade := range grades {
		if _, err := g.gradesStorage.NewGrade(grade); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	return nil
}

func (g *GradeService) UpdateGrades(grades []models.MinGrade) error {
	const op = "GradeService.UpdateGrades"
	for _, grade := range grades {
		if err := g.gradesStorage.Update(grade); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	return nil
}

func (g *GradeService) DeleteGrades(gradesId []int) error {
	const op = "GradeService.DeleteGrades"
	for _, id := range gradesId {
		if err := g.gradesStorage.Delete(id); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	return nil
}

/*
gradesMap:
	1:  "1",
	2:  "2",
	3:  "3",
	4:  "4",
	5:  "5",
	0:  "",
	-1: "Н",
	-2: "У",
	-3: "Зач",
	-4: "НеЗач",
*/
