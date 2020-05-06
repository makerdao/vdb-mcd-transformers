package backfill

import "github.com/makerdao/vulcanizedb/pkg/datastore/postgres"

type Urn struct {
	ID  int
	Ilk string
	Urn string
}

type UrnsRepository interface {
	GetUrns() ([]Urn, error)
	VatUrnArtExists(urnID, headerID int) (bool, error)
	VatUrnInkExists(urnID, headerID int) (bool, error)
}

type urnsRepository struct {
	db *postgres.DB
}

func NewUrnsRepository(db *postgres.DB) UrnsRepository {
	return urnsRepository{db: db}
}

func (u urnsRepository) GetUrns() ([]Urn, error) {
	var urns []Urn
	err := u.db.Select(&urns, `SELECT DISTINCT urns.id, ilks.ilk, urns.identifier AS urn
		FROM maker.urns
		JOIN maker.ilks on maker.ilks.id = maker.urns.ilk_id`)
	return urns, err
}

func (u urnsRepository) VatUrnArtExists(urnID, headerID int) (bool, error) {
	var exists bool
	err := u.db.Get(&exists, `SELECT EXISTS (SELECT 1 FROM maker.vat_urn_art WHERE urn_id = $1 and header_id = $2)`, urnID, headerID)
	return exists, err
}

func (u urnsRepository) VatUrnInkExists(urnID, headerID int) (bool, error) {
	var exists bool
	err := u.db.Get(&exists, `SELECT EXISTS (SELECT 1 FROM maker.vat_urn_ink WHERE urn_id = $1 and header_id = $2)`, urnID, headerID)
	return exists, err
}
