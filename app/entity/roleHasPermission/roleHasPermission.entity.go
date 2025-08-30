package roleHasPermission

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetPermissionCodeByRoleId(db *gorm.DB, roleId uuid.UUID)([]string, error){
	code := make([]string, 0)
	result := []map[string]interface{}{}

	query := db.Table("role_has_permission").
		Select("role_has_permission.*,permission.code").
		Joins("LEFT JOIN permission ON permission.id = role_has_permission.permission_id").
		Where("role_has_permission.role_id = ?", roleId).
		Where("role_has_permission.deleted_at IS NULL")

	if err := query.Scan(&result).Error; err != nil{
		log.Println("Query error:", err)
		return code, errors.New(err.Error())
	}

	for _, res := range result{
		code = append(code, res["code"].(string))
	}

	return code, nil
}