package models

import "time"

// AppGroupRole 应用组角色 - 属于特定 OAuth 应用的角色
type AppGroupRole struct {
	BaseModelWithUUID
	ClientID    string `gorm:"type:uuid;not null;index" json:"client_id"` // 所属的 OAuth 应用 ID
	Code        string `gorm:"size:100;not null;uniqueIndex" json:"code"` // 角色代码（全局唯一）
	Name        string `gorm:"size:100;not null" json:"name"`             // 角色名称
	Description string `gorm:"size:255" json:"description"`               // 角色描述
	IsDefault   bool   `gorm:"default:false" json:"is_default"`           // 是否为默认角色

	// Associations
	Client      OAuth2Client         `gorm:"foreignKey:ClientID;references:ID" json:"client,omitempty"`
	Permissions []AppGroupPermission `gorm:"many2many:app_group_roles_permissions;" json:"permissions,omitempty"`
	Users       []Users              `gorm:"many2many:app_group_users_roles;" json:"users,omitempty"`
}

func (AppGroupRole) TableName() string {
	return "app_group_roles"
}

// AppGroupRoleResponse 应用组角色响应
type AppGroupRoleResponse struct {
	ID              string `json:"id"`
	ClientID        string `json:"client_id"`
	Code            string `json:"code"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	IsDefault       bool   `json:"is_default"`
	UserCount       int64  `json:"user_count"`
	PermissionCount int64  `json:"permission_count"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

// ToResponse 转换为响应格式
func (r *AppGroupRole) ToResponse() *AppGroupRoleResponse {
	return &AppGroupRoleResponse{
		ID:          r.ID,
		ClientID:    r.ClientID,
		Code:        r.Code,
		Name:        r.Name,
		Description: r.Description,
		IsDefault:   r.IsDefault,
		CreatedAt:   r.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   r.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// AppGroupRolesPermissions 应用组角色与权限关联表
type AppGroupRolesPermissions struct {
	AppGroupRoleID       string `gorm:"type:uuid;primaryKey" json:"app_group_role_id"`
	AppGroupPermissionID string `gorm:"type:uuid;primaryKey" json:"app_group_permission_id"`

	// Associations
	Role       AppGroupRole       `gorm:"foreignKey:AppGroupRoleID;references:ID" json:"role,omitempty"`
	Permission AppGroupPermission `gorm:"foreignKey:AppGroupPermissionID;references:ID" json:"permission,omitempty"`
}

func (AppGroupRolesPermissions) TableName() string {
	return "app_group_roles_permissions"
}

// AppGroupUsersRoles 应用组用户角色关联表 - 用户在应用组中的角色分配
type AppGroupUsersRoles struct {
	ID             string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID         string    `gorm:"type:uuid;not null;index;uniqueIndex:idx_app_group_user_role,priority:1" json:"user_id"`
	AppGroupRoleID string    `gorm:"type:uuid;not null;index;uniqueIndex:idx_app_group_user_role,priority:2" json:"app_group_role_id"`
	AssignedAt     time.Time `gorm:"autoCreateTime" json:"assigned_at"`
	AssignedBy     string    `gorm:"type:uuid" json:"assigned_by"` // 分配者 ID

	// Associations
	User         Users        `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	AppGroupRole AppGroupRole `gorm:"foreignKey:AppGroupRoleID;references:ID" json:"app_group_role,omitempty"`
	Assigner     Users        `gorm:"foreignKey:AssignedBy;references:ID" json:"assigner,omitempty"`
}

func (AppGroupUsersRoles) TableName() string {
	return "app_group_users_roles"
}
