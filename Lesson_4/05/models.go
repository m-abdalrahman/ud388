package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"golang.org/x/crypto/bcrypt"
)

var DB = dbSetup()

func dbSetup() *gorm.DB {

	db, err := gorm.Open("sqlite3", "bagelShop.db")
	if err != nil {
		panic(err.Error())
	}

	if db.HasTable(&User{}) == false || db.HasTable(&Bagel{}) == false {
		db.CreateTable(&User{})
		db.CreateTable(&Bagel{})
	}

	return db
}

type User struct {
	ID           uint   `gorm:"primary_key"`
	Username     string `gorm:"size:32;index" json:"username"`
	PasswordHash string `gorm:"size:64" json:"password"`
}

func (u User) TableName() string {
	return "user"
}

func (u *User) HashPassword(password string) error {
	//encrypt password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPass)

	return nil
}

func (u *User) VerifyPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return err
	}

	return nil
}

type Bagel struct {
	ID          uint   `gorm:"primary_key" json:"-"`
	Name        string `json:"name"`
	Picture     string `json:"picture"`
	Description string `json:"description"`
	Price       string `json:"price"`
}

func (b Bagel) TableName() string {
	return "bagel"
}
