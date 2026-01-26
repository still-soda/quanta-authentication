package models

import "time"

// AppGroupAdminType 应用组管理员类型
type AppGroupAdminType string

const (
	// AppGroupAdminTypeOwner 应用创建者/所有者
	AppGroupAdminTypeOwner AppGroupAdminType = "owner"
	// AppGroupAdminTypeAdmin 应用组管理员
	AppGroupAdminTypeAdmin AppGroupAdminType = "admin"
	// AppGroupAdminTypeRoleManager 角色管理员 - 可以管理角色和分配角色
	AppGroupAdminTypeRoleManager AppGroupAdminType = "role_manager"
	// AppGroupAdminTypePermissionManager 权限管理员 - 可以管理权限和分配权限
	AppGroupAdminTypePermissionManager AppGroupAdminType = "permission_manager"
)

// AppGroupAdmin 应用组管理员 - 控制用户管理应用组的权限
type AppGroupAdmin struct {
	BaseModelWithUUID
	ClientID  string            `gorm:"type:uuid;not null;index;uniqueIndex:idx_app_group_admin,priority:1" json:"client_id"`   // 所属的 OAuth 应用 ID
	UserID    string            `gorm:"type:uuid;not null;index;uniqueIndex:idx_app_group_admin,priority:2" json:"user_id"`     // 用户 ID
	AdminType AppGroupAdminType `gorm:"type:varchar(50);not null;uniqueIndex:idx_app_group_admin,priority:3" json:"admin_type"` // 管理员类型
	GrantedAt time.Time         `gorm:"autoCreateTime" json:"granted_at"`                                                       // 授权时间
	GrantedBy string            `gorm:"type:uuid" json:"granted_by"`                                                            // 授权者 ID

	// Associations
	Client  OAuth2Client `gorm:"foreignKey:ClientID;references:ID" json:"client,omitempty"`
	User    Users        `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Granter Users        `gorm:"foreignKey:GrantedBy;references:ID" json:"granter,omitempty"`
}

func (AppGroupAdmin) TableName() string {
	return "app_group_admins"
}

// AppGroupAdminResponse 应用组管理员响应
type AppGroupAdminResponse struct {
	ID          string            `json:"id"`
	ClientID    string            `json:"client_id"`
	UserID      string            `json:"user_id"`
	UserName    string            `json:"user_name"`
	UserEmail   string            `json:"user_email"`
	AdminType   AppGroupAdminType `json:"admin_type"`
	GrantedAt   string            `json:"granted_at"`
	GrantedBy   string            `json:"granted_by"`
	GranterName string            `json:"granter_name,omitempty"`
	CreatedAt   string            `json:"created_at"`
}

// ToResponse 转换为响应格式
func (a *AppGroupAdmin) ToResponse() *AppGroupAdminResponse {
	resp := &AppGroupAdminResponse{
		ID:        a.ID,
		ClientID:  a.ClientID,
		UserID:    a.UserID,
		AdminType: a.AdminType,
		GrantedAt: a.GrantedAt.Format("2006-01-02T15:04:05Z07:00"),
		GrantedBy: a.GrantedBy,
		CreatedAt: a.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if a.User.ID != "" {
		resp.UserName = a.User.Name
		resp.UserEmail = a.User.Email
	}
	if a.Granter.ID != "" {
		resp.GranterName = a.Granter.Name
	}
	return resp
}

// GetAdminTypeDescription 获取管理员类型描述
func GetAdminTypeDescription(adminType AppGroupAdminType) string {
	switch adminType {
	case AppGroupAdminTypeOwner:
		return "应用所有者"
	case AppGroupAdminTypeAdmin:
		return "应用组管理员"
	case AppGroupAdminTypeRoleManager:
		return "角色管理员"
	case AppGroupAdminTypePermissionManager:
		return "权限管理员"
	default:
		return "未知类型"
	}
}

// IsHigherOrEqualAdminType 检查 a 是否具有高于或等于 b 的管理员权限
func IsHigherOrEqualAdminType(a, b AppGroupAdminType) bool {
	priority := map[AppGroupAdminType]int{
		AppGroupAdminTypeOwner:             100,
		AppGroupAdminTypeAdmin:             80,
		AppGroupAdminTypeRoleManager:       50,
		AppGroupAdminTypePermissionManager: 50,
	}
	return priority[a] >= priority[b]
}
