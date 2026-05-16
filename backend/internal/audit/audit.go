package audit

import (
	"encoding/json"

	"github.com/tofex/backend/internal/models"
	"gorm.io/gorm"
)

func Log(db *gorm.DB, userID *uint, action, resourceType, resourceID string, metadata map[string]any, ip, ua string) {
	meta := ""
	if metadata != nil {
		if b, err := json.Marshal(metadata); err == nil {
			meta = string(b)
		}
	}
	_ = db.Create(&models.AuditLog{
		UserID:       userID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		MetadataJSON: meta,
		IP:           ip,
		UserAgent:    ua,
	}).Error
}
