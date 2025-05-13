package roleHasPermission

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetPermissionCodeByRoleId(db *gorm.DB, roleId uuid.UUID)([]string, error){
	code := make([]string, 0)
	result := []map[string]interface{}{}
	
	query := `
		SELECT 	
			role_has_permission.*,
			permission.code
		FROM role_has_permission
		LEFT JOIN permission ON permission.id = role_has_permission.permission_id
		WHERE role_has_permission.role_id = ?
		AND role_has_permission.deleted_at IS NULL
	`

	if err := db.Raw(query, roleId).Scan(&result).Error; err != nil{
		fmt.Println("Query error:", err)
		return code, err
	}

	for _, res := range result{
		code = append(code, res["code"].(string))
	}

	return code, nil
}