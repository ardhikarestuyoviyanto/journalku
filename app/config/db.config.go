package config

import (
	"fmt"
	"fullstack-journal/app/entity"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error){
	env, err := GetEnv()
	
	if err != nil{
	  log.Fatal("Error load env file :", err.Error())
	  return nil, err
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", 
	env["dbHost"], env["dbUser"], env["dbPassword"], env["dbName"], env["dbPort"])

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil{
		log.Fatal("failed to connect to database :", err)
		return nil, err
	}

	startMigration(db)
	// seed.IndonesiaLocationSeed(db)
	// seed.MenuSeed(db)

	return db, nil

}

func startMigration(db *gorm.DB){
	db.AutoMigrate(
		entity.Users{},
		entity.Province{},
		entity.Regency{},
		entity.SubDistrict{},
		entity.Company{},
		entity.UserCompanyAccess{},
		entity.Menu{},
		entity.Role{},
		entity.Permission{},
		entity.RoleHasPermission{},
	)
}
