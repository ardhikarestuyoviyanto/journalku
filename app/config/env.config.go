package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnv() (map[string]interface{}, error) {
	err := godotenv.Load()
	if err != nil{
		log.Fatal("Error load .env file", err.Error())
		return nil, err
	}

	port := os.Getenv("PORT")
	appEnv := os.Getenv("APP_ENV")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPort := os.Getenv("DB_PORT")
	dbPassword := os.Getenv("DB_PASSWORD")
	appKey := os.Getenv("APP_KEY")
	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecrets := os.Getenv("GOOGLE_CLIENT_SECRETS")
	appUrl := os.Getenv("APP_URL")

	envValue := map[string]interface{}{
		"port": port,
		"appEnv": appEnv,
		"dbHost": dbHost,
		"dbName": dbName,
		"dbUser": dbUser,
		"dbPort": dbPort,
		"dbPassword": dbPassword,
		"appKey": appKey,
		"googleClientId": googleClientId,
		"googleClientSecrets":googleClientSecrets,
		"appUrl": appUrl,
	}

	return envValue, nil
}