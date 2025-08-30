package dropwdownHandler

import (
	"fullstack-journal/app/entity/metadataEntity"
	"fullstack-journal/app/filters"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

func GetCategoryAccounts(db *gorm.DB) echo.HandlerFunc{
	return func(c echo.Context) error {
		categoryAccount, err := metadataEntity.FindCategoryAccount(db)
		if err != nil {
			log.Error("Error Get CategoryAccount:", err)
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"error":   filters.Translate(c, "serverError"),
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"data": categoryAccount,
		})
	}
}
