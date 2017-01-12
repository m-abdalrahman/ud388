package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DB = dbSetup()

func dbSetup() *gorm.DB {

	db, err := gorm.Open("sqlite3", "restaruants.db")
	if err != nil {
		panic(err.Error())
	}

	if db.HasTable(&Restaurant{}) == false {
		db.CreateTable(&Restaurant{})
	}

	return db
}

type Restaurant struct {
	ID      uint   `gorm:"primary_key" json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Image   string `json:"image"`
}

func (r Restaurant) TableName() string {
	return "restaurant"
}
