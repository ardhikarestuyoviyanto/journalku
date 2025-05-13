package usersEntity

import (
	"fullstack-journal/app/entity"

	"gorm.io/gorm"
)

func FindFirstByEmail(db *gorm.DB, email string) entity.Users{
	var users entity.Users
	query := `
		SELECT *FROM users
		WHERE email=?
		AND deleted_at IS NULL
		LIMIT 1
	`
   	db.Raw(query, email).Scan(&users)
	return users
}