package kernel

import (
	"fullstack-journal/app/config"
	"fullstack-journal/app/filters"
	"fullstack-journal/app/routes"
	"log"

	"github.com/labstack/echo/v4"
)

func StartApplication(e *echo.Echo){
	// Load DB
	db, err := config.InitDB()
	if err != nil{
		log.Fatal("Error init db", err)
	}

	// Load Lang File
	if err := filters.LoadTranslations(); err != nil{
		log.Fatal("Error load file lang", err)
	}

	// Load Env File
	env, err := config.GetEnv()
	if err != nil{
		log.Fatal("Error load env", err)
	}
	port := env["port"].(string)

	// Register API
	routes.ApiV1(e, db)

	if env["appEnv"] == "development"{
		e.Debug = true
	}else{
		e.Debug = false
	}

	// Start App
	e.Logger.Fatal(e.Start(":"+port))
}