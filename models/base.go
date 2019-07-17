package models

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var db *gorm.DB

func init() {

	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)
	// fmt.Println(dbUri)

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	db.Debug().AutoMigrate(&User{}, &Key{})
	SeedFirstUser()
}

func GetDB() *gorm.DB {
	return db
}

func SeedFirstUser() {
	//check if user exists
	temp := &User{}
	err := GetDB().Table("users").Where("email = ?", "dev@hrshadhin.me").First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Println("Connection error. Please retry")
	}
	if temp.Email == "" {
		fmt.Println("Seeding defalult user: dev@hrshadhin.me dev#321")
		//create a default user
		user := User{Name: "H.R. Shadhin", Email: "dev@hrshadhin.me", Password: "dev#321"}
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		user.Password = string(hashedPassword)
		GetDB().Create(&user)
	}

}
