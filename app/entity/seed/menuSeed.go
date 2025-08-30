package seed

import (
	"encoding/json"
	"fullstack-journal/app/entity"
	"log"
	"os"

	"gorm.io/gorm"
)

type menuJson struct {
	ID         int64           	`json:"id"`
	ParentId   int64           	`json:"parentId"`
	NameId     string          	`json:"nameId"`
	NameEn     string          	`json:"nameEn"`
	Icon       *string         	`json:"icon"` 
	URL        string          	`json:"url"`
	Order      int64           	`json:"order"`
	Permission []permissionJson `json:"permission"`
	Child      []menuJson      	`json:"child"`
}

type permissionJson struct {
	MenuId int64  `json:"menuId"`
	Code   string `json:"code"`
	NameId string `json:"nameId"`
	NameEn string `json:"nameEn"`
	PermissionView int64 `json:"permissionView"`
}

func loadMenu() ([]menuJson, error) {
	var menus []menuJson
	data, err := os.ReadFile("storage/seed/menu.json")
	if err != nil {
		log.Println("Failed read file menu.json", err)
		return menus, err
	}
	if err := json.Unmarshal(data, &menus); err != nil {
		log.Println("Failed parse menu.json", err)
		return menus, err
	}
	return menus, nil
}

func insertMenu(db *gorm.DB, menu menuJson) error {
	menuEntity := entity.Menu{
		ID:       menu.ID,
		ParentId: menu.ParentId,
		NameId:   menu.NameId,
		NameEn:   menu.NameEn,
		Icon:     menu.Icon,
		Url:      menu.URL,
		Order:    menu.Order,
	}
	if err := db.Create(&menuEntity).Error; err != nil {
		return err
	}

	for _, perm := range menu.Permission {
		permEntity := entity.Permission{
			MenuId: perm.MenuId,
			Code:   perm.Code,
			NameId: perm.NameId,
			NameEn: perm.NameEn,
			PermissionView: perm.PermissionView,
		}
		if err := db.Create(&permEntity).Error; err != nil {
			return err
		}
	}

	// Recursively insert child menus
	for _, child := range menu.Child {
		if err := insertMenu(db, child); err != nil {
			return err
		}
	}

	return nil
}

func MenuSeed(db *gorm.DB) error {
	// Truncate Table 
	db.Exec("TRUNCATE TABLE menu RESTART IDENTITY CASCADE")
	db.Exec("TRUNCATE TABLE permission RESTART IDENTITY CASCADE")

	menus, err := loadMenu()
	if err != nil {
		return err
	}

	for _, m := range menus {
		if err := insertMenu(db, m); err != nil {
			log.Printf("Gagal insert menu ID %d: %v\n", m.ID, err)
			return err
		}
	}

	return nil
}
