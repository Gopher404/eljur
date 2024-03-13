package storage

import (
	"database/sql"
	"eljur/internal/config"
	"eljur/internal/domain/models"
	data "eljur/internal/storage/mysql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Users interface {
	NewUser(name, login string) (int, error)
	GetById(id int) (user *models.User, err error)
	GetId(login string) (int, error)
	GetAll() ([]*models.User, error)
	Update(user models.User) error
	Delete(id int) error
}

type Grades interface {
	NewGrade(grade *models.Grade) (int, error)
	GetAll() ([]models.Grade, error)
	Find(opts models.GradesFindOpts) ([]*models.Grade, error)
	Update(grade models.MinGrade) error
	Delete(id int) error
	DeleteByUser(userId int) error
	DeleteBySubject(subjectId int) error
}

type Subjects interface {
	GetById(id int) (string, error)
	GetAll() ([]models.Subject, error)
	GetBySemester(semester int8, course int8) ([]models.MinSubject, error)
	NewSubject(subject models.Subject) (int, error)
	Update(subject models.MinSubject) error
	Delete(id int) error
}

type Storage struct {
	Users    Users
	Grades   Grades
	Subjects Subjects
}

func New(cnf *config.DBConfig) (*Storage, error) {
	db, err := sql.Open(cnf.Driver, fmt.Sprintf("%s:%s@tcp(%s)/%s", cnf.User, cnf.Password, cnf.Host, cnf.Schema))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Storage{
		Users:    data.NewUsersStorage(db),
		Grades:   data.NewGradesStorage(db),
		Subjects: data.NewSubjectsStorage(db),
	}, nil
}
