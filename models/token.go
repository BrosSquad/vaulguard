package models

import "gorm.io/gorm"

type Token struct {
	gorm.Model
	Value         []byte `gorm:"not null"`
	ApplicationId uint   `gorm:"not null"`
	Application   Application `gorm:"foreignKey:ApplicationId"`
}
