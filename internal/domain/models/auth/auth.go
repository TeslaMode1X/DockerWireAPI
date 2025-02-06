package model

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
} // @name LoginModel

type Registration struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
} // @name RegistrationModel
