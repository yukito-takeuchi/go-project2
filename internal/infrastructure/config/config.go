package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DBUser     string
	DBPassword string
	DBHost     string
	DBName     string
	DBPort     string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  .envファイルが見つかりませんでした。")
	}

	DBUser = os.Getenv("DB_USER")
	DBPassword = os.Getenv("DB_PASSWORD")
	DBHost = os.Getenv("DB_HOST")
	DBPort = os.Getenv("DB_PORT")
	DBName = os.Getenv("DB_NAME")
}

// DB接続文字列を返す
func GetDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&loc=Local&sql_mode=TRADITIONAL",
		DBUser, DBPassword, DBHost, DBPort, DBName,
	)
}
