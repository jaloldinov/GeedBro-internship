package models

type CreateLike struct {
	UserId string `json:"user_id"`
	PostId string `json:"post_id"`
}

type Like struct {
	ID        string `json:"id"`
	UserId    string `json:"user_id"`
	PostId    string `json:"post_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

type DeleteLike struct {
	UserId string `json:"user_id"`
	PostId string `json:"post_id"`
}
