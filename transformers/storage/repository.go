// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package storage

import (
	"errors"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type Urn struct {
	Ilk        string
	Identifier string
}

var ErrNoFlips = errors.New("no flips exist in db")

type IMakerStorageRepository interface {
	GetDaiKeys() ([]string, error)
	GetGemKeys() ([]Urn, error)
	GetIlks() ([]string, error)
	GetVatSinKeys() ([]string, error)
	GetVowSinKeys() ([]string, error)
	GetUrns() ([]Urn, error)
	SetDB(db *postgres.DB)
}

type MakerStorageRepository struct {
	db *postgres.DB
}

func (repository *MakerStorageRepository) GetDaiKeys() ([]string, error) {
	var daiKeys []string
	err := repository.db.Select(&daiKeys, `
		SELECT DISTINCT src FROM maker.vat_move
		UNION
		SELECT DISTINCT dst FROM maker.vat_move
		UNION
		SELECT DISTINCT w FROM maker.vat_frob
		UNION
		SELECT DISTINCT v FROM maker.vat_suck
		UNION
		SELECT DISTINCT tx_from FROM public.header_sync_transactions AS transactions
			LEFT JOIN maker.vat_heal ON vat_heal.header_id = transactions.header_id
			WHERE vat_heal.tx_idx = transactions.tx_index
		UNION
		SELECT DISTINCT urns.identifier FROM maker.vat_fold
			INNER JOIN maker.urns on urns.id = maker.vat_fold.urn_id
	`)
	return daiKeys, err
}

func (repository *MakerStorageRepository) GetGemKeys() ([]Urn, error) {
	var gems []Urn
	err := repository.db.Select(&gems, `
		SELECT DISTINCT ilks.ilk, slip.usr AS identifier
		FROM maker.vat_slip slip
		INNER JOIN maker.ilks ilks ON ilks.id = slip.ilk_id
		UNION
		SELECT DISTINCT ilks.ilk, flux.src AS identifier
		FROM maker.vat_flux flux
		INNER JOIN maker.ilks ilks ON ilks.id = flux.ilk_id
		UNION
		SELECT DISTINCT ilks.ilk, flux.dst AS identifier
		FROM maker.vat_flux flux
		INNER JOIN maker.ilks ilks ON ilks.id = flux.ilk_id
		UNION
		SELECT DISTINCT ilks.ilk, frob.v AS identifier
		FROM maker.vat_frob frob
		INNER JOIN maker.urns on urns.id = frob.urn_id
		INNER JOIN maker.ilks ilks ON ilks.id = urns.ilk_id
		UNION
		SELECT DISTINCT ilks.ilk, grab.v AS identifier
		FROM maker.vat_grab grab
		INNER JOIN maker.urns on urns.id = grab.urn_id
		INNER JOIN maker.ilks ilks ON ilks.id = urns.ilk_id
	`)
	return gems, err
}

func (repository MakerStorageRepository) GetIlks() ([]string, error) {
	var ilks []string
	err := repository.db.Select(&ilks, `SELECT DISTINCT ilk FROM maker.ilks`)
	return ilks, err
}

func (repository *MakerStorageRepository) GetVatSinKeys() ([]string, error) {
	var sinKeys []string
	err := repository.db.Select(&sinKeys, `
		SELECT DISTINCT w FROM maker.vat_grab
		UNION
		SELECT DISTINCT u FROM maker.vat_suck
		UNION
		SELECT DISTINCT tx_from FROM public.header_sync_transactions AS transactions
			LEFT JOIN maker.vat_heal ON vat_heal.header_id = transactions.header_id
			WHERE vat_heal.tx_idx = transactions.tx_index`)
	return sinKeys, err
}

func (repository *MakerStorageRepository) GetVowSinKeys() ([]string, error) {
	var sinKeys []string
	err := repository.db.Select(&sinKeys, `
		SELECT DISTINCT era FROM maker.vow_flog
		UNION
		SELECT DISTINCT headers.block_timestamp
		FROM maker.vow_fess
		JOIN headers ON maker.vow_fess.header_id = headers.id`)
	return sinKeys, err
}

func (repository *MakerStorageRepository) GetUrns() ([]Urn, error) {
	var urns []Urn
	err := repository.db.Select(&urns, `
		SELECT DISTINCT ilks.ilk, urns.identifier
		FROM maker.urns
		JOIN maker.ilks on maker.ilks.id = maker.urns.ilk_id`)
	return urns, err
}

func (repository *MakerStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}
