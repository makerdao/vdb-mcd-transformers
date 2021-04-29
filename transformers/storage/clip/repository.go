package clip

import (
	"fmt"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	insertClipIlkQuery = `INSERT INTO maker.clip_ilk (diff_id, header_id, address_id, ilk_id) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
)

type StorageRepository struct {
	ContractAddress string
	db              *postgres.DB
}

func (repo *StorageRepository) Create(diffID, headerID int64, metadata types.ValueMetadata, value interface{}) error {
	switch metadata.Name {
	case storage.Ilk:
		return repo.insertIlk(diffID, headerID, value.(string))
	default:
		return fmt.Errorf("unrecognized clip contract storage name: %s", metadata.Name)
	}
}

func (repo *StorageRepository) SetDB(db *postgres.DB) {
	repo.db = db
}

func (repo *StorageRepository) insertIlk(diffID, headerID int64, ilk string) error {
	ilkID, ilkErr := shared.GetOrCreateIlk(ilk, repo.db)
	if ilkErr != nil {
		return fmt.Errorf("error getting or creating ilk for clip ilk: %w", ilkErr)
	}
	insertErr := shared.InsertRecordWithAddress(
		diffID,
		headerID,
		insertClipIlkQuery,
		strconv.FormatInt(ilkID, 10),
		repo.ContractAddress,
		repo.db)
	if insertErr != nil {
		msgToFormat := "error inserting clip %s ilk %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repo.ContractAddress, ilk, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}
