//go:generate mockgen -destination=user_mocks_test.go -package=user github.com/millbj92/nuboverflow-users/internal/user/service Service
package user

import (
	"errors"
	"log"

	"github.com/millbj92/nuboverflow-users/internal/repository"
	"github.com/millbj92/nuboverflow-users/internal/user"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	GetAllUsers() ([]user.User, error)
	GetUserByID(id int) (user.User, error)
	GetUserByEmail(email string) (user.User, error)
	CreateUser(user *user.User) (*user.User, error)
	UpdateUser(user user.User) (user.User, error)
	DeleteUser(id int) error
}

type service struct {
	Store repository.Store
}

func NewService(store repository.Store) Service {
	return &service{
		Store: store,
	}
}

func (s *service) GetAllUsers() ([]user.User, error) {
	users, err := s.Store.GetAllUsers()
	if err != nil {
		log.Printf("SERVICE ERROR: %s", err.Error())
		return []user.User{}, err
	}
	return users, nil
}

func (s *service) GetUserByID(id int) (user.User, error) {
	usr, err := s.Store.GetUserByID(id)
	if err != nil {
		return user.User{}, err
	}
	return usr, nil
}

func (s *service) GetUserByEmail(email string) (user.User, error) {
	usr, err := s.Store.GetUserByEmail(email)
	if err != nil {
		return user.User{}, err
	}
	return usr, nil
}

func (s *service) CreateUser(usr *user.User) (*user.User, error) {

	existingUser, err := s.GetUserByEmail(usr.Email)
	if err != nil {
		log.Println(err)
	}
	if existingUser.ID > 0 {
		log.Println("Returning error: User exists.")
		return &user.User{}, errors.New("user exists")
	}

	passBytes := []byte(usr.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(passBytes, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	usr.Password = string(hashedPassword)

	created, err := s.Store.CreateUser(usr)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (s *service) UpdateUser(usr user.User) (user.User, error) {
	usr, err := s.Store.UpdateUser(usr)
	if err != nil {
		return user.User{}, err
	}
	return usr, nil
}

func (s *service) DeleteUser(id int) error {
	err := s.Store.DeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}
