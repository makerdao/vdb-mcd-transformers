package clip

import (
	"fmt"

	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type StorageRepository struct {
	ContractAddress string
	db              *postgres.DB
}

func (repo *StorageRepository) Create(diffID, headerID int64, metadata types.ValueMetadata, value interface{}) error {
	switch metadata.Name {
	default:
		return fmt.Errorf("unrecognized cat contract storage name: %s", metadata.Name)
	}
}

func (repo *StorageRepository) SetDB(db *postgres.DB) {
	repo.db = db
}
