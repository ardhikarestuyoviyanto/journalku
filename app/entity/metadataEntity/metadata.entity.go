package metadataEntity

import (
	"errors"
	"fullstack-journal/app/response"
	"log"

	"gorm.io/gorm"
)

func FindCategoryAccount(db *gorm.DB)([]response.ResCategoryAccount, error){
	categoryAccount := make([]response.ResCategoryAccount, 0)
	query := db.Table("metadata").
		Where("deleted_at IS NULL").
		Where("key = ?", "category_account").
		Order("value ASC")

	if err := query.Scan(&categoryAccount).Error; err != nil{
		log.Println("Query error:", err)
		return categoryAccount, errors.New(err.Error())
	}

	return categoryAccount, nil
}

func FindFirstCategoryAccount(db *gorm.DB, id int64)(response.ResCategoryAccount, error) {
	var categoryAccount response.ResCategoryAccount
	query := db.Table("metadata").
		Where("deleted_at IS NULL").
		Where("key = ?", "category_account").
		Limit(1)
	
	if err := query.Scan(&categoryAccount).Error; err != nil{
		log.Println("Query error:", err)
		return categoryAccount, errors.New(err.Error())
	}

	return categoryAccount, nil
}