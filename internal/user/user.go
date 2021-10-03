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
	Interests  []Interest `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Awards     []Award    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Interest struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    int
	Interest  string
}

type Award struct {
	ID               int
	CreatedAt        time.Time
	UpdatedAt        time.Time
	UserID           int
	AwardName        string
	AwardDescription string
	AwardLevel       int8
	AwardScore       int32
}
