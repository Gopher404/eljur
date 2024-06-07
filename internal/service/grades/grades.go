package grades

import (
	"context"
	"eljur/internal/domain/models"
	"eljur/internal/storage"
	"eljur/pkg/tr"
	"errors"
)

type GradeService struct {
	storage *storage.Storage
}

func New(storage *storage.Storage) *GradeService {
	return &GradeService{
		storage: storage,
	}
}

type MinGradeWithDay struct {
	Id    int    `json:"id"`
	Value string `json:"value"`
	Day   int8   `json:"day"`
}

type MinGradeWithDayAndMonth struct {
	Id    int    `json:"id"`
	Value string `json:"value"`
	Day   int8   `json:"day"`
	Month int8   `json:"month"`
}

type UserGradesByMonth struct {
	SubjectsNames []string            `json:"subject_names"`
	Grades        [][]MinGradeWithDay `json:"grades"`
}

type UserGradesBySemester struct {
	SubjectsNames []string                    `json:"subject_names"`
	Grades        [][]MinGradeWithDayAndMonth `json:"grades"`
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

var (
	monthResDay    int8 = 100 // month result
	semesterResDay int8 = 101 // semester result
	courseResMonth int8 = 100 // course result
)

const resultSubjectSemester int8 = 3

func (g *GradeService) GetUserGradesByMonth(ctx context.Context, login string, month int8, course int8) (*UserGradesByMonth, error) {
	userId, err := g.storage.Users.GetId(ctx, login)
	if err != nil {
		return nil, tr.Trace(err)
	}

	var semester int8 = 1
	if month < 7 {
		semester = 2
	}
	subjects, err := g.storage.Subjects.GetBySemester(ctx, semester, course)
	if err != nil {
		return nil, tr.Trace(err)
	}

	var userGradesByMonth UserGradesByMonth

	for _, subject := range subjects {
		userGradesByMonth.SubjectsNames = append(userGradesByMonth.SubjectsNames, subject.Name)

		grades, err := g.storage.Grades.Find(ctx, models.GradesFindOpts{
			UserId:    &userId,
			SubjectId: &subject.Id,
			Month:     &month,
			Course:    &course,
		})
		if err != nil {
			return nil, tr.Trace(err)
		}
		newGradesSlice := make([]MinGradeWithDay, 0)
		for _, grade := range grades {
			newGradesSlice = append(newGradesSlice, MinGradeWithDay{
				Id:    grade.Id,
				Value: grade.Value,
				Day:   grade.Day,
			})
		}
		userGradesByMonth.Grades = append(userGradesByMonth.Grades, newGradesSlice)
	}

	return &userGradesByMonth, nil
}

func (g *GradeService) GetUserGradesBySemester(ctx context.Context, login string, semester int8, course int8) (*UserGradesBySemester, error) {
	userId, err := g.storage.Users.GetId(ctx, login)
	if err != nil {
		return nil, tr.Trace(err)
	}

	subjects, err := g.storage.Subjects.GetBySemester(ctx, semester, course)
	if err != nil {
		return nil, tr.Trace(err)
	}

	var (
		startMonth, endMonth int8
	)
	if semester == 1 {
		startMonth = 9
		endMonth = 12
	} else {
		startMonth = 1
		endMonth = 6
	}

	var userGradesByMonth UserGradesBySemester

	for _, subject := range subjects {
		userGradesByMonth.SubjectsNames = append(userGradesByMonth.SubjectsNames, subject.Name)
		var grades []*models.Grade

		for month := startMonth; month <= endMonth; month++ {
			monthGrades, err := g.storage.Grades.Find(ctx, models.GradesFindOpts{
				UserId:    &userId,
				SubjectId: &subject.Id,
				Month:     &month,
				Course:    &course,
			})
			if err != nil {
				return nil, tr.Trace(err)
			}

			grades = append(grades, monthGrades...)
		}

		newGradesSlice := make([]MinGradeWithDayAndMonth, 0)
		for _, grade := range grades {
			newGradesSlice = append(newGradesSlice, MinGradeWithDayAndMonth{
				Id:    grade.Id,
				Value: grade.Value,
				Day:   grade.Day,
				Month: grade.Month,
			})
		}
		userGradesByMonth.Grades = append(userGradesByMonth.Grades, newGradesSlice)
	}

	return &userGradesByMonth, nil
}

func (g *GradeService) GetByMonthAndSubject(ctx context.Context, month int8, subjectId int, course int8) (*SubjectGradesByMonth, error) {
	grades, err := g.storage.Grades.Find(ctx, models.GradesFindOpts{
		SubjectId: &subjectId,
		Month:     &month,
		Course:    &course,
	})
	if err != nil {
		return nil, tr.Trace(err)
	}

	users, err := g.storage.Users.GetAll(ctx)
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

type ResUserGradesLine struct {
	MonthRes    [12]models.MinGrade `json:"month_res"`
	SemesterRes [2]models.MinGrade  `json:"semester_res"`
	CourseRes   models.MinGrade     `json:"course_res"`
	UserName    string              `json:"user_name"`
}

type ResultGradesBySubject struct {
	Users       []ResUserGradesLine `json:"users"`
	SubjectName string              `json:"subject_name"`
}

func (g *GradeService) GetResultGradesBySubject(ctx context.Context, subjectId int, course int8) (*ResultGradesBySubject, error) {
	var res ResultGradesBySubject

	userList, err := g.storage.Users.GetAll(ctx)
	if err != nil {
		return nil, tr.Trace(err)
	}
	subject, err := g.storage.Subjects.GetById(ctx, subjectId)
	if err != nil {
		return nil, tr.Trace(err)
	}
	res.SubjectName = subject.Name

	if err != nil {
		return nil, tr.Trace(err)
	}

	for _, user := range userList {
		var userGrades ResUserGradesLine

		gradesMonth, err := g.storage.Grades.Find(ctx, models.GradesFindOpts{
			SubjectId: &subjectId,
			Day:       &monthResDay,
			Course:    &course,
			UserId:    &user.Id,
		})
		if err != nil {
			return nil, tr.Trace(err)
		}

		for _, grade := range gradesMonth {
			userGrades.MonthRes[grade.Month-1] = models.MinGrade{
				Id:    grade.Id,
				Value: grade.Value,
			}
		}

		gradesSemester, err := g.storage.Grades.Find(ctx, models.GradesFindOpts{
			SubjectId: &subjectId,
			Day:       &semesterResDay,
			Course:    &course,
			UserId:    &user.Id,
		})
		if err != nil {
			return nil, tr.Trace(err)
		}
		for _, grade := range gradesSemester {
			userGrades.SemesterRes[grade.Month-1] = models.MinGrade{Id: grade.Id, Value: grade.Value}
		}

		gradesCourse, err := g.storage.Grades.Find(ctx, models.GradesFindOpts{
			SubjectId: &subjectId,
			Month:     &courseResMonth,
			Course:    &course,
			UserId:    &user.Id,
		})
		if err != nil {
			return nil, tr.Trace(err)
		}
		userGrades.CourseRes = models.MinGrade{
			Id:    gradesCourse[0].Id,
			Value: gradesCourse[0].Value,
		}

		userGrades.UserName = user.FullName

		res.Users = append(res.Users, userGrades)
	}

	return &res, nil
}

type ResUserGradesLineBySubject struct {
	MonthRes    [12]models.MinGrade `json:"month_res"`
	SemesterRes [2]models.MinGrade  `json:"semester_res"`
	CourseRes   models.MinGrade     `json:"course_res"`
	SubjectName string              `json:"subject_name"`
}

func (g *GradeService) GetResultGradesByUser(ctx context.Context, login string, course int8) ([]*ResUserGradesLineBySubject, error) {
	userId, err := g.storage.Users.GetId(ctx, login)
	if err != nil {
		return nil, tr.Trace(err)
	}
	subjects, err := g.storage.Subjects.GetBySemester(ctx, resultSubjectSemester, course)
	if err != nil {
		return nil, tr.Trace(err)
	}

	var res []*ResUserGradesLineBySubject

	for _, subject := range subjects {
		var subjectGrades ResUserGradesLineBySubject

		subjectGrades.SubjectName = subject.Name

		gradesByMonth, err := g.storage.Grades.Find(ctx, models.GradesFindOpts{
			UserId:    &userId,
			SubjectId: &subject.Id,
			Day:       &monthResDay,
			Course:    &course,
		})
		if err != nil {
			return nil, tr.Trace(err)
		}

		for _, grade := range gradesByMonth {
			subjectGrades.MonthRes[grade.Month-1] = models.MinGrade{
				Id:    grade.Id,
				Value: grade.Value,
			}
		}

		gradesSemester, err := g.storage.Grades.Find(ctx, models.GradesFindOpts{
			SubjectId: &subject.Id,
			Day:       &semesterResDay,
			Course:    &course,
			UserId:    &userId,
		})
		if err != nil {
			return nil, tr.Trace(err)
		}
		for _, grade := range gradesSemester {
			subjectGrades.SemesterRes[grade.Month-1] = models.MinGrade{Id: grade.Id, Value: grade.Value}
		}

		gradesCourse, err := g.storage.Grades.Find(ctx, models.GradesFindOpts{
			SubjectId: &subject.Id,
			Month:     &courseResMonth,
			Course:    &course,
			UserId:    &userId,
		})
		if err != nil {
			return nil, tr.Trace(err)
		}
		if len(gradesCourse) < 1 {
			return nil, tr.Trace(errors.New("course result grade not found"))
		}
		subjectGrades.CourseRes = models.MinGrade{
			Id:    gradesCourse[0].Id,
			Value: gradesCourse[0].Value,
		}

		res = append(res, &subjectGrades)

	}
	return res, nil
}

const (
	GradeActionUpdate = int8(iota)
	GradeActionDelete
	GradeActionNew
)

func (g *GradeService) Save(ctx context.Context, grades []*models.GradeToSave) error {
	ctx, err := g.storage.Tx.Begin(ctx)
	if err != nil {
		return tr.Trace(err)
	}
	for _, grade := range grades {
		switch grade.Action {
		case GradeActionNew:
			if _, err := g.storage.Grades.NewGrade(ctx, &models.Grade{
				UserId:    grade.UserId,
				SubjectId: grade.SubjectId,
				Value:     grade.Value,
				Day:       grade.Day,
				Month:     grade.Month,
				Course:    grade.Course,
			}); err != nil {
				return tr.Trace(err)
			}
			break

		case GradeActionUpdate:
			if err := g.storage.Grades.Update(ctx, models.MinGrade{
				Id:    grade.Id,
				Value: grade.Value,
			}); err != nil {
				return tr.Trace(err)
			}
			break

		case GradeActionDelete:
			if err := g.storage.Grades.Delete(ctx, grade.Id); err != nil {
				return tr.Trace(err)
			}
			break
		}
	}

	if err := g.storage.Tx.Commit(ctx); err != nil {
		return tr.Trace(err)
	}
	return nil
}

func (g *GradeService) NewUserGrades(ctx context.Context, userId int) error {
	grades, err := g.storage.Grades.GetAll(ctx)
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
		if _, err := g.storage.Grades.NewGrade(ctx, &models.Grade{
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

func (g *GradeService) NewResGradesBySubject(ctx context.Context, subjectId int, course int8) error {
	users, err := g.storage.Users.GetAll(ctx)
	if err != nil {
		return tr.Trace(err)
	}

	var month int8
	for month = 1; month <= 12; month++ {
		for _, user := range users {
			if _, err := g.storage.Grades.NewGrade(ctx, &models.Grade{
				UserId:    user.Id,
				SubjectId: subjectId,
				Value:     "",
				Day:       monthResDay,
				Month:     month,
				Course:    course,
			}); err != nil {
				return tr.Trace(err)
			}
		}
	}

	for _, user := range users {
		// res grades by semesters
		var semester int8
		for semester = 1; semester <= 2; semester++ {
			if _, err := g.storage.Grades.NewGrade(ctx, &models.Grade{
				UserId:    user.Id,
				SubjectId: subjectId,
				Value:     "",
				Day:       semesterResDay,
				Month:     semester,
				Course:    course,
			}); err != nil {
				return tr.Trace(err)
			}
		}

		// res grade by course
		if _, err := g.storage.Grades.NewGrade(ctx, &models.Grade{
			UserId:    user.Id,
			SubjectId: subjectId,
			Value:     "",
			Day:       0,
			Month:     courseResMonth,
			Course:    course,
		}); err != nil {
			return tr.Trace(err)
		}
	}
	return nil
}

func (g *GradeService) DeleteByUser(ctx context.Context, userId int) error {
	if err := g.storage.Grades.DeleteByUser(ctx, userId); err != nil {
		return tr.Trace(err)
	}
	return nil
}

func (g *GradeService) DeleteBySubject(ctx context.Context, subjectId int) error {
	if err := g.storage.Grades.DeleteBySubject(ctx, subjectId); err != nil {
		return tr.Trace(err)
	}
	return nil
}
