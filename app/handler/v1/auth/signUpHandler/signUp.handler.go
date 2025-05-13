package signUpHandler

import (
	"errors"
	"fullstack-journal/app/entity"
	"fullstack-journal/app/entity/usersEntity"
	"fullstack-journal/app/filters"
	"fullstack-journal/app/helpers/globalFunc"
	"fullstack-journal/app/helpers/validationFunc"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SignUp(db *gorm.DB)echo.HandlerFunc {
	return func(c echo.Context) error {
		safeInput, err := vSignUp(db, c)
		if err != nil{
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"error": err.Error(),
			})
		}

		hashPass, err := globalFunc.Hash(safeInput["password"].(string))
		if err != nil{
			log.Fatal("error hash", err.Error())
			return err
		}

		user := entity.Users{
			ID: uuid.New(),
			Name: safeInput["name"].(string),
			Email: safeInput["email"].(string),
			Password: hashPass,
			PasswordView: safeInput["password"].(string),
			PhoneNumber: safeInput["phoneNumber"].(string),
		}

		db.Create(&user)
		
		return c.JSON(http.StatusCreated, map[string]interface{}{
			"success": true,
			"message": filters.Translate(c, "usersCreated"),
		})
	}
}

func vSignUp(db *gorm.DB, c echo.Context)(map[string]interface{}, error){
	safeInput := make(map[string]interface{})

	name := c.FormValue("name")
	email := c.FormValue("email")
	phoneNumber := c.FormValue("phoneNumber")
	password := c.FormValue("password")

	if name == ""{
		return safeInput, errors.New(
			filters.Translate(c, "nameRequired"),
		)
	}
		
	if err := validationFunc.VEmail(c, email); err != nil{
		return safeInput, err;
	}

	if err := validationFunc.VPhoneNumber(c, phoneNumber); err != nil{
		return safeInput, err
	}

	if len(password) < 6 {
		return safeInput, errors.New(
			filters.Translate(c, "passwordMin6"),
		)
	}

	users := usersEntity.FindFirstByEmail(db, email)
	if users.ID != uuid.Nil{
		return safeInput, errors.New(
			filters.Translate(c, "emailAlreadyUsed"),
		)
	}

	safeInput["name"] = name
	safeInput["email"] = email
	safeInput["phoneNumber"] = phoneNumber
	safeInput["password"] = password

	return safeInput, nil
}
