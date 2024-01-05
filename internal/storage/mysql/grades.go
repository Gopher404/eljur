package mysql

import (
	"database/sql"
	"eljur/internal/domain/models"
	"fmt"
)

type Grades struct {
	db *sql.DB
}

func NewGradesStorage(db *sql.DB) *Grades {
	return &Grades{
		db: db,
	}
}

// create new grade

func (g *Grades) NewGrade(grade models.Grade) (int, error) {
	const op = "mysql.NewGrade"

	r, err := g.db.Exec("INSERT INTO grades (user_id, subject_id, value, day, month, course) VALUES (?, ?, ?, ?, ?, ?)",
		grade.UserId, grade.SubjectId, grade.Value, grade.Day, grade.Month, grade.Course)
	if err != nil {
		return 0, err
	}
	id, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (g *Grades) GetGrade(id int) (*models.Grade, error) {
	const op = "mysql.GetGrade"

	var grade models.Grade
	if err := g.db.QueryRow("SELECT * FROM grades WHERE id=?", id).Scan(
		&grade.Id, &grade.UserId, &grade.SubjectId, &grade.Value, &grade.Day, &grade.Month, &grade.Course); err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	return &grade, nil
}

func (g *Grades) GetByMonth(month int8) ([]*models.Grade, error) {
	const op = "mysql.GetByMonth"

	var grades []*models.Grade

	rows, err := g.db.Query("SELECT * FROM grades WHERE month=?", month)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	for rows.Next() {
		var grade models.Grade
		if err := rows.Scan(&grade.Id, &grade.UserId, &grade.SubjectId,
			&grade.Value, &grade.Day, &grade.Month, &grade.Course); err != nil {
			return nil, fmt.Errorf("%s, %w", op, err)
		}
		grades = append(grades, &grade)
	}
	return grades, nil
}

func (g *Grades) GetByUser(userId int) ([]*models.Grade, error) {
	const op = "mysql.GetByUser"

	var grades []*models.Grade

	rows, err := g.db.Query("SELECT * FROM grades WHERE user_id=?", userId)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	for rows.Next() {
		var grade models.Grade
		if err := rows.Scan(&grade.Id, &grade.UserId, &grade.SubjectId,
			&grade.Value, &grade.Day, &grade.Month, &grade.Course); err != nil {
			return nil, fmt.Errorf("%s, %w", op, err)
		}
		grades = append(grades, &grade)
	}
	return grades, nil
}

func (g *Grades) Update(grade models.Grade) error {
	const op = "mysql.Update"

	if _, err := g.db.Exec("UPDATE grades SET user_id=?, subject_id=?, value=?, day=?, month=?, course=? WHERE id=?",
		grade.UserId, grade.SubjectId, grade.Value, grade.Day, grade.Month, grade.Course, grade.Id); err != nil {
		return fmt.Errorf("%s, %w", op, err)
	}
	return nil
}

func (g *Grades) Delete(id int) error {
	const op = "mysql.Delete"

	if _, err := g.db.Exec("DELETE FROM grades WHERE id=?", id); err != nil {
		return fmt.Errorf("%s, %w", op, err)
	}
	return nil
}
