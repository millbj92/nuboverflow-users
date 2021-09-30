package repository

import (
	"fmt"
	"log"
	"os"

	"github.com/millbj92/nuboverflow-users/internal/user"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func New() (Store, error) {
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbDatabase := os.Getenv("DB_DATABASE")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUsername,
		dbPassword,
		dbHost,
		dbPort,
		dbDatabase,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return Store{}, err
	}

	err = db.AutoMigrate(&user.User{}, &user.Award{})

	if err != nil {
		log.Println("Failed to migrate database.")
		return Store{}, err
	}

	return Store{
		db: db,
	}, nil
}

func (s Store) GetAllUsers() ([]user.User, error) {
	var users []user.User
	if result := s.db.Find(&users); result.Error != nil {
		return []user.User{}, result.Error
	}
	return users, nil
}

func (s Store) GetUserByID(id string) (user.User, error){
	var usr user.User
	if result := s.db.First(&usr, id); result.Error != nil {
		return user.User{}, result.Error
	}
	return usr, nil
}

func (s Store) GetUserByEmail(email string) (user.User, error) {
	var usr user.User
	if result := s.db.Where("Email = ", email).First(&usr); result.Error != nil {
		return user.User{}, result.Error
	}
	return usr, nil
}

func (s Store) CreateUser(usr user.User) (user.User, error) {
	if result := s.db.Create(usr); result.Error != nil {
		return user.User{}, result.Error
	}
	return usr, nil
}

func (s Store) UpdateUser(usr user.User) (user.User, error) {
	if result := s.db.UpdateColumns(usr); result.Error != nil {
		return user.User{}, result.Error
	}
	return usr, nil
}

func (s Store) DeleteUser(id string) error {
	var usr user.User
	if result := s.db.Delete(&usr, id); result.Error != nil {
		return result.Error
	}
	return nil
}

