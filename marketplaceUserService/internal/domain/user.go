package domain

import "time"

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Email        string `grom:"uniqueIndex;not null"`
	PasswordHash string
	CreatedAT    time.Time
	rote         byte
}
