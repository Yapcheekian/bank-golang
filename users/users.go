package users

import (
	"time"

	"github.com/Yapcheekian/bank-golang/helpers"
	"github.com/Yapcheekian/bank-golang/interfaces"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func prepareToken(user *interfaces.User) string {
	tokenContent := jwt.MapClaims{
		"user_id": user.ID,
		"expiry":  time.Now().Add(time.Minute ^ 60).Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte("TokenPassword"))
	helpers.HandleErr(err)

	return token
}

func prepareResponse(user *interfaces.User, accounts []interfaces.ResponseAccount) map[string]interface{} {
	responseUser := &interfaces.ResponseUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Accounts: accounts,
	}

	var token = prepareToken(user)
	var response = map[string]interface{}{"message": "all is fine"}
	response["jwt"] = token
	response["data"] = responseUser

	return response
}

func Login(username string, password string) map[string]interface{} {
	valid := helpers.Validate(
		[]interfaces.Validation{
			{Value: username, Valid: "username"},
			{Value: password, Valid: "password"},
		})

	if valid {
		db := helpers.ConnectDB()
		user := &interfaces.User{}

		if db.Where("username = ?", username).First(&user).RecordNotFound() {
			return map[string]interface{}{"message": "user not found"}
		}

		passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

		if passErr == bcrypt.ErrMismatchedHashAndPassword && passErr != nil {
			return map[string]interface{}{"message": "wrong password"}
		}

		accounts := []interfaces.ResponseAccount{}
		db.Table("accounts").Select("id, name, balance").Where("user_id = ?", user.ID).Scan(&accounts)

		defer db.Close()

		var response = prepareResponse(user, accounts)

		return response
	} else {
		return map[string]interface{}{"message": "not valid values"}
	}
}

func Register(username string, email string, password string) map[string]interface{} {
	valid := helpers.Validate(
		[]interfaces.Validation{
			{Value: username, Valid: "username"},
			{Value: password, Valid: "password"},
			{Value: email, Valid: "email"},
		})
	if valid {
		db := helpers.ConnectDB()
		defer db.Close()

		generatedPassword := helpers.HashAndSalt([]byte(password))
		user := interfaces.User{Username: username, Password: generatedPassword, Email: email}
		db.Create(&user)

		account := interfaces.Account{Type: "Daily Account", Name: string(username + "'s" + " account"), Balance: 0, UserID: user.ID}
		db.Create(&account)

		accounts := []interfaces.ResponseAccount{}
		respAccount := interfaces.ResponseAccount{
			ID:      account.ID,
			Name:    account.Name,
			Balance: int(account.Balance),
		}
		accounts = append(accounts, respAccount)
		response := prepareResponse(&user, accounts)

		return response
	} else {
		return map[string]interface{}{"message": "not valid values"}
	}
}
