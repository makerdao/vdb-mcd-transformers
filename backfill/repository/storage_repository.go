package repository

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type Urn struct {
	Ilk   string
	IlkID int64 `db:"ilk_id"`
	Urn   string
	UrnID int64 `db:"urn_id"`
}

type StorageRepository interface {
	GetOrCreateUrn(urn, ilk string) (int64, error)
	GetUrnByID(id int64) (Urn, error)
	InsertDiff(diff types.RawDiff) error
	VatIlkArtExists(ilkID, headerID int64) (bool, error)
	VatUrnArtExists(urnID, headerID int64) (bool, error)
	VatUrnInkExists(urnID, headerID int64) (bool, error)
}

type storageRepository struct {
	db *postgres.DB
}

func (repo storageRepository) GetOrCreateUrn(urn, ilk string) (int64, error) {
	return shared.GetOrCreateUrn(urn, ilk, repo.db)
}

func NewStorageRepository(db *postgres.DB) StorageRepository {
	return storageRepository{db: db}
}

func (repo storageRepository) GetUrnByID(id int64) (Urn, error) {
	var urn Urn
	err := repo.db.Get(&urn, `
		SELECT DISTINCT urns.id AS urn_id, ilks.ilk, ilks.id AS ilk_id, urns.identifier AS urn
		FROM maker.urns
		    JOIN maker.ilks on maker.ilks.id = maker.urns.ilk_id
		    WHERE urns.id = $1`, id)
	return urn, err
}

func (repo storageRepository) InsertDiff(diff types.RawDiff) error {
	_, err := repo.db.Exec(`INSERT INTO public.storage_diff (block_height, block_hash, address, storage_key,
                storage_value, from_backfill, eth_node_id) VALUES ($1, $2, $3, $4, $5, true, $6) ON CONFLICT DO NOTHING;`,
		diff.BlockHeight, diff.BlockHash.Bytes(), diff.Address.Bytes(), diff.StorageKey.Bytes(),
		diff.StorageValue.Bytes(), repo.db.NodeID)
	return err
}

func (repo storageRepository) VatIlkArtExists(ilkID, headerID int64) (bool, error) {
	var exists bool
	err := repo.db.Get(&exists, `SELECT EXISTS (SELECT 1 FROM maker.vat_ilk_art WHERE ilk_id = $1 and header_id = $2)`, ilkID, headerID)
	return exists, err
}

func (repo storageRepository) VatUrnArtExists(urnID, headerID int64) (bool, error) {
	var exists bool
	err := repo.db.Get(&exists, `SELECT EXISTS (SELECT 1 FROM maker.vat_urn_art WHERE urn_id = $1 and header_id = $2)`, urnID, headerID)
	return exists, err
}

func (repo storageRepository) VatUrnInkExists(urnID, headerID int64) (bool, error) {
	var exists bool
	err := repo.db.Get(&exists, `SELECT EXISTS (SELECT 1 FROM maker.vat_urn_ink WHERE urn_id = $1 and header_id = $2)`, urnID, headerID)
	return exists, err
}
