package median

import (
	"fmt"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
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
	db                *postgres.DB
	ContractAddress   string
	ContractAddressID int64
}

func (repo *MedianStorageRepository) Create(diffID, headerID int64, metadata types.ValueMetadata, value interface{}) error {
	switch metadata.Name {
	case storage.Packed:
		return repo.insertPackedValueRecord(diffID, headerID, metadata, value.(map[int]string))
	case wards.Wards:
		return wards.InsertWards(diffID, headerID, metadata, repo.ContractAddress, value.(string), repo.db)
	case Bar:
		return repo.insertMedianBar(diffID, headerID, value.(string))
	case Bud:
		return repo.insertBud(diffID, headerID, metadata, value.(string))
	case Orcl:
		return repo.insertOrcl(diffID, headerID, metadata, value.(string))
	case Slot:
		return repo.insertSlot(diffID, headerID, metadata, value.(string))
	default:
		return fmt.Errorf("unrecognized median contract storage name: %s", metadata.Name)
	}
}
func (repo *MedianStorageRepository) SetDB(db *postgres.DB) {
	repo.db = db
}

func (repo MedianStorageRepository) insertMedianVal(diffID, headerID int64, val string) error {
	tx, txErr := repo.db.Beginx()
	if txErr != nil {
		return txErr
	}

	addressID, addressErr := repo.getMemoizedAddressId()
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
		return fmt.Errorf("error inserting median val for address %s: %w", repo.ContractAddress, insertErr)

	}
	return tx.Commit()
}

func (repo MedianStorageRepository) insertMedianAge(diffID, headerID int64, age string) error {
	tx, txErr := repo.db.Beginx()
	if txErr != nil {
		return txErr
	}

	addressID, addressErr := repo.getMemoizedAddressId()
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
		return fmt.Errorf("error inserting median age for address %s: %w", repo.ContractAddress, insertErr)

	}
	return tx.Commit()
}

func (repo *MedianStorageRepository) insertMedianBar(diffID, headerID int64, bar string) error {
	tx, txErr := repo.db.Beginx()
	if txErr != nil {
		return txErr
	}

	addressID, addressErr := repo.getMemoizedAddressId()
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

func (repo *MedianStorageRepository) insertPackedValueRecord(diffID, headerID int64, metadata types.ValueMetadata, packedValues map[int]string) error {
	for order, value := range packedValues {
		var insertErr error
		switch metadata.PackedNames[order] {
		case Val:
			insertErr = repo.insertMedianVal(diffID, headerID, value)
		case Age:
			insertErr = repo.insertMedianAge(diffID, headerID, value)
		default:
			return fmt.Errorf("unrecognized median contract storage name: in packed values: %s", metadata.PackedNames[order])
		}
		if insertErr != nil {
			return insertErr
		}
	}
	return nil
}

func (repo *MedianStorageRepository) insertOrcl(diffID, headerID int64, metadata types.ValueMetadata, orcl string) error {
	orclAddress, orclErr := getOrclAddress(metadata.Keys)
	if orclErr != nil {
		return orclErr
	}

	tx, txErr := repo.db.Beginx()
	if txErr != nil {
		return txErr
	}

	addressID, addressErr := repo.getMemoizedAddressId()
	if addressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("contract address", addressErr)
		}
		return addressErr
	}

	orclAddressID, orclAddressErr := repository.GetOrCreateAddress(repo.db, orclAddress)
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

func (repo *MedianStorageRepository) insertBud(diffID, headerID int64, metadata types.ValueMetadata, bud string) error {
	budAddress, budErr := getBudAddress(metadata.Keys)
	if budErr != nil {
		return budErr
	}

	tx, txErr := repo.db.Beginx()
	if txErr != nil {
		return txErr
	}

	addressID, addressErr := repo.getMemoizedAddressId()
	if addressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("contract address", addressErr)
		}
		return addressErr
	}

	budAddressID, budAddressErr := repository.GetOrCreateAddress(repo.db, budAddress)
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

func (repo *MedianStorageRepository) insertSlot(diffID, headerID int64, metadata types.ValueMetadata, slot string) error {
	slotID, slotErr := getSlotID(metadata.Keys)
	if slotErr != nil {
		return slotErr
	}

	tx, txErr := repo.db.Beginx()
	if txErr != nil {
		return txErr
	}

	addressID, addressErr := repo.getMemoizedAddressId()
	if addressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("contract address", addressErr)
		}
		return addressErr
	}

	slotAddressID, slotAddressErr := repository.GetOrCreateAddress(repo.db, slot)
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

func (repo *MedianStorageRepository) getMemoizedAddressId() (int64, error) {
	if repo.ContractAddressID == 0 {
		var addressID int64
		addressID, addressErr := repository.GetOrCreateAddress(repo.db, repo.ContractAddress)
		repo.ContractAddressID = addressID
		return addressID, addressErr
	}
	return repo.ContractAddressID, nil
}
