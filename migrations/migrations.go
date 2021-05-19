package migrations

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/Yapcheekian/bank-golang/helpers"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username string
	Email    string
	Password string
}

type Account struct {
	gorm.Model
	Type    string
	Name    string
	Balance uint
	UserID  uint
}

func connectDB() *gorm.DB {
	db, err := gorm.Open("postgres", "host=127.0.0.1 port=5432 user=postgres dbname=bankapp password=mysecretpassword sslmode=disable")
	helpers.HandleErr(err)
	return db
}

func createAccounts() {
	db := connectDB()
	users := [2]User{
		{Username: "Yap", Email: "yap@test.com"},
		{Username: "Kath", Email: "kath@test.com"},
	}

	for i, v := range users {
		generatedPassword := helpers.HashAndSalt([]byte(v.Username))
		user := User{Username: v.Username, Password: generatedPassword, Email: v.Email}
		db.Create(&user)

		account := Account{Type: "Daily Account", Name: string(v.Username + "'s" + " account"), Balance: uint(10000 * int(i)), UserID: user.ID}
		db.Create(&account)
	}

	defer db.Close()
}

func Migrate() {
	db := connectDB()
	db.AutoMigrate(&User{}, &Account{})
	defer db.Close()
	createAccounts()
}
