package apihandler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"meshed/meshnode/model/categories"
	"meshed/meshnode/model/users"
	"net/http"
)

const UserLogin = "login"
const UserPassword = "pwd"
const UserEmail = "email"
const UserRegistrationToken = "regtoken"

const CategoryUserUnregistrated = "user-unregistrated"

type registrationSuccess struct {
	message string
	success bool
}

// look if essential categories are present
func init() {
	categories.CreateCategoryIfNotExists(CategoryUserUnregistrated)
}

// Register an user and send out registration link
func RegisterUser(writer http.ResponseWriter, request *http.Request) {
	log.Println("Register user called")
	requestVars := mux.Vars(request)
	log.Println("requestvars", requestVars)
	if userLogin, err := requestVars[UserLogin]; !err {
		log.Println("Could not find any user name from request")
		writeWrongUserRegistration(writer)
	} else if userPassword, err := requestVars[UserPassword]; !err {
		log.Println("Could not find any password from request")
		writeWrongUserRegistration(writer)
	} else if userEmail, err := requestVars[UserEmail]; !err {
		log.Println("Could not find any email from request")
		writeWrongUserRegistration(writer)
	} else {
		userNode := users.NewNodeFromRegistration(userLogin, userEmail, userPassword)
		if unregistrated, found := categories.FindCategoryByName(CategoryUserUnregistrated); found {
			userNode.AddChild(unregistrated)
			// generate token and send email
			writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
			writer.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(writer).Encode(registrationSuccess{
				message: "You are successfully registrated, please check your mail to complete the registration process",
				success: true,
			}); err != nil {
				log.Fatal("Error while encoding respose")
			}
		} else {
			// essential category not found
			writeUnsuccessfulUserRegistration(writer)
		}
	}
}

// Confirms user registration with confirmation token
func ConfirmRegistration(writer http.ResponseWriter, request *http.Request) {
	log.Println("Reset password called")
	writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
	writer.WriteHeader(http.StatusOK)
}

// Checks user account, generates a jwt
func LoginUser(writer http.ResponseWriter, request *http.Request) {
	log.Println("Login user called")
	writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
	writer.WriteHeader(http.StatusOK)
}

// Send out mail to users mail address reset password with generated password-reset-link
func RequestPasswordReset(writer http.ResponseWriter, request *http.Request) {
	log.Println("Reset password request called")
	writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
	writer.WriteHeader(http.StatusOK)
}

// Reset users password
// params
func PasswordReset(writer http.ResponseWriter, request *http.Request) {
	log.Println("Reset password called")
	writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
	writer.WriteHeader(http.StatusOK)
}

// renews an user token if the old one was expired
func RenewUserToken(writer http.ResponseWriter, request *http.Request) {
	log.Println("Renew user token called")
	writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
	writer.WriteHeader(http.StatusOK)
}

func writeWrongUserRegistration(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(writer).Encode(JSONError{Code: http.StatusBadRequest, Text: "Missing user information"}); err != nil {
		log.Fatal("Error while encoding respose")
	}
}

func writeUnsuccessfulUserRegistration(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusFailedDependency)
	if err := json.NewEncoder(writer).Encode(JSONError{Code: http.StatusFailedDependency, Text: "Missing other resources to register this user"}); err != nil {
		log.Fatal("Error while encoding respose")
	}
}

func writeUserNotAuthorized(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusUnauthorized)
	if err := json.NewEncoder(writer).Encode(JSONError{Code: http.StatusUnauthorized, Text: "User or password are wrong"}); err != nil {
		log.Fatal("Error while encoding respose")
	}
}
