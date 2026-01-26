package models

// AppGroupPermission 应用组权限 - 属于特定 OAuth 应用的自定义权限
type AppGroupPermission struct {
	BaseModelWithUUID
	ClientID    string `gorm:"type:uuid;not null;index" json:"client_id"`                         // 所属的 OAuth 应用 ID
	Resource    string `gorm:"size:50;not null" json:"resource"`                                  // 资源名称
	Action      int8   `gorm:"not null" json:"action"`                                            // 操作类型: 1=Create, 2=Read, 3=Update, 4=Delete
	Code        string `gorm:"size:100;not null;uniqueIndex:idx_app_group_perm_code" json:"code"` // 权限代码（全局唯一）
	Name        string `gorm:"size:100;not null" json:"name"`                                     // 权限名称
	Description string `gorm:"size:255" json:"description"`                                       // 权限描述

	// Associations
	Client OAuth2Client   `gorm:"foreignKey:ClientID;references:ID" json:"client,omitempty"`
	Roles  []AppGroupRole `gorm:"many2many:app_group_roles_permissions;" json:"roles,omitempty"`
}

func (AppGroupPermission) TableName() string {
	return "app_group_permissions"
}

// AppGroupPermissionResponse 应用组权限响应
type AppGroupPermissionResponse struct {
	ID          string `json:"id"`
	ClientID    string `json:"client_id"`
	Resource    string `json:"resource"`
	Action      int8   `json:"action"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ToResponse 转换为响应格式
func (p *AppGroupPermission) ToResponse() *AppGroupPermissionResponse {
	return &AppGroupPermissionResponse{
		ID:          p.ID,
		ClientID:    p.ClientID,
		Resource:    p.Resource,
		Action:      p.Action,
		Code:        p.Code,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
