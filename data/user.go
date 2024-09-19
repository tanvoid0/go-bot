package data

import (
	"fmt"
	"log"
)

type UserType string

const (
	AppUser UserType = "user"
	BotUser UserType = "bot"
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"type:varchar(100)"`
	Username string `gorm:"unique;type:varchar(100)"`
	Password string `gorm:"type:varchar(100)"`
}

func (user User) Create() (*User, error) {
	db, err := ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	newUser := db.Create(&user)
	if newUser.Error != nil {
		log.Fatalf("Failed to create record: %v", newUser.Error)
	} else {
		fmt.Printf("User %v created", user)
	}
	return &user, nil
}
