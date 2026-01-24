package services

import (
	"encoding/json"
	"qauth-server/internal/models"
	"qauth-server/internal/utilities"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// AuditService 审计日志服务
type AuditService struct {
	db          *gorm.DB
	userService *UserService
}

// NewAuditService 创建审计日志服务
func NewAuditService(db *gorm.DB, userService *UserService) *AuditService {
	return &AuditService{
		db:          db,
		userService: userService,
	}
}

// AuditContext 审计上下文信息
type AuditContext struct {
	OperatorID   string
	OperatorName string
	IP           string
	UserAgent    string
	ClientID     *string
	SessionID    string
	RequestID    string
}

// AuditEntry 审计日志条目
type AuditEntry struct {
	Module       models.AuditModule
	Action       models.AuditAction
	TargetID     string
	TargetType   string
	TargetName   string
	Detail       *models.AuditLogDetail
	Status       models.AuditStatus
	ErrorMessage string
	DurationMs   int64
}

// ExtractAuditContext 从 Gin 上下文中提取审计信息
func (s *AuditService) ExtractAuditContext(c *gin.Context) *AuditContext {
	ctx := &AuditContext{
		IP:        c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		RequestID: c.GetHeader("X-Request-ID"),
		SessionID: c.GetHeader("X-Session-ID"),
	}

	// 尝试从上下文中获取用户信息
	if userID, exists := c.Get("userID"); exists {
		ctx.OperatorID = userID.(string)

		// 获取用户名
		if user, err := s.userService.GetUserByID(ctx.OperatorID, false); err == nil {
			ctx.OperatorName = user.Name
		}
	}

	return ctx
}

// Log 记录审计日志
func (s *AuditService) Log(ctx *AuditContext, entry *AuditEntry) error {
	logger := utilities.GetLogger()

	// 构建详情 JSON
	var detailJSON datatypes.JSON
	if entry.Detail != nil {
		detailBytes, err := json.Marshal(entry.Detail)
		if err != nil {
			logger.Error("Failed to marshal audit detail", "error", err)
			detailJSON = datatypes.JSON([]byte("{}"))
		} else {
			detailJSON = datatypes.JSON(detailBytes)
		}
	} else {
		detailJSON = datatypes.JSON([]byte("{}"))
	}

	// 设置默认状态
	if entry.Status == "" {
		entry.Status = models.AuditStatusSuccess
	}

	auditLog := &models.AuditLog{
		OperatorID:   ctx.OperatorID,
		OperatorName: ctx.OperatorName,
		Module:       entry.Module,
		Action:       entry.Action,
		TargetID:     entry.TargetID,
		TargetType:   entry.TargetType,
		TargetName:   entry.TargetName,
		Detail:       detailJSON,
		IP:           ctx.IP,
		UserAgent:    ctx.UserAgent,
		Status:       entry.Status,
		ErrorMessage: entry.ErrorMessage,
		DurationMs:   entry.DurationMs,
		ClientID:     ctx.ClientID,
		SessionID:    ctx.SessionID,
		RequestID:    ctx.RequestID,
	}

	if err := s.db.Create(auditLog).Error; err != nil {
		logger.Error("Failed to create audit log", "error", err)
		return err
	}

	logger.Info("Audit log created",
		"module", entry.Module,
		"action", entry.Action,
		"operator", ctx.OperatorID,
		"target", entry.TargetID,
		"status", entry.Status,
	)

	return nil
}

// LogWithGinContext 使用 Gin 上下文直接记录审计日志
func (s *AuditService) LogWithGinContext(c *gin.Context, entry *AuditEntry) error {
	ctx := s.ExtractAuditContext(c)
	return s.Log(ctx, entry)
}

// LogLogin 记录登录审计日志
func (s *AuditService) LogLogin(c *gin.Context, userID, userName, loginType string, success bool, failReason string, durationMs int64) error {
	status := models.AuditStatusSuccess
	if !success {
		status = models.AuditStatusError
	}

	ctx := &AuditContext{
		OperatorID:   userID,
		OperatorName: userName,
		IP:           c.ClientIP(),
		UserAgent:    c.GetHeader("User-Agent"),
		RequestID:    c.GetHeader("X-Request-ID"),
	}

	return s.Log(ctx, &AuditEntry{
		Module:     models.AuditModuleAuth,
		Action:     models.AuditActionLogin,
		TargetID:   userID,
		TargetType: "user",
		TargetName: userName,
		Detail: &models.AuditLogDetail{
			LoginType:  loginType,
			FailReason: failReason,
		},
		Status:       status,
		ErrorMessage: failReason,
		DurationMs:   durationMs,
	})
}

// LogOAuthAuthorize 记录 OAuth 授权审计日志
func (s *AuditService) LogOAuthAuthorize(c *gin.Context, userID, userName string, clientID *string, scopes []string, grantType, responseType, redirectURI string, success bool, failReason string) error {
	status := models.AuditStatusSuccess
	if !success {
		status = models.AuditStatusError
	}

	clientIDStr := ""
	if clientID != nil {
		clientIDStr = *clientID
	}

	ctx := &AuditContext{
		OperatorID:   userID,
		OperatorName: userName,
		IP:           c.ClientIP(),
		UserAgent:    c.GetHeader("User-Agent"),
		ClientID:     clientID,
		RequestID:    c.GetHeader("X-Request-ID"),
	}

	return s.Log(ctx, &AuditEntry{
		Module:     models.AuditModuleOAuth,
		Action:     models.AuditActionOAuthAuthorize,
		TargetID:   clientIDStr,
		TargetType: "oauth_client",
		Detail: &models.AuditLogDetail{
			Scopes:       scopes,
			GrantType:    grantType,
			ResponseType: responseType,
			RedirectURI:  redirectURI,
			FailReason:   failReason,
		},
		Status:       status,
		ErrorMessage: failReason,
	})
}

// QueryAuditLogs 查询审计日志
func (s *AuditService) QueryAuditLogs(query *models.AuditLogQuery) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64

	db := s.db.Model(&models.AuditLog{})

	// 构建查询条件
	if query.OperatorID != "" {
		db = db.Where("operator_id = ?", query.OperatorID)
	}
	if query.Module != "" {
		db = db.Where("module = ?", query.Module)
	}
	if query.Action != "" {
		db = db.Where("action = ?", query.Action)
	}
	if query.TargetID != "" {
		db = db.Where("target_id = ?", query.TargetID)
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.ClientID != "" {
		db = db.Where("client_id = ?", query.ClientID)
	}
	if query.StartTime != nil {
		db = db.Where("created_at >= ?", query.StartTime)
	}
	if query.EndTime != nil {
		db = db.Where("created_at <= ?", query.EndTime)
	}

	// 计算总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}
	if query.PageSize > 100 {
		query.PageSize = 100
	}

	offset := (query.Page - 1) * query.PageSize
	if err := db.Preload("Operator").
		Order("created_at DESC").
		Offset(offset).
		Limit(query.PageSize).
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetAuditLogByID 根据 ID 获取审计日志
func (s *AuditService) GetAuditLogByID(id string) (*models.AuditLog, error) {
	var log models.AuditLog
	if err := s.db.Preload("Operator").Preload("Client").First(&log, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

// GetRecentActivities 获取最近活动列表
func (s *AuditService) GetRecentActivities(limit int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	if err := s.db.Preload("Operator").
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

// GetAuditStatsByModule 按模块统计审计日志
func (s *AuditService) GetAuditStatsByModule(startTime, endTime time.Time) (map[string]int64, error) {
	type Result struct {
		Module string
		Count  int64
	}
	var results []Result

	if err := s.db.Model(&models.AuditLog{}).
		Select("module, COUNT(*) as count").
		Where("created_at >= ? AND created_at <= ?", startTime, endTime).
		Group("module").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, r := range results {
		stats[r.Module] = r.Count
	}
	return stats, nil
}

// GetAuditStatsByAction 按操作类型统计审计日志
func (s *AuditService) GetAuditStatsByAction(startTime, endTime time.Time) (map[string]int64, error) {
	type Result struct {
		Action string
		Count  int64
	}
	var results []Result

	if err := s.db.Model(&models.AuditLog{}).
		Select("action, COUNT(*) as count").
		Where("created_at >= ? AND created_at <= ?", startTime, endTime).
		Group("action").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, r := range results {
		stats[r.Action] = r.Count
	}
	return stats, nil
}

// GetAuditStatsByStatus 按状态统计审计日志
func (s *AuditService) GetAuditStatsByStatus(startTime, endTime time.Time) (map[string]int64, error) {
	type Result struct {
		Status string
		Count  int64
	}
	var results []Result

	if err := s.db.Model(&models.AuditLog{}).
		Select("status, COUNT(*) as count").
		Where("created_at >= ? AND created_at <= ?", startTime, endTime).
		Group("status").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, r := range results {
		stats[r.Status] = r.Count
	}
	return stats, nil
}

// GetLoginStats 获取登录统计（用于替代单独的登录状态查询）
func (s *AuditService) GetLoginStats(startTime, endTime time.Time) (successCount, failCount int64, err error) {
	type Result struct {
		Status string
		Count  int64
	}
	var results []Result

	if err := s.db.Model(&models.AuditLog{}).
		Select("status, COUNT(*) as count").
		Where("action = ? AND created_at >= ? AND created_at <= ?", models.AuditActionLogin, startTime, endTime).
		Group("status").
		Scan(&results).Error; err != nil {
		return 0, 0, err
	}

	for _, r := range results {
		if r.Status == string(models.AuditStatusSuccess) {
			successCount = r.Count
		} else {
			failCount += r.Count
		}
	}

	return successCount, failCount, nil
}

// GetTopClients 获取热门客户端
func (s *AuditService) GetTopClients(startTime, endTime time.Time, limit int) ([]map[string]any, error) {
	type Result struct {
		ClientID string
		Count    int64
	}
	var results []Result

	if err := s.db.Model(&models.AuditLog{}).
		Select("client_id, COUNT(*) as count").
		Where("client_id IS NOT NULL AND created_at >= ? AND created_at <= ?", startTime, endTime).
		Group("client_id").
		Order("count DESC").
		Limit(limit).
		Scan(&results).Error; err != nil {
		return nil, err
	}

	// 获取客户端详情
	clients := make([]map[string]any, 0, len(results))
	for _, r := range results {
		var client models.OAuth2Client
		if err := s.db.First(&client, "id = ?", r.ClientID).Error; err == nil {
			clients = append(clients, map[string]any{
				"client_id":   r.ClientID,
				"client_name": client.Name,
				"count":       r.Count,
			})
		}
	}

	return clients, nil
}

// CleanupOldLogs 清理旧的审计日志
func (s *AuditService) CleanupOldLogs(retentionDays int) (int64, error) {
	cutoff := time.Now().AddDate(0, 0, -retentionDays)
	result := s.db.Where("created_at < ?", cutoff).Delete(&models.AuditLog{})
	return result.RowsAffected, result.Error
}
