package business

import (
	app_error "qauth-server/internal/errors"
	"qauth-server/internal/models"
	"qauth-server/internal/permissions"
	"qauth-server/internal/services"
	"qauth-server/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

// AuditHandler 审计日志处理器
type AuditHandler struct {
	auditService *services.AuditService
	roleService  *services.RoleService
}

// NewAuditHandler 创建审计日志处理器
func NewAuditHandler(
	auditService *services.AuditService,
	roleService *services.RoleService,
) *AuditHandler {
	return &AuditHandler{
		auditService: auditService,
		roleService:  roleService,
	}
}

// GetAuditLogs 获取审计日志列表
// GET /audit/logs
func (h *AuditHandler) GetAuditLogs(c *gin.Context) {
	// 验证权限
	if err := services.VerifyPermissions(c, h.roleService, []string{permissions.AuditView}); err != nil {
		response.HandlerError(c, err)
		return
	}

	var query models.AuditLogQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	logs, total, err := h.auditService.QueryAuditLogs(&query)
	if err != nil {
		response.HandlerError(c, app_error.ErrFailedToGetAuditLogs)
		return
	}

	response.HandlerSuccess(c, gin.H{
		"items":     logs,
		"total":     total,
		"page":      query.Page,
		"page_size": query.PageSize,
	})
}

// GetAuditLogDetail 获取审计日志详情
// GET /audit/logs/:id
func (h *AuditHandler) GetAuditLogDetail(c *gin.Context) {
	// 验证权限
	if err := services.VerifyPermissions(c, h.roleService, []string{permissions.AuditView}); err != nil {
		response.HandlerError(c, err)
		return
	}

	id := c.Param("id")
	log, err := h.auditService.GetAuditLogByID(id)
	if err != nil {
		response.HandlerError(c, app_error.ErrAuditLogNotFound)
		return
	}

	response.HandlerSuccess(c, log)
}

// GetRecentActivities 获取最近活动
// GET /audit/activities
func (h *AuditHandler) GetRecentActivities(c *gin.Context) {
	// 验证权限
	if err := services.VerifyPermissions(c, h.roleService, []string{permissions.AuditView}); err != nil {
		response.HandlerError(c, err)
		return
	}

	limit := 10
	logs, err := h.auditService.GetRecentActivities(limit)
	if err != nil {
		response.HandlerError(c, app_error.ErrFailedToGetAuditLogs)
		return
	}

	// 转换为前端期望的格式
	activities := make([]gin.H, 0, len(logs))
	for _, log := range logs {
		activity := gin.H{
			"id":          log.ID,
			"user":        log.OperatorName,
			"operator_id": log.OperatorID,
			"action":      log.Action,
			"module":      log.Module,
			"target_id":   log.TargetID,
			"target_name": log.TargetName,
			"ip":          log.IP,
			"time":        log.CreatedAt,
			"status":      log.Status,
			"duration_ms": log.DurationMs,
		}

		// 添加用户头像（如果有）
		if log.Operator.ID != "" {
			if log.Operator.Avatar != nil && log.Operator.Avatar.File != nil {
				activity["avatar"] = "/uploads/" + log.Operator.Avatar.File.StorageKey
			} else {
				activity["avatar"] = "https://api.dicebear.com/7.x/avataaars/svg?seed=" + log.OperatorID
			}
		}

		// 添加客户端信息
		if log.Client != nil {
			activity["client"] = log.Client.Name
		} else {
			activity["client"] = "System"
		}

		activities = append(activities, activity)
	}

	response.HandlerSuccess(c, activities)
}

// GetAuditStats 获取审计统计
// GET /audit/stats
func (h *AuditHandler) GetAuditStats(c *gin.Context) {
	// 验证权限
	if err := services.VerifyPermissions(c, h.roleService, []string{permissions.AuditView}); err != nil {
		response.HandlerError(c, err)
		return
	}

	// 默认获取最近 7 天的统计
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -7)

	// 按模块统计
	moduleStats, err := h.auditService.GetAuditStatsByModule(startTime, endTime)
	if err != nil {
		response.HandlerError(c, app_error.ErrFailedToGetAuditLogs)
		return
	}

	// 按操作统计
	actionStats, err := h.auditService.GetAuditStatsByAction(startTime, endTime)
	if err != nil {
		response.HandlerError(c, app_error.ErrFailedToGetAuditLogs)
		return
	}

	// 按状态统计
	statusStats, err := h.auditService.GetAuditStatsByStatus(startTime, endTime)
	if err != nil {
		response.HandlerError(c, app_error.ErrFailedToGetAuditLogs)
		return
	}

	// 登录统计
	successCount, failCount, err := h.auditService.GetLoginStats(startTime, endTime)
	if err != nil {
		response.HandlerError(c, app_error.ErrFailedToGetAuditLogs)
		return
	}

	response.HandlerSuccess(c, gin.H{
		"module_stats": moduleStats,
		"action_stats": actionStats,
		"status_stats": statusStats,
		"login_stats": gin.H{
			"success": successCount,
			"fail":    failCount,
		},
		"start_time": startTime,
		"end_time":   endTime,
	})
}

// GetTopClients 获取热门客户端
// GET /audit/top-clients
func (h *AuditHandler) GetTopClients(c *gin.Context) {
	// 验证权限
	if err := services.VerifyPermissions(c, h.roleService, []string{permissions.AuditView}); err != nil {
		response.HandlerError(c, err)
		return
	}

	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -30)

	clients, err := h.auditService.GetTopClients(startTime, endTime, 10)
	if err != nil {
		response.HandlerError(c, app_error.ErrFailedToGetAuditLogs)
		return
	}

	response.HandlerSuccess(c, clients)
}

// ExportAuditLogs 导出审计日志
// GET /audit/export
func (h *AuditHandler) ExportAuditLogs(c *gin.Context) {
	// 验证权限
	if err := services.VerifyPermissions(c, h.roleService, []string{permissions.AuditExport}); err != nil {
		response.HandlerError(c, err)
		return
	}

	var query models.AuditLogQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	// 限制导出数量
	query.PageSize = 10000

	logs, _, err := h.auditService.QueryAuditLogs(&query)
	if err != nil {
		response.HandlerError(c, app_error.ErrFailedToGetAuditLogs)
		return
	}

	// 返回 JSON 格式的导出数据
	c.Header("Content-Disposition", "attachment; filename=audit_logs.json")
	c.Header("Content-Type", "application/json")
	c.JSON(200, logs)
}
