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

func (s *Subjects) GetAll() ([]models.Subject, error) {
	rows, err := s.db.Query("SELECT * FROM subjects")
	if err != nil {
		return nil, tr.Trace(err)
	}

	var subjects []models.Subject
	for rows.Next() {
		var subject models.Subject
		if err := rows.Scan(&subject.Id, &subject.Name, &subject.Semester, &subject.Course); err != nil {
			return nil, tr.Trace(err)
		}
		subjects = append(subjects, subject)
	}

	return subjects, nil
}

func (s *Subjects) GetBySemester(semester int8, course int8) ([]models.MinSubject, error) {
	rows, err := s.db.Query("SELECT id, name FROM subjects WHERE semester=? AND course=?", semester, course)
	if err != nil {
		return nil, tr.Trace(err)
	}

	var subjects []models.MinSubject
	for rows.Next() {
		var subject models.MinSubject
		if err := rows.Scan(&subject.Id, &subject.Name); err != nil {
			return nil, tr.Trace(err)
		}
		subjects = append(subjects, subject)
	}

	return subjects, nil
}

// create new subject

func (s *Subjects) NewSubject(subject models.Subject) (int, error) {
	res, err := s.db.Exec("INSERT INTO subjects (name, semester, course) VALUES (?, ?, ?)", subject.Name, subject.Semester, subject.Course)
	if err != nil {
		return 0, tr.Trace(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, tr.Trace(err)
	}
	return int(id), nil
}

func (s *Subjects) Update(subject models.MinSubject) error {
	if _, err := s.db.Exec("UPDATE subjects SET name=? WHERE id=?", subject.Name, subject.Id); err != nil {
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
