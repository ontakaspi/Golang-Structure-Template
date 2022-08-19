package entity

import "gorm.io/gorm"

type ExampleData struct {
	gorm.Model
	Name    string `gorm:"type:varchar(255);not null"`
	Age     int    `gorm:"type:int;not null"`
	Address string `gorm:"type:varchar(255);not null"`
}
