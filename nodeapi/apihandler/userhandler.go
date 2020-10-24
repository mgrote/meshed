package apihandler

import (
	"log"
	"net/http"
)

// Register an user
func RegisterUser(writer http.ResponseWriter, request *http.Request) {
	log.Println("Register user called")
	writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
	writer.WriteHeader(http.StatusOK)
}

// Checks user account, generates a jwt
func LoginUser(writer http.ResponseWriter, request *http.Request) {
	log.Println("Login user called")
	writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
	writer.WriteHeader(http.StatusOK)
}

// renews an user token if the old one was expired
func RenewUserToken(writer http.ResponseWriter, request *http.Request) {
	log.Println("Renew user token called")
	writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
	writer.WriteHeader(http.StatusOK)
}
