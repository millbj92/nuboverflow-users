package service

import (
	"github.com/millbj92/nuboverflow-users/internal/user"
)

type Store interface {
	GetAllUsers() ([]user.User, error)
	GetUserByID(id string) (user.User, error)
	GetUserByEmail(email string) (user.User, error)
	CreateUser(user user.User) (user.User, error)
	UpdateUser(user user.User) (user.User, error)
	DeleteUser(id string) error
}

type UserService struct {
	Store Store
}

func New(store Store) UserService {
	return UserService{
		Store: store,
	}
}

func (s UserService) GetAllUsers() ([]user.User, error){
	users, err := s.Store.GetAllUsers()
	if err != nil {
		return []user.User{}, err
	}
	return users, nil
}

func (s UserService) GetUserByID(id string) (user.User, error) {
	usr, err := s.Store.GetUserByID(id)
	if err != nil {
		return user.User{}, err
	}
	return usr, nil
}

func (s UserService) GetUserByEmail(email string) (user.User, error) {
	usr, err := s.Store.GetUserByEmail(email)
	if err != nil {
		return user.User{}, err
	}
	return usr, nil
}

func (s UserService) CreateUser(usr user.User) (user.User, error) {
	usr, err := s.Store.CreateUser(usr)
	if err != nil {
		return user.User{}, err
	}
	return usr, nil
}

func (s UserService) UpdateUser(usr user.User) (user.User, error) {
	usr, err := s.Store.UpdateUser(usr)
	if err != nil {
		return user.User{}, err
	}
	return usr, nil
}

func (s UserService) DeleteUser(id string) error {
	err := s.Store.DeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}



