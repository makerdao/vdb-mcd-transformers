package clip

import (
	"fmt"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	insertClipDogQuery = `INSERT INTO maker.clip_dog (diff_id, header_id, address_id, dog) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
)

type StorageRepository struct {
	ContractAddress string
	db              *postgres.DB
}

func (repo *StorageRepository) Create(diffID, headerID int64, metadata types.ValueMetadata, value interface{}) error {
	switch metadata.Name {
	case storage.Dog:
		return repo.insertDog(diffID, headerID, value.(string))
	default:
		return fmt.Errorf("unrecognized clip contract storage name: %s", metadata.Name)
	}
}

func (repo *StorageRepository) SetDB(db *postgres.DB) {
	repo.db = db
}

func (repo *StorageRepository) insertDog(diffID, headerID int64, dog string) error {
	dogAddressID, addressErr := repository.GetOrCreateAddress(repo.db, dog)
	if addressErr != nil {
		return fmt.Errorf("error inserting clip dog: %w", addressErr)
	}
	insertErr := shared.InsertRecordWithAddress(
		diffID,
		headerID,
		insertClipDogQuery,
		strconv.FormatInt(dogAddressID, 10),
		repo.ContractAddress,
		repo.db)
	if insertErr != nil {
		msgToFormat := "error inserting clip %s dog %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repo.ContractAddress, dog, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}
