package clip

import (
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type StorageRepository struct {
	ContractAddress string
	db              *postgres.DB
}

func (repo *StorageRepository) Create(diffID, headerID int64, metadata types.ValueMetadata, value interface{}) error {
	panic("implment me")
}

func (repo *StorageRepository) SetDB(db *postgres.DB) {
	repo.db = db
}
