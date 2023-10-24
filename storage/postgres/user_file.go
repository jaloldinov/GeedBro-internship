package postgres

import (
	"auth/api/response"
	"auth/models"
	"auth/pkg/helper"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lib/pq"
)

type fileRepo struct {
	db          *pgxpool.Pool
	fileService *helper.Service
}

func NewFileRepo(db *pgxpool.Pool) *fileRepo {
	return &fileRepo{
		db:          db,
		fileService: helper.NewService(),
	}
}

func (r *fileRepo) CreateFile(ctx context.Context, req *models.CreateFile) (*models.CreateFileResponse, *response.ErrorResp) {
	var detail models.CreateFileResponse
	userInfo := ctx.Value("user_info").(helper.TokenInfo)

	link, err := r.fileService.Upload(ctx, req.File, "media")
	if err != nil {
		return &models.CreateFileResponse{}, err
	}
	fmt.Println(link)

	query := `		
		INSERT INTO "medias" (
			id,
			link,
			type,
			created_at,
			created_by
			)
			VALUES (
			$1, 
			$2,
			$3,
			NOW(),
			$4
			)`

	detail.Id = uuid.NewString()
	detail.Type = req.Type
	detail.Link = &link
	detail.CreatedAt = time.Now().String()
	detail.CreatedBy = userInfo.User_id

	_, er := r.db.Exec(context.Background(), query,
		detail.Id,
		link,
		detail.Type,
		userInfo.User_id,
	)
	if er != nil {
		return &models.CreateFileResponse{}, &response.ErrorResp{Message: er.Error()}
	}

	return &detail, nil
}

func (r *fileRepo) CreateFiles(ctx context.Context, req *models.CreateFiles) (*models.CreateFilesResponse, *response.ErrorResp) {
	var detail models.CreateFilesResponse
	userInfo := ctx.Value("user_info").(helper.TokenInfo)

	links, err := r.fileService.MultipleUpload(ctx, req.File, "user/")
	if err != nil {
		return &models.CreateFilesResponse{}, err
	}

	query := `		
		INSERT INTO "medias" (
			id,
			link,
			type,
			created_at,
			created_by
		)
		VALUES (
			$1, $2, $3, NOW(), $4
		)`

	detail.Id = uuid.NewString()
	detail.Type = req.Type
	detail.Link = links // Assign the array directly

	detail.CreatedAt = time.Now().String()
	detail.CreatedBy = userInfo.User_id

	_, er := r.db.Exec(context.Background(), query,
		detail.Id,
		pq.Array(detail.Link), // Convert array to pq.Array
		detail.Type,
		userInfo.User_id,
	)
	if er != nil {
		return &models.CreateFilesResponse{}, &response.ErrorResp{Message: er.Error()}
	}

	return &detail, nil
}
