package datasource

import (
	"fmt"
	"log"
	"quizer/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() *gorm.DB {
	if DB != nil {
		return DB
	}

	cfg := config.LoadConfig()

	host := cfg.Database.Host
	user := cfg.Database.User
	password := cfg.Database.Password
	dbname := cfg.Database.DBName
	port := cfg.Database.Port
	sslmode := cfg.Database.SSLMode

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Jakarta",
		host, user, password, dbname, port, sslmode)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	sqlDb, err := database.DB()
	if err != nil {
		log.Fatal("Failed to get database connection:", err)
	}

	log.Println("Database connected successfully")

	sqlDb.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDb.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDb.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	DB = database
	return database
}

func GetDB() *gorm.DB {
	return DB
}

func CloseDB() error {
	if DB == nil {
		return nil
	}
	sqlDb, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	if err := sqlDb.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	log.Println("Database connection closed successfully")
	DB = nil
	return nil
}

func CheckDB() bool {
	if DB == nil {
		return false
	}
	sqlDb, err := DB.DB()
	if err != nil {
		return false
	}
	return sqlDb.Ping() == nil
}
