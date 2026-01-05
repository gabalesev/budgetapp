package common

import (
	"gorm.io/gorm"
)

func WhereUserIDScope(userID uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userID)
	}
}
