package model

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Registration struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
