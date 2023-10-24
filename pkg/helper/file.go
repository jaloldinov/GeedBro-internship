package helper

import (
	"auth/api/response"
	"context"
	"errors"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s Service) Upload(ctx context.Context, file *multipart.FileHeader, folder string) (string, *response.ErrorResp) {
	if file == nil {
		return "", nil
	}

	//fixedFile := strings.Split(file.Filename, ".")
	//
	var chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321"
	//str := make([]rune, 30)
	//for i := range str {
	//	str[i] = chars[rand.Intn(len(chars))]
	//}

	//filename := filepath.Base(string(str) + "." + fixedFile[len(fixedFile)-1])
	i := 0
	randName := ""
	for i < 30 {
		var randSpeed = rand.New(rand.NewSource(time.Now().UnixNano()))
		randIdx := randSpeed.Intn(62)

		randName += string(chars[randIdx])

		i++
	}

	contentTypes := map[string]bool{
		"application/msword": true,
		"image/jpeg":         true,
		"image/jpg":          true,
		"image/png":          true,
		"application/pdf":    true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	}

	if len(file.Header.Values("Content-Type")) > 0 && !contentTypes[file.Header.Values("Content-Type")[0]] {
		return "", &response.ErrorResp{
			Message: "content-type of this file has not permission to upload into the server!",
			Code:    http.StatusInternalServerError,
		}
	}

	//filename := filepath.Base(strings.Join(strings.Split(file.Filename, " "), "-"))

	splitFileName := strings.Split(file.Filename, ".")

	//filename := filepath.Base(strconv.Itoa(int(time.Now().UnixMicro())) + "." + splitFileName[len(splitFileName)-1])
	filename := filepath.Base(randName + "." + splitFileName[len(splitFileName)-1])

	if _, err := os.Stat("./media/" + folder); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll("./media/"+folder, os.ModePerm)
		if err != nil {
			return "", &response.ErrorResp{
				Message: "file upload mkdir",
				Code:    http.StatusInternalServerError,
			}
		}
	}

	files, err := os.ReadDir("./media/" + folder)

	if err != nil {
		return "", &response.ErrorResp{
			Message: "file upload read",
			Code:    http.StatusInternalServerError,
		}
	}

	for _, f := range files {
		if !f.IsDir() && (f.Name() == filename) {
			splitString := strings.Split(filename, ".")
			extra := strconv.Itoa(int(time.Now().Unix()))
			splitString[len(splitString)-2] = splitString[len(splitString)-2] + "-" + extra
			filename = strings.Join(splitString, ".")
			break
		}
	}

	dst := "./media/" + folder + filename

	src, err := file.Open()
	if err != nil {
		return "", &response.ErrorResp{
			Message: "file upload open",
			Code:    http.StatusInternalServerError,
		}
	}
	defer log.Println("file upload src.Close() error: ", src.Close())

	out, err := os.Create(dst)
	if err != nil {
		return "", &response.ErrorResp{
			Message: "file upload create",
			Code:    http.StatusInternalServerError,
		}
	}
	//defer log.Println("file upload out.Close() error: ", out.Close())

	_, err = io.Copy(out, src)

	if err != nil {
		return "", &response.ErrorResp{
			Message: "file upload copy",
			Code:    http.StatusInternalServerError,
		}
	}

	return "/media/" + folder + filename, nil
}

func (s Service) Delete(ctx context.Context, url string) *response.ErrorResp {
	err := os.RemoveAll("." + url)
	if err != nil {
		return &response.ErrorResp{
			Message: "file  delete",
			Code:    http.StatusInternalServerError,
		}
	}

	return nil
}

func (s Service) MultipleUpload(ctx context.Context, files []*multipart.FileHeader, folder string) ([]string, *response.ErrorResp) {
	var links []string

	for _, f := range files {
		link, err := s.Upload(ctx, f, folder)

		if err != nil {
			return nil, err
		}

		links = append(links, link)
	}

	return links, nil
}
