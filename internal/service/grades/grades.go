package grades

import (
	"eljur/internal/domain/models"
	"eljur/internal/storage"
	"fmt"
	"log/slog"
)

type GradeService struct {
	gradesStorage   storage.Grades
	subjectsStorage storage.Subjects
	userStorage     storage.Users
	l               *slog.Logger
}

func New(gradesStorage storage.Grades, subjectsStorage storage.Subjects, l *slog.Logger) *GradeService {
	return &GradeService{
		gradesStorage:   gradesStorage,
		subjectsStorage: subjectsStorage,
		l:               l,
	}
}

type MinGrade struct {
	Id    int    `json:"id"`
	Value string `json:"value"`
}

type MinGradeWithDay struct {
	Value string `json:"value"`
	Day   int8   `json:"day"`
}

type UserGradesByMonth struct {
	SubjectsNames []string            `json:"subject_names"`
	Grades        [][]MinGradeWithDay `json:"grades"`
}

type SubjectGradesByMonth struct {
	Days   []int8       `json:"days"`
	Users  []string     `json:"users"`
	Grades [][]MinGrade `json:"grades"`
}

func (g *GradeService) GetUserGradesByMonth(userId int, month int8, course int8) (*UserGradesByMonth, error) {
	const op = "grades.GetGradesByMonthForUser"
	grades, err := g.gradesStorage.Find(models.GradesFindOpts{
		UserId: &userId,
		Month:  &month,
		Course: &course,
	})
	if err != nil {
		g.l.Error(fmt.Sprintf("%s: %s", op, err.Error()))
		return nil, err
	}

	subjects, err := g.subjectsStorage.GetAll()
	if err != nil {
		g.l.Error(fmt.Sprintf("%s: %s", op, err.Error()))
		return nil, err
	}

	var userGradesByMonth UserGradesByMonth

	for _, subject := range subjects {
		userGradesByMonth.SubjectsNames = append(userGradesByMonth.SubjectsNames, subject.Name)

		var newGradesSlice []MinGradeWithDay
		for _, grade := range grades {
			if subject.Id == grade.SubjectId {
				newGradesSlice = append(newGradesSlice, MinGradeWithDay{
					Value: gradesMap[grade.Value],
					Day:   grade.Day,
				})

			}

		}
		userGradesByMonth.Grades = append(userGradesByMonth.Grades, newGradesSlice)
	}

	return &userGradesByMonth, nil
}

func (g *GradeService) GetByMonthAndSubject(month int8, subjectId int, course int8) (*SubjectGradesByMonth, error) {
	const op = "grades.GetByMonth"
	grades, err := g.gradesStorage.Find(models.GradesFindOpts{
		SubjectId: &subjectId,
		Month:     &month,
		Course:    &course,
	})
	if err != nil {
		g.l.Error(fmt.Sprintf("%s: %s", op, err.Error()))
		return nil, err
	}

	users, err := g.userStorage.GetAll()
	if err != nil {
		g.l.Error(fmt.Sprintf("%s: %s", op, err.Error()))
		return nil, err
	}

	var subjectGradesByMonth SubjectGradesByMonth

	for i, user := range users {
		subjectGradesByMonth.Users = append(subjectGradesByMonth.Users, user.FullName)
		var newGradesSlice []MinGrade
		for _, grade := range grades {
			if grade.UserId == user.Id {
				if i == 0 {
					subjectGradesByMonth.Days = append(subjectGradesByMonth.Days, grade.Day)
				}
				newGradesSlice = append(newGradesSlice, MinGrade{
					Id:    grade.Id,
					Value: gradesMap[grade.Value],
				})
			}
		}
		subjectGradesByMonth.Grades = append(subjectGradesByMonth.Grades, newGradesSlice)
	}

	return &subjectGradesByMonth, nil
}

var gradesMap = map[int8]string{
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
}
