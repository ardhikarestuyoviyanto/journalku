package routes

import (
	"fullstack-journal/app/filters"
	"fullstack-journal/app/handler/v1/auth/signInHandler"
	"fullstack-journal/app/handler/v1/auth/signUpHandler"
	"fullstack-journal/app/handler/v1/companyHandler"
	"fullstack-journal/app/handler/v1/metaData/locationHandler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

func ApiV1(e *echo.Echo, db *gorm.DB) {
	// STATIC FILE (pastikan ini sebelum AuthFilters)
	e.Static("/storage", "storage")

	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(filters.LangFilters)
	
	// MODUL AUTH
	e.POST("/api/v1/auth/signUp", signUpHandler.SignUp(db))
	e.POST("/api/v1/auth/signIn", signInHandler.SignIn(db))
	e.GET("/api/v1/auth/google/signIn", signInHandler.SignInGoogle)
	e.GET("/api/v1/auth/google/callback", signInHandler.SignGoogleCallback(db))

	// ROUTE PROTECTED AUTH
   	routeAuth := e.Group("/api/v1")
	routeAuth.Use(filters.AuthFilters)
	routeAuth.GET("/metadata/province", locationHandler.GetProvince(db))
	routeAuth.GET("/metadata/regency", locationHandler.GetRegency(db))
	routeAuth.GET("/metadata/subDistrict", locationHandler.GetSubDistrict(db))
	// MODUL COMPANY BY USER
	routeAuth.POST("/company", companyHandler.Store(db))
	routeAuth.GET("/company", companyHandler.GetAll(db))
	// MODUL CHOOSE COMPANY
	routeAuth.POST("/choose-company", companyHandler.ChooseCompany(db))
	
}
