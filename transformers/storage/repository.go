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
	repository2 "github.com/vulcanize/vulcanizedb/libraries/shared/repository"
	"strconv"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type Urn struct {
	Ilk        string
	Identifier string
}

var ErrNoFlips = errors.New("no flips exist in db")

type IMakerStorageRepository interface {
	GetDaiKeys() ([]string, error)
	GetFlapBidIds(string) ([]string, error)
	GetGemKeys() ([]Urn, error)
	GetIlks() ([]string, error)
	GetVatSinKeys() ([]string, error)
	GetVowSinKeys() ([]string, error)
	GetUrns() ([]Urn, error)
	GetCdpis() ([]string, error)
	GetOwners() ([]string, error)
	GetFlipBidIds(contractAddress string) ([]string, error)
	GetFlopBidIds(contractAddress string) ([]string, error)
	SetDB(db *postgres.DB)
}

type MakerStorageRepository struct {
	db *postgres.DB
}

func (repository *MakerStorageRepository) GetFlapBidIds(contractAddress string) ([]string, error) {
	var bidIds []string
	addressId, addressErr := repository.GetOrCreateAddress(contractAddress)
	if addressErr != nil {
		return []string{}, addressErr
	}
	err := repository.db.Select(&bidIds, `
		SELECT bid_id FROM maker.flap_kick WHERE address_id = $1
		UNION
		SELECT kicks FROM maker.flap_kicks WHERE address_id = $1
		UNION
		SELECT bid_id from maker.tend WHERE address_id = $1
		UNION
		SELECT bid_id from maker.deal WHERE address_id = $1
		UNION
		SELECT bid_id from maker.yank WHERE address_id = $1`, addressId)
	return bidIds, err
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
			LEFT JOIN public.header_sync_logs ON header_sync_logs.id = vat_heal.log_id
			WHERE header_sync_logs.tx_index = transactions.tx_index
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
			LEFT JOIN public.header_sync_logs ON header_sync_logs.id = vat_heal.log_id
			WHERE header_sync_logs.tx_index = transactions.tx_index`)
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
		JOIN maker.ilks on maker.ilks.id = maker.urns.ilk_id
		UNION
		SELECT DISTINCT ilks.ilk, fork.src AS identifier
		FROM maker.vat_fork fork
		INNER JOIN maker.ilks ilks ON ilks.id = fork.ilk_id
		UNION
		SELECT DISTINCT ilks.ilk, fork.dst AS identifier
		FROM maker.vat_fork fork
		INNER JOIN maker.ilks ilks ON ilks.id = fork.ilk_id`)
	return urns, err
}

func (repository *MakerStorageRepository) GetCdpis() ([]string, error) {
	nullValue := 0
	var maxCdpi int
	readErr := repository.db.Get(&maxCdpi, `
		SELECT COALESCE(MAX(cdpi), $1)
		FROM maker.cdp_manager_cdpi`, nullValue)
	if readErr != nil {
		return nil, readErr
	}
	if maxCdpi == nullValue {
		return []string{}, nil
	}
	return rangeIntsAsStrings(1, maxCdpi), readErr
}

func (repository *MakerStorageRepository) GetOwners() ([]string, error) {
	var owners []string
	err := repository.db.Select(&owners, `
		SELECT DISTINCT owner
		FROM maker.cdp_manager_owns`)
	return owners, err
}

func (repository *MakerStorageRepository) GetFlipBidIds(contractAddress string) ([]string, error) {
	var bidIds []string
	addressId, addressErr := repository.GetOrCreateAddress(contractAddress)
	if addressErr != nil {
		return []string{}, addressErr
	}
	err := repository.db.Select(&bidIds, `
   		SELECT DISTINCT bid_id FROM maker.tick
		WHERE address_id = $1
		UNION
   		SELECT DISTINCT bid_id FROM maker.flip_kick
		WHERE address_id = $1
		UNION
		SELECT DISTINCT bid_id FROM maker.tend
		WHERE address_id = $1
		UNION
		SELECT DISTINCT bid_id FROM maker.dent
		WHERE address_id = $1
		UNION
		SELECT DISTINCT bid_id FROM maker.deal
		WHERE address_id = $1
		UNION
		SELECT DISTINCT bid_id FROM maker.yank
		WHERE address_id = $1
		UNION
		SELECT DISTINCT kicks FROM maker.flip_kicks
		WHERE address_id = $1`, addressId)
	return bidIds, err
}

func (repository *MakerStorageRepository) GetFlopBidIds(contractAddress string) ([]string, error) {
	var bidIds []string
	addressId, addressErr := repository.GetOrCreateAddress(contractAddress)
	if addressErr != nil {
		return []string{}, addressErr
	}
	err := repository.db.Select(&bidIds, `
		SELECT bid_id FROM maker.flop_kick
		WHERE address_id = $1
		UNION
		SELECT DISTINCT bid_id FROM maker.dent
		WHERE address_id = $1
		UNION
		SELECT DISTINCT bid_id FROM maker.deal
		WHERE address_id = $1
		UNION
		SELECT DISTINCT bid_id FROM maker.yank
		WHERE address_id = $1
		UNION
		SELECT DISTINCT kicks FROM maker.flop_kicks
		WHERE address_id = $1`, addressId)
	return bidIds, err
}

func (repository *MakerStorageRepository) GetOrCreateAddress(contractAddress string) (int64, error) {
	return repository2.GetOrCreateAddress(repository.db, contractAddress)
}

func (repository *MakerStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func rangeIntsAsStrings(start, end int) []string {
	var strSlice []string
	for i := start; i <= end; i++ {
		strSlice = append(strSlice, strconv.Itoa(i))
	}
	return strSlice
}
