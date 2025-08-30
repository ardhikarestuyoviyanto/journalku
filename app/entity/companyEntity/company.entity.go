package companyEntity

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func FindFirstByCompanyNameAndUserId(db *gorm.DB, companyName string, userId uuid.UUID) (map[string]interface{}, error) {
	result := map[string]interface{}{}

	query := db.Table("company").
		Select("company.*, user_company_access.is_owner").
		Joins("LEFT JOIN user_company_access ON company.id = user_company_access.company_id").
		Where("company.name = ? ", companyName).
		Where("user_company_access.user_id = ? ", userId).
		Where("user_company_access.is_owner = ?", 1).
		Where("company.deleted_at IS NULL").
		Limit(1)

	if err := query.Scan(&result).Error; err != nil {
		log.Println("Query error:", err)
		return result, errors.New(err.Error())
	}

	return result, nil
}

func GetByUserId(db *gorm.DB, userId uuid.UUID)([]map[string]interface{}, error){
	result := []map[string]interface{}{}

	query := db.Table("company").
		Select("company.*,user_company_access.is_owner,role.name AS role_name").
		Joins("LEFT JOIN user_company_access ON company.id = user_company_access.company_id").
		Joins("LEFT JOIN role ON role.id = user_company_access.role_id").
		Where("user_company_access.user_id = ? ", userId).
		Where("company.deleted_at IS NULL").
		Where("user_company_access.deleted_at IS NULL")

	if err := query.Scan(&result).Error; err != nil{
		log.Println("Query error:", err)
		return result, errors.New(err.Error())
	}

	return result, nil
}

func FindFirstByUserIdAndCompanyId(db *gorm.DB, userId uuid.UUID, companyId uuid.UUID)(map[string]interface{}, error){
	result := map[string]interface{}{}

	query := db.Table("company").
		Select("company.*,user_company_access.role_id,user_company_access.is_owner,user_company_access.id AS user_company_access_id,role.name AS role_name").
		Joins("LEFT JOIN user_company_access ON company.id = user_company_access.company_id").
		Joins("LEFT JOIN role ON role.id = user_company_access.role_id").
		Where("user_company_access.user_id = ? ", userId).
		Where(" user_company_access.company_id = ?", companyId).
		Where("company.deleted_at IS NULL").
		Where("user_company_access.deleted_at IS NULL").
		Limit(1)

	if err := query.Scan(&result).Error; err != nil{
		log.Println("Query error:", err)
		return result, errors.New(err.Error())
	}

	return result, nil
}