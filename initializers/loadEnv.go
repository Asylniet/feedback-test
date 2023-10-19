package initializers

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DBHost         string `mapstructure:"POSTGRES_HOST"`
	DBUserName     string `mapstructure:"POSTGRES_USER"`
	DBUserPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName         string `mapstructure:"POSTGRES_DB"`
	DBPort         string `mapstructure:"POSTGRES_PORT"`
	ServerPort     string `mapstructure:"PORT"`

	ClientOrigin string `mapstructure:"CLIENT_ORIGIN"`
	GinMode      string `mapstructure:"GIN_MODE"`

	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`

	EmailFrom string `mapstructure:"EMAIL_FROM"`
	SMTPHost  string `mapstructure:"SMTP_HOST"`
	SMTPPass  string `mapstructure:"SMTP_PASS"`
	SMTPPort  int    `mapstructure:"SMTP_PORT"`
	SMTPUser  string `mapstructure:"SMTP_USER"`
}

func LoadConfig() (config Config, err error) {
	if len(os.Getenv("GIN_MODE")) < 1 || os.Getenv("GIN_MODE") == "debug" {
		if err := godotenv.Load("./app.env"); err != nil {
			log.Fatal("Error loading app.env file")
		}
	}
	accessExpiresIn, _ := time.ParseDuration(os.Getenv("ACCESS_TOKEN_EXPIRED_IN"))
	refreshExpiresIn, _ := time.ParseDuration(os.Getenv("REFRESH_TOKEN_EXPIRED_IN"))
	accessMaxAge, _ := strconv.ParseInt(os.Getenv("ACCESS_TOKEN_MAXAGE"), 10, 64)
	refreshMaxAge, _ := strconv.ParseInt(os.Getenv("REFRESH_TOKEN_MAXAGE"), 10, 64)
	smtpPort, _ := strconv.ParseInt(os.Getenv("SMTP_PORT"), 10, 64)
	config = Config{
		DBHost:         os.Getenv("POSTGRES_HOST"),
		DBUserName:     os.Getenv("POSTGRES_USER"),
		DBUserPassword: os.Getenv("POSTGRES_PASSWORD"),
		DBName:         os.Getenv("POSTGRES_DB"),
		DBPort:         os.Getenv("POSTGRES_PORT"),
		ServerPort:     os.Getenv("PORT"),

		ClientOrigin: os.Getenv("CLIENT_ORIGIN"),
		GinMode:      os.Getenv("GIN_MODE"),

		AccessTokenPrivateKey:  os.Getenv("ACCESS_TOKEN_PRIVATE_KEY"),
		AccessTokenPublicKey:   os.Getenv("ACCESS_TOKEN_PUBLIC_KEY"),
		RefreshTokenPrivateKey: os.Getenv("REFRESH_TOKEN_PRIVATE_KEY"),
		RefreshTokenPublicKey:  os.Getenv("REFRESH_TOKEN_PUBLIC_KEY"),

		EmailFrom: os.Getenv("EMAIL_FROM"),
		SMTPHost:  os.Getenv("SMTP_HOST"),
		SMTPPass:  os.Getenv("SMTP_PASS"),
		SMTPUser:  os.Getenv("SMTP_USER"),
	}
	config.AccessTokenExpiresIn = accessExpiresIn
	config.RefreshTokenExpiresIn = refreshExpiresIn
	config.AccessTokenMaxAge = int(accessMaxAge)
	config.RefreshTokenMaxAge = int(refreshMaxAge)
	config.SMTPPort = int(smtpPort)
	return
}

func (c *Config) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.DBHost, c.DBPort, c.DBUserName, c.DBUserPassword, c.DBName)
}
