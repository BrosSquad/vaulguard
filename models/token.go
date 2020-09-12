package models

import (
	"time"
)

type Token struct {
	ID            uint        `gorm:"primarykey"`
	Value         []byte      `gorm:"not null"`
	ApplicationId uint        `gorm:"not null"`
	Application   Application `gorm:"foreignKey:ApplicationId"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type TokenDto struct {
	ID            interface{}
	Value         []byte
	ApplicationId interface{}
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Application   ApplicationDto
}
