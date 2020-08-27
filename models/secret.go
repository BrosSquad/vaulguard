package models

type Secret struct {
	Key   string `gorm:"not null"`
	Value string `gorm:"not null"`
}
