package users

import (
	"context"
	"eljur/internal/domain/models"
	"eljur/internal/storage"
	"eljur/pkg/tr"
	"errors"
	"fmt"
)

type UserService struct {
	storage *storage.Storage
	auth    AuthService
	grades  GradesUser
}

type AuthService interface {
	Register(ctx context.Context, login string, password string) error
	DeleteUser(ctx context.Context, login string) error
	UpdateLogin(ctx context.Context, login string, newLogin string) error
	ChangePassword(ctx context.Context, login string, newPassword string) error
	GetPermission(ctx context.Context, login string) (int32, error)
	SetPermission(ctx context.Context, login string, permission int32) error
}

type GradesUser interface {
	NewUserGrades(ctx context.Context, userId int) error
	DeleteByUser(ctx context.Context, userId int) error
}

func New(storage *storage.Storage, auth AuthService, grades GradesUser) *UserService {
	return &UserService{
		storage: storage,
		auth:    auth,
		grades:  grades,
	}
}

type SaveUsersIn struct {
	Action   string `json:"action"`
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	Group    int8   `json:"group"`
	Password string `json:"password"`
	Perm     int32  `json:"perm"`
}

var ErrUserIsExist = errors.New("user is exist")

func (u *UserService) Save(ctx context.Context, users []SaveUsersIn) error {
	ctx, err := u.storage.Tx.Begin(ctx)
	if err != nil {
		return tr.Trace(err)
	}
	for _, user := range users {
		fmt.Printf("%+v\n", user)
		switch user.Action {
		case "new":
			if _, err := u.storage.Users.GetId(ctx, user.Login); err == nil {
				return tr.Trace(ErrUserIsExist)
			}
			if err := u.auth.Register(ctx, user.Login, user.Password); err != nil {
				_ = u.storage.Tx.Rollback(ctx)
				return tr.Trace(err)
			}
			if err := u.auth.SetPermission(ctx, user.Login, user.Perm); err != nil {
				_ = u.storage.Tx.Rollback(ctx)
				return tr.Trace(err)
			}
			id, err := u.storage.Users.NewUser(ctx, models.User{
				FullName: user.Name,
				Login:    user.Login,
				Group:    user.Group,
			})
			if err != nil {
				return tr.Trace(err)
			}
			if err := u.grades.NewUserGrades(ctx, id); err != nil {
				return tr.Trace(err)
			}
			break
		case "update":
			if user.Password != "" {
				if err := u.auth.ChangePassword(ctx, user.Login, user.Password); err != nil {
					return tr.Trace(err)
				}
				if user.Id == 0 {
					continue
				}
			}
			userWithRealLogin, err := u.storage.Users.GetById(ctx, user.Id)
			if err != nil {
				return tr.Trace(err)
			}
			if id, err := u.storage.Users.GetId(ctx, user.Login); err == nil && userWithRealLogin.Id != id {
				return tr.Trace(ErrUserIsExist)
			}
			if err := u.auth.SetPermission(ctx, userWithRealLogin.Login, user.Perm); err != nil {
				_ = u.storage.Tx.Rollback(ctx)
				return tr.Trace(err)
			}
			if err := u.storage.Users.Update(ctx, models.User{
				Id:       user.Id,
				FullName: user.Name,
				Login:    user.Login,
				Group:    user.Group,
			}); err != nil {
				return tr.Trace(err)
			}
			if userWithRealLogin.Login != user.Login {
				if err := u.auth.UpdateLogin(ctx, userWithRealLogin.Login, user.Login); err != nil {
					_ = u.storage.Tx.Rollback(ctx)
					return tr.Trace(err)
				}
			}
			break
		case "del":
			if err := u.storage.Users.Delete(ctx, user.Id); err != nil {
				return tr.Trace(err)
			}
			if err := u.grades.DeleteByUser(ctx, user.Id); err != nil {
				return tr.Trace(err)
			}
			if err := u.auth.DeleteUser(ctx, user.Login); err != nil {
				_ = u.storage.Tx.Rollback(ctx)
				return tr.Trace(err)
			}
			break
		}
	}
	if err := u.storage.Tx.Commit(ctx); err != nil {
		return tr.Trace(err)
	}
	return nil
}

type UserWithPerm struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Login string `json:"login"`
	Group int8   `json:"group"`
	Perm  int32  `json:"perm"`
}

func (u *UserService) GetAll(ctx context.Context) ([]UserWithPerm, error) {
	users, err := u.storage.Users.GetAll(ctx)
	if err != nil {
		return nil, tr.Trace(err)
	}
	var fullUsers []UserWithPerm
	for _, user := range users {
		perm, err := u.auth.GetPermission(ctx, user.Login)
		if err != nil {
			return nil, tr.Trace(err)
		}
		fullUsers = append(fullUsers, UserWithPerm{
			Id:    user.Id,
			Name:  user.FullName,
			Login: user.Login,
			Group: user.Group,
			Perm:  perm,
		})
	}
	return fullUsers, nil
}

func (u *UserService) GetUser(ctx context.Context, login string) (*UserWithPerm, error) {
	userId, err := u.storage.Users.GetId(ctx, login)
	if err != nil {
		return nil, tr.Trace(err)
	}
	user, err := u.storage.Users.GetById(ctx, userId)
	if err != nil {
		return nil, tr.Trace(err)
	}
	perm, err := u.auth.GetPermission(ctx, login)
	if err != nil {
		return nil, tr.Trace(err)
	}
	userWithPerm := &UserWithPerm{
		Id:    user.Id,
		Name:  user.FullName,
		Login: user.Login,
		Group: user.Group,
		Perm:  perm,
	}
	return userWithPerm, nil
}
