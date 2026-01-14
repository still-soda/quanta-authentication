package services

import (
	"mime/multipart"
	"net/http"
	"qauth-server/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
	storageService *StorageService
	db             *gorm.DB
}

func NewFileService(storageService *StorageService, db *gorm.DB) *FileService {
	return &FileService{storageService: storageService, db: db}
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

	if err := s.storageService.Upload(ctx, uniqueName, file); err != nil {
		return "", err
	}

	fileRecord := &models.Files{
		StorageKey: uniqueName,
		Bucket:     *bucketName,
		MimeType:   mimeType,
		SizeBytes:  header.Size,
		CreatorID:  creatorID,
	}

	if err := s.db.Create(fileRecord).Error; err != nil {
		return "", err
	}

	return uniqueName, nil
}
