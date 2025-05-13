package companyEntity

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func FindFirstByCompanyNameAndUserId(db *gorm.DB, companyName string, userId uuid.UUID) (map[string]interface{}, error) {
	result := map[string]interface{}{}

	query := `
		SELECT 
			company.*, 
			user_company_access.is_owner 
		FROM company
		LEFT JOIN user_company_access ON company.id = user_company_access.company_id 
		WHERE company.name = ? 
		AND user_company_access.user_id = ? 
		AND user_company_access.is_owner = 1
		AND company.deleted_at IS NULL
		LIMIT 1
	`

	if err := db.Raw(query, companyName, userId).Scan(&result).Error; err != nil {
		fmt.Println("Query error:", err)
		return result, errors.New(err.Error())
	}

	return result, nil
}

func GetByUserId(db *gorm.DB, userId uuid.UUID)([]map[string]interface{}, error){
	result := []map[string]interface{}{}
	query := `
		SELECT 
			company.*, 
			user_company_access.is_owner, 
			role.name as role_name
		FROM company
		LEFT JOIN user_company_access ON company.id = user_company_access.company_id 
		LEFT JOIN role ON role.id = user_company_access.role_id
		WHERE user_company_access.user_id = ? 
		AND company.deleted_at IS NULL
		AND user_company_access.deleted_at IS NULL
	`
	if err := db.Raw(query, userId).Scan(&result).Error; err != nil{
		fmt.Println("Query error:", err)
		return result, errors.New(err.Error())
	}

	return result, nil
}

func FindFirstByUserIdAndCompanyId(db *gorm.DB, userId uuid.UUID, companyId uuid.UUID)(map[string]interface{}, error){
	result := map[string]interface{}{}
	query := `
		SELECT 
			company.*, 
			user_company_access.role_id,
			user_company_access.is_owner, 
			user_company_access.id AS user_company_access_id,
			role.name AS role_name
		FROM company
		LEFT JOIN user_company_access ON company.id = user_company_access.company_id 
		LEFT JOIN role ON role.id = user_company_access.role_id
		WHERE user_company_access.user_id = ? 
		AND user_company_access.company_id = ?
		AND company.deleted_at IS NULL
		AND user_company_access.deleted_at IS NULL
		LIMIT 1
	`

	if err := db.Raw(query, userId, companyId).Scan(&result).Error; err != nil{
		fmt.Println("Query error:", err)
		return result, errors.New(err.Error())
	}

	return result, nil
}