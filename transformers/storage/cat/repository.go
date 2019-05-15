package cat

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

const (
	InsertCatIlkChopQuery = `INSERT INTO maker.cat_ilk_chop (block_number, block_hash, ilk_id, chop) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertCatIlkFlipQuery = `INSERT INTO maker.cat_ilk_flip (block_number, block_hash, ilk_id, flip) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertCatIlkLumpQuery = `INSERT INTO maker.cat_ilk_lump (block_number, block_hash, ilk_id, lump) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`

	insertCatVowQuery     = `INSERT INTO maker.cat_vow (block_number, block_hash, vow) VALUES ($1, $2, $3 ) ON CONFLICT DO NOTHING`
	insertCatFlipIlkQuery = `INSERT INTO maker.cat_flip_ilk (block_number, block_hash, ilk_id, flip) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertCatNFlipQuery   = `INSERT INTO maker.cat_nflip (block_number, block_hash, nflip) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertCatLiveQuery    = `INSERT INTO maker.cat_live (block_number, block_hash, live) VALUES ($1, $2, $3 ) ON CONFLICT DO NOTHING`
	insertCatFlipUrnQuery = `INSERT INTO maker.cat_flip_urn (block_number, block_hash, flip, urn) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertCatVatQuery     = `INSERT INTO maker.cat_vat (block_number, block_hash, vat) VALUES ($1, $2, $3 ) ON CONFLICT DO NOTHING`
	insertCatFlipInkQuery = `INSERT INTO maker.cat_flip_ink (block_number, block_hash, flip, ink) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertCatFlipTabQuery = `INSERT INTO maker.cat_flip_tab (block_number, block_hash, flip, tab) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
)

type CatStorageRepository struct {
	db *postgres.DB
}

func (repository *CatStorageRepository) Create(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, value interface{}) error {
	switch metadata.Name {
	case NFlip:
		return repository.insertNFlip(blockNumber, blockHash, value.(string))
	case Live:
		return repository.insertLive(blockNumber, blockHash, value.(string))
	case Vat:
		return repository.insertVat(blockNumber, blockHash, value.(string))
	case Vow:
		return repository.insertVow(blockNumber, blockHash, value.(string))
	case IlkFlip:
		return repository.insertIlkFlip(blockNumber, blockHash, metadata, value.(string))
	case IlkChop:
		return repository.insertIlkChop(blockNumber, blockHash, metadata, value.(string))
	case IlkLump:
		return repository.insertIlkLump(blockNumber, blockHash, metadata, value.(string))
	case FlipIlk:
		return repository.insertFlipIlk(blockNumber, blockHash, metadata, value.(string))
	case FlipUrn:
		return repository.insertFlipUrn(blockNumber, blockHash, metadata, value.(string))
	case FlipInk:
		return repository.insertFlipInk(blockNumber, blockHash, metadata, value.(string))
	case FlipTab:
		return repository.insertFlipTab(blockNumber, blockHash, metadata, value.(string))
	default:
		panic(fmt.Sprintf("unrecognized cat contract storage name: %s", metadata.Name))
	}
}

func (repository *CatStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository *CatStorageRepository) insertNFlip(blockNumber int, blockHash string, nflip string) error {
	_, writeErr := repository.db.Exec(insertCatNFlipQuery, blockNumber, blockHash, nflip)
	return writeErr
}

func (repository *CatStorageRepository) insertLive(blockNumber int, blockHash string, live string) error {
	_, writeErr := repository.db.Exec(insertCatLiveQuery, blockNumber, blockHash, live)
	return writeErr
}

func (repository *CatStorageRepository) insertVat(blockNumber int, blockHash string, vat string) error {
	_, writeErr := repository.db.Exec(insertCatVatQuery, blockNumber, blockHash, vat)
	return writeErr
}

func (repository *CatStorageRepository) insertVow(blockNumber int, blockHash string, vow string) error {
	_, writeErr := repository.db.Exec(insertCatVowQuery, blockNumber, blockHash, vow)
	return writeErr
}

// Ilks mapping: bytes32 => flip address; chop (ray), lump (wad) uint256
func (repository *CatStorageRepository) insertIlkFlip(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, flip string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertFieldWithIlk(blockNumber, blockHash, ilk, IlkFlip, InsertCatIlkFlipQuery, flip)
}

func (repository *CatStorageRepository) insertIlkChop(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, chop string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertFieldWithIlk(blockNumber, blockHash, ilk, IlkChop, InsertCatIlkChopQuery, chop)
}

func (repository *CatStorageRepository) insertIlkLump(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, lump string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertFieldWithIlk(blockNumber, blockHash, ilk, IlkLump, InsertCatIlkLumpQuery, lump)
}

// Flips mapping: uint256 => ilk, urn bytes32; ink, tab uint256 (both wad)
func (repository *CatStorageRepository) insertFlipIlk(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, ilk string) error {
	flip, err := getFlip(metadata.Keys)
	if err != nil {
		return err
	}
	// Ilks are stored without prefix in the db, but bytes32 values are decoded with a prefix
	ilkWithoutPrefix := common.Bytes2Hex(common.FromHex(ilk))
	return repository.insertFieldWithIlk(blockNumber, blockHash, ilkWithoutPrefix, FlipIlk, insertCatFlipIlkQuery, flip)
}

func (repository *CatStorageRepository) insertFlipUrn(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, urn string) error {
	flip, err := getFlip(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(insertCatFlipUrnQuery, blockNumber, blockHash, flip, urn)
	return writeErr
}

func (repository *CatStorageRepository) insertFlipInk(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, ink string) error {
	flip, err := getFlip(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(insertCatFlipInkQuery, blockNumber, blockHash, flip, ink)
	return writeErr
}

func (repository *CatStorageRepository) insertFlipTab(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, tab string) error {
	flip, err := getFlip(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(insertCatFlipTabQuery, blockNumber, blockHash, flip, tab)
	return writeErr
}

func (repository *CatStorageRepository) insertFieldWithIlk(blockNumber int, blockHash, ilk, variableName, query, value string) error {
	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return txErr
	}
	ilkID, ilkErr := shared.GetOrCreateIlkInTransaction(ilk, tx)
	if ilkErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return formatRollbackError("ilk", ilkErr.Error())
		}
		return ilkErr
	}
	_, writeErr := tx.Exec(query, blockNumber, blockHash, ilkID, value)
	if writeErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return formatRollbackError(variableName, writeErr.Error())
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

func formatRollbackError(field, err string) error {
	return fmt.Errorf("failed to rollback transaction after failing to insert %s: %s", field, err)
}
