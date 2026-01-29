package repository

import (
	"context"
	e "qauth-server/internal/errors"
	"qauth-server/internal/models"
	"time"

	"gorm.io/gorm"
)

type OAuthRepository struct {
	db *gorm.DB
}

func NewOAuthRepository(db *gorm.DB) *OAuthRepository {
	return &OAuthRepository{
		db: db,
	}
}

// FindClientByID 根据 ID 获取客户端
func (r *OAuthRepository) FindClientByID(ctx context.Context, id string) (*models.OAuth2Client, error) {
	var client models.OAuth2Client
	if err := r.db.WithContext(ctx).First(&client, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, e.ErrClientNotFound.WithMessage(id)
		}
		return nil, e.ErrClientQueryFailed.Wrap(err)
	}
	return &client, nil
}

// CreateClient 创建客户端
func (r *OAuthRepository) CreateClient(client *models.OAuth2Client) error {
	if err := r.db.Create(client).Error; err != nil {
		return e.ErrClientCreationFailed.Wrap(err)
	}
	return nil
}

// UpdateClient 更新客户端（使用结构体）
func (r *OAuthRepository) UpdateClient(client *models.OAuth2Client) error {
	if err := r.db.Save(client).Error; err != nil {
		return e.ErrClientUpdateFailed.Wrap(err)
	}
	return nil
}

// UpdateClientFields 更新客户端字段（使用 map）
func (r *OAuthRepository) UpdateClientFields(clientID string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return e.ErrNoFieldsToUpdate
	}

	if err := r.db.Model(&models.OAuth2Client{}).
		Where("id = ?", clientID).
		Updates(updates).Error; err != nil {
		return e.ErrClientUpdateFailed.Wrap(err)
	}
	return nil
}

// DeleteClient 删除客户端
func (r *OAuthRepository) DeleteClient(clientID string) error {
	if err := r.db.Delete(&models.OAuth2Client{}, "id = ?", clientID).Error; err != nil {
		return e.ErrClientDeletionFailed.Wrap(err)
	}
	return nil
}

// CountClients 统计客户端总数
func (r *OAuthRepository) CountClients() (int64, error) {
	var count int64
	if err := r.db.Model(&models.OAuth2Client{}).Count(&count).Error; err != nil {
		return 0, e.ErrClientQueryFailed.Wrap(err)
	}
	return count, nil
}

// ListClients 获取客户端列表（分页）
func (r *OAuthRepository) ListClients(page, pageSize int) ([]models.OAuth2Client, int64, error) {
	var clients []models.OAuth2Client
	var total int64

	// 计算总数
	if err := r.db.Model(&models.OAuth2Client{}).Count(&total).Error; err != nil {
		return nil, 0, e.ErrClientQueryFailed.Wrap(err)
	}

	// 查询数据
	offset := (page - 1) * pageSize
	if err := r.db.Offset(offset).Limit(pageSize).Find(&clients).Error; err != nil {
		return nil, 0, e.ErrClientQueryFailed.Wrap(err)
	}

	return clients, total, nil
}

// ListClientsWithFilter 获取客户端列表（支持过滤和排序）
func (r *OAuthRepository) ListClientsWithFilter(search, status, sortBy string, sortDesc bool, page, pageSize int) ([]models.OAuth2Client, int64, error) {
	var clients []models.OAuth2Client
	var total int64

	query := r.db.Model(&models.OAuth2Client{})

	// 搜索条件
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("name LIKE ? OR domain LIKE ? OR description LIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	// 状态过滤
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, e.ErrClientQueryFailed.Wrap(err)
	}

	// 排序
	orderClause := "created_at DESC"
	if sortBy != "" {
		orderClause = sortBy
		if sortDesc {
			orderClause += " DESC"
		} else {
			orderClause += " ASC"
		}
	}
	query = query.Order(orderClause)

	// 分页
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&clients).Error; err != nil {
		return nil, 0, e.ErrClientQueryFailed.Wrap(err)
	}

	return clients, total, nil
}

// UpdateClientSecret 更新客户端密钥
func (r *OAuthRepository) UpdateClientSecret(clientID, newSecret string) error {
	if err := r.db.Model(&models.OAuth2Client{}).
		Where("id = ?", clientID).
		Update("secret", newSecret).Error; err != nil {
		return e.ErrClientUpdateFailed.Wrap(err)
	}
	return nil
}

// IncrementClientRequestCount 增加客户端请求计数
func (r *OAuthRepository) IncrementClientRequestCount(clientID string) error {
	if err := r.db.Model(&models.OAuth2Client{}).
		Where("id = ?", clientID).
		Updates(map[string]interface{}{
			"request_count": gorm.Expr("request_count + 1"),
			"last_used_at":  time.Now(),
		}).Error; err != nil {
		return e.ErrClientUpdateFailed.Wrap(err)
	}
	return nil
}

// ClientStatusCount 客户端状态统计结果
type ClientStatusCount struct {
	Status models.ClientStatus
	Count  int64
}

// CountClientsByStatus 按状态统计客户端数量
func (r *OAuthRepository) CountClientsByStatus() ([]ClientStatusCount, error) {
	var statusCounts []ClientStatusCount
	if err := r.db.Model(&models.OAuth2Client{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Find(&statusCounts).Error; err != nil {
		return nil, e.ErrClientQueryFailed.Wrap(err)
	}
	return statusCounts, nil
}

// CreateLoginState 创建登录状态记录
func (r *OAuthRepository) CreateLoginState(loginState *models.LoginState) error {
	if err := r.db.Create(loginState).Error; err != nil {
		return e.ErrLoginStateRecordFailed.Wrap(err)
	}
	return nil
}

// CreateErrorRecord 创建错误记录
func (r *OAuthRepository) CreateErrorRecord(errorRecord *models.ErrorRecord) error {
	if err := r.db.Create(errorRecord).Error; err != nil {
		return e.ErrErrorRecordFailed.Wrap(err)
	}
	return nil
}
