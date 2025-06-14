package initHandler

import (
	"fullstack-journal/app/entity/companyEntity"
	"fullstack-journal/app/filters"
	"fullstack-journal/app/response"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetAll(db *gorm.DB)echo.HandlerFunc{
	return func(c echo.Context) error {
		companyAccess, err := dropdownCompanyAccess(db, c)
		if err != nil{
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"error": filters.Translate(c, "serverError"),
			})
		}

		resInit := response.ResInit{
			CompanyAccess: companyAccess,
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"data": resInit,
			"success": true,
		})
	}
}

func dropdownCompanyAccess(db *gorm.DB, c echo.Context)([]response.ResCompanyAccess, error){
	user := c.Get("user").(map[string]interface{})
	userId, _ := uuid.Parse(user["id"].(string))

	companyAccess := []response.ResCompanyAccess{}

	company, err := companyEntity.GetByUserId(db, userId)
	if err != nil{
		return companyAccess, err
	}	

	for _, comp := range company{
		id, _ := uuid.Parse(comp["id"].(string))

		var photo *string
		if val, ok := comp["photo"].(string); ok {
			photo = &val
		}
		companyAccess = append(companyAccess, response.ResCompanyAccess{
			ID: id,
			Name: comp["name"].(string),
			Photo: photo,
			IsOwner: comp["is_owner"].(int64),
			Address: comp["address"].(string),
			Role: comp["role_name"].(string),
		})
	}

	return companyAccess, nil

}