package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
)

type Server struct {
	AppName string
	AppPort string
	AppEnv  string

	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
		SSLMode  string
	}

	SecretKey     string
	TokenDuration string
	TimeNow       time.Time
}

func (s *Server) Load() *Server {

	err := godotenv.Load(".env")

	if err != nil {
		log.Info("Error loading .env file")
	}

	s.Database.Host = os.Getenv("PGHOST")
	s.Database.Port = os.Getenv("PGPORT")
	s.Database.User = os.Getenv("PGUSER")
	s.Database.Password = os.Getenv("PGPASSWORD")
	s.Database.DBName = os.Getenv("PGDATABASE")
	s.Database.SSLMode = os.Getenv("PGSSLMODE")

	s.AppName = os.Getenv("APPNAME")
	s.AppPort = os.Getenv("APPPORT")
	s.AppEnv = os.Getenv("APPENV")

	s.SecretKey = os.Getenv("SECRETKEY")
	s.TokenDuration = os.Getenv("TOKENDURATION")
	s.TimeNow = time.Now()

	return s
}
