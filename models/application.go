package models

import "gorm.io/gorm"

type Application struct {
	gorm.Model
	Name   string  `gorm:"not null"`
	Tokens []Token `gorm:"foreignKey:ApplicationId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
