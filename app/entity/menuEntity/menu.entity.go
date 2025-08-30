package menuEntity

import (
	"errors"
	"fullstack-journal/app/entity"
	"fullstack-journal/app/response"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetMenuByRoleId(db *gorm.DB, roleId uuid.UUID) ([]response.ResMenu, error) {
	result := []response.ResMenu{}
	var parentMenuQry []entity.Menu

	// Ambil parent menu
	if err := db.Model(&entity.Menu{}).
		Where("deleted_at IS NULL").
		Where("parent_id = ?", 0).
		Order(`"order" ASC`).
		Scan(&parentMenuQry).Error; err != nil {
		log.Println("Exception get parentMenu:", err.Error())
		return result, err
	}

	// Iterasi parent menu dan cek permission
	for _, parent := range parentMenuQry {
		permissionView, err := findFirstPermissionView(db, roleId, parent.ID)
		if err != nil {
			log.Println("Exception get permission parentMenu:", err.Error())
			return result, err
		}

		if permissionView["id"] != nil {
			result = append(result, response.ResMenu{
				ID:       parent.ID,
				ParentId: parent.ParentId,
				NameId:   parent.NameId,
				NameEn:   parent.NameEn,
				Icon:     parent.Icon,
				Url:      parent.Url,
				Order:    parent.Order,
				Child:    []response.ResMenu{}, // initialize child empty
			})
		}
	}

	// Ambil dan masukkan child menu untuk setiap parent yang valid
	for i, parent := range result {
		var childMenuQry []entity.Menu // penting: reset setiap iterasi

		if err := db.Model(&entity.Menu{}).
			Where("deleted_at IS NULL").
			Where("parent_id = ?", parent.ID).
			Order(`"order" ASC`).
			Scan(&childMenuQry).Error; err != nil {
			log.Println("Exception get childMenuQry:", err.Error())
			return result, err
		}

		childMenus := make([]response.ResMenu, 0)

		for _, child := range childMenuQry {
			permissionView, err := findFirstPermissionView(db, roleId, child.ID)
			if err != nil {
				log.Println("Exception get permission childMenu:", err.Error())
				return result, err
			}

			if permissionView["id"] != nil {
					childMenus = append(childMenus, response.ResMenu{
					ID:       child.ID,
					ParentId: child.ParentId,
					NameId:   child.NameId,
					NameEn:   child.NameEn,
					Icon:     child.Icon,
					Url:      child.Url,
					Order:    child.Order,
					Child:    []response.ResMenu{}, // initialize child empty
				})
			}
		}

		// set child menu hanya kalau ada
		if len(childMenus) > 0 {
			result[i].Child = childMenus
		}
	}

	return result, nil
}


func findFirstPermissionView(db *gorm.DB, roleId uuid.UUID, menuId int64)(map[string]interface{}, error){
	result := map[string]interface{}{}
	query := db.Table("role_has_permission").
		Select("role_has_permission.*,permission.menu_id").
		Joins("LEFT JOIN permission ON permission.id = role_has_permission.permission_id").
		Where("role_has_permission.role_id = ?", roleId).
		Where("permission.menu_id = ?", menuId).
		Where("permission.permission_view = ?", 1).
		Where("role_has_permission.deleted_at IS NULL").
		Limit(1)

	if err := query.Scan(&result).Error; err != nil{
		log.Println("Query error:", err)
		return result, errors.New(err.Error())
	}

	return result, nil
}