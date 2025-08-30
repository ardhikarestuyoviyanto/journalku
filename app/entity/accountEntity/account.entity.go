package accountEntity

import (
	"errors"
	"fmt"
	"fullstack-journal/app/response"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaginateComponent struct {
	Limit  int
	Offset int
	Order  string
	SortBy string
	Where  map[string]interface{}
	Search string
}


func FindAll(db *gorm.DB, pagination PaginateComponent) ([]response.ResAccount, int64, error) {
	resAccount := make([]response.ResAccount, 0)
	var totalRows int64

	query := db.Table("account").
		Select("account.*, tb_category_account.description AS category_account").
		Joins("LEFT JOIN metadata AS tb_category_account ON tb_category_account.id = account.category_account_id").
		Where("account.deleted_at IS NULL").
		Where("account.company_id = ?", pagination.Where["companyId"])

	if pagination.Search != "" {
		like := "%" + pagination.Search + "%"
		query = query.Where(`
			account.number_account ILIKE ? OR
			account.name ILIKE ? OR
			account.description ILIKE ? OR
			tb_category_account.description ILIKE ?`,
			like, like, like, like,
		)
	}

	if err := query.Count(&totalRows).Error; err != nil {
		return resAccount, 0, err
	}

	order := fmt.Sprintf("%s %s", pagination.SortBy, pagination.Order)
	query = query.Order(order).Limit(pagination.Limit).Offset(pagination.Offset)

	if err := query.Scan(&resAccount).Error; err != nil {
		log.Println("Query error:", err)
		return resAccount, totalRows, errors.New(err.Error())
	}

	return resAccount, totalRows, nil
}

func FindFirstByNumberAccount(db *gorm.DB, numberAccount string, companyId uuid.UUID)(response.ResAccount, error){
	var resAccount response.ResAccount
	
	query := db.Table("account").
		Where("deleted_at IS NULL").
		Where("number_account = ?", numberAccount).
		Where("company_id = ?", companyId).
		Limit(1)

	if err := query.Scan(&resAccount).Error; err != nil {
		log.Println("Query error:", err)
		return resAccount,  errors.New(err.Error())
	}

	return resAccount, nil
}
