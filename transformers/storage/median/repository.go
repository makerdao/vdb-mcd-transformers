package median

import (
	"fmt"

	"github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	insertMedianValQuery = `INSERT INTO maker.median_val (diff_id, header_id, val) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertMedianAgeQuery = `INSERT INTO maker.median_age (diff_id, header_id, age) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertMedianBarQuery = `INSERT INTO maker.median_bar (diff_id, header_id, bar) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
)

type MedianStorageRepository struct {
	db              *postgres.DB
	ContractAddress string
}

func (repository MedianStorageRepository) Create(diffID, headerID int64, metadata types.ValueMetadata, value interface{}) error {
	switch metadata.Name {
	case storage.Packed:
		return repository.insertPackedValueRecord(diffID, headerID, metadata, value.(map[int]string))
	case wards.Wards:
		return wards.InsertWards(diffID, headerID, metadata, repository.ContractAddress, value.(string), repository.db)
	case Bar:
		return repository.insertMedianBar(diffID, headerID, value.(string))
	default:
		panic(fmt.Sprintf("unrecognized median contract storage name: %s", metadata.Name))
	}
}
func (repository *MedianStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository MedianStorageRepository) insertMedianVal(diffID, headerID int64, val string) error {
	_, err := repository.db.Exec(insertMedianValQuery, diffID, headerID, val)
	return err
}

func (repository MedianStorageRepository) insertMedianAge(diffID, headerID int64, age string) error {
	_, err := repository.db.Exec(insertMedianAgeQuery, diffID, headerID, age)
	return err
}

func (repository MedianStorageRepository) insertMedianBar(diffID, headerID int64, bar string) error {
	_, err := repository.db.Exec(insertMedianBarQuery, diffID, headerID, bar)
	return err
}
func (repository *MedianStorageRepository) insertPackedValueRecord(diffID, headerID int64, metadata types.ValueMetadata, packedValues map[int]string) error {
	for order, value := range packedValues {
		var insertErr error
		switch metadata.PackedNames[order] {
		case Val:
			insertErr = repository.insertMedianVal(diffID, headerID, value)
		case Age:
			insertErr = repository.insertMedianAge(diffID, headerID, value)
		default:
			panic(fmt.Sprintf("unrecognized median contract storage name in packed values: %s", metadata.Name))
		}
		if insertErr != nil {
			return insertErr
		}
	}
	return nil
}
