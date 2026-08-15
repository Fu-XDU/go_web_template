package model

import (
	"time"
)

// User
// https://gorm.io/docs/scopes.html#Dynamically-Table
type User struct {
	ID           uint           `gorm:"primary_key" json:"id"`
	Name         string         `gorm:"type:varchar(255);not null" json:"name"`
	CreatedAt    time.Time      `json:"-"`
}

func NewUser(name string) *User {
	return &User{
		Name: name,
	}
}
