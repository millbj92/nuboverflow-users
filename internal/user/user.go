package user

import (
	"time"
)

type User struct {
	ID         int
	CreatedAt  time.Time 
	UpdatedAt  time.Time
	UserName   string
	Password   string `json:"-"`
	Email      string
	Github     string 
	Linkedin   string 
	UserScore  int
	Bio        string 
	Profession string
	WorkPlace  string 
}