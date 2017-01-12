package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	DB = dbSetup()
)

func dbSetup() *gorm.DB {
	db, err := gorm.Open("sqlite3", "bargainMart.db")
	if err != nil {
		panic(err.Error())
	}

	if db.HasTable(&Item{}) == false {
		db.CreateTable(&Item{})
	}

	return db
}

type Item struct {
	ID          uint   `gorm:"primary_key" json:"id"`
	Name        string `json:"name"`
	Picture     string `json:"picture"`
	Price       string `json:"price"`
	Description string `json:"description"`
}

func (i Item) TableName() string {
	return "item"
}
