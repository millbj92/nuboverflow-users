package user

import (
	"time"
)

type User struct {
	ID uint
	CreatedAt time.Time
	UpdatedAt time.Time
    UserName string
	Password string
    Email string
    Github string
    Linkedin string
    UserScore int
    Bio string
    Profession string
    WorkPlace string
    Interests *[]Interests
    Awards *[]Award
}

type Interests struct {
	ID uint
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID uint
	Interest string
}

type Award struct {
	ID uint
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID uint
    AwardName string
    AwardDescription string
    AwardLevel int8
    AwardScore int32
}
