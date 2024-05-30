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
		err := rows.Scan(&lesson.Id, &lesson.Week, &lesson.WeekDay, &lesson.Number,
			&lesson.Auditorium, &lesson.Name, &lesson.Teacher, &lesson.Group)
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
		err := rows.Scan(&lesson.Id, &lesson.Week, &lesson.WeekDay, &lesson.Number,
			&lesson.Auditorium, &lesson.Name, &lesson.Teacher, &lesson.Group)
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

func (s *Schedule) New(ctx context.Context, lesson *models.Lesson) error {
	res, err := s.Exec(ctx, "INSERT INTO lessons (week, week_day, number, auditorium, name, teacher, user_group) VALUES (?, ?, ?, ?, ?, ?, ?)",
		lesson.Week, lesson.WeekDay, lesson.Number, lesson.Auditorium, lesson.Name, lesson.Teacher, lesson.Group)
	if err != nil {
		return tr.Trace(err)
	}
	id, _ := res.LastInsertId()
	lesson.Id = int(id)
	return nil
}

func (s *Schedule) Update(ctx context.Context, lesson *models.Lesson) error {
	if _, err := s.Exec(ctx, "UPDATE lessons SET week=?, week_day=?, number=?, auditorium=?, name=?, teacher=?, user_group=? WHERE id=?",
		lesson.Week, lesson.WeekDay, lesson.Number, lesson.Auditorium, lesson.Name, lesson.Teacher, lesson.Group, lesson.Id); err != nil {
		return tr.Trace(err)
	}
	return nil
}

func (s *Schedule) Delete(ctx context.Context, id int) error {
	if _, err := s.Exec(ctx, "DELETE FROM lessons WHERE id=?", id); err != nil {
		return err
	}
	return nil
}
