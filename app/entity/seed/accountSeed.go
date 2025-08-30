package seed

import (
	"encoding/json"
	"fullstack-journal/app/entity"
	"log"
	"os"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type accountJson struct {
	CategoryId    int64  `json:"category_id"`
	Name          string  `json:"name"`
	NumberAccount string  `json:"number_account"`
	Description   *string `json:"description"`
	IsPrimary     int64   `json:"is_primary"`
	IsArchive     int64   `json:"is_archive"`
}

// loadAccount reads and parses the default account JSON file.
func loadAccount() ([]accountJson, error) {
	var accounts []accountJson

	data, err := os.ReadFile("storage/seed/account_default.json")
	if err != nil {
		log.Println("❌ Failed to read account_default.json:", err)
		return accounts, err
	}

	if err := json.Unmarshal(data, &accounts); err != nil {
		log.Println("❌ Failed to parse account_default.json:", err)
		return accounts, err
	}

	return accounts, nil
}

// AccountSeed seeds the account data into the database.
func AccountSeed(db *gorm.DB, companyId uuid.UUID) error {
	accounts, err := loadAccount()
	if err != nil {
		return err
	}

	accountList := make([]entity.Account, 0)

	for _, a := range accounts {
		accountList = append(accountList, entity.Account{
			ID:                uuid.New(),
			CompanyId:         companyId,
			CategoryAccountId: a.CategoryId,
			NumberAccount:     a.NumberAccount,
			Name:              a.Name,
			Description:       a.Description,
			IsPrimary:         a.IsPrimary,
			IsArchive:         a.IsArchive,
		})
	}

	if len(accountList) == 0 {
		log.Println("⚠️  No valid accounts to insert.")
		return nil
	}

	if err := db.Create(&accountList).Error; err != nil {
		log.Println("❌ Failed to insert accounts:", err)
		return err
	}

	log.Printf("✅ Seeded %d accounts successfully.\n", len(accountList))
	return nil
}
