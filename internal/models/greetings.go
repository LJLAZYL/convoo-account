package models

import (
	"time"
)

// Greetings 是数据库中的模型结构，数据库实体（表结构映射）
type Greetings struct {
	ID         int64     `gorm:"primaryKey"`        // 自动增长的主键
	UserName   string    `gorm:"size:100;not null"` // 用户名，非空
	Greeting   string    `gorm:"size:255;not null"` // 问候语，非空
	ModifyTime time.Time `gorm:"autoUpdateTime"`    // 修改时间，自动更新时间
}
