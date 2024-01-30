package config

import (
	"os"
	"strconv"
	"github.com/joho/godotenv"
)

type ConfigType struct {
	DATABASE_URL string
	JWT_SECRET_KEY string
	DATABASE_USER string
	DATABASE_PASSWORD string
	DATABASE_NAME string
	ACCESS_TOKEN_EXPIRY int
	REFRESH_TOKEN_EXPIRY int
}

func GetConfig () ConfigType {
	godotenv.Load(".env")

	var DATABASE_USER = os.Getenv("POSTGRES_USER")
	var DATABASE_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	var DATABASE_NAME = os.Getenv("POSTGRES_DB")
	var DATABASE_URL = os.Getenv("DATABASE_URL")
	var JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")
	var ACCESS_TOKEN_EXPIRY, _ = strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRY"))
	var REFRESH_TOKEN_EXPIRY, _ = strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRY"))

	return ConfigType{
		DATABASE_URL: DATABASE_URL,
		JWT_SECRET_KEY: JWT_SECRET_KEY,
		ACCESS_TOKEN_EXPIRY: ACCESS_TOKEN_EXPIRY,
		REFRESH_TOKEN_EXPIRY: REFRESH_TOKEN_EXPIRY,
		DATABASE_USER: DATABASE_USER,
		DATABASE_PASSWORD: DATABASE_PASSWORD,
		DATABASE_NAME: DATABASE_NAME,
	}
}