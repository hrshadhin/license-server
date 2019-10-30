package models

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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
	dbType := os.Getenv("db_type")
	dbUri := ""

	if dbType == "postgres" {
		dbUri = fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)
	} else if dbType == "mysql" {
		dbUri = fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local",username, password, dbName)
	} else if dbType == "sqlite3" {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		dbUri = pwd+"/license-server.db"
	} else {
		fmt.Println("\nDatabase type not selected!")
		os.Exit(1)
	}

	conn, err := gorm.Open(dbType, dbUri)
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	db.AutoMigrate(&User{}, &Key{}, &KeyAccessLog{})
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
