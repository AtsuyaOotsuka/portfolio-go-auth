package models

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

type UserRefreshToken struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	UserID       uint      `gorm:"not null"`
	RefreshToken string    `gorm:"type:varchar(512);uniqueIndex;not null"`
	ExpiresAt    time.Time `gorm:"type:datetime;not null"`
	IsUsed       bool      `gorm:"default:false"`
	UseIP        string    `gorm:"type:varchar(45)"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func CreateRefreshToken() string {
	bytes := make([]byte, 256)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}
