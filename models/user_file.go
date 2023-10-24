package models

import (
	"mime/multipart"
)

type CreateFile struct {
	Link *string               `json:"-" form:"-"`
	Type int                   `json:"type" form:"type"`
	File *multipart.FileHeader `json:"-" form:"file"`
}

type CreateFiles struct {
	Link *string                 `json:"-" form:"-"`
	Type int                     `json:"type" form:"type"`
	File []*multipart.FileHeader `json:"-" form:"files"`
}

type CreateFilesResponse struct {
	Id        string   `json:"id" bun:"id,pk,autoincrement"`
	Link      []string `json:"link" bun:"link"`
	Type      int      `json:"type" bun:"type"`
	CreatedAt string   `json:"-" bun:"created_at"`
	CreatedBy string   `json:"-" bun:"created_by"`
}

type CreateFileResponse struct {
	Id        string  `json:"id" bun:"id,pk,autoincrement"`
	Link      *string `json:"link" bun:"link"`
	Type      int     `json:"type" bun:"type"`
	CreatedAt string  `json:"-" bun:"created_at"`
	CreatedBy string  `json:"-" bun:"created_by"`
}
