package mysql

import (
	"context"
	"database/sql"
	"eljur/internal/domain/models"
	"eljur/internal/storage/transaction"
	"eljur/pkg/tr"
)

type UsersTest struct {
	transaction.TxStorage
}

func NewUsersTestStorage(db *sql.DB) *UsersTest {
	var s UsersTest
	s.DB = db
	return &s
}

// create new user

func (u *UsersTest) NewUser(ctx context.Context, user models.User) (int, error) {
	res, err := u.Exec(ctx, "INSERT INTO users_test (name, login, user_group) VALUES (?, ?, ?);", user.FullName, user.Login, user.Group)
	if err != nil {
		return 0, tr.Trace(err)
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

func (u *UsersTest) GetById(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	row := u.QueryRow(ctx, "SELECT * FROM users_test WHERE id=?;", id)
	if err := row.Scan(&user.Id, &user.FullName, &user.Login, &user.Group); err != nil {
		return nil, tr.Trace(err)
	}
	return &user, nil
}

func (u *UsersTest) GetId(ctx context.Context, login string) (int, error) {
	var id int
	if err := u.QueryRow(ctx, "SELECT id FROM users_test WHERE login=?;", login).Scan(&id); err != nil {
		return 0, tr.Trace(err)
	}
	return id, nil
}

func (u *UsersTest) GetGroup(ctx context.Context, login string) (int8, error) {
	var group int8
	if err := u.QueryRow(ctx, "SELECT user_group FROM users_test WHERE login=?;", login).Scan(&group); err != nil {
		return 0, tr.Trace(err)
	}
	return group, nil
}

func (u *UsersTest) GetAll(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	rows, err := u.Query(ctx, "SELECT * FROM users_test ORDER BY name;")
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

func (u *UsersTest) Update(ctx context.Context, user models.User) error {
	if _, err := u.Exec(ctx, "UPDATE users_test SET name=?, login=?, user_group=? WHERE id=?;", user.FullName, user.Login, user.Group, user.Id); err != nil {
		return tr.Trace(err)
	}
	return nil
}

func (u *UsersTest) Delete(ctx context.Context, id int) error {
	if _, err := u.Exec(ctx, "DELETE FROM users_test WHERE id=?;", id); err != nil {
		return tr.Trace(err)
	}
	return nil
}
