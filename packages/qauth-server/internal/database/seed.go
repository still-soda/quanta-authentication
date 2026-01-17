package database

import (
	"qauth-server/internal/models"
	"qauth-server/internal/utilities"

	"gorm.io/gorm"
)

func SeedingDB(db *gorm.DB) error {
	superAdmin := &models.Roles{Code: "system_super_admin", Name: "系统超级管理员", Description: "拥有系统内所有权限", IsSystem: true}
	db.Create(&superAdmin)
	admin := models.Roles{Code: "system_admin", Name: "系统管理员", Description: "拥有系统内大部分权限", IsSystem: true}
	db.Create(&admin)
	user := models.Roles{Code: "system_user", Name: "系统普通用户", Description: "拥有系统内基本权限", IsSystem: true}
	db.Create(&user)

	salt, _ := utilities.GenerateSalt(16)
	hash := utilities.HashPassword("123456", salt)
	adminUser := &models.Users{Name: "超级管理员", StudentId: "20231003059", Email: "951040628@qq.com", Salt: salt, PasswordHash: hash}
	db.Create(&adminUser)

	adminUserRole := &models.UsersRoles{UserID: adminUser.ID, RoleID: superAdmin.ID}
	db.Create(&adminUserRole)

	return nil
}
