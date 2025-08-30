package accountHandler

import (
	"errors"
	"fullstack-journal/app/entity"
	"fullstack-journal/app/entity/accountEntity"
	"fullstack-journal/app/entity/metadataEntity"
	"fullstack-journal/app/filters"
	"fullstack-journal/app/helpers/globalFunc"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

func GetAll(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Ambil query param
		pageStr := c.QueryParam("page")
		perPageStr := c.QueryParam("perPage")
		orderStr := strings.ToLower(c.QueryParam("order"))
		sortByStr := c.QueryParam("sortBy")
		search := c.QueryParam("search")


		// Ambil company ID dari context
		user := c.Get("user").(map[string]interface{})
		currentCompany := user["currentCompany"].(map[string]interface{})
		currentCompanyIdStr, ok := currentCompany["id"].(string)
		if !ok {
			log.Error("CurrentCompany ID is not a string")
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"error":   filters.Translate(c, "serverError"),
			})
		}

		currentCompanyId, err := uuid.Parse(currentCompanyIdStr)
		if err != nil {
			log.Error("Failed to parse company UUID:", err)
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"error":   filters.Translate(c, "serverError"),
			})
		}

		// Validasi dan set default order dan sortBy
		allowedOrder := []string{"asc", "desc"}
		allowedSortBy := []string{
			"account.created_at", "account.number_account", "account.name",
			"account.category_account_id", "account.description", "account.status_archive",
		}

		order := "asc"
		if globalFunc.BuildSet(allowedOrder)[orderStr] {
			order = orderStr
		}

		sortBy := "account.number_account"
		if globalFunc.BuildSet(allowedSortBy)[sortByStr] {
			sortBy = sortByStr
		}

		// Convert perPage
		perPage := 20
		if perPageStr != "" {
			if v, err := strconv.Atoi(perPageStr); err == nil && v > 0 {
				perPage = v
			}
		}

		// Convert page to offset
		offset := 0
		if pageStr != "" {
			if page, err := strconv.Atoi(pageStr); err == nil && page > 1 {
				offset = (page - 1) * perPage
			}
		}

		// Buat komponen paginasi
		paginateComponent := accountEntity.PaginateComponent{
			Limit:  perPage,
			Offset: offset,
			Order:  order,
			SortBy: sortBy,
			Where: map[string]interface{}{
				"companyId": currentCompanyId,
			},
			Search: search,
		}

		// Ambil data dari entity
		accounts, totalRows, err := accountEntity.FindAll(db, paginateComponent)
		if err != nil {
			log.Error("Error FindAll:", err)
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"error":   filters.Translate(c, "serverError"),
			})
		}

		lastPage := int(math.Ceil(float64(totalRows) / float64(perPage)))


		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"data":    map[string]interface{}{
				"accounts": accounts,
				"lastPage": lastPage,
				"totalRows": totalRows,
			},
		})
	}
}

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

		// CompanyId
		currentCompany := user["currentCompany"].(map[string]interface{})
		companyId, err := uuid.Parse(currentCompany["id"].(string))

		if err != nil{
			log.Print("companyId is not string", err)
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"error": err.Error(),
			})
		}

		var description *string
		if safeInput["description"] != ""{
			descriptionStr, ok := safeInput["description"].(string)
			if !ok{
				log.Print("description is not string")
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"success": false,
					"error": "description is not string",
				})
			}

			description = &descriptionStr
		}

		account := entity.Account{
			CompanyId: companyId,
			CategoryAccountId: safeInput["categoryAccountId"].(int64),
			NumberAccount: safeInput["numberAccount"].(string),
			Name: safeInput["name"].(string),
			Description: description,
			IsPrimary: 0,
			IsArchive: safeInput["isArchive"].(int64),
		}

		if err := db.Create(&account).Error; err != nil{
			log.Print("failed create account", err.Error())
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"error": err.Error(),
			})

		}

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"success": true,
			"message": filters.Translate(c, "accountCreated"),
		}) 
	}
}


func vStore(db *gorm.DB, c echo.Context)(map[string]interface{}, error){
	safeInput := make(map[string]interface{})
	categoryAccountIdStr := c.FormValue("categoryAccountId")
	numberAccount := c.FormValue("numberAccount")
	name := c.FormValue("name")
	statusArchiveStr := c.FormValue("statusArchive")
	description := c.FormValue("description")

	if categoryAccountIdStr == ""{
		return safeInput, errors.New(filters.Translate(c, "categoryAccountRequired"))
	}

	if numberAccount == ""{
		return safeInput, errors.New(filters.Translate(c, "numberAccountRequired"))
	}

	if name == ""{
		return safeInput, errors.New(filters.Translate(c, "nameAccountRequired"))
	}

	if statusArchiveStr == ""{
		return safeInput, errors.New(filters.Translate(c, "statusArchiveRequired"))
	}

	categoryAccountId, err := strconv.Atoi(categoryAccountIdStr)
	if err != nil{
		return safeInput, errors.New("Category Account Must Be Numeric")
	}

	statusArchive, err := strconv.Atoi(statusArchiveStr)
	if err != nil{
		return safeInput, errors.New("Status Archive Must Be Numeric")
	}

	if !globalFunc.Contains([]int{1, 0}, statusArchive){
		return safeInput, errors.New("Status Archive Must Be 1 or 0")
	}

	categoryAccount, _ := metadataEntity.FindFirstCategoryAccount(db, int64(categoryAccountId))
	if categoryAccount.ID == 0{
		return safeInput, errors.New("Category Account Not Found")
	}
	
	user := c.Get("user").(map[string]interface{})
	// CompanyId
	currentCompany := user["currentCompany"].(map[string]interface{})
	companyId, _ := uuid.Parse(currentCompany["id"].(string))
	resAccount, _ := accountEntity.FindFirstByNumberAccount(db, numberAccount, companyId)

	if resAccount.Name != ""{
		return safeInput, errors.New(filters.Translate(c, "numberAccountUnique"))
	}

	safeInput["categoryAccountId"] = int64(categoryAccountId)
	safeInput["numberAccount"] = numberAccount
	safeInput["name"] = name
	safeInput["statusArchive"] = int64(statusArchive)
	safeInput["description"] = description

	return safeInput, nil
}