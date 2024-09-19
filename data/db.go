package data

import (
	"fmt"
	"github.com/tanvoid0/dev-bot/util"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
)

func SetupDatabase() (*gorm.DB, error) {
	db, err := ConnectDatabase()
	err = db.AutoMigrate(&User{}, &ChatMessage{}, &ChatRoom{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func ConnectDatabase() (*gorm.DB, error) {
	dbPath := util.ReadEnvWithDefault("DATABASE", "database.sqlite")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		file, err := os.Create(dbPath)
		if err != nil {
			return nil, fmt.Errorf("error creating database file: %w", err)
		}
		file.Close()
		fmt.Println("Database file created: ", dbPath)
	} else if err != nil {
		return nil, fmt.Errorf("error opening database file: %w", err)
	}
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		return nil, err
	}
	return db, nil
}
