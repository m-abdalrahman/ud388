package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DB = dbSetup()

func dbSetup() *gorm.DB {

	db, err := gorm.Open("sqlite3", "puppies.db")
	if err != nil {
		panic(err.Error())
	}

	if db.HasTable(&Puppy{}) == false {
		db.CreateTable(&Puppy{})
	}

	return db
}

type Puppy struct {
	ID          uint   `gorm:"primary_key"`
	Name        string `gorm:"size:80;not null"`
	Description string `gorm:"size:250"`
}

func (p Puppy) TableName() string {
	return "puppy"
}
