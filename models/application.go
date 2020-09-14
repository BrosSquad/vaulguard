package models

import (
	"time"
)

type Application struct {
	ID        uint   `gorm:"primarykey"`
	Name      string `gorm:"not null;uniqueIndex"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Tokens    []Token `gorm:"foreignKey:ApplicationId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type ApplicationDto struct {
	ID        interface{}
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
