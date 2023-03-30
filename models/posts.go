package models

import "time"

type Post struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"desc"`
	Image       string    `json:"image"`
	AuthorID    string    `json:"-"`
	User        User      `gorm:"foreignKey:AuthorID"`
	CreateAt    time.Time `json:"created_at gorm:autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at"`
}
