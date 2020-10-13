package cat

import (
	"fmt"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
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

	InsertCatIlkChopQuery = `INSERT INTO maker.cat_ilk_chop (diff_id, header_id, address_id, ilk_id, chop) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	insertCatIlkDunkQuery = `INSERT INTO maker.cat_ilk_dunk (diff_id, header_id, address_id, ilk_id, dunk) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertCatIlkFlipQuery = `INSERT INTO maker.cat_ilk_flip (diff_id, header_id, address_id, ilk_id, flip) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertCatIlkLumpQuery = `INSERT INTO maker.cat_ilk_lump (diff_id, header_id, address_id, ilk_id, lump) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`

	insertCatBoxQuery    = `INSERT INTO maker.cat_box (diff_id, header_id, address_id, box) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertCatLitterQuery = `INSERT INTO maker.cat_litter (diff_id, header_id, address_id, litter) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertCatLiveQuery   = `INSERT INTO maker.cat_live (diff_id, header_id, address_id, live) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertCatVatQuery    = `INSERT INTO maker.cat_vat (diff_id, header_id, address_id, vat) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertCatVowQuery    = `INSERT INTO maker.cat_vow (diff_id, header_id, address_id, vow) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
)

type StorageRepository struct {
	db                *postgres.DB
	ContractAddress   string
	contractAddressID int64
}

func (repo *StorageRepository) Create(diffID, headerID int64, metadata types.ValueMetadata, value interface{}) error {
	switch metadata.Name {
	case Live:
		return repo.insertLive(diffID, headerID, value.(string))
	case Vat:
		return repo.insertVat(diffID, headerID, value.(string))
	case Vow:
		return repo.insertVow(diffID, headerID, value.(string))
	case Box:
		return repo.insertBox(diffID, headerID, value.(string))
	case Litter:
		return repo.insertLitter(diffID, headerID, value.(string))
	case wards.Wards:
		return wards.InsertWards(diffID, headerID, metadata, repo.ContractAddress, value.(string), repo.db)
	case IlkChop:
		return repo.insertIlkChop(diffID, headerID, metadata, value.(string))
	case IlkFlip:
		return repo.insertIlkFlip(diffID, headerID, metadata, value.(string))
	case IlkLump:
		return repo.insertIlkLump(diffID, headerID, metadata, value.(string))
	case IlkDunk:
		return repo.insertIlkDunk(diffID, headerID, metadata, value.(string))
	default:
		panic(fmt.Sprintf("unrecognized cat contract storage name: %s", metadata.Name))
	}
}

func (repo *StorageRepository) SetDB(db *postgres.DB) {
	repo.db = db
}

func (repo *StorageRepository) insertLive(diffID, headerID int64, live string) error {
	addressID, addressErr := repo.ContractAddressID()
	if addressErr != nil {
		return fmt.Errorf("could not retrieve address id for %s, error: %w", repo.ContractAddress, addressErr)
	}

	_, err := repo.db.Exec(insertCatLiveQuery, diffID, headerID, addressID, live)
	if err != nil {
		return fmt.Errorf("error inserting cat live %s from diff ID %d: %w", live, diffID, err)
	}
	return nil
}

func (repo *StorageRepository) insertVat(diffID, headerID int64, vat string) error {
	addressID, addressErr := repo.ContractAddressID()
	if addressErr != nil {
		return fmt.Errorf("could not retrieve address id for %s, error: %w", repo.ContractAddress, addressErr)
	}

	_, err := repo.db.Exec(insertCatVatQuery, diffID, headerID, addressID, vat)
	if err != nil {
		return fmt.Errorf("error inserting cat vat %s from diff ID %d: %w", vat, diffID, err)
	}
	return nil
}

func (repo *StorageRepository) insertVow(diffID, headerID int64, vow string) error {
	addressID, addressErr := repo.ContractAddressID()
	if addressErr != nil {
		return fmt.Errorf("could not retrieve address id for %s, error: %w", repo.ContractAddress, addressErr)
	}

	_, err := repo.db.Exec(insertCatVowQuery, diffID, headerID, addressID, vow)
	if err != nil {
		return fmt.Errorf("error inserting cat vow %s from diff ID %d: %w", vow, diffID, err)
	}
	return nil
}

func (repo *StorageRepository) insertBox(diffID, headerID int64, box string) error {
	addressID, addressErr := repo.ContractAddressID()
	if addressErr != nil {
		return fmt.Errorf("could not retrieve address id for %s, error: %w", repo.ContractAddress, addressErr)
	}

	_, err := repo.db.Exec(insertCatBoxQuery, diffID, headerID, addressID, box)
	if err != nil {
		return fmt.Errorf("error inserting cat box %s from diff ID %d: %w", box, diffID, err)
	}
	return nil
}

func (repo *StorageRepository) insertLitter(diffID, headerID int64, litter string) error {
	addressID, addressErr := repo.ContractAddressID()
	if addressErr != nil {
		return fmt.Errorf("could not retrieve address id for %s, error: %w", repo.ContractAddress, addressErr)
	}

	_, err := repo.db.Exec(insertCatLitterQuery, diffID, headerID, addressID, litter)
	if err != nil {
		return fmt.Errorf("error inserting cat litter %s from diff ID %d: %w", litter, diffID, err)
	}
	return nil
}

func (repo *StorageRepository) insertIlkFlip(diffID, headerID int64, metadata types.ValueMetadata, flip string) error {
	addressID, addressErr := repo.ContractAddressID()
	if addressErr != nil {
		return fmt.Errorf("could not retrieve address id for %s, error: %w", repo.ContractAddress, addressErr)
	}

	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting ilk for ilk flip: %w", err)
	}
	insertErr := shared.InsertFieldWithIlkAndAddress(diffID, headerID, addressID, ilk, IlkFlip, InsertCatIlkFlipQuery, flip, repo.db)
	if insertErr != nil {
		return fmt.Errorf("error inserting ilk %s flip %s from diff ID %d: %w", insertErr, flip, diffID, insertErr)
	}
	return nil
}

func (repo *StorageRepository) insertIlkChop(diffID, headerID int64, metadata types.ValueMetadata, chop string) error {
	addressID, addressErr := repo.ContractAddressID()
	if addressErr != nil {
		return fmt.Errorf("could not retrieve address id for %s, error: %w", repo.ContractAddress, addressErr)
	}

	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting ilk for ilk chop: %w", err)
	}
	insertErr := shared.InsertFieldWithIlkAndAddress(diffID, headerID, addressID, ilk, IlkChop, InsertCatIlkChopQuery, chop, repo.db)
	if insertErr != nil {
		return fmt.Errorf("error inserting ilk %s chop %s from diff Id %d: %w", ilk, chop, diffID, insertErr)
	}
	return nil
}

func (repo *StorageRepository) insertIlkLump(diffID, headerID int64, metadata types.ValueMetadata, lump string) error {
	addressID, addressErr := repo.ContractAddressID()
	if addressErr != nil {
		return fmt.Errorf("could not retrieve address id for %s, error: %w", repo.ContractAddress, addressErr)
	}

	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting ilk for ilk lump: %w", err)
	}

	insertErr := shared.InsertFieldWithIlkAndAddress(diffID, headerID, addressID, ilk, IlkLump, InsertCatIlkLumpQuery, lump, repo.db)
	if insertErr != nil {
		return fmt.Errorf("error inserting ilk %s lump %s from diff ID %d: %w", ilk, lump, diffID, insertErr)
	}
	return nil
}

func (repo *StorageRepository) insertIlkDunk(diffID, headerID int64, metadata types.ValueMetadata, dunk string) error {
	addressID, addressErr := repo.ContractAddressID()
	if addressErr != nil {
		return fmt.Errorf("could not retrieve address id for %s, error: %w", repo.ContractAddress, addressErr)
	}

	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting ilk for ilk dunk: %w", err)
	}
	insertErr := shared.InsertFieldWithIlkAndAddress(diffID, headerID, addressID, ilk, IlkDunk, insertCatIlkDunkQuery, dunk, repo.db)
	if insertErr != nil {
		return fmt.Errorf("error inserting ilk %s dunk %s from diff ID %d: %w", ilk, dunk, diffID, insertErr)
	}
	return nil
}

func (repo *StorageRepository) ContractAddressID() (int64, error) {
	if repo.contractAddressID == 0 {
		addressID, addressErr := repository.GetOrCreateAddress(repo.db, repo.ContractAddress)
		repo.contractAddressID = addressID
		return repo.contractAddressID, addressErr
	}
	return repo.contractAddressID, nil
}

func getIlk(keys map[types.Key]string) (string, error) {
	ilk, ok := keys[constants.Ilk]
	if !ok {
		return "", types.ErrMetadataMalformed{MissingData: constants.Ilk}
	}
	return ilk, nil
}
