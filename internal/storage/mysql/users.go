package mysql

import (
	"database/sql"
	"eljur/internal/domain/models"
	"eljur/pkg/tr"
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
	if _, err := u.db.Exec("INSERT INTO users (name) VALUES (?);", fullName); err != nil {
		return tr.Trace(err)
	}
	return nil
}

// return username

func (u *Users) GetById(id int) (user models.User, err error) {
	// todo

	return user, nil
}

func (u *Users) GetId(login string) (int, error) {
	// todo
	var id int
	if err := u.db.QueryRow("SELECT id FROM users WHERE login=?;", login).Scan(&id); err != nil {
		return 0, tr.Trace(err)
	}
	return id, nil
}

// return [user1, user2, ...] sorted by name

func (u *Users) GetAll() ([]*models.User, error) {
	var users []*models.User
	rows, err := u.db.Query("SELECT * FROM users;")
	if err != nil {
		return nil, tr.Trace(err)
	}

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.FullName, &user.Login); err != nil {
			return nil, tr.Trace(err)
		}
		users = append(users, &user)
	}

	return users, nil
}

func (u *Users) Delete(id int) error {
	// todo

	return nil
}
