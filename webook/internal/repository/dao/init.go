package dao

import "gorm.io/gorm"

func InitTables(db *gorm.DB) error {
	// 理论上应该走db结构更改审批流程，这个不是优秀实践
	return db.AutoMigrate(&User{})
}
