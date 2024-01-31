package mysql

import (
	"database/sql"
	"eljur/internal/domain/models"
	"eljur/pkg/tr"
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
	var name string
	if err := s.db.QueryRow("SELECT name from subjects WHERE id=?", id).Scan(&name); err != nil {
		return "", tr.Trace(err)
	}
	return name, nil
}

// return [subject1, subject2, ...]

func (s *Subjects) GetAll() ([]*models.Subject, error) {
	rows, err := s.db.Query("SELECT * FROM subjects")
	if err != nil {
		return nil, tr.Trace(err)
	}

	var subjects []*models.Subject
	for rows.Next() {
		var subject models.Subject
		if err := rows.Scan(&subject.Id, &subject.Name); err != nil {
			return nil, tr.Trace(err)
		}
		subjects = append(subjects, &subject)
	}

	return subjects, nil
}

// create new subject

func (s *Subjects) NewSubject(name string) error {
	if _, err := s.db.Exec("INSERT INTO subjects (name) VALUES (?)", name); err != nil {
		return tr.Trace(err)
	}
	return nil
}

// delete subject

func (s *Subjects) Delete(id int) error {
	if _, err := s.db.Exec("DELETE FROM subjects WHERE id=?", id); err != nil {
		return tr.Trace(err)
	}
	return nil
}
