package mysql

import (
	"database/sql"
	"eljur/internal/domain/models"
)

type Users struct {
	db *sql.DB
}

func NewUsersStorage(db *sql.DB) *Users {
	return &Users{
		db: db,
	}
}

// create new user

func (u *Users) NewUser(fullName string) error {
	if _, err := u.db.Exec("INSERT INTO users (name) VALUES (?)", fullName); err != nil {
		return err
	}
	return nil
}

// return username

func (u *Users) GetById(id int) (string, error) {
	// todo

	return "", nil
}

// return [user1, user2, ...] sorted by name

func (u *Users) GetAll() ([]*models.User, error) {
	// todo

	return nil, nil
}

func (u *Users) Delete(id int) error {
	// todo

	return nil
}
