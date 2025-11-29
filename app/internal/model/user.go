package model

type User struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	UUID         string `gorm:"type:char(36);uniqueIndex;not null"`
	Username     string `gorm:"type:varchar(255);uniqueIndex;not null"`
	Email        string `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash string `gorm:"type:varchar(255);not null"`
	CreatedAt    int64  `gorm:"autoCreateTime"`
	UpdatedAt    int64  `gorm:"autoUpdateTime"`
}
