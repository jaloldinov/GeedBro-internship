package models

type CreateCommentLike struct {
	CommentId string `json:"comment_id"`
}

type CommentLike struct {
	Id        string `json:"id"`
	CommentId string `json:"post_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

type DeleteCommentLike struct {
	CommentId *string `json:"comment_id"`
}
