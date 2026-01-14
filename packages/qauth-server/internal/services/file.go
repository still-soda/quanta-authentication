package services

import (
	"mime/multipart"
	"net/http"
	"qauth-server/internal/database"
	"qauth-server/internal/models"
	"qauth-server/internal/utilities"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func _getMimeType(file multipart.File) string {
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return "application/octet-stream"
	}
	mimeType := http.DetectContentType(buffer)
	if _, err := file.Seek(0, 0); err != nil {
		return "application/octet-stream"
	}
	return mimeType
}

type FileService struct {
	storage utilities.Storage
}

func NewFileService(storage utilities.Storage) *FileService {
	return &FileService{storage: storage}
}

func (s *FileService) SaveFile(
	ctx *gin.Context,
	header *multipart.FileHeader,
	creatorID *string,
	bucketName *string,
) (string, error) {
	file, err := header.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	uniqueName := uuid.New().String() + "_" + header.Filename

	mimeType := _getMimeType(file)

	if err := s.storage.Upload(ctx, uniqueName, file); err != nil {
		return "", err
	}

	db := database.GetDB(ctx)
	fileRecord := &models.Files{
		StorageKey: uniqueName,
		Bucket:     *bucketName,
		MimeType:   mimeType,
		SizeBytes:  header.Size,
		CreatorID:  creatorID,
	}

	if err := db.Create(fileRecord).Error; err != nil {
		return "", err
	}

	return uniqueName, nil
}
