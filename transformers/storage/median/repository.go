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
	insertMedianValQuery  = `INSERT INTO maker.median_val (diff_id, header_id, address_id,  val) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertMedianAgeQuery  = `INSERT INTO maker.median_age (diff_id, header_id, address_id, age) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertMedianBarQuery  = `INSERT INTO maker.median_bar (diff_id, header_id, address_id, bar) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertMedianBudQuery  = `INSERT INTO maker.median_bud (diff_id, header_id, address_id, a, bud) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	insertMedianOrclQuery = `INSERT INTO maker.median_orcl (diff_id, header_id, address_id, a, orcl) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	insertMedianSlotQuery = `INSERT INTO maker.median_slot (diff_id, header_id, address_id, slot_id, slot) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
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
		return repository.insertBud(diffID, headerID, metadata, value.(string))
	case Orcl:
		return repository.insertOrcl(diffID, headerID, metadata, value.(string))
	case Slot:
		return repository.insertSlot(diffID, headerID, metadata, value.(string))
	default:
		return fmt.Errorf("unrecognized median contract storage name: %s", metadata.Name)
	}
}
func (repository *MedianStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository MedianStorageRepository) insertMedianVal(diffID, headerID int64, contractAddress, val string) error {
	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return txErr
	}

	addressID, addressErr := shared.GetOrCreateAddress(contractAddress, repository.db)
	if addressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("contract address", addressErr)
		}
		return addressErr
	}
	_, insertErr := tx.Exec(insertMedianValQuery, diffID, headerID, addressID, val)
	if insertErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("val record with address", insertErr)
		}
		return fmt.Errorf("error inserting median val for address %s: %w", contractAddress, insertErr)

	}
	return tx.Commit()
}

func (repository MedianStorageRepository) insertMedianAge(diffID, headerID int64, contractAddress, age string) error {
	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return txErr
	}

	addressID, addressErr := shared.GetOrCreateAddress(contractAddress, repository.db)
	if addressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("contract address", addressErr)
		}
		return addressErr
	}
	_, insertErr := tx.Exec(insertMedianAgeQuery, diffID, headerID, addressID, age)
	if insertErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("age record with address", insertErr)
		}
		return fmt.Errorf("error inserting median age for address %s: %w", contractAddress, insertErr)

	}
	return tx.Commit()
}

func (repository MedianStorageRepository) insertMedianBar(diffID, headerID int64, bar string) error {
	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return txErr
	}

	addressID, addressErr := shared.GetOrCreateAddress(repository.ContractAddress, repository.db)
	if addressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("contract address", addressErr)
		}
		return addressErr
	}
	_, insertErr := tx.Exec(insertMedianBarQuery, diffID, headerID, addressID, bar)
	if insertErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("bar record with address", insertErr)
		}
		return fmt.Errorf("error inserting median bar for address %s: %w", bar, insertErr)

	}
	return tx.Commit()
}

func (repository *MedianStorageRepository) insertPackedValueRecord(diffID, headerID int64, metadata types.ValueMetadata, packedValues map[int]string) error {
	for order, value := range packedValues {
		var insertErr error
		switch metadata.PackedNames[order] {
		case Val:
			insertErr = repository.insertMedianVal(diffID, headerID, repository.ContractAddress, value)
		case Age:
			insertErr = repository.insertMedianAge(diffID, headerID, repository.ContractAddress, value)
		default:
			return fmt.Errorf("unrecognized median contract storage name: in packed values: %s", metadata.PackedNames[order])
		}
		if insertErr != nil {
			return insertErr
		}
	}
	return nil
}

func (repository *MedianStorageRepository) insertOrcl(diffID, headerID int64, metadata types.ValueMetadata, orcl string) error {
	orclAddress, orclErr := getOrclAddress(metadata.Keys)
	if orclErr != nil {
		return orclErr
	}

	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return txErr
	}

	addressID, addressErr := shared.GetOrCreateAddress(repository.ContractAddress, repository.db)
	if addressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("contract address", addressErr)
		}
		return addressErr
	}

	orclAddressID, orclAddressErr := shared.GetOrCreateAddress(orclAddress, repository.db)
	if orclAddressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("median orcl address", orclAddressErr)
		}
		return addressErr
	}

	_, insertErr := tx.Exec(insertMedianOrclQuery, diffID, headerID, addressID, orclAddressID, orcl)
	if insertErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("bud record with address", insertErr)
		}
		return fmt.Errorf("error inserting median orcl for address %s: %w", orclAddress, insertErr)

	}
	return tx.Commit()
}

func getOrclAddress(keys map[types.Key]string) (string, error) {
	user, ok := keys[constants.Address]
	if !ok {
		return "", types.ErrMetadataMalformed{MissingData: constants.A}
	}
	return user, nil
}

func (repository *MedianStorageRepository) insertBud(diffID, headerID int64, metadata types.ValueMetadata, bud string) error {
	budAddress, budErr := getBudAddress(metadata.Keys)
	if budErr != nil {
		return budErr
	}

	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return txErr
	}

	addressID, addressErr := shared.GetOrCreateAddress(repository.ContractAddress, repository.db)
	if addressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("contract address", addressErr)
		}
		return addressErr
	}

	budAddressID, budAddressErr := shared.GetOrCreateAddress(budAddress, repository.db)
	if budAddressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("median bud address", budAddressErr)
		}
		return addressErr
	}

	_, insertErr := tx.Exec(insertMedianBudQuery, diffID, headerID, addressID, budAddressID, bud)
	if insertErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("bud record with address", insertErr)
		}
		return fmt.Errorf("error inserting median bud for address %s: %w", budAddress, insertErr)
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

func (repository *MedianStorageRepository) insertSlot(diffID, headerID int64, metadata types.ValueMetadata, slot string) error {
	slotID, slotErr := getSlotID(metadata.Keys)
	if slotErr != nil {
		return slotErr
	}

	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return txErr
	}

	addressID, addressErr := shared.GetOrCreateAddress(repository.ContractAddress, repository.db)
	if addressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("contract address", addressErr)
		}
		return addressErr
	}

	slotAddressID, slotAddressErr := shared.GetOrCreateAddress(slot, repository.db)
	if slotAddressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("median slot address", slotAddressErr)
		}
		return slotAddressErr
	}

	_, insertErr := tx.Exec(insertMedianSlotQuery, diffID, headerID, addressID, slotID, slotAddressID)
	if insertErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("slot record with address", insertErr)
		}
		return insertErr
	}
	return tx.Commit()
}

func getSlotID(keys map[types.Key]string) (string, error) {
	slotIDstr, ok := keys[constants.SlotId]
	if !ok {
		return "", types.ErrMetadataMalformed{MissingData: constants.SlotId}
	}
	return slotIDstr, nil
}
