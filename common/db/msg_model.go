package db

import (
	"context"

	"github.com/n8sPxD/cowIM/common/db/models"
)

func (db *DB) InsertSingleMessage(ctx context.Context, msg *models.SingleMessage) error {
	db.rwMutex.Lock()
	defer db.rwMutex.Unlock()
	return db.client.WithContext(ctx).Create(msg).Error
}
