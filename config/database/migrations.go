package database

import "golang-example/app/models/entity"

func Migrate() {
	db := PostgreDB
	err := db.AutoMigrate(&entity.ExampleData{})
	if err != nil {
		return
	}
}
