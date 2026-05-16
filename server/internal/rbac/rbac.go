package rbac

import (
	"github.com/tofex/backend/internal/models"
	"github.com/tofex/backend/internal/permissions"
)

func UserHasPermission(u *models.User, perm string) bool {
	if u == nil {
		return false
	}
	for _, r := range u.Roles {
		for _, p := range r.Permissions {
			if p.Permission == permissions.Wildcard {
				return true
			}
			if p.Permission == perm {
				return true
			}
		}
	}
	return false
}

func IsSuperAdmin(u *models.User) bool {
	if u == nil {
		return false
	}
	for _, r := range u.Roles {
		if r.Name == "super_admin" {
			return true
		}
	}
	return false
}

func RolesIncludeSuperAdmin(roles []models.Role) bool {
	for _, r := range roles {
		if r.Name == "super_admin" {
			return true
		}
	}
	return false
}
