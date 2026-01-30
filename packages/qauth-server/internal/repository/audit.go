package repository

import (
	"errors"
	"qauth-server/internal/models"
	"time"

	"gorm.io/gorm"
)

type AuditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) *AuditRepository {
	return &AuditRepository{
		db: db,
	}
}

// Create 创建审计日志
func (r *AuditRepository) Create(log *models.AuditLog) error {
	return r.db.Create(log).Error
}

// FindByID 根据 ID 获取审计日志
func (r *AuditRepository) FindByID(id string) (*models.AuditLog, error) {
	var log models.AuditLog
	if err := r.db.Preload("Operator").Preload("Client").First(&log, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

// Query 查询审计日志
func (r *AuditRepository) Query(query *models.AuditLogQuery) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64

	db := r.db.Model(&models.AuditLog{})

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

	// 排序
	orderClause := "created_at DESC"
	if query.SortBy != "" {
		allowedSortFields := map[string]string{
			"operator_name": "operator_name",
			"module":        "module",
			"action":        "action",
			"target_name":   "target_name",
			"status":        "status",
			"duration_ms":   "duration_ms",
			"created_at":    "created_at",
		}
		if field, ok := allowedSortFields[query.SortBy]; ok {
			direction := "ASC"
			if query.SortDesc {
				direction = "DESC"
			}
			orderClause = field + " " + direction
		}
	}

	offset := (query.Page - 1) * query.PageSize
	if err := db.Preload("Operator").
		Order(orderClause).
		Offset(offset).
		Limit(query.PageSize).
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// FindRecent 获取最近活动列表
func (r *AuditRepository) FindRecent(limit int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	if err := r.db.Preload("Operator").
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

// GetStatsByModule 按模块统计审计日志
func (r *AuditRepository) GetStatsByModule(startTime, endTime time.Time) (map[string]int64, error) {
	type Result struct {
		Module string
		Count  int64
	}
	var results []Result

	if err := r.db.Model(&models.AuditLog{}).
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

// GetStatsByAction 按操作类型统计审计日志
func (r *AuditRepository) GetStatsByAction(startTime, endTime time.Time) (map[string]int64, error) {
	type Result struct {
		Action string
		Count  int64
	}
	var results []Result

	if err := r.db.Model(&models.AuditLog{}).
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

// GetStatsByStatus 按状态统计审计日志
func (r *AuditRepository) GetStatsByStatus(startTime, endTime time.Time) (map[string]int64, error) {
	type Result struct {
		Status string
		Count  int64
	}
	var results []Result

	if err := r.db.Model(&models.AuditLog{}).
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

// GetLoginStats 获取登录统计
func (r *AuditRepository) GetLoginStats(startTime, endTime time.Time) (successCount, failCount int64, err error) {
	type Result struct {
		Status string
		Count  int64
	}
	var results []Result

	if err := r.db.Model(&models.AuditLog{}).
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
func (r *AuditRepository) GetTopClients(startTime, endTime time.Time, limit int) ([]TopClientResult, error) {
	var results []TopClientResult

	if err := r.db.Model(&models.AuditLog{}).
		Select("client_id, COUNT(*) as count").
		Where("client_id IS NOT NULL AND created_at >= ? AND created_at <= ?", startTime, endTime).
		Group("client_id").
		Order("count DESC").
		Limit(limit).
		Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

// TopClientResult 热门客户端结果
type TopClientResult struct {
	ClientID string
	Count    int64
}

// FindClientByID 根据 ID 查找客户端（用于 GetTopClients）
func (r *AuditRepository) FindClientByID(clientID string) (*models.OAuth2Client, error) {
	var client models.OAuth2Client
	if err := r.db.First(&client, "id = ?", clientID).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

// DeleteOldLogs 删除旧的审计日志
func (r *AuditRepository) DeleteOldLogs(cutoffTime time.Time) (int64, error) {
	result := r.db.Where("created_at < ?", cutoffTime).Delete(&models.AuditLog{})
	return result.RowsAffected, result.Error
}

// FindLastLoginByUserID 获取用户最后登录时间
func (r *AuditRepository) FindLastLoginByUserID(userID string) (*time.Time, error) {
	var auditLog models.AuditLog
	if err := r.db.Where("operator_id = ? AND action = ?", userID, models.AuditActionLogin).
		Order("created_at DESC").
		First(&auditLog).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &auditLog.CreatedAt, nil
}
