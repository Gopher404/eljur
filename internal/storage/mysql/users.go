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

func (u *Users) NewUser(name, login string) (int, error) {
	res, err := u.db.Exec("INSERT INTO users (name, login) VALUES (?, ?);", name, login)
	if err != nil {
		return 0, tr.Trace(err)
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

func (u *Users) GetById(id int) (*models.User, error) {
	var user models.User
	row := u.db.QueryRow("SELECT * FROM users WHERE id=?;", id)
	if err := row.Scan(&user.Id, &user.FullName, &user.Login); err != nil {
		return nil, tr.Trace(err)
	}
	return &user, nil
}

func (u *Users) GetId(login string) (int, error) {
	var id int
	if err := u.db.QueryRow("SELECT id FROM users WHERE login=?;", login).Scan(&id); err != nil {
		return 0, tr.Trace(err)
	}
	return id, nil
}

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

func (u *Users) Update(user models.User) error {
	if _, err := u.db.Exec("UPDATE users SET name=?, login=? WHERE id=?;", user.FullName, user.Login, user.Id); err != nil {
		return tr.Trace(err)
	}
	return nil
}

func (u *Users) Delete(id int) error {
	if _, err := u.db.Exec("DELETE FROM users WHERE id=?;", id); err != nil {
		return tr.Trace(err)
	}
	return nil
}
