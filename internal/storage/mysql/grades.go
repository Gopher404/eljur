package mysql

import (
	"database/sql"
	"eljur/internal/domain/models"
	"eljur/pkg/tr"
	"errors"
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

func (g *Grades) NewGrade(grade *models.Grade) (int, error) {
	r, err := g.db.Exec("INSERT INTO grades (user_id, subject_id, value, day, month, course) VALUES (?, ?, ?, ?, ?, ?)",
		grade.UserId, grade.SubjectId, grade.Value, grade.Day, grade.Month, grade.Course)
	if err != nil {
		return 0, tr.Trace(err)
	}
	id, err := r.LastInsertId()
	if err != nil {
		return 0, tr.Trace(err)
	}
	return int(id), nil
}

func (g *Grades) GetAll() ([]models.Grade, error) {
	rows, err := g.db.Query("SELECT * FROM grades ORDER BY day;")
	var grades []models.Grade
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return grades, nil
		}
		return nil, tr.Trace(err)
	}
	for rows.Next() {
		var grade models.Grade
		if err := rows.Scan(&grade.Id, &grade.UserId, &grade.SubjectId,
			&grade.Value, &grade.Day, &grade.Month, &grade.Course); err != nil {
			return nil, tr.Trace(err)
		}
		grades = append(grades, grade)
	}
	return grades, nil

}

func (g *Grades) Find(opts models.GradesFindOpts) ([]*models.Grade, error) {
	query := "SELECT * FROM grades WHERE " // начало sql запроса
	var args []any                         // то что надо подставить в sql запрос

	// генерируем запрос относительно того что передано в opts также не нулевые поля opts добавляем в args
	if opts.Id != nil {
		query += "id=? AND "
		args = append(args, *opts.Id)
	}
	if opts.UserId != nil {
		query += "user_id=? AND "
		args = append(args, *opts.UserId)
	}
	if opts.SubjectId != nil {
		query += "subject_id=? AND "
		args = append(args, *opts.SubjectId)
	}
	if opts.Value != nil {
		query += "value=? AND "
		args = append(args, *opts.Value)
	}
	if opts.Day != nil {
		query += "day=? AND "
		args = append(args, *opts.Day)
	}
	if opts.Month != nil {
		query += "month=? AND "
		args = append(args, *opts.Month)
	}
	if opts.Course != nil {
		query += "course=? AND "
		args = append(args, *opts.Course)
	}
	query = query + " month<13 AND course<5 ORDER BY day ;"

	var grades []*models.Grade

	rows, err := g.db.Query(query, args...)
	if err != nil {
		// если нет строк (rows) то возвращаем пустой список []*grades БЕЗ ОШИБКИ
		if errors.Is(err, sql.ErrNoRows) {
			return grades, nil
		}
		return nil, tr.Trace(err)
	}

	// считываем данные с rows
	for rows.Next() {
		var grade models.Grade
		if err := rows.Scan(&grade.Id, &grade.UserId, &grade.SubjectId,
			&grade.Value, &grade.Day, &grade.Month, &grade.Course); err != nil {
			return nil, tr.Trace(err)
		}
		grades = append(grades, &grade)
	}

	return grades, nil
}

func (g *Grades) FindResultGrades() {

}

func (g *Grades) Update(grade models.MinGrade) error {
	if _, err := g.db.Exec("UPDATE grades SET value=? WHERE id=?", grade.Value, grade.Id); err != nil {
		return tr.Trace(err)
	}
	return nil
}

func (g *Grades) Delete(id int) error {
	if _, err := g.db.Exec("DELETE FROM grades WHERE id=?", id); err != nil {
		return tr.Trace(err)
	}
	return nil
}

func (g *Grades) DeleteByUser(userId int) error {
	if _, err := g.db.Exec("DELETE FROM grades WHERE user_id=?", userId); err != nil {
		return tr.Trace(err)
	}
	return nil
}
func (g *Grades) DeleteBySubject(subjectId int) error {
	if _, err := g.db.Exec("DELETE FROM grades WHERE subject_id=?", subjectId); err != nil {
		return tr.Trace(err)
	}
	return nil
}
