package services

import (
	"encoding/json"
	e "qauth-server/internal/errors"
	"qauth-server/internal/models"
	"qauth-server/internal/providers"
	"qauth-server/internal/repository"
	"qauth-server/pkg/jwt"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

// AuditService 审计日志服务
type AuditService struct {
	repo    *repository.AuditRepository
	userSrv *UserService
	logger  providers.ILogger
}

// NewAuditService 创建审计日志服务
func NewAuditService(
	repo *repository.AuditRepository,
	userSrv *UserService,
	logger providers.ILogger,
) *AuditService {
	return &AuditService{
		repo:    repo,
		userSrv: userSrv,
		logger:  logger.With("service", "AuditService"),
	}
}

// AuditContext 审计上下文信息
type AuditContext struct {
	OperatorID   *string
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
	if userInfo, exists := c.Get("userInfo"); exists {
		userInfo := userInfo.(*jwt.UserJWTClaims)
		ctx.OperatorID = &userInfo.UserID

		// 获取用户名
		if user, err := s.userSrv.GetUserByID(userInfo.UserID, false); err == nil {
			ctx.OperatorName = user.Name
		}
	}

	return ctx
}

// Log 记录审计日志
func (s *AuditService) Log(ctx *AuditContext, entry *AuditEntry) error {
	// 构建详情 JSON
	var detailJSON datatypes.JSON
	if entry.Detail != nil {
		detailBytes, err := json.Marshal(entry.Detail)
		if err != nil {
			s.logger.Error("Failed to marshal audit detail", "error", err)
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

	if err := s.repo.Create(auditLog); err != nil {
		s.logger.Error("Failed to create audit log", "error", err)
		return err
	}

	s.logger.Info("Audit log created",
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

	var operatorID *string
	if userID != "" {
		operatorID = &userID
	}

	ctx := &AuditContext{
		OperatorID:   operatorID,
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

	var operatorID *string
	if userID != "" {
		operatorID = &userID
	}

	ctx := &AuditContext{
		OperatorID:   operatorID,
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
	return s.repo.Query(query)
}

// GetAuditLogByID 根据 ID 获取审计日志
func (s *AuditService) GetAuditLogByID(id string) (*models.AuditLog, error) {
	return s.repo.FindByID(id)
}

// GetRecentActivities 获取最近活动列表
func (s *AuditService) GetRecentActivities(limit int) ([]models.AuditLog, error) {
	if limit <= 0 {
		return nil, e.ErrInvalidParameter.Wrapf("limit must be greater than 0, got %d", limit).WithScope("GetRecentActivities")
	}

	return s.repo.FindRecent(limit)
}

// GetAuditStatsByModule 按模块统计审计日志
func (s *AuditService) GetAuditStatsByModule(startTime, endTime time.Time) (map[string]int64, error) {
	if startTime.IsZero() || endTime.IsZero() || !endTime.After(startTime) {
		return nil, e.ErrInvalidTimeRange.Wrapf("startTime: %v, endTime: %v", startTime, endTime).WithScope("GetAuditStatsByModule")
	}

	return s.repo.GetStatsByModule(startTime, endTime)
}

// GetAuditStatsByAction 按操作类型统计审计日志
func (s *AuditService) GetAuditStatsByAction(startTime, endTime time.Time) (map[string]int64, error) {
	if startTime.IsZero() || endTime.IsZero() || !endTime.After(startTime) {
		return nil, e.ErrInvalidTimeRange.Wrapf("startTime: %v, endTime: %v", startTime, endTime).WithScope("GetAuditStatsByAction")
	}

	return s.repo.GetStatsByAction(startTime, endTime)
}

// GetAuditStatsByStatus 按状态统计审计日志
func (s *AuditService) GetAuditStatsByStatus(startTime, endTime time.Time) (map[string]int64, error) {
	if startTime.IsZero() || endTime.IsZero() || !endTime.After(startTime) {
		return nil, e.ErrInvalidTimeRange.Wrapf("startTime: %v, endTime: %v", startTime, endTime).WithScope("GetAuditStatsByStatus")
	}

	return s.repo.GetStatsByStatus(startTime, endTime)
}

// GetLoginStats 获取登录统计（用于替代单独的登录状态查询）
func (s *AuditService) GetLoginStats(startTime, endTime time.Time) (successCount, failCount int64, err error) {
	if startTime.IsZero() || endTime.IsZero() || !endTime.After(startTime) {
		return 0, 0, e.ErrInvalidTimeRange.Wrapf("startTime: %v, endTime: %v", startTime, endTime).WithScope("GetLoginStats")
	}

	return s.repo.GetLoginStats(startTime, endTime)
}

// GetTopClients 获取热门客户端
func (s *AuditService) GetTopClients(startTime, endTime time.Time, limit int) ([]map[string]any, error) {
	results, err := s.repo.GetTopClients(startTime, endTime, limit)
	if err != nil {
		return nil, err
	}

	// 获取客户端详情
	clients := make([]map[string]any, 0, len(results))
	for _, r := range results {
		client, err := s.repo.FindClientByID(r.ClientID)
		if err == nil {
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
	return s.repo.DeleteOldLogs(cutoff)
}
