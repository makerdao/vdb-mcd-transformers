package clip

import (
	"fmt"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	Dog     = "dog"
	Vow     = "vow"
	Spotter = "spotter"
	Calc    = "calc"
	Buf     = "buf"
	Tail    = "tail"
	Cusp    = "cusp"
	Chip    = "chip"
	Tip     = "tip"
	Chost   = "chost"
	Active  = "active"

	insertClipDogQuery     = `INSERT INTO maker.clip_dog (diff_id, header_id, address_id, dog) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertClipVowQuery     = `INSERT INTO maker.clip_vow (diff_id, header_id, address_id, vow) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertClipSpotterQuery = `INSERT INTO maker.clip_spotter (diff_id, header_id, address_id, spotter) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
)

type StorageRepository struct {
	ContractAddress string
	db              *postgres.DB
}

func (repo *StorageRepository) Create(diffID, headerID int64, metadata types.ValueMetadata, value interface{}) error {
	switch metadata.Name {
	case Dog:
		return repo.insertDog(diffID, headerID, value.(string))
	case Vow:
		return repo.insertVow(diffID, headerID, value.(string))
	case Spotter:
		return repo.insertSpotter(diffID, headerID, value.(string))
	case wards.Wards:
		return wards.InsertWards(diffID, headerID, metadata, repo.ContractAddress, value.(string), repo.db)
	default:
		return fmt.Errorf("unrecognized clip contract storage name: %s", metadata.Name)
	}
}

func (repo *StorageRepository) SetDB(db *postgres.DB) {
	repo.db = db
}

func (repo *StorageRepository) insertDog(diffID, headerID int64, dog string) error {
	dogAddressID, addressErr := repository.GetOrCreateAddress(repo.db, dog)
	if addressErr != nil {
		return fmt.Errorf("error inserting clip dog: %w", addressErr)
	}
	insertErr := shared.InsertRecordWithAddress(
		diffID,
		headerID,
		insertClipDogQuery,
		strconv.FormatInt(dogAddressID, 10),
		repo.ContractAddress,
		repo.db)
	if insertErr != nil {
		msgToFormat := "error inserting clip %s dog %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repo.ContractAddress, dog, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repo *StorageRepository) insertVow(diffID, headerID int64, vow string) error {
	vowAddressID, addressErr := repository.GetOrCreateAddress(repo.db, vow)
	if addressErr != nil {
		return fmt.Errorf("error inserting clip vow: %w", addressErr)
	}
	insertErr := shared.InsertRecordWithAddress(
		diffID,
		headerID,
		insertClipVowQuery,
		strconv.FormatInt(vowAddressID, 10),
		repo.ContractAddress,
		repo.db)
	if insertErr != nil {
		msgToFormat := "error inserting clip %s vow %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repo.ContractAddress, vow, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repo *StorageRepository) insertSpotter(diffID, headerID int64, spotter string) error {
	spotterAddressID, addressErr := repository.GetOrCreateAddress(repo.db, spotter)
	if addressErr != nil {
		return fmt.Errorf("error inserting clip spotter: %w", addressErr)
	}
	insertErr := shared.InsertRecordWithAddress(
		diffID,
		headerID,
		insertClipSpotterQuery,
		strconv.FormatInt(spotterAddressID, 10),
		repo.ContractAddress,
		repo.db)
	if insertErr != nil {
		msgToFormat := "error inserting clip %s spotter %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repo.ContractAddress, spotter, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}
