package mysql

import (
	"context"
	"database/sql"
	"eljur/internal/domain/models"
	"eljur/internal/storage/transaction"
	"eljur/pkg/tr"
)

type Users struct {
	transaction.TxStorage
}

func NewUsersStorage(db *sql.DB) *Users {
	var s Users
	s.DB = db
	return &s
}

// create new user

func (u *Users) NewUser(ctx context.Context, user models.User) (int, error) {
	res, err := u.Exec(ctx, "INSERT INTO users (name, login, user_group) VALUES (?, ?, ?);", user.FullName, user.Login, user.Group)
	if err != nil {
		return 0, tr.Trace(err)
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

func (u *Users) GetById(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	row := u.QueryRow(ctx, "SELECT * FROM users WHERE id=?;", id)
	if err := row.Scan(&user.Id, &user.FullName, &user.Login, &user.Group); err != nil {
		return nil, tr.Trace(err)
	}
	return &user, nil
}

func (u *Users) GetId(ctx context.Context, login string) (int, error) {
	var id int
	if err := u.QueryRow(ctx, "SELECT id FROM users WHERE login=?;", login).Scan(&id); err != nil {
		return 0, tr.Trace(err)
	}
	return id, nil
}

func (u *Users) GetGroup(ctx context.Context, login string) (int8, error) {
	var group int8
	if err := u.QueryRow(ctx, "SELECT user_group FROM users WHERE login=?;", login).Scan(&group); err != nil {
		return 0, tr.Trace(err)
	}
	return group, nil
}

func (u *Users) GetAll(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	rows, err := u.Query(ctx, "SELECT * FROM users ORDER BY name;")
	if err != nil {
		return nil, tr.Trace(err)
	}

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.FullName, &user.Login, &user.Group); err != nil {
			return nil, tr.Trace(err)
		}
		users = append(users, &user)
	}

	return users, nil
}

func (u *Users) Update(ctx context.Context, user models.User) error {
	if _, err := u.Exec(ctx, "UPDATE users SET name=?, login=?, user_group=? WHERE id=?;", user.FullName, user.Login, user.Group, user.Id); err != nil {
		return tr.Trace(err)
	}
	return nil
}

func (u *Users) Delete(ctx context.Context, id int) error {
	if _, err := u.Exec(ctx, "DELETE FROM users WHERE id=?;", id); err != nil {
		return tr.Trace(err)
	}
	return nil
}
