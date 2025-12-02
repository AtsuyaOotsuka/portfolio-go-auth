package seeder

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UsersSeeder struct {
	UUID      string
	UserName  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func GetUsersSeeders(count int, isDelete bool) []UsersSeeder {
	users := make([]UsersSeeder, count)
	for i := 0; i < count; i++ {
		users[i] = UsersSeeder{
			UUID:      fmt.Sprintf("uuid_%d", i+1),
			UserName:  fmt.Sprintf("Test User %d", i+1),
			Email:     fmt.Sprintf("user%d@example.com", i+1),
			Password:  fmt.Sprintf("password%d", i+1),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}
	return users
}
