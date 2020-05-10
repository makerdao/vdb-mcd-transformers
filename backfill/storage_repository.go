package backfill

import (
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type Urn struct {
	Ilk   string
	IlkID int `db:"ilk_id"`
	Urn   string
	UrnID int `db:"urn_id"`
}

type StorageRepository interface {
	GetUrnByID(id int) (Urn, error)
	GetUrns() ([]Urn, error)
	InsertDiff(diff types.RawDiff) error
	VatIlkArtExists(ilkID, headerID int) (bool, error)
	VatUrnArtExists(urnID, headerID int) (bool, error)
	VatUrnInkExists(urnID, headerID int) (bool, error)
}

type storageRepository struct {
	db *postgres.DB
}

func NewStorageRepository(db *postgres.DB) StorageRepository {
	return storageRepository{db: db}
}

func (repo storageRepository) GetUrnByID(id int) (Urn, error) {
	var urn Urn
	err := repo.db.Get(&urn, `
		SELECT DISTINCT urns.id AS urn_id, ilks.ilk, ilks.id AS ilk_id, urns.identifier AS urn
		FROM maker.urns
		    JOIN maker.ilks on maker.ilks.id = maker.urns.ilk_id
		    WHERE urns.id = $1`, id)
	return urn, err
}

func (repo storageRepository) GetUrns() ([]Urn, error) {
	var urns []Urn
	err := repo.db.Select(&urns, `
		SELECT DISTINCT urns.id AS urn_id, ilks.ilk, ilks.id AS ilk_id, urns.identifier AS urn
		FROM maker.urns
		    JOIN maker.ilks on maker.ilks.id = maker.urns.ilk_id`)
	return urns, err
}

func (repo storageRepository) InsertDiff(diff types.RawDiff) error {
	_, err := repo.db.Exec(`INSERT INTO public.storage_diff (block_height, block_hash, hashed_address, storage_key,
                storage_value, from_backfill) VALUES ($1, $2, $3, $4, $5, true) ON CONFLICT DO NOTHING;`,
		diff.BlockHeight, diff.BlockHash.Bytes(), diff.HashedAddress.Bytes(), diff.StorageKey.Bytes(),
		diff.StorageValue.Bytes())
	return err
}

func (repo storageRepository) VatIlkArtExists(ilkID, headerID int) (bool, error) {
	var exists bool
	err := repo.db.Get(&exists, `SELECT EXISTS (SELECT 1 FROM maker.vat_ilk_art WHERE ilk_id = $1 and header_id = $2)`, ilkID, headerID)
	return exists, err
}

func (repo storageRepository) VatUrnArtExists(urnID, headerID int) (bool, error) {
	var exists bool
	err := repo.db.Get(&exists, `SELECT EXISTS (SELECT 1 FROM maker.vat_urn_art WHERE urn_id = $1 and header_id = $2)`, urnID, headerID)
	return exists, err
}

func (repo storageRepository) VatUrnInkExists(urnID, headerID int) (bool, error) {
	var exists bool
	err := repo.db.Get(&exists, `SELECT EXISTS (SELECT 1 FROM maker.vat_urn_ink WHERE urn_id = $1 and header_id = $2)`, urnID, headerID)
	return exists, err
}
