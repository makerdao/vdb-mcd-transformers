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

package storage_test

import (
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("Maker storage repository", func() {
	var (
		db                  *postgres.DB
		repository          storage.IMakerStorageRepository
		ilk1                = "494c4b31" // ILK1
		ilk2                = "494c4b32" // ILK2
		guy1                = "47555931" // GUY1
		guy2                = "47555932" // GUY2
		guy3                = "47555933" // GUY3
		era                 = big.NewInt(0).SetBytes(common.FromHex("0x000000000000000000000000000000000000000000000000000000005bb48864")).String()
		tab                 = big.NewInt(0).SetBytes(common.FromHex("0x0000000000000000000000000000000000000000000002544faa778090e00000")).String()
		timestamp           = int64(1538558053)
		transactionFromGuy1 = core.TransactionModel{
			From:    guy1,
			TxIndex: 3,
			Value:   "0",
		}
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repository = &storage.MakerStorageRepository{}
		repository.SetDB(db)
	})

	Describe("getting dai keys", func() {
		It("fetches guy from both src and dst field on vat_move", func() {
			insertVatMove(guy1, guy2, 1, db)

			keys, err := repository.GetDaiKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(keys)).To(Equal(2))
			Expect(keys).To(ConsistOf(guy1, guy2))
		})

		It("fetches guy from w field on vat_frob", func() {
			insertVatFrob(ilk1, guy1, guy1, guy2, 1, db)

			keys, err := repository.GetDaiKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(keys)).To(Equal(1))
			Expect(keys).To(ConsistOf(guy2))
		})

		It("fetches guy from vat_heal transaction", func() {
			insertVatHeal(1, transactionFromGuy1, db)

			keys, err := repository.GetDaiKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(keys)).To(Equal(1))
			Expect(keys).To(ConsistOf(guy1))
		})

		It("fetches guy from v field on vat_suck", func() {
			insertVatSuck(guy1, guy2, 0, 1, db)

			daiKeys , repoErr := repository.GetDaiKeys()

			Expect(repoErr).NotTo(HaveOccurred())
			Expect(len(daiKeys)).To(Equal(1))
			Expect(daiKeys).To(ConsistOf(guy2))
		})

		It("fetches guy from u field on vat_fold", func() {
			insertVatFold(guy1, 1, db)

			daiKeys , repoErr := repository.GetDaiKeys()

			Expect(repoErr).NotTo(HaveOccurred())
			Expect(len(daiKeys)).To(Equal(1))
			Expect(daiKeys).To(ConsistOf(guy1))
		})

		It("fetches unique guys from vat_move + vat_frob + vat_heal + vat_fold + vat_suck", func() {
			guy4 := "47555934"
			guy5 := "47555935"
			guy6 := "47555936"
			transactionFromGuy4 := core.TransactionModel{From: guy4, TxIndex: 4, Value: "0"}
			insertVatMove(guy1, guy2, 1, db)
			insertVatFrob(ilk1, guy1, guy1, guy3, 2, db)
			insertVatHeal(3, transactionFromGuy4, db)
			insertVatFold(guy5, 4, db)
			insertVatSuck(guy1, guy6, 0, 5, db)

			// duplicates
			insertVatMove(guy3, guy1, 6, db)
			insertVatFrob(ilk2, guy2, guy2, guy5, 7, db)
			insertVatHeal(8, transactionFromGuy1, db)
			insertVatFold(guy4, 9, db)
			insertVatSuck(guy1, guy1, 0, 10, db)

			keys, err := repository.GetDaiKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(keys)).To(Equal(6))
			Expect(keys).To(ConsistOf(guy1, guy2, guy3, guy4, guy5, guy6))
		})

		It("fetches the correct guy when there are multiple transactions in a block", func() {
			insertVatHeal(1, transactionFromGuy1, db)
			unrelatedTransaction := core.TransactionModel{From: "unrelated guy", TxIndex: 15, Value: "0"}
			insertTransaction(1, unrelatedTransaction, db)

			sinKeys, err := repository.GetDaiKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(sinKeys)).To(Equal(1))
			Expect(sinKeys).To(ConsistOf(guy1))
		})

		It("does not return error if no matching rows", func() {
			daiKeys, err := repository.GetDaiKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(daiKeys)).To(BeZero())
		})
	})

	Describe("getting gem keys", func() {
		It("fetches guy from both src and dst field on vat_flux", func() {
			insertVatFlux(ilk1, guy1, guy2, 1, db)

			gems, err := repository.GetGemKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(gems)).To(Equal(2))
			Expect(gems).To(ConsistOf([]storage.Urn{{
				Ilk:        ilk1,
				Identifier: guy1,
			}, {
				Ilk:        ilk1,
				Identifier: guy2,
			}}))
		})

		It("fetches guy from v field on vat_frob + vat_grab", func() {
			insertVatFrob(ilk1, guy1, guy2, guy1, 1, db)
			insertVatGrab(ilk1, guy1, guy3, guy1, 2, db)

			gems, err := repository.GetGemKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(gems)).To(Equal(2))
			Expect(gems).To(ConsistOf([]storage.Urn{{
				Ilk:        ilk1,
				Identifier: guy2,
			}, {
				Ilk:        ilk1,
				Identifier: guy3,
			}}))
		})

		It("fetches unique urns from vat_slip + vat_flux + vat_frob + vat_grab events", func() {
			insertVatSlip(ilk1, guy1, 1, db)
			insertVatFlux(ilk1, guy2, guy3, 2, db)
			insertVatFrob(ilk2, guy1, guy1, guy1, 3, db)
			insertVatGrab(ilk2, guy1, guy2, guy1, 4, db)
			// duplicates
			insertVatSlip(ilk1, guy2, 6, db)
			insertVatFlux(ilk2, guy2, guy3, 7, db)
			insertVatFrob(ilk2, guy1, guy1, guy1, 8, db)
			insertVatGrab(ilk1, guy1, guy1, guy1, 9, db)

			gems, err := repository.GetGemKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(gems)).To(Equal(6))
			Expect(gems).To(ConsistOf([]storage.Urn{{
				Ilk:        ilk1,
				Identifier: guy1,
			}, {
				Ilk:        ilk1,
				Identifier: guy2,
			}, {
				Ilk:        ilk1,
				Identifier: guy3,
			}, {
				Ilk:        ilk2,
				Identifier: guy1,
			}, {
				Ilk:        ilk2,
				Identifier: guy2,
			}, {
				Ilk:        ilk2,
				Identifier: guy3,
			}}))
		})

		It("does not return error if no matching rows", func() {
			gemKeys, err := repository.GetGemKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(gemKeys)).To(BeZero())
		})
	})

	Describe("getting ilks", func() {
		It("fetches unique ilks from vat init events", func() {
			insertVatInit(ilk1, 1, db)
			insertVatInit(ilk2, 2, db)
			insertVatInit(ilk2, 3, db)

			ilks, err := repository.GetIlks()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(ilks)).To(Equal(2))
			Expect(ilks).To(ConsistOf(ilk1, ilk2))
		})

		It("does not return error if no matching rows", func() {
			ilks, err := repository.GetIlks()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(ilks)).To(BeZero())
		})
	})

	Describe("getting vat sin keys", func() {
		It("fetches guy from w field of vat grab", func() {
			insertVatGrab(guy1, guy1, guy1, guy2, 1, db)

			sinKeys, err := repository.GetVatSinKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(sinKeys)).To(Equal(1))
			Expect(sinKeys).To(ConsistOf(guy2))
		})

		It("fetches guy from vat heal transaction", func() {
			insertVatHeal(1, transactionFromGuy1, db)

			sinKeys, err := repository.GetVatSinKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(sinKeys)).To(Equal(1))
			Expect(sinKeys).To(ConsistOf(guy1))
		})

		It("fetches guy from u field of vat suck", func() {
			insertVatSuck(guy1, guy2, 0, 1, db)
			sinKeys, repoErr := repository.GetVatSinKeys()

			Expect(repoErr).NotTo(HaveOccurred())
			Expect(len(sinKeys)).To(Equal(1))
			Expect(sinKeys).To(ConsistOf(guy1))
		})

		It("fetches the correct guy when there are multiple transactions in a block", func() {
			insertVatHeal(1, transactionFromGuy1, db)
			unrelatedTransaction := core.TransactionModel{From: "unrelated guy", TxIndex: 15, Value: "0"}
			insertTransaction(1, unrelatedTransaction, db)

			sinKeys, err := repository.GetVatSinKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(sinKeys)).To(Equal(1))
			Expect(sinKeys).To(ConsistOf(guy1))
		})

		It("fetches unique sin keys from vat_grab + vat_heal + vat_suck", func() {
			transactionFromGuy2 := core.TransactionModel{From: guy2, TxIndex: 2, Value: "0"}
			insertVatGrab(guy3, guy3, guy3, guy1, 1, db)
			insertVatHeal(2, transactionFromGuy2, db)
			insertVatSuck(guy3, guy3, 0, 3, db)
			// duplicate
			insertVatGrab(guy2, guy2, guy2, guy2, 4, db)
			insertVatHeal(5, transactionFromGuy2, db)
			insertVatSuck(guy1, guy2, 0, 6, db)

			sinKeys, err := repository.GetVatSinKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(sinKeys)).To(Equal(3))
			Expect(sinKeys).To(ConsistOf(guy1, guy2, guy3))
		})

		It("does not return error if no matching rows", func() {
			sinKeys, err := repository.GetVatSinKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(sinKeys)).To(BeZero())
		})
	})

	Describe("getting vow sin keys", func() {
		It("fetches timestamp from era field of vow flog", func() {
			insertVowFlog(era, 1, db)

			sinKeys, err := repository.GetVowSinKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(sinKeys)).To(Equal(1))
			Expect(sinKeys).To(ConsistOf(era))
		})

		It("fetches timestamp from header of vow fess event", func() {
			insertVowFess(tab, timestamp, 1, db)

			sinKeys, err := repository.GetVowSinKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(sinKeys)).To(Equal(1))
			Expect(sinKeys).To(ConsistOf(strconv.FormatInt(timestamp, 10)))
		})

		It("fetches unique sin keys from vow flog and vow fess header", func() {
			insertVowFlog(era, 1, db)
			insertVowFess(tab, timestamp, 2, db)
			// duplicates
			insertVowFlog(era, 3, db)
			insertVowFess(tab, timestamp, 4, db)

			sinKeys, err := repository.GetVowSinKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(sinKeys)).To(Equal(2))
			Expect(sinKeys).To(ConsistOf(era, strconv.FormatInt(timestamp, 10)))
		})

		It("does not return error if no matching rows", func() {
			sinKeys, err := repository.GetVowSinKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(sinKeys)).To(BeZero())
		})
	})

	Describe("getting urns", func() {
		It("fetches unique urns from vat_frob + vat_grab + vat_fork events", func() {
			insertVatFrob(ilk1, guy1, guy1, guy1, 1, db)
			insertVatFrob(ilk1, guy2, guy1, guy1, 2, db)
			insertVatFrob(ilk2, guy1, guy1, guy1, 3, db)
			insertVatFrob(ilk1, guy1, guy1, guy1, 4, db)
			insertVatGrab(ilk1, guy1, guy1, guy1, 5, db)
			insertVatGrab(ilk1, guy3, guy1, guy1, 6, db)
			insertVatFork(ilk2, guy2, guy3, 7, db)

			urns, err := repository.GetUrns()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(urns)).To(Equal(6))
			Expect(urns).To(ConsistOf([]storage.Urn{{
				Ilk:        ilk1,
				Identifier: guy1,
			}, {
				Ilk:        ilk1,
				Identifier: guy2,
			}, {
				Ilk:        ilk2,
				Identifier: guy1,
			}, {
				Ilk:        ilk1,
				Identifier: guy3,
			}, {
				Ilk:        ilk2,
				Identifier: guy2,
			}, {
				Ilk:        ilk2,
				Identifier: guy3,
			}}))
		})

		It("does not return error if no matching rows", func() {
			urns, err := repository.GetUrns()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(urns)).To(BeZero())
		})
	})
})

func insertVatFold(urn string, blockNumber int64, db *postgres.DB) {
	headerRepository := repositories.NewHeaderRepository(db)
	headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(blockNumber))
	Expect(err).NotTo(HaveOccurred())
	ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
	Expect(err).NotTo(HaveOccurred())
	urnID, err := shared.GetOrCreateUrn(urn, ilkID, db)
	Expect(err).NotTo(HaveOccurred())

	_, execErr := db.Exec(
		`INSERT INTO maker.vat_fold (header_id, urn_id, log_idx, tx_idx)
			VALUES($1, $2, $3, $4)`,
		headerID, urnID, 0, 0,
	)
	Expect(execErr).NotTo(HaveOccurred())
}

func insertVowFlog(era string, blockNumber int64, db *postgres.DB) {
	headerRepository := repositories.NewHeaderRepository(db)
	headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(blockNumber))

	Expect(err).NotTo(HaveOccurred())
	_, execErr := db.Exec(
		`INSERT INTO maker.vow_flog (header_id, era, log_idx, tx_idx)
			VALUES($1, $2, $3, $4)`,
		headerID, era, 0, 0,
	)
	Expect(execErr).NotTo(HaveOccurred())
}

func insertVowFess(tab string, timestamp, blockNumber int64, db *postgres.DB) {
	headerRepository := repositories.NewHeaderRepository(db)
	fakeHeader := fakes.GetFakeHeader(blockNumber)
	fakeHeader.Timestamp = strconv.FormatInt(timestamp, 10)
	// TODO: replace above 2 lines with fakes.GetFakeHeaderWithTimestamp once it's in a versioned release
	headerID, err := headerRepository.CreateOrUpdateHeader(fakeHeader)

	Expect(err).NotTo(HaveOccurred())
	_, execErr := db.Exec(
		`INSERT INTO maker.vow_fess (header_id, tab, log_idx, tx_idx)
			VALUES($1, $2, $3, $4)`,
		headerID, tab, 0, 0,
	)
	Expect(execErr).NotTo(HaveOccurred())
}

func insertTransaction(blockNumber int64, transaction core.TransactionModel, db *postgres.DB) {
	var headerID int64
	err := db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, blockNumber)
	Expect(err).NotTo(HaveOccurred())

	headerRepository := repositories.NewHeaderRepository(db)
	err = headerRepository.CreateTransactions(headerID, []core.TransactionModel{transaction})
	Expect(err).NotTo(HaveOccurred())
}

func insertVatInit(ilk string, blockNumber int64, db *postgres.DB) {
	headerRepository := repositories.NewHeaderRepository(db)
	headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(blockNumber))
	Expect(err).NotTo(HaveOccurred())
	ilkID, err := shared.GetOrCreateIlk(ilk, db)
	Expect(err).NotTo(HaveOccurred())
	_, execErr := db.Exec(
		`INSERT INTO maker.vat_init (header_id, ilk_id, log_idx, tx_idx)
			VALUES($1, $2, $3, $4)`,
		headerID, ilkID, 0, 0,
	)
	Expect(execErr).NotTo(HaveOccurred())
}

func insertVatFlux(ilk, src, dst string, blockNumber int64, db *postgres.DB) {
	headerRepository := repositories.NewHeaderRepository(db)
	headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(blockNumber))
	Expect(err).NotTo(HaveOccurred())
	ilkID, err := shared.GetOrCreateIlk(ilk, db)
	Expect(err).NotTo(HaveOccurred())

	_, execErr := db.Exec(
		`INSERT INTO maker.vat_flux (header_id, ilk_id, src, dst, log_idx, tx_idx)
			VALUES($1, $2, $3, $4, $5, $6)`,
		headerID, ilkID, src, dst, 0, 0,
	)
	Expect(execErr).NotTo(HaveOccurred())
}

func insertVatFork(ilk, src, dst string, blockNumber int64, db *postgres.DB) {
	headerRepository := repositories.NewHeaderRepository(db)
	headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(blockNumber))
	Expect(err).NotTo(HaveOccurred())
	ilkID, err := shared.GetOrCreateIlk(ilk, db)
	Expect(err).NotTo(HaveOccurred())

	_, execErr := db.Exec(
		`INSERT INTO maker.vat_fork (header_id, ilk_id, src, dst, log_idx, tx_idx)
			VALUES($1, $2, $3, $4, $5, $6)`,
		headerID, ilkID, src, dst, 0, 0,
	)
	Expect(execErr).NotTo(HaveOccurred())
}

func insertVatFrob(ilk, urn, v, w string, blockNumber int64, db *postgres.DB) {
	headerRepository := repositories.NewHeaderRepository(db)
	headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(blockNumber))
	Expect(err).NotTo(HaveOccurred())
	ilkID, err := shared.GetOrCreateIlk(ilk, db)
	Expect(err).NotTo(HaveOccurred())
	urnID, err := shared.GetOrCreateUrn(urn, ilkID, db)
	Expect(err).NotTo(HaveOccurred())

	_, execErr := db.Exec(
		`INSERT INTO maker.vat_frob (header_id, urn_id, v, w, log_idx, tx_idx)
			VALUES($1, $2, $3, $4, $5, $6)`,
		headerID, urnID, v, w, 0, 0,
	)
	Expect(execErr).NotTo(HaveOccurred())
}

func insertVatGrab(ilk, urn, v, w string, blockNumber int64, db *postgres.DB) {
	headerRepository := repositories.NewHeaderRepository(db)
	headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(blockNumber))
	Expect(err).NotTo(HaveOccurred())
	ilkID, err := shared.GetOrCreateIlk(ilk, db)
	Expect(err).NotTo(HaveOccurred())
	urnID, err := shared.GetOrCreateUrn(urn, ilkID, db)
	Expect(err).NotTo(HaveOccurred())

	_, execErr := db.Exec(
		`INSERT INTO maker.vat_grab (header_id, urn_id, v, w, log_idx, tx_idx)
			VALUES($1, $2, $3, $4, $5, $6)`,
		headerID, urnID, v, w, 0, 0,
	)
	Expect(execErr).NotTo(HaveOccurred())
}

func insertVatSuck(u, v string, rad int, blockNumber int64, db *postgres.DB) {
	headerRepository := repositories.NewHeaderRepository(db)
	headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(blockNumber))
	Expect(err).NotTo(HaveOccurred())

	_, execErr := db.Exec(
		`INSERT INTO maker.vat_suck (header_id, u, v, rad, log_idx, tx_idx)
			VALUES($1, $2, $3, $4, $5, $6)`,
		headerID, u, v, rad, 0, 0,
	)
	Expect(execErr).NotTo(HaveOccurred())
}

func insertVatHeal(blockNumber int64, transaction core.TransactionModel, db *postgres.DB) {
	headerRepository := repositories.NewHeaderRepository(db)
	headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(blockNumber))
	Expect(err).NotTo(HaveOccurred())
	err = headerRepository.CreateTransactions(headerID, []core.TransactionModel{transaction})
	Expect(err).NotTo(HaveOccurred())
	_, execErr := db.Exec(
		`INSERT INTO maker.vat_heal (header_id, log_idx, tx_idx)
			VALUES($1, $2, $3)`,
		headerID, 0, transaction.TxIndex,
	)
	Expect(execErr).NotTo(HaveOccurred())
}

func insertVatMove(src, dst string, blockNumber int64, db *postgres.DB) {
	headerRepository := repositories.NewHeaderRepository(db)
	headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(blockNumber))
	Expect(err).NotTo(HaveOccurred())
	_, execErr := db.Exec(
		`INSERT INTO maker.vat_move (header_id, src, dst, rad, log_idx, tx_idx)
			VALUES($1, $2, $3, $4, $5, $6)`,
		headerID, src, dst, 0, 0, 0,
	)
	Expect(execErr).NotTo(HaveOccurred())
}

func insertVatSlip(ilk, usr string, blockNumber int64, db *postgres.DB) {
	headerRepository := repositories.NewHeaderRepository(db)
	headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(blockNumber))
	Expect(err).NotTo(HaveOccurred())
	ilkID, err := shared.GetOrCreateIlk(ilk, db)
	Expect(err).NotTo(HaveOccurred())
	_, execErr := db.Exec(
		`INSERT INTO maker.vat_slip (header_id, ilk_id, usr, log_idx, tx_idx)
			VALUES($1, $2, $3, $4, $5)`,
		headerID, ilkID, usr, 0, 0,
	)
	Expect(execErr).NotTo(HaveOccurred())
}
