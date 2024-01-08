package storage

import (
	"database/sql"
	"eljur/internal/domain/models"
	data "eljur/internal/storage/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type Users interface {
	NewUser(fullName string) error
	GetById(id int) (string, error)
	GetAll() ([]*models.User, error)
	Delete(id int) error
}

type Grades interface {
	NewGrade(grade *models.Grade) (int, error)
	Find(opts models.GradesFindOpts) ([]*models.Grade, error)
	Update(grade models.MinGrade) error
	Delete(id int) error
}

type Subjects interface {
	GetById(id int) (string, error)
	GetAll() ([]*models.Subject, error)
	NewSubject(name string) error
	Delete(id int) error
}

type Storage struct {
	Users    Users
	Grades   Grades
	Subjects Subjects
}

func New(driver string, dataSource string) (*Storage, error) {
	db, err := sql.Open(driver, dataSource)
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
