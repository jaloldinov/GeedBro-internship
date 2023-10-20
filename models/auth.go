package models

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRespond struct {
	Token string `json:"token"`
}

type Login struct {
	Login string `json:"login"`
}

type LoginDataRespond struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
