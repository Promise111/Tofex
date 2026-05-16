package permissions

// Action keys used for RBAC and audit naming.
const (
	Wildcard = "*"

	UsersRead   = "users.read"
	UsersCreate = "users.create"
	UsersUpdate = "users.update"
	UsersDelete = "users.delete"

	RolesRead   = "roles.read"
	RolesUpdate = "roles.update"

	ProductsRead   = "products.read"
	ProductsCreate = "products.create"
	ProductsUpdate = "products.update"
	ProductsDelete = "products.delete"

	PaymentAccountsRead   = "payment_accounts.read"
	PaymentAccountsCreate = "payment_accounts.create"
	PaymentAccountsUpdate = "payment_accounts.update"
	PaymentAccountsDelete = "payment_accounts.delete"

	OrdersRead   = "orders.read"
	OrdersCreate = "orders.create"
	OrdersUpdate = "orders.update"

	AuditRead = "audit.read"
)

var All = []string{
	UsersRead, UsersCreate, UsersUpdate, UsersDelete,
	RolesRead, RolesUpdate,
	ProductsRead, ProductsCreate, ProductsUpdate, ProductsDelete,
	PaymentAccountsRead, PaymentAccountsCreate, PaymentAccountsUpdate, PaymentAccountsDelete,
	OrdersRead, OrdersCreate, OrdersUpdate,
	AuditRead,
}
