package user

import (
	"time"
)

type User struct {
	ID         int `json:"id" example:"1" format:"int64"`
	CreatedAt  time.Time `json:"createdAt" example:"2021-10-03T20:54:53.144Z" format:"date"`
	UpdatedAt  time.Time `json:"updatedAt" example:"2021-10-03T20:54:53.144Z" format:"date"`
	UserName   string `json:"userName" example:"nuboverflow_user" format:"string"`
	Password   string `json:"-"`
	Email      string `json:"email" example:"test@testemail.com" format:"email"`
	Github     string `json:"github" example:"http://github.com/millbj92" format:"string"`
	Linkedin   string `json:"linkedIn" example:"http://linkedin.com/userName" format:"string"`
	UserScore  int `json:"userScore" example:"1800" format:"int64"`
	Bio        string `json:"bio" example:"A simple bio." format:"string"`
	Profession string `json:"profession" example:"Software Developer" format:"string"`
	WorkPlace  string `json:"workPlace" example:"NASA" format:"string"`
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
