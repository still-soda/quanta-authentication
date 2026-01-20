package database

import (
	"qauth-server/internal/models"
	"qauth-server/internal/permissions"
	"qauth-server/internal/services"
	"qauth-server/internal/utilities"

	"gorm.io/gorm"
)

func SeedingDB(db *gorm.DB) error {
	// 创建权限
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
	roleService := services.NewRoleService(db, services.NewPermissionService(db))
	roleService.GrantPermissionToRole(superAdmin.ID, []string{
		permissions.OAuthClientCreate,
		permissions.OAuthClientDelete,
		permissions.OAuthClientView,
		permissions.OAuthClientList,
		permissions.OAuthClientUpdate,
	})
	roleService.GrantPermissionToRole(admin.ID, []string{
		permissions.OAuthClientView,
		permissions.OAuthClientList,
	})

	salt, _ := utilities.GenerateSalt(16)
	hash := utilities.HashPassword("123456", salt)
	adminUser := &models.Users{Name: "超级管理员", StudentID: "20231003059", Email: "951040628@qq.com", Salt: salt, PasswordHash: hash}
	db.FirstOrCreate(&adminUser)

	adminUserRole := &models.UsersRoles{UserID: adminUser.ID, RoleID: superAdmin.ID}
	db.FirstOrCreate(&adminUserRole)

	return nil
}
