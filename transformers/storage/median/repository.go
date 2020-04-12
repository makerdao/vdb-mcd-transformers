package median

import (
	"fmt"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	insertMedianValQuery = `INSERT INTO maker.median_val (diff_id, header_id, val) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertMedianAgeQuery = `INSERT INTO maker.median_age (diff_id, header_id, age) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertMedianBarQuery = `INSERT INTO maker.median_bar (diff_id, header_id, bar) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertMedianBudQuery = `INSERT INTO maker.median_bud (diff_id, header_id, address_id, a, bud) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
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
	case Bud:
		return repository.insertBud(diffID, headerID, metadata, repository.ContractAddress, value.(string))
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

func (repository *MedianStorageRepository) insertBud(diffID, headerID int64, metadata types.ValueMetadata, contractAddress, bud string) error {
	budAddress, budErr := getBudAddress(metadata.Keys)
	if budErr != nil {
		return budErr
	}

	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return txErr
	}

	addressID, addressErr := shared.GetOrCreateAddress(contractAddress, repository.db)
	if addressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("contract address", addressErr.Error())
		}
		return addressErr
	}

	budAddressID, budAddressErr := shared.GetOrCreateAddress(budAddress, repository.db)
	if budAddressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("median bud address", budAddressErr.Error())
		}
		return addressErr
	}

	_, insertErr := tx.Exec(insertMedianBudQuery, diffID, headerID, addressID, budAddressID, bud)
	if insertErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("bud record with address", insertErr.Error())
		}
		return insertErr
	}
	return tx.Commit()
}

func getBudAddress(keys map[types.Key]string) (string, error) {
	user, ok := keys[constants.A]
	if !ok {
		return "", types.ErrMetadataMalformed{MissingData: constants.A}
	}
	return user, nil
}
