// Package apidocs holds types referenced by Swagger annotations.
package apidocs

import "github.com/tofex/backend/internal/models"

// ErrorResponse is the standard API error envelope.
type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

// ErrorBody describes an error.
type ErrorBody struct {
	Code    string `json:"code" example:"validation_error"`
	Message string `json:"message" example:"invalid request"`
}

// HealthResponse health check payload.
type HealthResponse struct {
	OK bool `json:"ok" example:"true"`
}

// BootstrapRequest creates the first super admin.
type BootstrapRequest struct {
	Email       string `json:"email" example:"admin@tofex.com"`
	Username    string `json:"username" example:"superadmin"`
	Password    string `json:"password" example:"changeme1234"`
	DisplayName string `json:"display_name" example:"Super Admin"`
}

// MessageResponse generic message payload.
type MessageResponse struct {
	Message string `json:"message" example:"ok"`
}

// BootstrapResponse bootstrap success.
type BootstrapResponse struct {
	Message string `json:"message" example:"first admin created"`
	UserID  uint   `json:"user_id" example:"1"`
}

// LoginRequest staff login.
type LoginRequest struct {
	Login    string `json:"login" example:"admin@tofex.com"`
	Password string `json:"password" example:"changeme1234"`
}

// LoginResponse JWT login success.
type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type" example:"Bearer"`
	ExpiresIn   int          `json:"expires_in" example:"86400"`
	User        models.User  `json:"user"`
}

// ForgotPasswordRequest request reset email/token.
type ForgotPasswordRequest struct {
	Email string `json:"email" example:"admin@tofex.com"`
}

// ForgotPasswordDevResponse dev-only when LOG_PASSWORD_RESET_LINK=true.
type ForgotPasswordDevResponse struct {
	Message    string `json:"message"`
	ResetToken string `json:"reset_token"`
	ExpiresIn  int    `json:"expires_in"`
}

// ResetPasswordRequest complete password reset.
type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password" example:"newpassword1234"`
}

// CreateUserRequest admin creates staff user.
type CreateUserRequest struct {
	Email       string `json:"email" example:"staff@tofex.com"`
	Username    string `json:"username" example:"staff1"`
	Password    string `json:"password" example:"changeme1234"`
	DisplayName string `json:"display_name" example:"Staff User"`
	RoleIDs     []uint `json:"role_ids" example:"2"`
}

// PatchUserRequest partial user update.
type PatchUserRequest struct {
	DisplayName *string `json:"display_name"`
	Password    *string `json:"password"`
	RoleIDs     *[]uint `json:"role_ids"`
}

// CreateUserResponse user created.
type CreateUserResponse struct {
	UserID uint `json:"user_id" example:"2"`
}

// UsersListResponse paginated users.
type UsersListResponse struct {
	Users   []models.User `json:"users"`
	Total   int64         `json:"total"`
	Page    int           `json:"page"`
	PerPage int           `json:"per_page"`
}

// PatchRoleRequest update role permissions.
type PatchRoleRequest struct {
	Permissions []string `json:"permissions" example:"products.read,orders.read"`
}

// RolesListResponse role list.
type RolesListResponse struct {
	Roles []models.Role `json:"roles"`
}

// CreateProductRequest new product.
type CreateProductRequest struct {
	Name        string `json:"name" example:"Chocolate cake"`
	Slug        string `json:"slug" example:"chocolate-cake"`
	Description string `json:"description" example:"Rich cocoa layers"`
	PriceCents  int64  `json:"price_cents" example:"450000"`
	Active      *bool  `json:"active" example:"true"`
}

// PatchProductRequest partial product update.
type PatchProductRequest struct {
	Name        *string `json:"name"`
	Slug        *string `json:"slug"`
	Description *string `json:"description"`
	PriceCents  *int64  `json:"price_cents"`
	Active      *bool   `json:"active"`
}

// ProductsListResponse product list wrapper.
type ProductsListResponse struct {
	Products []models.Product `json:"products"`
}

// PublicProductsListResponse public catalog.
type PublicProductsListResponse struct {
	Products []models.Product `json:"products"`
}

// CreatePaymentAccountRequest new bank account for customers.
type CreatePaymentAccountRequest struct {
	BankName      string `json:"bank_name" example:"Example Bank"`
	AccountName   string `json:"account_name" example:"Tofex Ltd"`
	AccountNumber string `json:"account_number" example:"0123456789"`
	Branch        string `json:"branch" example:"Lagos Main"`
	Currency      string `json:"currency" example:"NGN"`
	DisplayLabel  string `json:"display_label" example:"Main account"`
	Active        *bool  `json:"active" example:"true"`
}

// PatchPaymentAccountRequest partial payment account update.
type PatchPaymentAccountRequest struct {
	BankName      *string `json:"bank_name"`
	AccountName   *string `json:"account_name"`
	AccountNumber *string `json:"account_number"`
	Branch        *string `json:"branch"`
	Currency      *string `json:"currency"`
	DisplayLabel  *string `json:"display_label"`
	Active        *bool   `json:"active"`
}

// PaymentAccountsListResponse payment accounts wrapper.
type PaymentAccountsListResponse struct {
	PaymentAccounts []models.PaymentAccount `json:"payment_accounts"`
}

// PatchOrderRequest update order status.
type PatchOrderRequest struct {
	Status string `json:"status" example:"confirmed" enums:"pending,confirmed,ready,completed,cancelled"`
}

// OrdersListResponse paginated orders.
type OrdersListResponse struct {
	Orders  []models.Order `json:"orders"`
	Total   int64          `json:"total"`
	Page    int            `json:"page"`
	PerPage int            `json:"per_page"`
}

// AuditLogsListResponse paginated audit logs.
type AuditLogsListResponse struct {
	AuditLogs []models.AuditLog `json:"audit_logs"`
	Total     int64             `json:"total"`
	Page      int               `json:"page"`
	PerPage   int               `json:"per_page"`
}

// PublicCreateOrderResponse guest order created.
type PublicCreateOrderResponse struct {
	OrderID    string `json:"order_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	TotalCents int64  `json:"total_cents" example:"900000"`
	Status     string `json:"status" example:"pending"`
}

// OrderLineItem one line in the items JSON for guest checkout.
type OrderLineItem struct {
	ProductID uint `json:"product_id" example:"1"`
	Quantity  int  `json:"quantity" example:"2"`
}
