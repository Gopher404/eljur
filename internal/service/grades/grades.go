package grades

import (
	"eljur/internal/domain/models"
	"eljur/internal/storage"
	"eljur/pkg/tr"
	"errors"
	"fmt"
)

type GradeService struct {
	gradesStorage   storage.Grades
	subjectsStorage storage.Subjects
	userStorage     storage.Users
}

func New(gradesStorage storage.Grades, subjectsStorage storage.Subjects, userStorage storage.Users) *GradeService {
	return &GradeService{
		gradesStorage:   gradesStorage,
		subjectsStorage: subjectsStorage,
		userStorage:     userStorage,
	}
}

type MinGradeWithDay struct {
	Id    int    `json:"id"`
	Value string `json:"value"`
	Day   int8   `json:"day"`
}

type UserGradesByMonth struct {
	SubjectsNames []string            `json:"subject_names"`
	Grades        [][]MinGradeWithDay `json:"grades"`
}

type SubjectGradesByMonth struct {
	Days   []int8              `json:"days"`
	Users  []MinUser           `json:"users"`
	Grades [][]models.MinGrade `json:"grades"`
}

type MinUser struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (g *GradeService) GetUserGradesByMonth(login string, month int8, course int8) (*UserGradesByMonth, error) {
	userId, err := g.userStorage.GetId(login)
	if err != nil {
		return nil, tr.Trace(err)
	}
	grades, err := g.gradesStorage.Find(models.GradesFindOpts{
		UserId: &userId,
		Month:  &month,
		Course: &course,
	})
	if err != nil {
		return nil, tr.Trace(err)
	}

	var semester int8 = 1
	if month < 7 {
		semester = 2
	}
	subjects, err := g.subjectsStorage.GetBySemester(semester, course)
	if err != nil {
		return nil, tr.Trace(err)
	}

	var userGradesByMonth UserGradesByMonth

	for _, subject := range subjects {
		userGradesByMonth.SubjectsNames = append(userGradesByMonth.SubjectsNames, subject.Name)

		var newGradesSlice []MinGradeWithDay
		for _, grade := range grades {
			if subject.Id == grade.SubjectId {
				newGradesSlice = append(newGradesSlice, MinGradeWithDay{
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
	grades, err := g.gradesStorage.Find(models.GradesFindOpts{
		SubjectId: &subjectId,
		Month:     &month,
		Course:    &course,
	})
	if err != nil {
		return nil, tr.Trace(err)
	}

	users, err := g.userStorage.GetAll()
	if err != nil {
		return nil, tr.Trace(err)
	}

	var subjectGradesByMonth SubjectGradesByMonth

	for i, user := range users {
		subjectGradesByMonth.Users = append(subjectGradesByMonth.Users, MinUser{user.Id, user.FullName})

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

func (g *GradeService) GetResultGradesBySubject(subjectId int, course int8) {

}

func (g *GradeService) GetResultGradesByUser(userid int, course int8) {

}

func (g *GradeService) Save(grades []*models.GradeToSave) error {
	var AllErr string

	for _, grade := range grades {
		switch grade.Action {
		case models.GradeActionNew:
			if _, err := g.gradesStorage.NewGrade(&models.Grade{
				UserId:    grade.UserId,
				SubjectId: grade.SubjectId,
				Value:     grade.Value,
				Day:       grade.Day,
				Month:     grade.Month,
				Course:    grade.Course,
			}); err != nil {
				AllErr += err.Error() + "; "
			}
			break

		case models.GradeActionUpdate:
			if err := g.gradesStorage.Update(models.MinGrade{
				Id:    grade.Id,
				Value: grade.Value,
			}); err != nil {
				AllErr += err.Error() + "; "
			}
			break

		case models.GradeActionDelete:
			if err := g.gradesStorage.Delete(grade.Id); err != nil {
				AllErr += err.Error() + "; "
			}
			break
		}
	}

	if AllErr != "" {
		return tr.Trace(errors.New(AllErr))
	}
	return nil
}

func (g *GradeService) NewUserGrades(userId int) error {
	grades, err := g.gradesStorage.GetAll()
	if err != nil {
		return tr.Trace(err)
	}

	ignoreGrades := make(map[[4]int]struct{})

	for _, grade := range grades {
		day := [4]int{grade.SubjectId, int(grade.Day), int(grade.Month), int(grade.Course)}
		_, ok := ignoreGrades[day]
		if ok {
			continue
		}
		if _, err := g.gradesStorage.NewGrade(&models.Grade{
			UserId:    userId,
			SubjectId: grade.SubjectId,
			Value:     "",
			Day:       grade.Day,
			Month:     grade.Month,
			Course:    grade.Course,
		}); err != nil {
			return tr.Trace(err)
		}
		ignoreGrades[day] = struct{}{}
	}
	return nil
}

func (g *GradeService) NewResGradesBySubject(subjectId int, semester int8, course int8) error {
	var (
		startMonth int8
		endMonth   int8
	)
	if semester == 1 {
		startMonth = 9
		endMonth = 12
	} else {
		startMonth = 1
		endMonth = 6
	}
	users, err := g.userStorage.GetAll()
	if err != nil {
		return tr.Trace(err)
	}

	for month := startMonth; month <= endMonth; month++ {
		fmt.Println(month)
		for _, user := range users {
			if _, err := g.gradesStorage.NewGrade(&models.Grade{
				UserId:    user.Id,
				SubjectId: subjectId,
				Value:     "",
				Day:       100,
				Month:     month,
				Course:    course,
			}); err != nil {
				return tr.Trace(err)
			}
		}
	}
	for _, user := range users {
		fmt.Println(user)
		if _, err := g.gradesStorage.NewGrade(&models.Grade{
			UserId:    user.Id,
			SubjectId: subjectId,
			Value:     "",
			Day:       100,
			Month:     semester + 100,
			Course:    course,
		}); err != nil {
			return tr.Trace(err)
		}
	}
	return nil
}

func (g *GradeService) DeleteByUser(userId int) error {
	if err := g.gradesStorage.DeleteByUser(userId); err != nil {
		return tr.Trace(err)
	}
	return nil
}

func (g *GradeService) DeleteBySubject(subjectId int) error {
	if err := g.gradesStorage.DeleteBySubject(subjectId); err != nil {
		return tr.Trace(err)
	}
	return nil
}
