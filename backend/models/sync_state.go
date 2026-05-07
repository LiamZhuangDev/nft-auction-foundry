package models

import "time"

type SyncState struct {
	ID                     uint   `gorm:"primaryKey"` // always 1, the table only have one row
	LastProcessedBlock     uint64 `gorm:"not null"`
	LastProcessedBlockHash string `gorm:"type:varchar(66)"` // optional but useful for debugging / reorg safety
	CreatedAt              time.Time
	UpdatedAt              time.Time
}
