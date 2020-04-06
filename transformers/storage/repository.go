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
	"strconv"

	vdbRepository "github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type Urn struct {
	Ilk        string
	Identifier string
}

type IMakerStorageRepository interface {
	GetCdpis() ([]string, error)
	GetDaiKeys() ([]string, error)
	GetFlapBidIDs(string) ([]string, error)
	GetFlipBidIDs(contractAddress string) ([]string, error)
	GetFlopBidIDs(contractAddress string) ([]string, error)
	GetGemKeys() ([]Urn, error)
	GetIlks() ([]string, error)
	GetMedianBudAddresses(string) ([]string, error)
	GetOwners() ([]string, error)
	GetPotPieUsers() ([]string, error)
	GetUrns() ([]Urn, error)
	GetVatSinKeys() ([]string, error)
	GetVowSinKeys() ([]string, error)
	GetVatWardsAddresses() ([]string, error)
	GetWardsAddresses(string) ([]string, error)
	SetDB(db *postgres.DB)
}

type MakerStorageRepository struct {
	db *postgres.DB
}

func (repository *MakerStorageRepository) GetFlapBidIDs(contractAddress string) ([]string, error) {
	var bidIDs []string
	addressID, addressErr := repository.GetOrCreateAddress(contractAddress)
	if addressErr != nil {
		return []string{}, addressErr
	}
	err := repository.db.Select(&bidIDs, `
		SELECT bid_id FROM maker.flap_kick WHERE address_id = $1
		UNION
		SELECT kicks FROM maker.flap_kicks WHERE address_id = $1
		UNION
		SELECT bid_id from maker.tend WHERE address_id = $1
		UNION
		SELECT bid_id from maker.deal WHERE address_id = $1
		UNION
		SELECT bid_id from maker.yank WHERE address_id = $1`, addressID)
	return bidIDs, err
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
		SELECT DISTINCT u FROM maker.vat_fold
		UNION
		SELECT DISTINCT tx_from FROM public.transactions AS transactions
			LEFT JOIN maker.vat_heal ON vat_heal.header_id = transactions.header_id
			LEFT JOIN public.event_logs ON event_logs.id = vat_heal.log_id
			WHERE event_logs.tx_index = transactions.tx_index
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

func (repository *MakerStorageRepository) GetMedianBudAddresses(contractAddress string) ([]string, error) {
	addressID, addressErr := repository.GetOrCreateAddress(contractAddress)
	if addressErr != nil {
		return nil, addressErr
	}
	var budAddresses []string
	err := repository.db.Select(&budAddresses, `
		SELECT addresses.address
		FROM maker.median_kiss_single
		LEFT JOIN public.addresses ON median_kiss_single.a = addresses.id
		WHERE median_kiss_single.address_id = $1
		UNION
		SELECT addresses.address
		FROM maker.median_diss_single
		LEFT JOIN public.addresses ON median_diss_single.a = addresses.id
		WHERE median_diss_single.address_id = $1
		UNION
		SELECT a_address
		FROM maker.median_kiss_batch, UNNEST(a) a_address
		WHERE address_id = $1
		UNION
		SELECT a_address
		FROM maker.median_diss_batch, UNNEST(a) a_address
		WHERE address_id = $1
	`, addressID)
	return budAddresses, err
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
		SELECT DISTINCT tx_from FROM public.transactions AS transactions
			LEFT JOIN maker.vat_heal ON vat_heal.header_id = transactions.header_id
			LEFT JOIN public.event_logs ON event_logs.id = vat_heal.log_id
			WHERE event_logs.tx_index = transactions.tx_index`)
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

func (repository *MakerStorageRepository) GetFlipBidIDs(contractAddress string) ([]string, error) {
	var bidIDs []string
	addressID, addressErr := repository.GetOrCreateAddress(contractAddress)
	if addressErr != nil {
		return []string{}, addressErr
	}
	err := repository.db.Select(&bidIDs, `
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
		WHERE address_id = $1`, addressID)
	return bidIDs, err
}

func (repository *MakerStorageRepository) GetPotPieUsers() ([]string, error) {
	var userAddresses []string
	err := repository.db.Select(&userAddresses, `
		SELECT addresses.address
		FROM maker.pot_join
		    LEFT JOIN public.addresses ON pot_join.msg_sender = addresses.id
		UNION
		SELECT addresses.address
		FROM maker.pot_exit
		    LEFT JOIN public.addresses ON pot_exit.msg_sender = addresses.id`)
	return userAddresses, err
}

func (repository *MakerStorageRepository) GetFlopBidIDs(contractAddress string) ([]string, error) {
	var bidIDs []string
	addressID, addressErr := repository.GetOrCreateAddress(contractAddress)
	if addressErr != nil {
		return []string{}, addressErr
	}
	err := repository.db.Select(&bidIDs, `
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
		WHERE address_id = $1`, addressID)
	return bidIDs, err
}

func (repository *MakerStorageRepository) GetVatWardsAddresses() ([]string, error) {
	var wardsKeys []string
	selectErr := repository.db.Select(&wardsKeys, `
		SELECT addresses.address
		FROM maker.vat_rely
		    LEFT JOIN public.addresses ON vat_rely.usr = addresses.id
		UNION
		SELECT addresses.address
		FROM maker.vat_deny
		    LEFT JOIN public.addresses ON vat_deny.usr = addresses.id`)
	return wardsKeys, selectErr
}

func (repository *MakerStorageRepository) GetWardsAddresses(contractAddress string) ([]string, error) {
	contractAddressID, addressErr := repository.GetOrCreateAddress(contractAddress)
	if addressErr != nil {
		return nil, addressErr
	}
	var wardsKeys []string
	selectErr := repository.db.Select(&wardsKeys, `
		SELECT addresses.address
		FROM maker.rely
		    LEFT JOIN public.addresses ON rely.usr = addresses.id
		WHERE rely.address_id = $1
		UNION
		SELECT addresses.address
		FROM maker.rely
		    LEFT JOIN public.addresses ON rely.msg_sender = addresses.id
		WHERE rely.address_id = $1
		UNION
		SELECT addresses.address
		FROM maker.deny
		    LEFT JOIN public.addresses ON deny.usr = addresses.id
		WHERE deny.address_id = $1
		UNION
		SELECT addresses.address
		FROM maker.deny
		LEFT JOIN public.addresses ON deny.msg_sender = addresses.id
		WHERE deny.address_id = $1`, contractAddressID)
	return wardsKeys, selectErr
}

func (repository *MakerStorageRepository) GetOrCreateAddress(contractAddress string) (int64, error) {
	return vdbRepository.GetOrCreateAddress(repository.db, contractAddress)
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
