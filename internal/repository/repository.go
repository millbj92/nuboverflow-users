package repository

import (
	"fmt"
	"log"
	"os"

	"github.com/millbj92/nuboverflow-users/internal/user"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Store interface {
	GetAllUsers() ([]user.User, error)
	GetUserByID(id int) (user.User, error)
	GetUserByEmail(email string) (user.User, error)
	CreateUser(user *user.User) (*user.User, error)
	UpdateUser(user user.User) (user.User, error)
	DeleteUser(id int) error
}

type store struct {
	DB *gorm.DB
}

func New() (Store, error) {
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbDatabase := os.Getenv("DB_DATABASE")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		dbUsername,
		dbPassword,
		dbHost,
		dbPort,
		dbDatabase,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&user.User{})

	if err != nil {
		log.Println("Failed to migrate database.")
		return nil, err
	}
	log.Println("Connection to database successful.")
	return &store{
		DB: db,
	}, nil
}

func (s *store) GetAllUsers() ([]user.User, error) {
	var users []user.User
	if result := s.DB.Find(&users); result.Error != nil {
		log.Printf("GORM ERROR: %s", result.Error.Error())
		return []user.User{}, result.Error
	}
	return users, nil
}

func (s *store) GetUserByID(id int) (user.User, error) {
	var usr user.User
	if result := s.DB.First(&usr, id); result.Error != nil {
		return user.User{}, result.Error
	}
	return usr, nil
}

func (s *store) GetUserByEmail(email string) (user.User, error) {
	var usr user.User
	if result := s.DB.Where("email = ?", email).First(&usr); result.Error != nil {
		log.Println(result.Error.Error())
		return user.User{}, result.Error
	}
	return usr, nil
}

func (s *store) CreateUser(usr *user.User) (*user.User, error) {
	if result := s.DB.Create(usr); result.Error != nil {
		return nil, result.Error
	}
	return usr, nil
}

func (s *store) UpdateUser(usr user.User) (user.User, error) {
	if result := s.DB.UpdateColumns(usr); result.Error != nil {
		return user.User{}, result.Error
	}
	return usr, nil
}

func (s *store) DeleteUser(id int) error {
	var usr user.User
	if result := s.DB.Delete(&usr, id); result.Error != nil {
		return result.Error
	}
	return nil
}
