package migrations

import (
	"github.com/Yapcheekian/bank-golang/helpers"
	"github.com/Yapcheekian/bank-golang/interfaces"
)

func createAccounts() {
	db := helpers.ConnectDB()
	users := [2]interfaces.User{
		{Username: "Yap", Email: "yap@test.com"},
		{Username: "Kath", Email: "kath@test.com"},
	}

	for i, v := range users {
		generatedPassword := helpers.HashAndSalt([]byte(v.Username))
		user := interfaces.User{Username: v.Username, Password: generatedPassword, Email: v.Email}
		db.Create(&user)

		account := interfaces.Account{Type: "Daily Account", Name: string(v.Username + "'s" + " account"), Balance: uint(10000 * int(i)), UserID: user.ID}
		db.Create(&account)
	}

	defer db.Close()
}

func Migrate() {
	db := helpers.ConnectDB()
	db.AutoMigrate(&interfaces.User{}, &interfaces.Account{})
	defer db.Close()
	createAccounts()
}
