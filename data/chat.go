package data

import (
	"fmt"
	"log"
)

type ChatType string
type ChatStatus string

const (
	Text  ChatType = "text"
	Audio ChatType = "audio"
	Image ChatType = "image"
	Video ChatType = "video"

	Sent ChatStatus = "sent"
	Seen ChatStatus = "seen"
)

type ChatMessage struct {
	ID         uint       `gorm:"primaryKey"`
	Data       string     `gorm:"type:text;default='text'"`
	Type       ChatType   `gorm:"type:varchar(20)"`
	Status     ChatStatus `gorm:"type:varchar(20)"`
	Sender     uint       `gorm:"not null"` // foreign key for user table
	User       User       `gorm:"foreignKey:Sender;constraint:OnDelete:CASCADE"`
	ChatRoomID uint       // Foreign key for ChatRoom
	ChatRoom   ChatRoom   `gorm:"foreignKey:ChatRoomID;constraint:OnDelete:CASCADE"`
}

type ChatRoom struct {
	ID       uint          `gorm:"primaryKey"`
	Messages []ChatMessage `gorm:"foreignKey:ChatRoomID"` // Define the foreign key
}

func (message ChatMessage) Create() (*ChatMessage, error) {
	db, err := ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	newMessage := db.Create(&message)
	if newMessage.Error != nil {
		log.Fatalf("failed to create record: %v", newMessage.Error)
	} else {
		fmt.Printf("Record created: %+v\n", newMessage)
	}
	return &message, nil

}
