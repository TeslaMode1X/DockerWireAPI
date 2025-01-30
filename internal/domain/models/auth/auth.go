package model

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Registration struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
