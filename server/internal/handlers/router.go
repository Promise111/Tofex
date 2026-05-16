package handlers

import (
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/tofex/backend/internal/config"
	"github.com/tofex/backend/internal/middleware"
	"github.com/tofex/backend/internal/permissions"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	s := NewServer(db, cfg)
	r.Use(middleware.CORS(os.Getenv("CORS_ORIGIN")))

	v1 := r.Group("/api/v1")
	{
		v1.POST("/bootstrap", s.Bootstrap)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", s.Login)
			auth.POST("/forgot-password", s.ForgotPassword)
			auth.POST("/reset-password", s.ResetPassword)
		}

		pub := v1.Group("/public")
		{
			pub.GET("/products", s.PublicListProducts)
			pub.GET("/products/:id", s.PublicGetProduct)
			pub.GET("/product-images/:id", s.PublicProductImageFile)
			pub.GET("/branches", s.PublicListStoreBranches)
			pub.GET("/payment-accounts", s.PublicListPaymentAccounts)
			pub.POST("/orders", s.PublicCreateOrder)
		}

		adm := v1.Group("/admin")
		adm.Use(middleware.Auth(cfg, db))
		{
			adm.GET("/me", s.Me)

			adm.GET("/users", middleware.RequirePermission(permissions.UsersRead), s.AdminListUsers)
			adm.GET("/users/:id", middleware.RequirePermission(permissions.UsersRead), s.AdminGetUser)
			adm.POST("/users", middleware.RequirePermission(permissions.UsersCreate), s.AdminCreateUser)
			adm.PATCH("/users/:id", middleware.RequirePermission(permissions.UsersUpdate), s.AdminPatchUser)
			adm.DELETE("/users/:id", middleware.RequirePermission(permissions.UsersDelete), s.AdminDeleteUser)

			adm.GET("/roles", middleware.RequirePermission(permissions.RolesRead), s.AdminListRoles)
			adm.PATCH("/roles/:id", middleware.RequirePermission(permissions.RolesUpdate), s.AdminPatchRole)

			adm.GET("/products", middleware.RequirePermission(permissions.ProductsRead), s.AdminListProducts)
			adm.POST("/products", middleware.RequirePermission(permissions.ProductsCreate), s.AdminCreateProduct)
			adm.PATCH("/products/:id", middleware.RequirePermission(permissions.ProductsUpdate), s.AdminPatchProduct)
			adm.DELETE("/products/:id", middleware.RequirePermission(permissions.ProductsDelete), s.AdminDeleteProduct)
			adm.POST("/products/:id/images", middleware.RequirePermission(permissions.ProductsUpdate), s.AdminUploadProductImage)

			adm.GET("/branches", middleware.RequirePermission(permissions.BranchesRead), s.AdminListStoreBranches)
			adm.POST("/branches", middleware.RequirePermission(permissions.BranchesCreate), s.AdminCreateStoreBranch)
			adm.PATCH("/branches/:id", middleware.RequirePermission(permissions.BranchesUpdate), s.AdminPatchStoreBranch)
			adm.DELETE("/branches/:id", middleware.RequirePermission(permissions.BranchesDelete), s.AdminDeleteStoreBranch)

			adm.GET("/payment-accounts", middleware.RequirePermission(permissions.PaymentAccountsRead), s.AdminListPaymentAccounts)
			adm.POST("/payment-accounts", middleware.RequirePermission(permissions.PaymentAccountsCreate), s.AdminCreatePaymentAccount)
			adm.PATCH("/payment-accounts/:id", middleware.RequirePermission(permissions.PaymentAccountsUpdate), s.AdminPatchPaymentAccount)
			adm.DELETE("/payment-accounts/:id", middleware.RequirePermission(permissions.PaymentAccountsDelete), s.AdminDeletePaymentAccount)

			adm.GET("/orders", middleware.RequirePermission(permissions.OrdersRead), s.AdminListOrders)
			adm.GET("/orders/:id", middleware.RequirePermission(permissions.OrdersRead), s.AdminGetOrder)
			adm.PATCH("/orders/:id", middleware.RequirePermission(permissions.OrdersUpdate), s.AdminPatchOrder)
			adm.GET("/orders/:id/receipts/:rid/file", middleware.RequirePermission(permissions.OrdersRead), s.AdminDownloadOrderReceipt)

			adm.GET("/audit-logs", middleware.RequirePermission(permissions.AuditRead), s.AdminListAuditLogs)
		}
	}

	r.GET("/healthz", s.Health)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
