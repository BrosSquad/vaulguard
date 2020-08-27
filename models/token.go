package models

import "gorm.io/gorm"

type Token struct {
	gorm.Model
	Value         [64]byte `gorm:"not null"`
	ApplicationId uint   `gorm:"not null"`
	Application   Application `gorm:"foreignKey:ApplicationId"`
}
