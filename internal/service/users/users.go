package users

import "eljur/internal/storage"

type UserService struct {
	userStorage storage.Users
}

func New(userStorage storage.Users) *UserService {
	return &UserService{
		userStorage: userStorage,
	}
}
