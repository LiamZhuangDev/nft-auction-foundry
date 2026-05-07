package models

import "time"

type ProcessedEvent struct {
	TxHash      string `gorm:"primaryKey;type:varchar(66)"`
	LogIndex    uint   `gorm:"primaryKey"`
	EventType   string `gorm:"type:varchar(50)"` // optional but nice for debugging
	BlockNumber uint64 `gorm:"index"`            // optional but useful
	CreatedAt   time.Time
}
