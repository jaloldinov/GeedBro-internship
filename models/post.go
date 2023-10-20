package models

type CreatePost struct {
	Description string   `json:"description"`
	Photos      []string `json:"photos"`
	CreatedBy   string   `json:"created_by"`
}

type Post struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Photos      []string `json:"photos"`
	CreatedAt   string   `json:"created_at"`
	CreatedBy   string   `json:"created_by"`
	UpdatedAt   string   `json:"updated_at"`
	UpdatedBy   string   `json:"updated_by"`
	DeletedAt   string   `json:"deleted_at"`
	DeletedBy   string   `json:"deleted_by"`
}

type DeletePost struct {
	Id     string `json:"id"`
	UserId string `json:"user_id"`
}

type UpdatePost struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Photos      []string `json:"photos"`
	CreatedBy   string   `json:"created_by"`
}

type GetAllPostRequest struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Search string `json:"username"`
}

type GetAllPost struct {
	Posts []Post `json:"Posts"`
	Count int    `json:"count"`
}
