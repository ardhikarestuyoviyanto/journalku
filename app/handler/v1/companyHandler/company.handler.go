package companyHandler

import (
	"errors"
	"fullstack-journal/app/entity"
	"fullstack-journal/app/entity/companyEntity"
	"fullstack-journal/app/entity/menuEntity"
	"fullstack-journal/app/entity/roleHasPermission"
	"fullstack-journal/app/filters"
	"fullstack-journal/app/helpers/globalFunc"
	"fullstack-journal/app/helpers/validationFunc"
	"fullstack-journal/app/response"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)


func Store(db *gorm.DB)echo.HandlerFunc{
	return func(c echo.Context) error {
		user := c.Get("user").(map[string]interface{})
		safeInput, err := vStore(db, c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"error": err.Error(),
			})
		}

		var photo *string
		if val, ok := safeInput["photo"].(string); ok {
			photo = &val
		}

		// Start Trx
	   	err = db.Transaction(func(tx *gorm.DB) error {
			// Insert Company
			company := entity.Company{
				ID:            uuid.New(),
				ProvinceId:    safeInput["provinceId"].(int64),
				RegencyId:     safeInput["regencyId"].(int64),
				SubDistrictId: safeInput["subDistrictId"].(int64),
				Name:          safeInput["name"].(string),
				Address:       safeInput["address"].(string),
				Photo:         photo,
			}

			if err := tx.Create(&company).Error; err != nil {
				return err
			}

			// Create Role Owner
			role := entity.Role{
				ID:        uuid.New(),
				CompanyId: company.ID,
				Name:      "Owner",
			}

			if err := tx.Create(&role).Error; err != nil {
				return err
			}

			// Get all permissions
			var permissions []entity.Permission
			if err := tx.Where("deleted_at IS NULL").Find(&permissions).Error; err != nil {
				return err
			}

			// Create RoleHasPermission (bulk insert recommended)
			var rolePerms []entity.RoleHasPermission
			for _, perm := range permissions {
				rolePerms = append(rolePerms, entity.RoleHasPermission{
					RoleId:       role.ID,
					PermissionId: perm.ID,
				})
			}

			if len(rolePerms) > 0 {
				if err := tx.Create(&rolePerms).Error; err != nil {
					return err
				}
			}

			// Insert UserCompanyAccess
			userId, err := uuid.Parse(user["id"].(string))
			if err != nil {
				return err
			}

			userCompanyAccess := entity.UserCompanyAccess{
				UserId:         userId,
				CompanyId:      company.ID,
				RoleId:         role.ID,
				IsOwner:        1,
				CurrentCompany: 0,
			}

			if err := tx.Create(&userCompanyAccess).Error; err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			log.Println("Transaction failed:", err)
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"success": true,
				"message": filters.Translate(c, "serverError"),
			})
		}

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"success": true,
			"message": filters.Translate(c, "companyCreated"),
		})
	}
}

func GetAll(db *gorm.DB)echo.HandlerFunc{
	return func(c echo.Context) error {
		user := c.Get("user").(map[string]interface{})
		userId, _ := uuid.Parse(user["id"].(string))

		company, err := companyEntity.GetByUserId(db, userId)
		if err != nil{
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"error": filters.Translate(c, "serverError"),
			})
		}

		companyList := []response.ResCompanyAccess{}
		for _, comp := range company{
			id, _ := uuid.Parse(comp["id"].(string))

			var photo *string
			if val, ok := comp["photo"].(string); ok {
				photo = &val
			}
			companyList = append(companyList, response.ResCompanyAccess{
				ID: id,
				Name: comp["name"].(string),
				Photo: photo,
				IsOwner: comp["is_owner"].(int64),
				Address: comp["address"].(string),
				Role: comp["role_name"].(string),
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"data": companyList,
			"success": true,
		})
	}
}

func ChooseCompany(db *gorm.DB)echo.HandlerFunc{
	return func(c echo.Context) error {
		user := c.Get("user").(map[string]interface{})
		userId, _ := uuid.Parse(user["id"].(string))

		safeInput, err := vChooseCompany(c)
		if err != nil{
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"error": err.Error(),
			})
		}
		
		companyAccess, err := companyEntity.FindFirstByUserIdAndCompanyId(
			db,
			userId,
			safeInput["companyId"],
		)

		if err != nil{
			log.Println("Error get company access", err.Error())
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"error": filters.Translate(c, "forhibiden"),
			})	
		}

		if companyAccess["id"] == nil{
			log.Println("companyAccess nil (logic exception)", companyAccess)
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"error": filters.Translate(c, "serverError"),
			})	
		} 

		roleId, _ := uuid.Parse(companyAccess["role_id"].(string))

		permission, err := roleHasPermission.GetPermissionCodeByRoleId(db, roleId)
		if err != nil{
			log.Println("error get permission", err)
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"error": filters.Translate(c, "serverError"),
			})	
		}

		menu, err := menuEntity.GetMenuByRoleId(db, roleId)
		if err != nil{
			log.Println("error get menu", err)
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"error": filters.Translate(c, "serverError"),
			})
		}

		companyId, _ := uuid.Parse(companyAccess["id"].(string))

		var photo *string
		if val, ok := companyAccess["photo"].(string); ok {
			photo = &val
		}

		currentCompany := response.ResCurrentCompany{
			ID:         companyId,
			Name:       companyAccess["name"].(string),
			Photo:      photo,
			IsOwner:    companyAccess["is_owner"].(int64),
			Address:    companyAccess["address"].(string),
			Role:       companyAccess["role_name"].(string),
			Permission: permission,
			ResMenu: menu,
		}

		// Payload New Token
		exp := time.Now().Add(time.Hour * 24).Unix() 
		payload := response.ResToken{
			User: response.ResUserToken{
				ID: userId,
				Name: user["name"].(string),
				Email: user["email"].(string),
				CurrentCompany: currentCompany,
			},
			Exp: exp,

		}

		token, err := globalFunc.GetJwt(payload)
		if err != nil{
			log.Println("Error generate jwt token", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"error": filters.Translate(c, "networkError"),
			})
		}

		// Update token in user
		db.Model(&entity.Users{}).Where("id = ?", user["id"]).Update("token", token)
		// Update CurrentCompany
		db.Model(&entity.UserCompanyAccess{}).Where("user_id = ?", user["id"]).Update("current_company", 0)
		db.Model(&entity.UserCompanyAccess{}).Where("id = ?", companyAccess["user_company_access_id"]).Update("current_company", 1)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"data": map[string]interface{}{
				"token": token,
				"user": payload.User,			
			},
			"success": true,
		}); 
	}
}

func vChooseCompany(c echo.Context)(map[string]uuid.UUID, error){
	safeInput := make(map[string]uuid.UUID)

	companyIdStr := c.FormValue("companyId")
	companyId, err := uuid.Parse(companyIdStr)
	if err != nil{
		return safeInput, errors.New("Invalid companyId")
	}

	safeInput["companyId"] = companyId
	return safeInput,nil
}

func vStore(db *gorm.DB, c echo.Context)(map[string]interface{}, error){
	safeInput := make(map[string]interface{})

	user := c.Get("user").(map[string]interface{})
	name := c.FormValue("name")
	provinceIdx := c.FormValue("provinceId")
	regencyIdx := c.FormValue("regencyId")
	subDistrictIdx := c.FormValue("subDistrictId")
	address := c.FormValue("address")
	photo, errPhoto := c.FormFile("photo")
	
	if errPhoto == nil{
		// Validate Photo
		err := validationFunc.VFile([]string{".png", ".jpg", ".jpeg"}, photo, 3, c)
		if err != nil{
			return safeInput, err
		}
	}

	if name == ""{
		return safeInput, errors.New(filters.Translate(c, "companyRequired"))
	}

	if address == ""{
		return safeInput, errors.New(filters.Translate(c, "addressRequired"))
	}

	provinceId, err := globalFunc.Decrypt(provinceIdx)
	if err != nil{
		return safeInput, errors.New("provinceId is invalid")
	}

	regencyId, err := globalFunc.Decrypt(regencyIdx)
	if err != nil{
		return safeInput, errors.New("regencyId is invalid")
	}

	subDistrictId, err := globalFunc.Decrypt(subDistrictIdx)
	if err != nil{
		return safeInput, errors.New("subDistrictId is invalid")
	}

	// if company exists by user id
	userId, _ := uuid.Parse(user["id"].(string))
	company, err := companyEntity.FindFirstByCompanyNameAndUserId(db, name, userId)
	if err != nil{
		return safeInput, err
	}

	if company["id"] != nil{
		return safeInput, errors.New(filters.Translate(c, "companyNameAlreadyExists"))
	}

	safeInput["name"] = name
	safeInput["provinceId"] = provinceId
	safeInput["regencyId"] = regencyId
	safeInput["subDistrictId"] = subDistrictId
	safeInput["address"] = address
	if errPhoto == nil{
		// Ada fotonya
		dstDir := "./storage/image/"
		fileName, err := globalFunc.UploadFile(dstDir, photo)
		if err != nil{
			return safeInput, err
		}
		safeInput["photo"] = fileName
	}else{
		// Nil
		safeInput["photo"] = nil
	}

	return safeInput, nil
}