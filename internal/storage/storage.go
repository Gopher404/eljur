package storage

import (
	"context"
	"database/sql"
	"eljur/internal/config"
	"eljur/internal/domain/models"
	"eljur/internal/storage/file"
	data "eljur/internal/storage/mysql"
	"eljur/internal/storage/transaction"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
)

type Users interface {
	NewUser(ctx context.Context, user models.User) (int, error)
	GetById(ctx context.Context, id int) (user *models.User, err error)
	GetId(ctx context.Context, login string) (int, error)
	GetGroup(ctx context.Context, login string) (int8, error)
	GetAll(ctx context.Context) ([]*models.User, error)
	Update(ctx context.Context, user models.User) error
	Delete(ctx context.Context, id int) error
}

type Grades interface {
	NewGrade(ctx context.Context, grade *models.Grade) (int, error)
	GetAll(ctx context.Context) ([]models.Grade, error)
	Find(ctx context.Context, opts models.GradesFindOpts) ([]*models.Grade, error)
	Update(ctx context.Context, grade models.MinGrade) error
	Delete(ctx context.Context, id int) error
	DeleteByUser(ctx context.Context, userId int) error
	DeleteBySubject(ctx context.Context, subjectId int) error
}

type Subjects interface {
	GetById(ctx context.Context, id int) (*models.Subject, error)
	GetAll(ctx context.Context) ([]models.Subject, error)
	Find(ctx context.Context, opts models.SubjectFindOpts) ([]models.Subject, error)
	GetBySemester(ctx context.Context, semester int8, course int8) ([]models.MinSubject, error)
	NewSubject(ctx context.Context, subject models.Subject) (int, error)
	Update(ctx context.Context, subject models.MinSubject) error
	Delete(ctx context.Context, id int) error
}

type Schedule interface {
	GetAll(ctx context.Context) ([]models.Lesson, error)
	GetByWeek(ctx context.Context, week int8) ([]models.Lesson, error)
	GetByWeekAndDay(ctx context.Context, week int8, day int8) ([]models.Lesson, error)
	New(ctx context.Context, lesson *models.Lesson) error
	Update(ctx context.Context, lesson *models.Lesson) error
	Delete(ctx context.Context, id int) error
}

type Files interface {
	Get(path string) (*models.File, error)
	ListDir(path string) ([]models.FileInfo, error)
	Save(path string, data io.Reader) error
	CreateDir(path string) error
}

type Storage struct {
	Users    Users
	Grades   Grades
	Subjects Subjects
	Schedule Schedule
	Files    Files
	Tx       *transaction.TxManager
}

func New(cnf *config.DBConfig) (*Storage, error) {
	db, err := sql.Open(cnf.Driver, fmt.Sprintf("%s:%s@tcp(%s)/%s", cnf.User, cnf.Password, cnf.Host, cnf.Schema))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	txManager := transaction.NewTxManager(db)

	return &Storage{
		Users:    data.NewUsersStorage(db),
		Grades:   data.NewGradesStorage(db),
		Subjects: data.NewSubjectsStorage(db),
		Schedule: data.NewScheduleStorage(db),
		Files:    file.NewFilesStorage(cnf.FileStoragePath),
		Tx:       txManager,
	}, nil
}
