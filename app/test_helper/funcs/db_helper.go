package funcs

import (
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/AtsuyaOotsuka/portfolio-go-auth/test_helper/seeder"
)

func truncateTable(db *sql.DB, tableName string) error {
	query := "TRUNCATE TABLE " + tableName
	_, err := db.Exec(query)
	return err
}

func DbCleanup(db *sql.DB) ([]DbRecords, error) {
	truncateTable(db, "users")
	truncateTable(db, "user_refresh_tokens")
	dbRecords, err := CreateSeeders(db)
	if err != nil {
		return nil, err
	}
	return dbRecords, nil
}

type DbRecords struct {
	TableName string
	Count     int
	Data      []map[string]interface{}
}

func FilterRecordsByTableName(dbRecords []DbRecords, tableName string) []DbRecords {
	filtered := []DbRecords{}
	for _, record := range dbRecords {
		if record.TableName == tableName {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

func CreateSeeders(db *sql.DB) ([]DbRecords, error) {
	dbRecords := []DbRecords{}

	users := seeder.GetUsersSeeders(5, false)
	for _, user := range users {
		password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		InsertUser, err := db.Exec("INSERT INTO users (uuid, username, email, password_hash, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
			user.UUID, user.UserName, user.Email, password, user.CreatedAt, user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		InsertUserId, err := InsertUser.LastInsertId()
		if err != nil {
			return nil, err
		}
		dbRecords = append(dbRecords, DbRecords{
			TableName: "users",
			Count:     len(users),
			Data: []map[string]interface{}{
				{
					"id":         InsertUserId,
					"uuid":       user.UUID,
					"username":   user.UserName,
					"email":      user.Email,
					"password":   user.Password, // 平文のまま保存
					"created_at": user.CreatedAt,
					"updated_at": user.UpdatedAt,
				},
			},
		})

		refreshTokenString := fmt.Sprintf("refresh_token_sample%d", InsertUserId)
		InsertUserRefreshToken, err := db.Exec("INSERT INTO user_refresh_tokens (user_id, refresh_token, expires_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
			InsertUserId, refreshTokenString, user.CreatedAt.Add(24*7*time.Hour), user.CreatedAt, user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		_, err = InsertUserRefreshToken.LastInsertId()
		if err != nil {
			return nil, err
		}

		dbRecords = append(dbRecords, DbRecords{
			TableName: "user_refresh_tokens",
			Count:     len(users),
			Data: []map[string]interface{}{
				{
					"user_id":       InsertUserId,
					"refresh_token": refreshTokenString,
					"expires_at":    user.CreatedAt.Add(24 * 7 * time.Hour),
					"created_at":    user.CreatedAt,
					"updated_at":    user.UpdatedAt,
				},
			},
		})

	}
	return dbRecords, nil
}

func ExistsRecord(db *sql.DB, table string, filter map[string]interface{}) bool {
	record := GetRecords(db, table, filter)
	return len(record) > 0
}

func GetRecords(db *sql.DB, table string, filter map[string]interface{}) []DbRecords {
	records := []DbRecords{}
	query := "SELECT * FROM " + table + " WHERE "
	args := []interface{}{}
	i := 0
	for k, v := range filter {
		if i > 0 {
			query += " AND "
		}
		query += k + " = ?"
		args = append(args, v)
		i++
	}
	rows, err := db.Query(query, args...)
	if err != nil {
		return records
	}
	cols, err := rows.Columns()
	if err != nil {
		return records
	}
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		if err := rows.Scan(columnPointers...); err != nil {
			return records
		}
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}
		records = append(records, DbRecords{
			TableName: table,
			Data:      []map[string]interface{}{m},
		})
	}
	return records
}
