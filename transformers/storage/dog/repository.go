package dog

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
	Hole = "Hole"
	Dirt = "Dirt"
	Live = "live"
	Vat  = "vat"
	Vow  = "vow"

	IlkClip = "clip"
	IlkChop = "chop"
	IlkHole = "hole"
	IlkDirt = "dirt"

	InsertDogIlkClipQuery = `INSERT INTO maker.dog_ilk_clip (diff_id, header_id, address_id, ilk_id, clip) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`

	insertDogDirtQuery = `INSERT INTO maker.dog_dirt (diff_id, header_id, address_id, dirt) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertDogHoleQuery = `INSERT INTO maker.dog_hole (diff_id, header_id, address_id, hole) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertDogLiveQuery = `INSERT INTO maker.dog_live (diff_id, header_id, address_id, live) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertDogVatQuery  = `INSERT INTO maker.dog_vat (diff_id, header_id, address_id, vat) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertDogVowQuery  = `INSERT INTO maker.dog_vow (diff_id, header_id, address_id, vow) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
)

type StorageRepository struct {
	db                *postgres.DB
	ContractAddress   string
	contractAddressID int64
}

func (repo *StorageRepository) Create(diffID, headerID int64, metadata types.ValueMetadata, value interface{}) error {
	switch metadata.Name {
	case Dirt:
		return repo.insertDirt(diffID, headerID, value.(string))
	case Hole:
		return repo.insertHole(diffID, headerID, value.(string))
	case Live:
		return repo.insertLive(diffID, headerID, value.(string))
	case Vat:
		return repo.insertVat(diffID, headerID, value.(string))
	case Vow:
		return repo.insertVow(diffID, headerID, value.(string))
	case wards.Wards:
		return wards.InsertWards(diffID, headerID, metadata, repo.ContractAddress, value.(string), repo.db)
	case IlkClip:
		return repo.insertIlkClip(diffID, headerID, metadata, value.(string))
	default:
		return fmt.Errorf("unrecognized dog contract storage name: %s", metadata.Name)
	}
}

func (repo *StorageRepository) SetDB(db *postgres.DB) {
	repo.db = db
}

func (repo *StorageRepository) insertDirt(diffID, headerID int64, dirt string) error {
	addressID, addressErr := repo.ContractAddressID()
	if addressErr != nil {
		return fmt.Errorf("could not retrieve address id for %s, error: %w", repo.ContractAddress, addressErr)
	}

	_, err := repo.db.Exec(insertDogDirtQuery, diffID, headerID, addressID, dirt)
	if err != nil {
		return fmt.Errorf("error inserting dog dirt %s from diff ID %d: %w", dirt, diffID, err)
	}
	return nil
}

func (repo *StorageRepository) insertHole(diffID, headerID int64, hole string) error {
	addressID, addressErr := repo.ContractAddressID()
	if addressErr != nil {
		return fmt.Errorf("could not retrieve address id for %s, error: %w", repo.ContractAddress, addressErr)
	}

	_, err := repo.db.Exec(insertDogHoleQuery, diffID, headerID, addressID, hole)
	if err != nil {
		return fmt.Errorf("error inserting dog hole %s from diff ID %d: %w", hole, diffID, err)
	}
	return nil
}

func (repo *StorageRepository) insertLive(diffID, headerID int64, live string) error {
	addressID, addressErr := repo.ContractAddressID()
	if addressErr != nil {
		return fmt.Errorf("could not retrieve address id for %s, error: %w", repo.ContractAddress, addressErr)
	}

	_, err := repo.db.Exec(insertDogLiveQuery, diffID, headerID, addressID, live)
	if err != nil {
		return fmt.Errorf("error inserting dog live %s from diff ID %d: %w", live, diffID, err)
	}
	return nil
}

func (repo *StorageRepository) insertVat(diffID, headerID int64, vat string) error {
	addressID, addressErr := repo.ContractAddressID()
	if addressErr != nil {
		return fmt.Errorf("could not retrieve address id for %s, error: %w", repo.ContractAddress, addressErr)
	}

	vatAddressID, vatAddressErr := repository.GetOrCreateAddress(repo.db, vat)
	if vatAddressErr != nil {
		return fmt.Errorf("error persisting vat address: %w", vatAddressErr)
	}

	_, err := repo.db.Exec(insertDogVatQuery, diffID, headerID, addressID, vatAddressID)
	if err != nil {
		return fmt.Errorf("error inserting dog vat %s from diff ID %d: %w", vat, diffID, err)
	}
	return nil
}

func (repo *StorageRepository) insertVow(diffID, headerID int64, vow string) error {
	addressID, addressErr := repo.ContractAddressID()
	if addressErr != nil {
		return fmt.Errorf("could not retrieve address id for %s, error: %w", repo.ContractAddress, addressErr)
	}

	vowAddressID, vowAddressErr := repository.GetOrCreateAddress(repo.db, vow)
	if vowAddressErr != nil {
		return fmt.Errorf("error persisting vow address: %w", vowAddressErr)
	}

	_, err := repo.db.Exec(insertDogVowQuery, diffID, headerID, addressID, vowAddressID)
	if err != nil {
		return fmt.Errorf("error inserting dog vat %s from diff ID %d: %w", vow, diffID, err)
	}
	return nil
}

func (repo *StorageRepository) insertIlkClip(diffID, headerID int64, metadata types.ValueMetadata, clip string) error {
	addressID, addressErr := repo.ContractAddressID()
	if addressErr != nil {
		return fmt.Errorf("could not retrieve address id for %s, error: %w", repo.ContractAddress, addressErr)
	}

	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting ilk for ilk flip: %w", err)
	}
	insertErr := shared.InsertFieldWithIlkAndAddress(diffID, headerID, addressID, ilk, IlkClip, InsertDogIlkClipQuery, clip, repo.db)
	if insertErr != nil {
		return fmt.Errorf("error inserting ilk %s clip %s from diff ID %d: %w", insertErr, clip, diffID, insertErr)
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
