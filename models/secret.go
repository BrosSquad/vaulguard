package models

type Secret struct {
	ID            uint   `gorm:"primaryKey"`
	Key           string `gorm:"uniqueIndex:application_id_key_idx;not null;"`
	ApplicationId uint   `gorm:"not null;uniqueIndex:application_id_key_idx;"`
	Value         []byte `gorm:"not null;"`
}
