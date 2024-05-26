package mysql

import (
	"context"
	"database/sql"
	"eljur/internal/domain/models"
	"eljur/internal/storage/transaction"
	"eljur/pkg/tr"
)

type Schedule struct {
	transaction.TxStorage
}

func NewScheduleStorage(db *sql.DB) *Schedule {
	var s Schedule
	s.DB = db
	return &s
}

func (s *Schedule) GetAll(ctx context.Context) ([]models.Lesson, error) {
	var lessons []models.Lesson

	rows, err := s.Query(ctx, "SELECT * FROM lessons;")
	if err != nil {
		return nil, tr.Trace(err)
	}

	for rows.Next() {
		var lesson models.Lesson
		err := rows.Scan(&lesson.Id, &lesson.Week, &lesson.Number, &lesson.WeekDay,
			&lesson.Auditorium, &lesson.Name, &lesson.Teacher)
		if err == sql.ErrNoRows {
			return lessons, nil
		}
		if err != nil {
			return nil, tr.Trace(err)
		}
		lessons = append(lessons, lesson)
	}

	return lessons, nil
}

func (s *Schedule) GetByWeek(ctx context.Context, week int8) ([]models.Lesson, error) {
	var lessons []models.Lesson

	rows, err := s.Query(ctx, "SELECT * FROM lessons WHERE week=?;", week)
	if err != nil {
		return nil, tr.Trace(err)
	}

	for rows.Next() {
		var lesson models.Lesson
		err := rows.Scan(&lesson.Id, &lesson.Week, &lesson.Number, &lesson.WeekDay,
			&lesson.Auditorium, &lesson.Name, &lesson.Teacher)
		if err == sql.ErrNoRows {
			return lessons, nil
		}
		if err != nil {
			return nil, tr.Trace(err)
		}
		lessons = append(lessons, lesson)
	}

	return lessons, nil
}

func (s *Schedule) Update(ctx context.Context, lesson *models.LessonForUpdate) error {
	if _, err := s.Exec(ctx, "UPDATE lessons SET auditorium=?, name=?, teacher=? WHERE id=?",
		&lesson.Auditorium, lesson.Name, lesson.Teacher, lesson.Id); err != nil {
		return tr.Trace(err)
	}
	return nil
}

func (s *Schedule) Delete(ctx context.Context, id int) error {
	if _, err := s.Exec(ctx, "DELETE FROM lesson WHERE id=?", id); err != nil {
		return err
	}
	return nil
}
