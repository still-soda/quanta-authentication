package handlers

import (
	app_error "qauth-server/internal/errors"
	"qauth-server/internal/services"
	"qauth-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	fileService services.FileService
}

func NewFileHandler(fileService *services.FileService) *FileHandler {
	return &FileHandler{fileService: *fileService}
}

func (h *FileHandler) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}
	defer file.Close()

	bucketName := "uploads"
	fileName, err := h.fileService.SaveFile(c, header, nil, &bucketName)

	response.HandlerSuccess(c, gin.H{
		"fileName": fileName,
	})
}
