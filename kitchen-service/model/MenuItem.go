package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MenuItem struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name" gorm:"not null"`
	RestaurantID   string    `json:"restaurant"`
	TicketLineItem []TicketLineItem
}

func (menuItem *MenuItem) BeforeCreate(scope *gorm.DB) error {
	menuItem.ID = uuid.New()
	return nil
}
