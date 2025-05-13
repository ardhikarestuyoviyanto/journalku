package filters

import (
	"fullstack-journal/app/helpers/globalFunc"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func AuthFilters(next echo.HandlerFunc)echo.HandlerFunc{
	return func(c echo.Context) error {
		// Get Token
		authHeader  := c.Request().Header.Get("Authorization")
		// Cek Apakah Token Kosong
		if authHeader  == ""{
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"error": "Unauthorized",
			})
		
		}
		// Cek Struktur Token
		// Baerer xxxxx
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer"{
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"error": "Unauthorized",
			})
		} 

		// Get Token
		token := tokenParts[1]
		decodedToken, err := globalFunc.DecodeJwt(token)
		if err !=nil{
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"error": "Unauthorized",
			})
		}

		c.Set("user", decodedToken["user"])
		return next(c)
	}
}