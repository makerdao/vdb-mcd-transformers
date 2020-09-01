package cat

import (
	"fmt"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	Live   = "live"
	Vat    = "vat"
	Vow    = "vow"
	Box    = "box"
	Litter = "litter"

	IlkFlip = "flip"
	IlkChop = "chop"
	IlkLump = "lump"
	IlkDunk = "dunk"

	InsertCatIlkChopQuery = `INSERT INTO maker.cat_ilk_chop (diff_id, header_id, ilk_id, chop) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertCatIlkDunkQuery = `INSERT INTO maker.cat_ilk_dunk (diff_id, header_id, ilk_id, dunk) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertCatIlkFlipQuery = `INSERT INTO maker.cat_ilk_flip (diff_id, header_id, ilk_id, flip) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertCatIlkLumpQuery = `INSERT INTO maker.cat_ilk_lump (diff_id, header_id, ilk_id, lump) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`

	insertCatBoxQuery    = `INSERT INTO maker.cat_box (diff_id, header_id, box) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertCatLitterQuery = `INSERT INTO maker.cat_litter (diff_id, header_id, litter) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertCatLiveQuery   = `INSERT INTO maker.cat_live (diff_id, header_id, live) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertCatVatQuery    = `INSERT INTO maker.cat_vat (diff_id, header_id, vat) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertCatVowQuery    = `INSERT INTO maker.cat_vow (diff_id, header_id, vow) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
)

type StorageRepository struct {
	db              *postgres.DB
	ContractAddress string
}

func (repository *StorageRepository) Create(diffID, headerID int64, metadata types.ValueMetadata, value interface{}) error {
	switch metadata.Name {
	case Live:
		return repository.insertLive(diffID, headerID, value.(string))
	case Vat:
		return repository.insertVat(diffID, headerID, value.(string))
	case Vow:
		return repository.insertVow(diffID, headerID, value.(string))
	case Box:
		return repository.insertBox(diffID, headerID, value.(string))
	case Litter:
		return repository.insertLitter(diffID, headerID, value.(string))
	case wards.Wards:
		return wards.InsertWards(diffID, headerID, metadata, repository.ContractAddress, value.(string), repository.db)
	case IlkChop:
		return repository.insertIlkChop(diffID, headerID, metadata, value.(string))
	case IlkFlip:
		return repository.insertIlkFlip(diffID, headerID, metadata, value.(string))
	case IlkLump:
		return repository.insertIlkLump(diffID, headerID, metadata, value.(string))
	case IlkDunk:
		return repository.insertIlkDunk(diffID, headerID, metadata, value.(string))
	default:
		panic(fmt.Sprintf("unrecognized cat contract storage name: %s", metadata.Name))
	}
}

func (repository *StorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository *StorageRepository) insertLive(diffID, headerID int64, live string) error {
	_, err := repository.db.Exec(insertCatLiveQuery, diffID, headerID, live)
	if err != nil {
		return fmt.Errorf("error inserting cat live %s from diff ID %d: %w", live, diffID, err)
	}
	return nil
}

func (repository *StorageRepository) insertVat(diffID, headerID int64, vat string) error {
	_, err := repository.db.Exec(insertCatVatQuery, diffID, headerID, vat)
	if err != nil {
		return fmt.Errorf("error inserting cat vat %s from diff ID %d: %w", vat, diffID, err)
	}
	return nil
}

func (repository *StorageRepository) insertVow(diffID, headerID int64, vow string) error {
	_, err := repository.db.Exec(insertCatVowQuery, diffID, headerID, vow)
	if err != nil {
		return fmt.Errorf("error inserting cat vow %s from diff ID %d: %w", vow, diffID, err)
	}
	return nil
}

func (repository *StorageRepository) insertBox(diffID, headerID int64, box string) error {
	_, err := repository.db.Exec(insertCatBoxQuery, diffID, headerID, box)
	if err != nil {
		return fmt.Errorf("error inserting cat box %s from diff ID %d: %w", box, diffID, err)
	}
	return nil
}

func (repository *StorageRepository) insertLitter(diffID, headerID int64, litter string) error {
	_, err := repository.db.Exec(insertCatLitterQuery, diffID, headerID, litter)
	if err != nil {
		return fmt.Errorf("error inserting cat litter %s from diff ID %d: %w", litter, diffID, err)
	}
	return nil
}

func (repository *StorageRepository) insertIlkFlip(diffID, headerID int64, metadata types.ValueMetadata, flip string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting ilk for ilk flip: %w", err)
	}
	insertErr := shared.InsertFieldWithIlk(diffID, headerID, ilk, IlkFlip, InsertCatIlkFlipQuery, flip, repository.db)
	if insertErr != nil {
		return fmt.Errorf("error inserting ilk %s flip %s from diff ID %d: %w", insertErr, flip, diffID, insertErr)
	}
	return nil
}

func (repository *StorageRepository) insertIlkChop(diffID, headerID int64, metadata types.ValueMetadata, chop string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting ilk for ilk chop: %w", err)
	}
	insertErr := shared.InsertFieldWithIlk(diffID, headerID, ilk, IlkChop, InsertCatIlkChopQuery, chop, repository.db)
	if insertErr != nil {
		return fmt.Errorf("error inserting ilk %s chop %s from diff Id %d: %w", ilk, chop, diffID, insertErr)
	}
	return nil
}

func (repository *StorageRepository) insertIlkLump(diffID, headerID int64, metadata types.ValueMetadata, lump string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting ilk for ilk lump: %w", err)
	}
	insertErr := shared.InsertFieldWithIlk(diffID, headerID, ilk, IlkLump, InsertCatIlkLumpQuery, lump, repository.db)
	if insertErr != nil {
		return fmt.Errorf("error inserting ilk %s lump %s from diff ID %d: %w", ilk, lump, diffID, insertErr)
	}
	return nil
}

func (repository *StorageRepository) insertIlkDunk(diffID, headerID int64, metadata types.ValueMetadata, dunk string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting ilk for ilk dunk: %w", err)
	}
	insertErr := shared.InsertFieldWithIlk(diffID, headerID, ilk, IlkDunk, insertCatIlkDunkQuery, dunk, repository.db)
	if insertErr != nil {
		return fmt.Errorf("error inserting ilk %s dunk %s from diff ID %d: %w", ilk, dunk, diffID, insertErr)
	}
	return nil
}

func getIlk(keys map[types.Key]string) (string, error) {
	ilk, ok := keys[constants.Ilk]
	if !ok {
		return "", types.ErrMetadataMalformed{MissingData: constants.Ilk}
	}
	return ilk, nil
}
