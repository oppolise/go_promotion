package database

import (
	"fmt"
	databaseconfig "goPromotion/config/database_config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDatabase() (*gorm.DB, error) {

	config := databaseconfig.DatabaseConfig()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.USERNAME,
		config.PASSWORD,
		config.HOST,
		config.PORT,
		config.DATABASE)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying *sql.DB: %v", err)
	}

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("could not connect to the database: %v", err)
	}
	fmt.Println("connect database succesfully")
	fmt.Println("DB host:", config.HOST)
	fmt.Println("DB name:", config.DATABASE)

	return db, nil

}
