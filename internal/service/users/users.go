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
	userStorage storage.Users
	auth        AuthService
	grades      GradesUser
}

type AuthService interface {
	Register(ctx context.Context, login string, password string) error
	DeleteUser(ctx context.Context, login string) error
	UpdateLogin(ctx context.Context, login string, newLogin string) error
	GetPermission(ctx context.Context, login string) (int32, error)
	SetPermission(ctx context.Context, login string, permission int32) error
}

type GradesUser interface {
	NewUserGrades(userId int) error
	DeleteByUser(userId int) error
}

func New(userStorage storage.Users, auth AuthService, grades GradesUser) *UserService {
	return &UserService{
		userStorage: userStorage,
		auth:        auth,
		grades:      grades,
	}
}

type SaveUsersIn struct {
	Action   string `json:"action"`
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Perm     int32  `json:"perm"`
}

var ErrUserIsExist = errors.New("user is exist")

func (u *UserService) Save(ctx context.Context, users []SaveUsersIn) error {

	for _, user := range users {
		fmt.Printf("%+v\n", user)
		switch user.Action {
		case "new":
			if _, err := u.userStorage.GetId(user.Login); err == nil {
				return tr.Trace(ErrUserIsExist)
			}
			if err := u.auth.Register(ctx, user.Login, user.Password); err != nil {
				return tr.Trace(err)
			}
			if err := u.auth.SetPermission(ctx, user.Login, user.Perm); err != nil {
				return tr.Trace(err)
			}
			id, err := u.userStorage.NewUser(user.Name, user.Login)
			if err != nil {
				return tr.Trace(err)
			}
			if err := u.grades.NewUserGrades(id); err != nil {
				return tr.Trace(err)
			}
			break
		case "update":
			userWithRealLogin, err := u.userStorage.GetById(user.Id)
			if err != nil {
				return tr.Trace(err)
			}
			if id, err := u.userStorage.GetId(user.Login); err == nil && userWithRealLogin.Id != id {
				return tr.Trace(ErrUserIsExist)
			}
			if err := u.auth.SetPermission(ctx, userWithRealLogin.Login, user.Perm); err != nil {
				return tr.Trace(err)
			}
			if err := u.userStorage.Update(models.User{
				Id:       user.Id,
				FullName: user.Name,
				Login:    user.Login,
			}); err != nil {
				return tr.Trace(err)
			}
			if userWithRealLogin.Login != user.Login {
				if err := u.auth.UpdateLogin(ctx, userWithRealLogin.Login, user.Login); err != nil {
					return tr.Trace(err)
				}
			}
			break
		case "del":
			if err := u.userStorage.Delete(user.Id); err != nil {
				return tr.Trace(err)
			}
			if err := u.grades.DeleteByUser(user.Id); err != nil {
				return tr.Trace(err)
			}
			if err := u.auth.DeleteUser(ctx, user.Login); err != nil {
				return tr.Trace(err)
			}
			break
		}
	}
	return nil
}

type UserWithPerm struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Login string `json:"login"`
	Perm  int32  `json:"perm"`
}

func (u *UserService) GetAll(ctx context.Context) ([]UserWithPerm, error) {
	users, err := u.userStorage.GetAll()
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
			Perm:  perm,
		})
	}
	return fullUsers, nil
}
