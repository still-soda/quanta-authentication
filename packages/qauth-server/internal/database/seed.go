package database

import (
	"qauth-server/internal/models"
	"qauth-server/internal/permissions"
	"qauth-server/internal/services"
	"qauth-server/internal/utilities"

	"gorm.io/gorm"
)

func SeedingDB(db *gorm.DB) error {
	// 创建 OAuth 客户端相关权限
	createClientPermission := &models.Permissions{Resource: "oauth_clients", Action: int8(permissions.Create), Code: permissions.OAuthClientCreate, Description: "创建 OAuth2 客户端"}
	db.Where(models.Permissions{Code: permissions.OAuthClientCreate}).
		FirstOrCreate(&createClientPermission)
	deleteClientPermission := &models.Permissions{Resource: "oauth_clients", Action: int8(permissions.Delete), Code: permissions.OAuthClientDelete, Description: "删除 OAuth2 客户端"}
	db.Where(models.Permissions{Code: permissions.OAuthClientDelete}).
		FirstOrCreate(&deleteClientPermission)
	viewClientPermission := &models.Permissions{Resource: "oauth_clients", Action: int8(permissions.Read), Code: permissions.OAuthClientView, Description: "查看 OAuth2 客户端"}
	db.Where(models.Permissions{Code: permissions.OAuthClientView}).
		FirstOrCreate(&viewClientPermission)
	listClientPermission := &models.Permissions{Resource: "oauth_clients", Action: int8(permissions.Read), Code: permissions.OAuthClientList, Description: "列出 OAuth2 客户端"}
	db.Where(models.Permissions{Code: permissions.OAuthClientList}).
		FirstOrCreate(&listClientPermission)
	updateClientPermission := &models.Permissions{Resource: "oauth_clients", Action: int8(permissions.Update), Code: permissions.OAuthClientUpdate, Description: "更新 OAuth2 客户端"}
	db.Where(models.Permissions{Code: permissions.OAuthClientUpdate}).
		FirstOrCreate(&updateClientPermission)

	// 创建系统管理相关权限
	systemKeyRotationViewPermission := &models.Permissions{Resource: "system", Action: int8(permissions.Read), Code: permissions.SystemKeyRotationView, Description: "查看密钥轮换信息"}
	db.Where(models.Permissions{Code: permissions.SystemKeyRotationView}).
		FirstOrCreate(&systemKeyRotationViewPermission)
	systemKeyRotationExecutePermission := &models.Permissions{Resource: "system", Action: int8(permissions.Update), Code: permissions.SystemKeyRotationExecute, Description: "执行密钥轮换"}
	db.Where(models.Permissions{Code: permissions.SystemKeyRotationExecute}).
		FirstOrCreate(&systemKeyRotationExecutePermission)

	// 创建角色相关权限
	roleCreatePermission := &models.Permissions{Resource: "roles", Action: int8(permissions.Create), Code: permissions.RoleCreate, Description: "创建角色"}
	db.Where(models.Permissions{Code: permissions.RoleCreate}).
		FirstOrCreate(&roleCreatePermission)
	roleDeletePermission := &models.Permissions{Resource: "roles", Action: int8(permissions.Delete), Code: permissions.RoleDelete, Description: "删除角色"}
	db.Where(models.Permissions{Code: permissions.RoleDelete}).
		FirstOrCreate(&roleDeletePermission)
	roleViewPermission := &models.Permissions{Resource: "roles", Action: int8(permissions.Read), Code: permissions.RoleView, Description: "查看角色"}
	db.Where(models.Permissions{Code: permissions.RoleView}).
		FirstOrCreate(&roleViewPermission)
	roleUpdatePermission := &models.Permissions{Resource: "roles", Action: int8(permissions.Update), Code: permissions.RoleUpdate, Description: "更新角色"}
	db.Where(models.Permissions{Code: permissions.RoleUpdate}).
		FirstOrCreate(&roleUpdatePermission)
	roleAssignPermissionsPermission := &models.Permissions{Resource: "roles", Action: int8(permissions.Update), Code: permissions.RoleAssignPermissions, Description: "为角色分配权限"}
	db.Where(models.Permissions{Code: permissions.RoleAssignPermissions}).
		FirstOrCreate(&roleAssignPermissionsPermission)
	roleRevokePermissionsPermission := &models.Permissions{Resource: "roles", Action: int8(permissions.Update), Code: permissions.RoleRevokePermissions, Description: "撤销角色权限"}
	db.Where(models.Permissions{Code: permissions.RoleRevokePermissions}).
		FirstOrCreate(&roleRevokePermissionsPermission)

	// 创建审计日志相关权限
	auditViewPermission := &models.Permissions{Resource: "audit", Action: int8(permissions.Read), Code: permissions.AuditView, Description: "查看审计日志"}
	db.Where(models.Permissions{Code: permissions.AuditView}).
		FirstOrCreate(&auditViewPermission)
	auditExportPermission := &models.Permissions{Resource: "audit", Action: int8(permissions.Read), Code: permissions.AuditExport, Description: "导出审计日志"}
	db.Where(models.Permissions{Code: permissions.AuditExport}).
		FirstOrCreate(&auditExportPermission)

	// 创建用户相关权限
	userListPermission := &models.Permissions{Resource: "users", Action: int8(permissions.Read), Code: permissions.UserList, Description: "列出用户"}
	db.Where(models.Permissions{Code: permissions.UserList}).
		FirstOrCreate(&userListPermission)
	userCreatePermission := &models.Permissions{Resource: "users", Action: int8(permissions.Create), Code: permissions.UserCreate, Description: "创建用户"}
	db.Where(models.Permissions{Code: permissions.UserCreate}).
		FirstOrCreate(&userCreatePermission)
	userDeletePermission := &models.Permissions{Resource: "users", Action: int8(permissions.Delete), Code: permissions.UserDelete, Description: "删除用户"}
	db.Where(models.Permissions{Code: permissions.UserDelete}).
		FirstOrCreate(&userDeletePermission)
	userViewPermission := &models.Permissions{Resource: "users", Action: int8(permissions.Read), Code: permissions.UserView, Description: "查看用户"}
	db.Where(models.Permissions{Code: permissions.UserView}).
		FirstOrCreate(&userViewPermission)
	userUpdatePermission := &models.Permissions{Resource: "users", Action: int8(permissions.Update), Code: permissions.UserUpdate, Description: "更新用户"}
	db.Where(models.Permissions{Code: permissions.UserUpdate}).
		FirstOrCreate(&userUpdatePermission)
	userAssignRolesPermission := &models.Permissions{Resource: "users", Action: int8(permissions.Update), Code: permissions.UserAssignRoles, Description: "为用户分配角色"}
	db.Where(models.Permissions{Code: permissions.UserAssignRoles}).
		FirstOrCreate(&userAssignRolesPermission)
	userRevokeRolesPermission := &models.Permissions{Resource: "users", Action: int8(permissions.Update), Code: permissions.UserRevokeRoles, Description: "撤销用户角色"}
	db.Where(models.Permissions{Code: permissions.UserRevokeRoles}).
		FirstOrCreate(&userRevokeRolesPermission)
	userResetPasswordPermission := &models.Permissions{Resource: "users", Action: int8(permissions.Update), Code: permissions.UserResetPassword, Description: "重置用户密码"}
	db.Where(models.Permissions{Code: permissions.UserResetPassword}).
		FirstOrCreate(&userResetPasswordPermission)

	// 创建角色
	superAdmin := &models.Roles{Code: "system_super_admin", Name: "系统超级管理员", Description: "拥有系统内所有权限", IsSystem: true}
	db.Where(&models.Roles{Code: "system_super_admin"}).
		FirstOrCreate(&superAdmin)
	admin := models.Roles{Code: "system_admin", Name: "系统管理员", Description: "拥有系统内大部分权限", IsSystem: true}
	db.Where(&models.Roles{Code: "system_admin"}).
		FirstOrCreate(&admin)
	user := models.Roles{Code: "system_user", Name: "系统普通用户", Description: "拥有系统内基本权限", IsSystem: true}
	db.Where(&models.Roles{Code: "system_user"}).
		FirstOrCreate(&user)

	// 分配权限到角色
	roleService := services.NewRoleService(db, services.NewPermissionService(db), services.NewUserService(db))

	// 超级管理员：拥有所有权限
	roleService.GrantPermissionToRole(superAdmin.ID, []string{
		// OAuth 客户端权限
		permissions.OAuthClientCreate,
		permissions.OAuthClientDelete,
		permissions.OAuthClientView,
		permissions.OAuthClientList,
		permissions.OAuthClientUpdate,
		// 系统管理权限
		permissions.SystemKeyRotationView,
		permissions.SystemKeyRotationExecute,
		// 角色管理权限
		permissions.RoleCreate,
		permissions.RoleDelete,
		permissions.RoleView,
		permissions.RoleUpdate,
		permissions.RoleAssignPermissions,
		permissions.RoleRevokePermissions,
		// 审计日志权限
		permissions.AuditView,
		permissions.AuditExport,
		// 用户管理权限
		permissions.UserList,
		permissions.UserCreate,
		permissions.UserDelete,
		permissions.UserView,
		permissions.UserUpdate,
		permissions.UserAssignRoles,
		permissions.UserRevokeRoles,
		permissions.UserResetPassword,
	})

	// 系统管理员：拥有大部分权限，不包括敏感操作
	roleService.GrantPermissionToRole(admin.ID, []string{
		// OAuth 客户端权限（不包括删除）
		permissions.OAuthClientCreate,
		permissions.OAuthClientView,
		permissions.OAuthClientList,
		permissions.OAuthClientUpdate,
		// 系统管理权限（仅查看）
		permissions.SystemKeyRotationView,
		// 角色管理权限（不包括创建和删除）
		permissions.RoleView,
		permissions.RoleUpdate,
		permissions.RoleAssignPermissions,
		permissions.RoleRevokePermissions,
		// 审计日志权限
		permissions.AuditView,
		permissions.AuditExport,
		// 用户管理权限（不包括创建和删除用户）
		permissions.UserList,
		permissions.UserView,
		permissions.UserUpdate,
		permissions.UserAssignRoles,
		permissions.UserRevokeRoles,
		permissions.UserResetPassword,
	})

	// 普通用户：仅基本查看权限
	roleService.GrantPermissionToRole(user.ID, []string{
		// OAuth 客户端权限（仅查看）
		permissions.OAuthClientView,
		permissions.OAuthClientList,
		// 审计日志权限（仅查看）
		permissions.AuditView,
		// 用户管理权限（仅查看自己信息）
		permissions.UserView,
	})

	// 创建超级管理员用户
	salt, _ := utilities.GenerateSalt(16)
	hash := utilities.HashPassword("123456", salt)
	adminUser := &models.Users{Name: "超级管理员", StudentID: "20231003059", Email: "951040628@qq.com", Salt: salt, PasswordHash: hash}
	db.FirstOrCreate(&adminUser)

	// 分配超级管理员角色给用户
	roleService.AssignRolesToUserByCode(adminUser.ID, []string{
		"system_super_admin",
	})

	return nil
}
