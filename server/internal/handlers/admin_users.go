package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tofex/backend/internal/audit"
	"github.com/tofex/backend/internal/httpx"
	"github.com/tofex/backend/internal/middleware"
	"github.com/tofex/backend/internal/models"
	"github.com/tofex/backend/internal/rbac"
	"github.com/tofex/backend/internal/security/password"
)

// Me returns the authenticated staff profile.
//
//	@Summary		Current user
//	@Tags			admin
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{object}	models.User
//	@Failure		401	{object}	apidocs.ErrorResponse
//	@Failure		500	{object}	apidocs.ErrorResponse
//	@Router			/api/v1/admin/me [get]
func (s *Server) Me(c *gin.Context) {
	u := middleware.CurrentUser(c)
	if u == nil {
		httpx.JSONError(c, 401, "unauthorized", "missing user")
		return
	}
	if err := s.DB.Preload("Roles.Permissions").First(u, u.ID).Error; err != nil {
		httpx.JSONError(c, 500, "server_error", "could not load user")
		return
	}
	c.JSON(http.StatusOK, u)
}

// AdminListUsers lists staff users (paginated).
//
//	@Summary		List users
//	@Tags			admin-users
//	@Security		BearerAuth
//	@Produce		json
//	@Param			page		query	int	false	"Page"		default(1)
//	@Param			per_page	query	int	false	"Per page"	default(20)
//	@Success		200			{object}	apidocs.UsersListResponse
//	@Failure		401			{object}	apidocs.ErrorResponse
//	@Failure		403			{object}	apidocs.ErrorResponse
//	@Failure		500			{object}	apidocs.ErrorResponse
//	@Router			/api/v1/admin/users [get]
func (s *Server) AdminListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	per, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	if page < 1 {
		page = 1
	}
	if per < 1 || per > 100 {
		per = 20
	}
	offset := (page - 1) * per
	var total int64
	_ = s.DB.Model(&models.User{}).Count(&total).Error
	var users []models.User
	if err := s.DB.Preload("Roles.Permissions").Order("id asc").Limit(per).Offset(offset).Find(&users).Error; err != nil {
		httpx.JSONError(c, 500, "server_error", "could not list users")
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users, "total": total, "page": page, "per_page": per})
}

// AdminGetUser returns one staff user.
//
//	@Summary		Get user
//	@Tags			admin-users
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path	int	true	"User ID"
//	@Success		200	{object}	models.User
//	@Failure		400	{object}	apidocs.ErrorResponse
//	@Failure		401	{object}	apidocs.ErrorResponse
//	@Failure		403	{object}	apidocs.ErrorResponse
//	@Failure		404	{object}	apidocs.ErrorResponse
//	@Router			/api/v1/admin/users/{id} [get]
func (s *Server) AdminGetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpx.JSONError(c, 400, "validation_error", "invalid id")
		return
	}
	var u models.User
	if err := s.DB.Preload("Roles.Permissions").First(&u, id).Error; err != nil {
		httpx.JSONError(c, 404, "not_found", "user not found")
		return
	}
	c.JSON(http.StatusOK, u)
}

type createUserBody struct {
	Email       string `json:"email" binding:"required,email"`
	Username    string `json:"username" binding:"required,min=2,max=64"`
	Password    string `json:"password" binding:"required,min=10,max=128"`
	DisplayName string `json:"display_name" binding:"max=255"`
	RoleIDs     []uint `json:"role_ids" binding:"required"`
}

// AdminCreateUser creates a staff user (no public signup).
//
//	@Summary		Create user
//	@Tags			admin-users
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		apidocs.CreateUserRequest	true	"New user"
//	@Success		201		{object}	apidocs.CreateUserResponse
//	@Failure		400		{object}	apidocs.ErrorResponse
//	@Failure		401		{object}	apidocs.ErrorResponse
//	@Failure		403		{object}	apidocs.ErrorResponse
//	@Failure		500		{object}	apidocs.ErrorResponse
//	@Router			/api/v1/admin/users [post]
func (s *Server) AdminCreateUser(c *gin.Context) {
	cur := middleware.CurrentUser(c)
	var body createUserBody
	if err := c.ShouldBindJSON(&body); err != nil {
		httpx.JSONError(c, 400, "validation_error", err.Error())
		return
	}
	if len(body.RoleIDs) == 0 {
		httpx.JSONError(c, 400, "validation_error", "role_ids must include at least one role")
		return
	}
	var roles []models.Role
	if err := s.DB.Where("id IN ?", body.RoleIDs).Find(&roles).Error; err != nil || len(roles) != len(body.RoleIDs) {
		httpx.JSONError(c, 400, "validation_error", "one or more role_ids are invalid")
		return
	}
	if rbac.RolesIncludeSuperAdmin(roles) && !rbac.IsSuperAdmin(cur) {
		httpx.JSONError(c, 403, "forbidden", "only a super admin can assign the super_admin role")
		return
	}
	hash, err := password.Hash(body.Password)
	if err != nil {
		httpx.JSONError(c, 500, "server_error", "password hashing failed")
		return
	}
	u := models.User{
		Email:        body.Email,
		Username:     body.Username,
		PasswordHash: hash,
		DisplayName:  body.DisplayName,
	}
	if err := s.DB.Create(&u).Error; err != nil {
		httpx.JSONError(c, 400, "create_failed", "could not create user (unique email/username)")
		return
	}
	if err := s.DB.Model(&u).Association("Roles").Replace(roles); err != nil {
		httpx.JSONError(c, 500, "server_error", "could not assign roles")
		return
	}
	audit.Log(s.DB, &cur.ID, "user.create", "user", itoa(u.ID), map[string]any{"email": u.Email}, httpx.ClientIP(c), c.Request.UserAgent())
	c.JSON(http.StatusCreated, gin.H{"user_id": u.ID})
}

type patchUserBody struct {
	DisplayName *string `json:"display_name"`
	Password    *string `json:"password"`
	RoleIDs     *[]uint `json:"role_ids"`
}

// AdminPatchUser updates a staff user.
//
//	@Summary		Update user
//	@Tags			admin-users
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path	int						true	"User ID"
//	@Param			body	body	apidocs.PatchUserRequest	true	"Fields to update"
//	@Success		200		{object}	apidocs.MessageResponse
//	@Failure		400		{object}	apidocs.ErrorResponse
//	@Failure		401		{object}	apidocs.ErrorResponse
//	@Failure		403		{object}	apidocs.ErrorResponse
//	@Failure		404		{object}	apidocs.ErrorResponse
//	@Failure		500		{object}	apidocs.ErrorResponse
//	@Router			/api/v1/admin/users/{id} [patch]
func (s *Server) AdminPatchUser(c *gin.Context) {
	cur := middleware.CurrentUser(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpx.JSONError(c, 400, "validation_error", "invalid id")
		return
	}
	var u models.User
	if err := s.DB.Preload("Roles").First(&u, id).Error; err != nil {
		httpx.JSONError(c, 404, "not_found", "user not found")
		return
	}
	var body patchUserBody
	if err := c.ShouldBindJSON(&body); err != nil {
		httpx.JSONError(c, 400, "validation_error", err.Error())
		return
	}
	if body.RoleIDs != nil {
		if len(*body.RoleIDs) == 0 {
			httpx.JSONError(c, 400, "validation_error", "role_ids must include at least one role")
			return
		}
		var roles []models.Role
		if err := s.DB.Where("id IN ?", *body.RoleIDs).Find(&roles).Error; err != nil || len(roles) != len(*body.RoleIDs) {
			httpx.JSONError(c, 400, "validation_error", "one or more role_ids are invalid")
			return
		}
		if rbac.RolesIncludeSuperAdmin(roles) && !rbac.IsSuperAdmin(cur) {
			httpx.JSONError(c, 403, "forbidden", "only a super admin can assign the super_admin role")
			return
		}
		if err := s.DB.Model(&u).Association("Roles").Replace(roles); err != nil {
			httpx.JSONError(c, 500, "server_error", "could not update roles")
			return
		}
	}
	if body.DisplayName != nil {
		u.DisplayName = *body.DisplayName
	}
	if body.Password != nil {
		if len(*body.Password) < 10 {
			httpx.JSONError(c, 400, "validation_error", "password too short")
			return
		}
		hash, err := password.Hash(*body.Password)
		if err != nil {
			httpx.JSONError(c, 500, "server_error", "password hashing failed")
			return
		}
		u.PasswordHash = hash
	}
	if err := s.DB.Save(&u).Error; err != nil {
		httpx.JSONError(c, 500, "server_error", "could not update user")
		return
	}
	audit.Log(s.DB, &cur.ID, "user.update", "user", itoa(u.ID), nil, httpx.ClientIP(c), c.Request.UserAgent())
	httpx.JSONOK(c, gin.H{"message": "updated"})
}

// AdminDeleteUser soft-deletes a staff user.
//
//	@Summary		Delete user
//	@Tags			admin-users
//	@Security		BearerAuth
//	@Param			id	path	int	true	"User ID"
//	@Success		204	"no content"
//	@Failure		400	{object}	apidocs.ErrorResponse
//	@Failure		401	{object}	apidocs.ErrorResponse
//	@Failure		403	{object}	apidocs.ErrorResponse
//	@Failure		404	{object}	apidocs.ErrorResponse
//	@Failure		500	{object}	apidocs.ErrorResponse
//	@Router			/api/v1/admin/users/{id} [delete]
func (s *Server) AdminDeleteUser(c *gin.Context) {
	cur := middleware.CurrentUser(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpx.JSONError(c, 400, "validation_error", "invalid id")
		return
	}
	if uint(id) == cur.ID {
		httpx.JSONError(c, 400, "validation_error", "cannot delete your own account")
		return
	}
	res := s.DB.Delete(&models.User{}, id)
	if res.Error != nil {
		httpx.JSONError(c, 500, "server_error", "delete failed")
		return
	}
	if res.RowsAffected == 0 {
		httpx.JSONError(c, 404, "not_found", "user not found")
		return
	}
	audit.Log(s.DB, &cur.ID, "user.delete", "user", strconv.FormatUint(id, 10), nil, httpx.ClientIP(c), c.Request.UserAgent())
	c.Status(http.StatusNoContent)
}