package models

type CreatePost struct {
	Description string   `json:"description"`
	Photos      []string `json:"photos" form:"photos"`
}

type Post struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Photos      []string `json:"photos"`
	LikeCount   int      `json:"likes_count"`
	CreatedAt   string   `json:"created_at"`
}

type DeletePost struct {
	Id string `json:"id"`
}

type UpdatePost struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Photos      []string `json:"photos"`
}

type GetAllPostRequest struct {
	Page   *int    `json:"page"`
	Limit  *int    `json:"limit"`
	Search *string `json:"description"`
}

type GetAllMyPostRequest struct {
	Page   *int    `json:"page"`
	Limit  *int    `json:"limit"`
	Search *string `json:"description"`
}

type GetAllPost struct {
	Posts []Post `json:"Posts"`
	Count int    `json:"count"`
}
