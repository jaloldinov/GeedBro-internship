package models

// CREATE TABLE "post_comments" (
// 	"id" varchar(36) PRIMARY KEY,
// 	"post_id" varchar(36) NOT NULL NOT NULL REFERENCES "post"("id"),
// 	"comment" varchar,
// 	"created_at" timestamp NOT NULL DEFAULT NOW(),
// 	"created_by" varchar(36) NOT NULL REFERENCES "users" ("id"),
// 	"updated_at" timestamp,
// 	"updated_by" varchar(36) REFERENCES "users" ("id"),
// 	"deleted_at" timestamp,
// 	"deleted_by" varchar(36) REFERENCES "users" ("id"),
// 	FOREIGN KEY ("post_id") REFERENCES "post"("id") ON DELETE CASCADE,
// 	FOREIGN KEY ("created_by") REFERENCES "users"("id") ON DELETE CASCADE
//   );

type CreateComment struct {
	PostId  string `json:"post_id"`
	Comment string `json:"comment"`
}

type Comment struct {
	ID        string `json:"id"`
	PostId    string `json:"post_id"`
	Comment   string `json:"comment"`
	LikeCount int    `json:"likes_count"`
	CreatedAt string `json:"created_at"`
}

type DeleteComment struct {
	Id string `json:"id"`
}

type UpdateComment struct {
	ID      string `json:"id"`
	PostId  string `json:"post_id"`
	Comment string `json:"comment"`
}

type GetAllPostComments struct {
	Page   *int    `json:"page"`
	Limit  *int    `json:"limit"`
	PostId *string `json:"post_id"`
}

type GetAllCommentResponse struct {
	Comments []Comment `json:"comments"`
	Count    int       `json:"count"`
}
