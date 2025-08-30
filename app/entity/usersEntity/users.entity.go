package usersEntity

import (
	"fullstack-journal/app/entity"

	"gorm.io/gorm"
)

func FindFirstByEmail(db *gorm.DB, email string) entity.Users{
	var users entity.Users
	query := db.Table("users").
		Where("email=?", email).
		Where("deleted_at IS NULL").
		Limit(1)

   	query.Scan(&users)
	return users
}