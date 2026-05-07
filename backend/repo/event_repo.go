package repo

import (
	"context"
	"errors"
	"nft-auction-backend/models"

	"gorm.io/gorm"
)

type EventRepo struct {
	DB *gorm.DB
}

func NewEventRepo(db *gorm.DB) *EventRepo {
	return &EventRepo{DB: db}
}

func (r *EventRepo) GetLastProcessedBlock(ctx context.Context) (uint64, error) {
	var state models.SyncState

	err := r.DB.First(&state, "id = ?", 1).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {

		state = models.SyncState{
			ID:                 1,
			LastProcessedBlock: 0,
		}

		if err := r.DB.Create(&state).Error; err != nil {
			return 0, err
		}

		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return state.LastProcessedBlock, nil
}

func (r *EventRepo) SetLastProcessedBlock(ctx context.Context, block uint64) error {
	return r.DB.Model(&models.SyncState{}).
		Where("id = ?", 1).
		Update("last_processed_block", block).Error
}

func (r *EventRepo) EventExists(ctx context.Context, txHash string, logIndex uint) (bool, error) {
	var count int64

	err := r.DB.Model(&models.ProcessedEvent{}).
		Where("tx_hash = ? AND log_index = ?", txHash, logIndex).
		Count(&count).Error

	return count > 0, err
}

func (r *EventRepo) MarkEventProcessed(ctx context.Context, txHash string, logIndex uint) error {
	return r.DB.Create(&models.ProcessedEvent{
		TxHash:   txHash,
		LogIndex: logIndex,
	}).Error
}
