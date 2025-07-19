package config

import (
	"fmt"
	"log"
	"os"

	"github.com/Claudio712005/go-microservices-architecture/auth-service/internal/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConectarBanco() {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName)

	log.Println("Conectando ao banco de dados:", dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados", err)
	}

	if os.Getenv("APP_ENV") == "test" {
		if err := db.AutoMigrate(&domain.Usuario{}); err != nil {
			log.Fatal("Erro ao migrar o banco de dados", err)
		}
	}

	DB = db
}

// ResetTestDB limpa o banco de dados de testes, removendo todos os dados
func ResetTestDB() {
    if DB == nil {
        return
    }

    DB.Exec("SET FOREIGN_KEY_CHECKS = 0")

    var tables []string
    DB.Raw("SHOW TABLES").Scan(&tables)

    for _, tbl := range tables {
        DB.Exec("TRUNCATE TABLE " + tbl)
    }

    DB.Exec("SET FOREIGN_KEY_CHECKS = 1")
}