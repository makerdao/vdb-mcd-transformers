package cat

import (
	"fmt"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/utils"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	InsertCatIlkChopQuery = `INSERT INTO maker.cat_ilk_chop (diff_id, header_id, ilk_id, chop) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertCatIlkFlipQuery = `INSERT INTO maker.cat_ilk_flip (diff_id, header_id, ilk_id, flip) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertCatIlkLumpQuery = `INSERT INTO maker.cat_ilk_lump (diff_id, header_id, ilk_id, lump) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`

	insertCatLiveQuery = `INSERT INTO maker.cat_live (diff_id, header_id, live) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertCatVatQuery  = `INSERT INTO maker.cat_vat (diff_id, header_id, vat) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertCatVowQuery  = `INSERT INTO maker.cat_vow (diff_id, header_id, vow) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
)

type CatStorageRepository struct {
	db *postgres.DB
}

func (repository *CatStorageRepository) Create(diffID, headerID int64, metadata utils.StorageValueMetadata, value interface{}) error {
	switch metadata.Name {
	case Live:
		return repository.insertLive(diffID, headerID, value.(string))
	case Vat:
		return repository.insertVat(diffID, headerID, value.(string))
	case Vow:
		return repository.insertVow(diffID, headerID, value.(string))
	case IlkChop:
		return repository.insertIlkChop(diffID, headerID, metadata, value.(string))
	case IlkFlip:
		return repository.insertIlkFlip(diffID, headerID, metadata, value.(string))
	case IlkLump:
		return repository.insertIlkLump(diffID, headerID, metadata, value.(string))
	default:
		panic(fmt.Sprintf("unrecognized cat contract storage name: %s", metadata.Name))
	}
}

func (repository *CatStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository *CatStorageRepository) insertLive(diffID, headerID int64, live string) error {
	_, writeErr := repository.db.Exec(insertCatLiveQuery, diffID, headerID, live)
	return writeErr
}

func (repository *CatStorageRepository) insertVat(diffID, headerID int64, vat string) error {
	_, writeErr := repository.db.Exec(insertCatVatQuery, diffID, headerID, vat)
	return writeErr
}

func (repository *CatStorageRepository) insertVow(diffID, headerID int64, vow string) error {
	_, writeErr := repository.db.Exec(insertCatVowQuery, diffID, headerID, vow)
	return writeErr
}

// Ilks mapping: bytes32 => flip address; chop (ray), lump (wad) uint256
func (repository *CatStorageRepository) insertIlkFlip(diffID, headerID int64, metadata utils.StorageValueMetadata, flip string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertFieldWithIlk(diffID, headerID, ilk, IlkFlip, InsertCatIlkFlipQuery, flip)
}

func (repository *CatStorageRepository) insertIlkChop(diffID, headerID int64, metadata utils.StorageValueMetadata, chop string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertFieldWithIlk(diffID, headerID, ilk, IlkChop, InsertCatIlkChopQuery, chop)
}

func (repository *CatStorageRepository) insertIlkLump(diffID, headerID int64, metadata utils.StorageValueMetadata, lump string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertFieldWithIlk(diffID, headerID, ilk, IlkLump, InsertCatIlkLumpQuery, lump)
}

func (repository *CatStorageRepository) insertFieldWithIlk(diffID, headerID int64, ilk, variableName, query, value string) error {
	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return txErr
	}
	ilkID, ilkErr := shared.GetOrCreateIlkInTransaction(ilk, tx)
	if ilkErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("ilk", ilkErr.Error())
		}
		return ilkErr
	}
	_, writeErr := tx.Exec(query, diffID, headerID, ilkID, value)
	if writeErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError(variableName, writeErr.Error())
		}
		return writeErr
	}
	return tx.Commit()
}

func getIlk(keys map[utils.Key]string) (string, error) {
	ilk, ok := keys[constants.Ilk]
	if !ok {
		return "", utils.ErrMetadataMalformed{MissingData: constants.Ilk}
	}
	return ilk, nil
}

func getFlip(keys map[utils.Key]string) (string, error) {
	flip, ok := keys[constants.Flip]
	if !ok {
		return "", utils.ErrMetadataMalformed{MissingData: constants.Flip}
	}
	return flip, nil
}
