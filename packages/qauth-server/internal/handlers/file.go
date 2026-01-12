package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"qauth-server/internal/database"
	"qauth-server/internal/models"
	"qauth-server/internal/storage"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// FileHandler 文件处理器
type FileHandler struct {
	storage storage.Service
}

// NewFileHandler 创建文件处理器
func NewFileHandler(storage storage.Service) *FileHandler {
	return &FileHandler{
		storage: storage,
	}
}

// Upload 上传文件
func (h *FileHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件上传失败"})
		return
	}

	ext := filepath.Ext(file.Filename)
	fileKey := fmt.Sprintf("%s/%s%s", time.Now().Format("2006/01/02"), uuid.New().String(), ext)

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法读取文件"})
		return
	}
	defer src.Close()

	if err := h.storage.Upload(c.Request.Context(), fileKey, src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件上传失败"})
		return
	}

	fileURL, _ := h.storage.GetURL(c.Request.Context(), fileKey)

	fileRecord := &models.File{
		UserID:      1,
		FileName:    file.Filename,
		FileKey:     fileKey,
		FileSize:    file.Size,
		ContentType: file.Header.Get("Content-Type"),
		URL:         fileURL,
	}

	if err := database.GetDB().Create(fileRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "文件上传成功",
		"file":    fileRecord,
	})
}

// Download 下载文件
func (h *FileHandler) Download(c *gin.Context) {
	fileKey := c.Param("key")

	data, err := h.storage.Download(c.Request.Context(), fileKey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	c.Data(http.StatusOK, "application/octet-stream", data)
}

// Delete 删除文件
func (h *FileHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	var file models.File
	if err := database.GetDB().First(&file, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	if err := h.storage.Delete(c.Request.Context(), file.FileKey); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文件失败"})
		return
	}

	if err := database.GetDB().Delete(&file).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文件删除成功"})
}

// List 文件列表
func (h *FileHandler) List(c *gin.Context) {
	var files []models.File

	if err := database.GetDB().Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"files": files,
		"total": len(files),
	})
}
