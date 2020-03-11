package medianizer

import (
	"fmt"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	insertMedianizerValQuery = `INSERT INTO maker.medianizer_val (diff_id, header_id, val) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertMedianizerAgeQuery = `INSERT INTO maker.medianizer_age (diff_id, header_id, age) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertMedianizerBarQuery = `INSERT INTO maker.medianizer_bar (diff_id, header_id, bar) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
)

type MedianizerStorageRepository struct {
	db              *postgres.DB
	ContractAddress string
}

func (repository MedianizerStorageRepository) Create(diffID, headerID int64, metadata types.ValueMetadata, value interface{}) error {
	switch metadata.Name {
	case storage.Packed:
		return repository.insertPackedValueRecord(diffID, headerID, metadata, value.(map[int]string))
	case Bar:
		return repository.insertMedianizerBar(diffID, headerID, value.(string))
	default:
		panic(fmt.Sprintf("unrecognized medianizer contract storage name: %s", metadata.Name))
	}
}
func (repository *MedianizerStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository MedianizerStorageRepository) insertMedianizerVal(diffID, headerID int64, val string) error {
	_, err := repository.db.Exec(insertMedianizerValQuery, diffID, headerID, val)
	return err
}

func (repository MedianizerStorageRepository) insertMedianizerAge(diffID, headerID int64, age string) error {
	_, err := repository.db.Exec(insertMedianizerAgeQuery, diffID, headerID, age)
	return err
}

func (repository MedianizerStorageRepository) insertMedianizerBar(diffID, headerID int64, bar string) error {
	_, err := repository.db.Exec(insertMedianizerBarQuery, diffID, headerID, bar)
	return err
}
func (repository *MedianizerStorageRepository) insertPackedValueRecord(diffID, headerID int64, metadata types.ValueMetadata, packedValues map[int]string) error {
	for order, value := range packedValues {
		var insertErr error
		switch metadata.PackedNames[order] {
		case Val:
			insertErr = repository.insertMedianizerVal(diffID, headerID, value)
		case Age:
			insertErr = repository.insertMedianizerAge(diffID, headerID, value)
		default:
			panic(fmt.Sprintf("unrecognized medianizer contract storage name in packed values: %s", metadata.Name))
		}
		if insertErr != nil {
			return insertErr
		}
	}
	return nil
}
