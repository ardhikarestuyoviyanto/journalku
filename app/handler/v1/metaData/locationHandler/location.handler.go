package locationHandler

import (
	"errors"
	"fullstack-journal/app/entity"
	"fullstack-journal/app/helpers/globalFunc"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetProvince(db *gorm.DB)echo.HandlerFunc{
	return func(c echo.Context) error {
		var province []entity.Province
		provinceList := make([]map[string]interface{}, 0)
		// Get provinsi
		db.Model(&entity.Province{}).Order("name asc").Scan(&province)
		for _, prov := range province{
			id, _ := globalFunc.Encrypt(strconv.FormatInt(prov.ID, 10))
			provinceList = append(provinceList, map[string]interface{}{
					"id": id,
					"name": prov.Name,
				})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"data": provinceList,
		})
	}
}

func GetRegency(db *gorm.DB)echo.HandlerFunc{
	return func(c echo.Context) error {	
		safeInput, err := vGetRegency(c)

		if err != nil{
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"error": err.Error(),
			})
		}

		var regency []entity.Regency
		regencyList := make([]map[string]interface{}, 0)
		// Get Regencies
		db.Model(&entity.Regency{}).Where("province_id", safeInput["provinceId"]).Order("name asc").Scan(&regency)
		for _, reg := range regency{
			id, _ := globalFunc.Encrypt(strconv.FormatInt(reg.ID, 10))
			regencyList = append(regencyList, map[string]interface{}{
				"id": id,
				"name": reg.Name,
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"data": regencyList,
		})	
	}
}

func vGetRegency(c echo.Context)(map[string]interface{}, error){
	safeInput := make(map[string]interface{})
	
	provinceIdx := c.QueryParam("provinceId")
	if provinceIdx == ""{
		return safeInput, errors.New("query params provinceId is required")
	}

	provinceId, err := globalFunc.Decrypt(provinceIdx)
	if err != nil{
		return safeInput, errors.New("provinceId is invalid")
	}

	safeInput["provinceId"] = provinceId
	return safeInput, nil
}

func GetSubDistrict(db *gorm.DB)echo.HandlerFunc{
	return func(c echo.Context) error {
		safeInput, err := vGetSubDistrict(c)
		if err != nil{
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"error": err.Error(),
			})
		}

		var subDistrict []entity.SubDistrict
		subDistrictList := make([]map[string]interface{}, 0)
		// Get Sub Districts
		db.Model(&entity.SubDistrict{}).Where("regency_id", safeInput["regencyId"]).Order("name asc").Scan(&subDistrict)
		for _, sub := range subDistrict{
			id, _ := globalFunc.Encrypt(strconv.FormatInt(sub.ID, 10))
			subDistrictList = append(subDistrictList, map[string]interface{}{
				"id": id,
				"name": sub.Name,
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"data": subDistrictList,
		})		
	}
}

func vGetSubDistrict(c echo.Context)(map[string]interface{}, error){
	safeInput := make(map[string]interface{})

	regencyx := c.QueryParam("regencyId")
	if regencyx == ""{
		return safeInput, errors.New("query params regencyId is required")
	}

	regencyId, err := globalFunc.Decrypt(regencyx)
	if err != nil{
		return safeInput, errors.New("regencyId is invalid")
	}

	safeInput["regencyId"] = regencyId
	return safeInput, nil

}
