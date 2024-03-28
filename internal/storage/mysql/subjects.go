package mysql

import (
	"context"
	"database/sql"
	"eljur/internal/domain/models"
	"eljur/internal/storage/transaction"
	"eljur/pkg/tr"
	"errors"
	"fmt"
)

type Subjects struct {
	transaction.TxStorage
}

func NewSubjectsStorage(db *sql.DB) *Subjects {
	var s Subjects
	s.DB = db
	return &s
}

// return subject name by id

func (s *Subjects) GetById(ctx context.Context, id int) (*models.Subject, error) {

	var subject models.Subject
	if err := s.QueryRow(ctx, "SELECT * FROM subjects WHERE id=?", id).Scan(
		&subject.Id, &subject.Name, &subject.Semester, &subject.Course); err != nil {
		return nil, tr.Trace(err)
	}
	return &subject, nil
}

func (s *Subjects) Find(ctx context.Context, opts models.SubjectFindOpts) ([]models.Subject, error) {
	query := "SELECT * FROM subjects WHERE "
	var args []any

	if opts.Id != nil {
		query += "id=? AND "
		args = append(args, *opts.Id)
	}
	if opts.Name != nil {
		query += "name=? AND "
		args = append(args, *opts.Name)
	}
	if opts.Semester != nil {
		query += "semester=? AND "
		args = append(args, *opts.Semester)
	}
	if opts.Course != nil {
		query += "course=? AND "
		args = append(args, *opts.Course)
	}
	query = query[:len(query)-5]
	fmt.Println(query, args)
	var subjects []models.Subject
	rows, err := s.Query(ctx, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return subjects, nil
		}
		return nil, tr.Trace(err)
	}
	for rows.Next() {
		var subject models.Subject
		if err := rows.Scan(&subject.Id, &subject.Name, &subject.Semester, &subject.Course); err != nil {
			return nil, tr.Trace(err)
		}
	}

	return subjects, nil
}

// return [subject1, subject2, ...]

func (s *Subjects) GetAll(ctx context.Context) ([]models.Subject, error) {
	rows, err := s.Query(ctx, "SELECT * FROM subjects")
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

func (s *Subjects) GetBySemester(ctx context.Context, semester int8, course int8) ([]models.MinSubject, error) {
	rows, err := s.Query(ctx, "SELECT id, name FROM subjects WHERE semester=? AND course=?", semester, course)
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

func (s *Subjects) NewSubject(ctx context.Context, subject models.Subject) (int, error) {
	res, err := s.Exec(ctx, "INSERT INTO subjects (name, semester, course) VALUES (?, ?, ?)", subject.Name, subject.Semester, subject.Course)
	if err != nil {
		return 0, tr.Trace(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, tr.Trace(err)
	}
	return int(id), nil
}

func (s *Subjects) Update(ctx context.Context, subject models.MinSubject) error {
	if _, err := s.Exec(ctx, "UPDATE subjects SET name=? WHERE id=?", subject.Name, subject.Id); err != nil {
		return tr.Trace(err)
	}
	return nil
}

// delete subject

func (s *Subjects) Delete(ctx context.Context, id int) error {
	if _, err := s.Exec(ctx, "DELETE FROM subjects WHERE id=?", id); err != nil {
		return tr.Trace(err)
	}
	return nil
}
