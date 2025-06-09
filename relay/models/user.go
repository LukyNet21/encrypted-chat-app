package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"primarykey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	UserName  string `json:"username" gorm:"uniqueIndex"`
	PublicKey string `json:"public" `

	SentMessages     []Message `json:"-" gorm:"foreignKey:SenderID"`
	RecievedMessages []Message `json:"-" gorm:"foreignKey:ReceiverID"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}
