package database

import (
	"fmt"
	"log"

	"github.com/tofex/backend/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.RolePermission{},
		&models.Product{},
		&models.ProductImage{},
		&models.PaymentAccount{},
		&models.Order{},
		&models.OrderItem{},
		&models.OrderReceipt{},
		&models.AuditLog{},
		&models.PasswordReset{},
	)
}

func CountUsers(db *gorm.DB) (int64, error) {
	var n int64
	err := db.Model(&models.User{}).Count(&n).Error
	return n, err
}

func SeedRBAC(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.Role{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	log.Println("seeding default roles and permissions")

	roles := []struct {
		role        models.Role
		permissions []string
	}{
		{
			role: models.Role{
				Name:        "super_admin",
				Label:       "Super Admin",
				Description: "Full access including payment account management",
			},
			permissions: []string{"*"},
		},
		{
			role: models.Role{
				Name:        "admin",
				Label:       "Admin",
				Description: "Manage catalog, orders, and staff users",
			},
			permissions: []string{
				"users.read", "users.create", "users.update", "users.delete",
				"roles.read",
				"products.read", "products.create", "products.update", "products.delete",
				"payment_accounts.read",
				"orders.read", "orders.update",
				"audit.read",
			},
		},
		{
			role: models.Role{
				Name:        "finance",
				Label:       "Finance",
				Description: "View orders and payment instructions",
			},
			permissions: []string{
				"payment_accounts.read",
				"orders.read",
				"audit.read",
			},
		},
	}

	for _, r := range roles {
		if err := db.Create(&r.role).Error; err != nil {
			return err
		}
		for _, p := range r.permissions {
			if err := db.Create(&models.RolePermission{RoleID: r.role.ID, Permission: p}).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
