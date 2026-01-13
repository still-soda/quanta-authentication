package handlers

import (
	"qauth-server/internal/services"
	"qauth-server/internal/utilities"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	fileService services.FileService
}

func NewFileHandler(storage utilities.Storage) *FileHandler {
	fileService := services.NewFileService(storage)
	return &FileHandler{fileService: *fileService}
}

func (h *FileHandler) Upload(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to get file from request"})
		return
	}
	defer file.Close()

}
