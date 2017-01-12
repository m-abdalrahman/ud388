package main

import (
	"time"

	"github.com/dchest/uniuri"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"golang.org/x/crypto/bcrypt"
)

var (
	DB        = dbSetup()
	secretKey = uniuri.NewLen(32)
)

func dbSetup() *gorm.DB {

	db, err := gorm.Open("sqlite3", "paleKale.db")
	if err != nil {
		panic(err.Error())
	}

	if db.HasTable(&User{}) == false {
		db.CreateTable(&User{})
	}

	return db
}

type User struct {
	ID           uint   `gorm:"primary_key"`
	Username     string `gorm:"size:32;index" json:"username"`
	Picture      string
	Email        string
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

type MyCustomClaims struct {
	ID uint `json:"id"`
	jwt.StandardClaims
}

func (u *User) GenerateAuthToken(expiration time.Duration) (string, error) {
	claims := MyCustomClaims{
		u.ID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * expiration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (u *User) VerifyAuthToken(tokenString string) (uint, error) {
	var userID uint
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return 0, err
	} else {
		if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
			userID = claims.ID
		}
	}

	return userID, nil
}
