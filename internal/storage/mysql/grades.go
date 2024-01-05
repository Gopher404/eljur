package mysql

import (
	"database/sql"
	"eljur/internal/domain/models"
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

func (g *Grades) NewGrade(grade models.Grade) error {
	// todo

	return nil
}

func (g *Grades) GetGrade(id int) (*models.Grade, error) {
	// todo

	return nil, nil
}

func (g *Grades) GetByMonth(month int8) ([]*models.Grade, error) {
	// todo

	return nil, nil
}

func (g *Grades) GetByUser(userId int) ([]*models.Grade, error) {
	// todo

	return nil, nil
}

func (g *Grades) Update(grade models.Grade) error {
	// todo

	return nil
}

func (g *Grades) Delete(id int) error {
	// todo

	return nil
}
