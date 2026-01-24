package models

import (
	"time"

	"gorm.io/datatypes"
)

// AuditAction 审计操作类型
type AuditAction string

const (
	// 认证相关操作
	AuditActionLogin          AuditAction = "LOGIN"
	AuditActionLogout         AuditAction = "LOGOUT"
	AuditActionRegister       AuditAction = "REGISTER"
	AuditActionPasswordReset  AuditAction = "PASSWORD_RESET"
	AuditActionPasswordChange AuditAction = "PASSWORD_CHANGE"
	AuditActionTokenRefresh   AuditAction = "TOKEN_REFRESH"

	// OAuth2 相关操作
	AuditActionOAuthAuthorize AuditAction = "OAUTH_AUTHORIZE"
	AuditActionOAuthToken     AuditAction = "OAUTH_TOKEN"
	AuditActionOAuthRevoke    AuditAction = "OAUTH_REVOKE"
	AuditActionClientCreate   AuditAction = "CLIENT_CREATE"
	AuditActionClientUpdate   AuditAction = "CLIENT_UPDATE"
	AuditActionClientDelete   AuditAction = "CLIENT_DELETE"

	// 角色权限相关操作
	AuditActionRoleCreate       AuditAction = "ROLE_CREATE"
	AuditActionRoleUpdate       AuditAction = "ROLE_UPDATE"
	AuditActionRoleDelete       AuditAction = "ROLE_DELETE"
	AuditActionPermissionCreate AuditAction = "PERMISSION_CREATE"
	AuditActionPermissionUpdate AuditAction = "PERMISSION_UPDATE"
	AuditActionPermissionDelete AuditAction = "PERMISSION_DELETE"
	AuditActionPermissionGrant  AuditAction = "PERMISSION_GRANT"
	AuditActionPermissionRevoke AuditAction = "PERMISSION_REVOKE"

	// 用户管理相关操作
	AuditActionUserCreate AuditAction = "USER_CREATE"
	AuditActionUserUpdate AuditAction = "USER_UPDATE"
	AuditActionUserDelete AuditAction = "USER_DELETE"

	// 系统管理相关操作
	AuditActionKeyRotation    AuditAction = "KEY_ROTATION"
	AuditActionSettingsChange AuditAction = "SETTINGS_CHANGE"
)

// AuditModule 审计模块
type AuditModule string

const (
	AuditModuleAuth       AuditModule = "AUTH"
	AuditModuleOAuth      AuditModule = "OAUTH"
	AuditModuleUser       AuditModule = "USER"
	AuditModuleRole       AuditModule = "ROLE"
	AuditModulePermission AuditModule = "PERMISSION"
	AuditModuleClient     AuditModule = "CLIENT"
	AuditModuleSystem     AuditModule = "SYSTEM"
)

// AuditStatus 审计状态
type AuditStatus string

const (
	AuditStatusSuccess AuditStatus = "SUCCESS"
	AuditStatusWarning AuditStatus = "WARNING"
	AuditStatusError   AuditStatus = "ERROR"
)

// AuditLog 审计日志模型
type AuditLog struct {
	ID           string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	OperatorID   string         `gorm:"type:uuid;index" json:"operator_id"`
	OperatorName string         `gorm:"size:100" json:"operator_name"`
	Module       AuditModule    `gorm:"size:50;not null;index" json:"module"`
	Action       AuditAction    `gorm:"size:50;not null;index" json:"action"`
	TargetID     string         `gorm:"size:255;index" json:"target_id"`
	TargetType   string         `gorm:"size:50" json:"target_type"`
	TargetName   string         `gorm:"size:255" json:"target_name"`
	Detail       datatypes.JSON `gorm:"type:jsonb" json:"detail"`
	IP           string         `gorm:"size:45" json:"ip"`
	UserAgent    string         `gorm:"size:500" json:"user_agent"`
	Location     string         `gorm:"size:200" json:"location"`
	Status       AuditStatus    `gorm:"size:20;not null;default:'SUCCESS'" json:"status"`
	ErrorMessage string         `gorm:"size:500" json:"error_message"`
	DurationMs   int64          `gorm:"default:0" json:"duration_ms"`
	ClientID     *string        `gorm:"type:uuid;index" json:"client_id,omitempty"`
	SessionID    string         `gorm:"size:100" json:"session_id"`
	RequestID    string         `gorm:"size:100" json:"request_id"`
	CreatedAt    time.Time      `gorm:"autoCreateTime;index" json:"created_at"`

	// Associations
	Operator Users         `gorm:"foreignKey:OperatorID;references:ID" json:"operator,omitempty"`
	Client   *OAuth2Client `gorm:"foreignKey:ClientID;references:ID" json:"client,omitempty"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

// AuditLogDetail 审计日志详情结构（用于 JSON 存储）
type AuditLogDetail struct {
	// 认证相关
	LoginType  string `json:"login_type,omitempty"`
	AuthMethod string `json:"auth_method,omitempty"`
	FailReason string `json:"fail_reason,omitempty"`

	// OAuth 相关
	Scopes       []string `json:"scopes,omitempty"`
	GrantType    string   `json:"grant_type,omitempty"`
	RedirectURI  string   `json:"redirect_uri,omitempty"`
	ResponseType string   `json:"response_type,omitempty"`

	// 变更相关
	OldValue      any      `json:"old_value,omitempty"`
	NewValue      any      `json:"new_value,omitempty"`
	ChangedFields []string `json:"changed_fields,omitempty"`

	// 其他
	Metadata map[string]any `json:"metadata,omitempty"`
}

// AuditLogQuery 审计日志查询参数
type AuditLogQuery struct {
	OperatorID string      `form:"operator_id"`
	Module     AuditModule `form:"module"`
	Action     AuditAction `form:"action"`
	TargetID   string      `form:"target_id"`
	Status     AuditStatus `form:"status"`
	ClientID   string      `form:"client_id"`
	StartTime  *time.Time  `form:"start_time"`
	EndTime    *time.Time  `form:"end_time"`
	Page       int         `form:"page"`
	PageSize   int         `form:"page_size"`
}
