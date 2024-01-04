package mysql

import (
	"database/sql"
	"eljur/internal/domain/models"
)

type Subjects struct {
	db *sql.DB
}

func NewSubjectsStorage(db *sql.DB) *Subjects {
	return &Subjects{
		db: db,
	}
}

// return subject name by id

func (s *Subjects) GetById(id int) (string, error) {
	// todo

	return "", nil
}

// return [subject1, subject2, ...]

func (s *Subjects) GetAll() ([]*models.Subject, error) {
	// todo

	return nil, nil
}

// create new subject

func (s *Subjects) NewSubject(name string) error {
	// todo

	return nil
}

// delete subject

func (s *Subjects) Delete(id int) error {
	// todo

	return nil
}
