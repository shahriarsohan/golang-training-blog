package models

import (
	"time"
)

//A recommended approach is to define the Golang structs
//with singular names since GORM will pluralize them in the database by default.

type User struct {
	ID               uint
	Name             string    `gorm:"type:varchar(255);not null"`
	Email            string    `json:"email" gorm:"uniqueIndex;not null"`
	Password         string    `gorm:"not null"`
	Role             string    `gorm:"type:varchar(255);not null"`
	Provider         string    `gorm:"not null"`
	Photo            string    `json:"photo"`
	VerificationCode string    `json:"verification_code"`
	Verified         bool      `gorm:"not null"`
	CreateAt         time.Time `json:"created_at gorm:autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at"`
}
