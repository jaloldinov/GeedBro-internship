package models

type CreateUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Is_active bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

type UpdateUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type IdRequest struct {
	Id string `json:"id"`
}

type GetAllUserRequest struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Search string `json:"username"`
}

type GetAllUser struct {
	Users []User `json:"Users"`
	Count int    `json:"count"`
}
