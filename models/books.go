package models

import "gorm.io/gorm"

type Books struct {
	ID      uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Yazar   *string `json:"yazar"`
	Adi     *string `json:"adi"`
	Yayinci *string `json:"yayinci"`
}

func MigrateBooks(db *gorm.DB) error {
	err := db.AutoMigrate(&Books{})
	return err
}
