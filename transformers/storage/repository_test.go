// VulcanizeDB
// Copyright © 2019 Vulcanize

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
	"fmt"
	"math/big"
	"math/rand"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lib/pq"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	query_helper "github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	mcdShared "github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flap"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flop"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	insertFlapKickQuery = `INSERT into maker.flap_kick
		(header_id, bid_id, lot, bid, address_id, log_id)
		VALUES($1, $2::NUMERIC, $3::NUMERIC, $4::NUMERIC, $5, $6)
		ON CONFLICT (header_id, log_id)
		DO UPDATE SET bid_id = $2, lot = $3, bid = $4, address_id = $5;`

	insertClipKickQuery = `INSERT into maker.clip_kick
		(header_id, address_id, log_id, sale_id, top, tab, lot, usr, kpr, coin )
		VALUES($1, $2, $3, $4::NUMERIC, $5::NUMERIC, $6::NUMERIC, $7::NUMERIC, $8, $9, $10::NUMERIC)
		ON CONFLICT (header_id, log_id)
		DO NOTHING;`

	insertFlipKickQuery = `INSERT into maker.flip_kick
		(header_id, bid_id, lot, bid, tab, usr, gal, address_id, log_id)
		VALUES($1, $2::NUMERIC, $3::NUMERIC, $4::NUMERIC, $5::NUMERIC, $6, $7, $8, $9)
		ON CONFLICT (header_id, log_id)
		DO UPDATE SET bid_id = $2, lot = $3, bid = $4, tab = $5, usr = $6, gal = $7, address_id = $8;`

	insertFlopKickQuery = `INSERT INTO maker.flop_kick
		(header_id, bid_id, lot, bid, gal, address_id, log_id)
		VALUES($1, $2::NUMERIC, 0, 0, '', $3, $4)
		ON CONFLICT (header_id, log_id)
		DO NOTHING;`
)

var _ = Describe("Maker storage repository", func() {
	var (
		address             = fakes.FakeAddress.Hex()
		addressID           int64
		addressErr          error
		db                  = test_config.NewTestDB(test_config.NewTestNode())
		repo                storage.IMakerStorageRepository
		ilk1                = common.HexToHash("0x494c4b31").Hex()
		ilk2                = common.HexToHash("0x494c4b32").Hex()
		guy1                = "0x47555931"
		guy2                = "0x47555932"
		guy3                = "0x47555933"
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
		test_config.CleanTestDB(db)
		repo = &storage.MakerStorageRepository{}
		repo.SetDB(db)
		addressID, addressErr = repository.GetOrCreateAddress(db, address)
		Expect(addressErr).NotTo(HaveOccurred())
	})

	Describe("getting flap bid ids", func() {
		var (
			bidID1, bidID2, bidID3, bidID4, bidID5, bidID6 string
		)
		BeforeEach(func() {
			bidID1 = strconv.FormatInt(rand.Int63(), 10)
			bidID2 = strconv.FormatInt(rand.Int63(), 10)
			bidID3 = strconv.FormatInt(rand.Int63(), 10)
			bidID4 = strconv.FormatInt(rand.Int63(), 10)
			bidID5 = strconv.FormatInt(rand.Int63(), 10)
			bidID6 = strconv.FormatInt(rand.Int63(), 10)
		})

		It("fetches unique bid ids from Flap methods", func() {
			insertFlapKick(1, bidID1, addressID, db)
			insertFlapKick(2, bidID1, addressID, db)

			bidIDs, err := repo.GetFlapBidIDs(address)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bidIDs)).To(Equal(1))
			Expect(bidIDs[0]).To(Equal(bidID1))
		})

		It("fetches unique bid ids from flap_kick, tend, deal and yank", func() {
			duplicateBidID := bidID1
			insertFlapKick(1, bidID1, addressID, db)
			insertFlapKicks(2, bidID2, addressID, db)
			insertTend(3, bidID3, addressID, db)
			insertDeal(4, bidID4, addressID, db)
			insertYank(5, bidID5, addressID, db)
			insertYank(6, duplicateBidID, addressID, db)

			bidIDs, err := repo.GetFlapBidIDs(address)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bidIDs)).To(Equal(5))
			Expect(bidIDs).To(ConsistOf(bidID1, bidID2, bidID3, bidID4, bidID5))
		})

		It("fetches bid ids only for the given contract address", func() {
			anotherAddress := address + "1"
			anotherAddressID, addressErr := repository.GetOrCreateAddress(db, anotherAddress)
			Expect(addressErr).NotTo(HaveOccurred())
			insertFlapKick(1, bidID1, addressID, db)
			insertFlapKick(2, bidID2, addressID, db)
			insertTend(3, bidID3, addressID, db)
			insertDeal(4, bidID4, addressID, db)
			insertYank(5, bidID5, addressID, db)
			insertYank(6, bidID6, anotherAddressID, db)

			bidIDs, err := repo.GetFlapBidIDs(address)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bidIDs)).To(Equal(5))
			Expect(bidIDs).To(ConsistOf(bidID1, bidID2, bidID3, bidID4, bidID5))
		})

		It("does not return error if no matching rows", func() {
			bidIDs, err := repo.GetFlapBidIDs(fakes.FakeAddress.Hex())

			Expect(err).NotTo(HaveOccurred())
			Expect(len(bidIDs)).To(BeZero())
		})
	})

	Describe("getting dai keys", func() {
		It("fetches guy from both src and dst field on vat_move", func() {
			insertVatMove(guy1, guy2, 1, db)

			keys, err := repo.GetDaiKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(keys)).To(Equal(2))
			Expect(keys).To(ConsistOf(guy1, guy2))
		})

		It("fetches guy from w field on vat_frob", func() {
			insertVatFrob(ilk1, guy1, guy1, guy2, 1, db)

			keys, err := repo.GetDaiKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(keys)).To(Equal(1))
			Expect(keys).To(ConsistOf(guy2))
		})

		It("fetches guy from vat_heal transaction", func() {
			insertVatHeal(1, transactionFromGuy1, db)

			keys, err := repo.GetDaiKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(keys)).To(Equal(1))
			Expect(keys).To(ConsistOf(guy1))
		})

		It("fetches guy from v field on vat_suck", func() {
			insertVatSuck(guy1, guy2, 0, 1, db)

			daiKeys, repoErr := repo.GetDaiKeys()

			Expect(repoErr).NotTo(HaveOccurred())
			Expect(len(daiKeys)).To(Equal(1))
			Expect(daiKeys).To(ConsistOf(guy2))
		})

		It("fetches guy from u field on vat_fold", func() {
			insertVatFold(guy1, 1, db)

			daiKeys, repoErr := repo.GetDaiKeys()

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

			keys, err := repo.GetDaiKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(keys)).To(Equal(6))
			Expect(keys).To(ConsistOf(guy1, guy2, guy3, guy4, guy5, guy6))
		})

		It("fetches the correct guy when there are multiple transactions in a block", func() {
			insertVatHeal(1, transactionFromGuy1, db)
			unrelatedTransaction := core.TransactionModel{From: "unrelated guy", TxIndex: 15, Value: "0"}
			insertTransaction(1, unrelatedTransaction, db)

			sinKeys, err := repo.GetDaiKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(sinKeys)).To(Equal(1))
			Expect(sinKeys).To(ConsistOf(guy1))
		})

		It("does not return error if no matching rows", func() {
			daiKeys, err := repo.GetDaiKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(daiKeys)).To(BeZero())
		})
	})

	Describe("getting gem keys", func() {
		It("fetches guy from both src and dst field on vat_flux", func() {
			insertVatFlux(ilk1, guy1, guy2, 1, db)

			gems, err := repo.GetGemKeys()

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

			gems, err := repo.GetGemKeys()

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

			gems, err := repo.GetGemKeys()

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
			gemKeys, err := repo.GetGemKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(gemKeys)).To(BeZero())
		})
	})

	Describe("getting ilks", func() {
		It("fetches unique ilks from vat init events", func() {
			insertVatInit(ilk1, 1, db)
			insertVatInit(ilk2, 2, db)
			insertVatInit(ilk2, 3, db)

			ilks, err := repo.GetIlks()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(ilks)).To(Equal(2))
			Expect(ilks).To(ConsistOf(ilk1, ilk2))
		})

		It("does not return error if no matching rows", func() {
			ilks, err := repo.GetIlks()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(ilks)).To(BeZero())
		})
	})

	Describe("getting vat sin keys", func() {
		It("fetches guy from w field of vat grab", func() {
			insertVatGrab(guy1, guy1, guy1, guy2, 1, db)

			sinKeys, err := repo.GetVatSinKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(sinKeys)).To(Equal(1))
			Expect(sinKeys).To(ConsistOf(guy2))
		})

		It("fetches guy from vat heal transaction", func() {
			insertVatHeal(1, transactionFromGuy1, db)

			sinKeys, err := repo.GetVatSinKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(sinKeys)).To(Equal(1))
			Expect(sinKeys).To(ConsistOf(guy1))
		})

		It("fetches guy from u field of vat suck", func() {
			insertVatSuck(guy1, guy2, 0, 1, db)
			sinKeys, repoErr := repo.GetVatSinKeys()

			Expect(repoErr).NotTo(HaveOccurred())
			Expect(len(sinKeys)).To(Equal(1))
			Expect(sinKeys).To(ConsistOf(guy1))
		})

		It("fetches the correct guy when there are multiple transactions in a block", func() {
			insertVatHeal(1, transactionFromGuy1, db)
			unrelatedTransaction := core.TransactionModel{From: "unrelated guy", TxIndex: 15, Value: "0"}
			insertTransaction(1, unrelatedTransaction, db)

			sinKeys, err := repo.GetVatSinKeys()

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

			sinKeys, err := repo.GetVatSinKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(sinKeys)).To(Equal(3))
			Expect(sinKeys).To(ConsistOf(guy1, guy2, guy3))
		})

		It("does not return error if no matching rows", func() {
			sinKeys, err := repo.GetVatSinKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(sinKeys)).To(BeZero())
		})
	})

	Describe("getting vow sin keys", func() {
		It("fetches timestamp from era field of vow flog", func() {
			insertVowFlog(era, 1, db)

			sinKeys, err := repo.GetVowSinKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(sinKeys)).To(Equal(1))
			Expect(sinKeys).To(ConsistOf(era))
		})

		It("fetches timestamp from header of vow fess event", func() {
			insertVowFess(tab, timestamp, 1, db)

			sinKeys, err := repo.GetVowSinKeys()

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

			sinKeys, err := repo.GetVowSinKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(sinKeys)).To(Equal(2))
			Expect(sinKeys).To(ConsistOf(era, strconv.FormatInt(timestamp, 10)))
		})

		It("does not return error if no matching rows", func() {
			sinKeys, err := repo.GetVowSinKeys()

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

			urns, err := repo.GetUrns()

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
			urns, err := repo.GetUrns()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(urns)).To(BeZero())
		})
	})

	Describe("getting clip sales ids", func() {
		var (
			saleID1 string
		)

		BeforeEach(func() {
			saleID1 = strconv.FormatInt(rand.Int63(), 10)
		})

		It("fetches sale ids from a clip kick event", func() {
			insertClipKick(1, saleID1, addressID, db)

			saleIDs, err := repo.GetClipSalesIDs(address)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(saleIDs)).To(Equal(1))
			Expect(saleIDs[0]).To(Equal(saleID1))
		})
	})

	Describe("getting flip bid ids", func() {
		var (
			bidID1, bidID2, bidID3, bidID4, bidID5, bidID6, bidID7 string
			address                                                = fakes.FakeAddress.Hex()
		)

		BeforeEach(func() {
			bidID1 = strconv.FormatInt(rand.Int63(), 10)
			bidID2 = strconv.FormatInt(rand.Int63(), 10)
			bidID3 = strconv.FormatInt(rand.Int63(), 10)
			bidID4 = strconv.FormatInt(rand.Int63(), 10)
			bidID5 = strconv.FormatInt(rand.Int63(), 10)
			bidID6 = strconv.FormatInt(rand.Int63(), 10)
			bidID7 = strconv.FormatInt(rand.Int63(), 10)
		})

		It("fetches unique bid ids from flip methods", func() {
			insertFlipKick(1, bidID1, addressID, db)
			insertFlipKick(2, bidID1, addressID, db)

			bidIDs, err := repo.GetFlipBidIDs(address)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bidIDs)).To(Equal(1))
			Expect(bidIDs[0]).To(Equal(bidID1))
		})

		It("fetches unique bid ids from tick, flip_kick, flip_kicks, tend, dent, deal and yank", func() {
			duplicateBidID := bidID1
			insertTick(1, bidID1, addressID, db)
			insertFlipKick(2, bidID2, addressID, db)
			insertFlipKicks(3, bidID3, addressID, db)
			insertTend(4, bidID4, addressID, db)
			insertDent(5, bidID5, addressID, db)
			insertDeal(6, bidID6, addressID, db)
			insertYank(7, bidID7, addressID, db)
			insertYank(8, duplicateBidID, addressID, db)

			bidIDs, err := repo.GetFlipBidIDs(address)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bidIDs)).To(Equal(7))
			Expect(bidIDs).To(ConsistOf(bidID1, bidID2, bidID3, bidID4, bidID5, bidID6, bidID7))
		})
	})

	Describe("getting flop bid ids", func() {
		var (
			bidID1, bidID2, bidID3, bidID4, bidID5, bidID6 string
			address                                        = fakes.FakeAddress.Hex()
		)

		BeforeEach(func() {
			bidID1 = strconv.FormatInt(rand.Int63(), 10)
			bidID2 = strconv.FormatInt(rand.Int63(), 10)
			bidID3 = strconv.FormatInt(rand.Int63(), 10)
			bidID4 = strconv.FormatInt(rand.Int63(), 10)
			bidID5 = strconv.FormatInt(rand.Int63(), 10)
			bidID6 = strconv.FormatInt(rand.Int63(), 10)
		})

		It("fetches unique flop bid ids from flop methods", func() {
			insertFlopKick(1, bidID1, addressID, db)
			insertFlopKick(2, bidID1, addressID, db)

			bidIDs, err := repo.GetFlopBidIDs(address)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bidIDs)).To(Equal(1))
			Expect(bidIDs[0]).To(Equal(bidID1))
		})

		It("fetches unique bid ids from flop_kick, dent, deal, and yank", func() {
			duplicateBidID := bidID1
			insertFlopKick(1, bidID1, addressID, db)
			insertFlopKicks(2, bidID2, addressID, db)
			insertDent(3, bidID3, addressID, db)
			insertDeal(4, bidID4, addressID, db)
			insertYank(5, bidID5, addressID, db)
			insertYank(6, duplicateBidID, addressID, db)

			bidIDs, err := repo.GetFlopBidIDs(address)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bidIDs)).To(Equal(5))
			Expect(bidIDs).To(ConsistOf(bidID1, bidID2, bidID3, bidID4, bidID5))
		})

		It("fetches bid ids only for the given contract address", func() {
			anotherAddress := address + "1"
			anotherAddressID, addressErr := repository.GetOrCreateAddress(db, anotherAddress)
			Expect(addressErr).NotTo(HaveOccurred())
			insertFlopKick(1, bidID1, addressID, db)
			insertFlopKick(2, bidID2, addressID, db)
			insertDent(3, bidID3, addressID, db)
			insertDeal(4, bidID4, addressID, db)
			insertYank(5, bidID5, addressID, db)
			insertYank(6, bidID6, anotherAddressID, db)

			bidIDs, err := repo.GetFlopBidIDs(address)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bidIDs)).To(Equal(5))
			Expect(bidIDs).To(ConsistOf(bidID1, bidID2, bidID3, bidID4, bidID5))
		})

		It("does not return error if no matching rows", func() {
			bidIDs, err := repo.GetFlopBidIDs(fakes.FakeAddress.Hex())

			Expect(err).NotTo(HaveOccurred())
			Expect(len(bidIDs)).To(BeZero())
		})
	})

	Describe("getting bud keys", func() {
		It("fetches unique addresses from kiss and diss (single and batch) events", func() {
			a1 := common.HexToAddress(test_data.RandomString(40)).Hex()
			a2 := common.HexToAddress(test_data.RandomString(40)).Hex()
			a3 := common.HexToAddress(test_data.RandomString(40)).Hex()
			a4 := common.HexToAddress(test_data.RandomString(40)).Hex()
			a5 := common.HexToAddress(test_data.RandomString(40)).Hex()
			medianAddressID, addressErr := repository.GetOrCreateAddress(db, test_data.MedianEthAddress())
			Expect(addressErr).NotTo(HaveOccurred())
			insertMedianKissSingle(a1, medianAddressID, db)
			insertMedianDissSingle(a2, medianAddressID, db)
			insertMedianKissBatch([]string{a1, a3}, medianAddressID, db)
			insertMedianDissBatch([]string{a2, a4}, medianAddressID, db)
			insertMedianKissSingle(a5, addressID, db)

			aAddresses, err := repo.GetMedianBudAddresses(test_data.MedianEthAddress())
			Expect(err).NotTo(HaveOccurred())

			Expect(aAddresses).To(ConsistOf(a1, a2, a3, a4))
		})
	})

	Describe("getting orcl keys", func() {
		It("fetches unique addresses from lift and drop events", func() {
			a1 := common.HexToAddress(test_data.RandomString(40)).Hex()
			a2 := common.HexToAddress(test_data.RandomString(40)).Hex()
			a3 := common.HexToAddress(test_data.RandomString(40)).Hex()
			a4 := common.HexToAddress(test_data.RandomString(40)).Hex()
			a5 := common.HexToAddress(test_data.RandomString(40)).Hex()
			medianAddressID, addressErr := repository.GetOrCreateAddress(db, test_data.MedianEthAddress())
			Expect(addressErr).NotTo(HaveOccurred())
			insertMedianLiftAddresses([]string{a1, a2, a3}, medianAddressID, db)
			insertMedianDropAddresses([]string{a4, a5}, medianAddressID, db)

			aAddresses, err := repo.GetMedianOrclAddresses(test_data.MedianEthAddress())
			Expect(err).NotTo(HaveOccurred())

			Expect(aAddresses).To(ConsistOf(a1, a2, a3, a4, a5))
		})
	})

	Describe("getting Pot pie users", func() {
		It("gets unique msg senders from Pot join and exit events", func() {
			userAddressOne := common.HexToAddress(test_data.RandomString(40)).Hex()
			userAddressTwo := common.HexToAddress(test_data.RandomString(40)).Hex()
			insertPotPieUser(1, userAddressOne, "maker.pot_join", db)
			insertPotPieUser(2, userAddressTwo, "maker.pot_exit", db)
			insertPotPieUser(3, userAddressTwo, "maker.pot_join", db)

			userAddresses, err := repo.GetPotPieUsers()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(userAddresses)).To(Equal(2))
			Expect(userAddresses).To(ConsistOf(userAddressOne, userAddressTwo))
		})
	})

	Describe("getting vat auth keys", func() {
		It("gets unique vat users", func() {
			userAddressOne := common.HexToAddress(test_data.RandomString(40)).Hex()
			userAddressTwo := common.HexToAddress(test_data.RandomString(40)).Hex()
			insertVatAuthEvent(1, userAddressOne, "maker.vat_rely", db)
			insertVatAuthEvent(2, userAddressTwo, "maker.vat_rely", db)
			insertVatAuthEvent(3, userAddressTwo, "maker.vat_deny", db)

			allUsers, userErr := repo.GetVatWardsAddresses()
			Expect(userErr).NotTo(HaveOccurred())
			Expect(allUsers).To(ConsistOf(userAddressOne, userAddressTwo))
		})
	})

	Describe("getting auth keys", func() {
		It("gets unique rely and deny users and msg.senders for a given contract", func() {
			msgSenderAddressOne := common.HexToAddress(test_data.RandomString(40)).Hex()
			msgSenderAddressTwo := common.HexToAddress(test_data.RandomString(40)).Hex()
			userAddressOne := common.HexToAddress(test_data.RandomString(40)).Hex()
			userAddressTwo := common.HexToAddress(test_data.RandomString(40)).Hex()
			insertAuthEvent(1, test_data.Cat100Address(), msgSenderAddressOne, userAddressOne, "maker.rely", db)
			insertAuthEvent(2, test_data.VowAddress(), msgSenderAddressOne, userAddressTwo, "maker.rely", db)
			insertAuthEvent(3, test_data.VowAddress(), msgSenderAddressTwo, userAddressTwo, "maker.deny", db)

			catUserAddresses, catUserErr := repo.GetWardsAddresses(test_data.Cat100Address())
			Expect(catUserErr).NotTo(HaveOccurred())
			Expect(catUserAddresses).To(ConsistOf(msgSenderAddressOne, userAddressOne))

			vowUserAddresses, vowUserErr := repo.GetWardsAddresses(test_data.VowAddress())
			Expect(vowUserErr).NotTo(HaveOccurred())
			Expect(vowUserAddresses).To(ConsistOf(msgSenderAddressOne, msgSenderAddressTwo, userAddressTwo))
		})
	})

	Describe("getting CDPIs", func() {
		It("returns string version of ints ranging from 1 to the max CDPI in the table", func() {
			insertCdpManagerCdpi(int64(rand.Int()), 2, db)
			insertCdpManagerCdpi(int64(rand.Int()), 5, db)
			insertCdpManagerCdpi(int64(rand.Int()), 3, db)

			cdpis, err := repo.GetCdpis()
			Expect(err).NotTo(HaveOccurred())

			Expect(len(cdpis)).To(Equal(5))
			Expect(cdpis).To(ConsistOf([]string{"1", "2", "3", "4", "5"}))
		})

		It("returns empty slice if table is empty", func() {
			cdpis, err := repo.GetCdpis()
			Expect(err).NotTo(HaveOccurred())

			Expect(cdpis).To(BeEmpty())
		})
	})
})

func insertFlapKick(blockNumber int64, bidID string, contractAddressID int64, db *postgres.DB) {
	//inserting a flap kick log event record
	headerID := insertHeader(db, blockNumber)

	flapKickLog := test_data.CreateTestLog(headerID, db)
	_, insertErr := db.Exec(insertFlapKickQuery,
		headerID, bidID, 0, 0, contractAddressID, flapKickLog.ID,
	)
	Expect(insertErr).NotTo(HaveOccurred())
}

func insertFlapKicks(blockNumber int64, kicks string, contractAddressID int64, db *postgres.DB) {
	//inserting a flap kicks storage record
	headerID := insertHeader(db, blockNumber)
	diffID := test_helpers.CreateFakeDiffRecord(db)
	_, insertErr := db.Exec(flap.InsertKicksQuery,
		diffID, headerID, contractAddressID, kicks,
	)
	Expect(insertErr).NotTo(HaveOccurred())
}

func insertTick(blockNumber int64, bidID string, contractAddressID int64, db *postgres.DB) {
	// tick event record
	headerID := insertHeader(db, blockNumber)
	flapTickLog := test_data.CreateTestLog(headerID, db)

	msgSender := shared.GetChecksumAddressString(test_data.FlipTickEventLog.Log.Topics[1].Hex())
	msgSenderID, msgSenderErr := repository.GetOrCreateAddress(db, msgSender)
	Expect(msgSenderErr).NotTo(HaveOccurred())

	_, insertErr := db.Exec(`INSERT INTO maker.tick (header_id, bid_id, address_id, log_id, msg_sender)
				VALUES($1, $2::NUMERIC, $3, $4, $5)`,
		headerID, bidID, contractAddressID, flapTickLog.ID, msgSenderID,
	)
	Expect(insertErr).NotTo(HaveOccurred())
}

func insertClipKick(blockNumber int64, saleID string, contractAddressID int64, db *postgres.DB) {
	// flip kick event record
	headerID := insertHeader(db, blockNumber)
	log := test_data.CreateTestLog(headerID, db)
	_, insertErr := db.Exec(insertClipKickQuery,
		headerID, contractAddressID, log.ID, saleID, 0, 0, 0, 0, 0, 0,
	)
	Expect(insertErr).NotTo(HaveOccurred())
}

func insertFlipKick(blockNumber int64, bidID string, contractAddressID int64, db *postgres.DB) {
	// flip kick event record
	headerID := insertHeader(db, blockNumber)
	log := test_data.CreateTestLog(headerID, db)
	_, insertErr := db.Exec(insertFlipKickQuery,
		headerID, bidID, 0, 0, 0, "", "", contractAddressID, log.ID,
	)
	Expect(insertErr).NotTo(HaveOccurred())
}

func insertFlipKicks(blockNumber int64, kicks string, contractAddressID int64, db *postgres.DB) {
	// flip kicks storage record
	headerID := insertHeader(db, blockNumber)
	diffID := test_helpers.CreateFakeDiffRecord(db)
	_, insertErr := db.Exec(flip.InsertFlipKicksQuery,
		diffID, headerID, contractAddressID, kicks,
	)
	Expect(insertErr).NotTo(HaveOccurred())
}

func insertFlopKick(blockNumber int64, bidID string, contractAddressID int64, db *postgres.DB) {
	// inserting a flop kick log event record
	headerID := insertHeader(db, blockNumber)
	flopKickLog := test_data.CreateTestLog(headerID, db)
	_, insertErr := db.Exec(insertFlopKickQuery, headerID, bidID, contractAddressID, flopKickLog.ID)
	Expect(insertErr).NotTo(HaveOccurred())
}

func insertFlopKicks(blockNumber int64, kicks string, contractAddressID int64, db *postgres.DB) {
	// inserting a flop kicks storage record
	diffID := test_helpers.CreateFakeDiffRecord(db)
	headerID := insertHeader(db, blockNumber)
	_, insertErr := db.Exec(flop.InsertFlopKicksQuery, diffID, headerID, contractAddressID, kicks)
	Expect(insertErr).NotTo(HaveOccurred())
}

func insertTend(blockNumber int64, bidID string, contractAddressID int64, db *postgres.DB) {
	headerID := insertHeader(db, blockNumber)
	tendLog := test_data.CreateTestLogFromEventLog(headerID, test_data.TendEventLog.Log, db)
	tendModel := test_data.TendModel()
	tendModel.ColumnValues[event.HeaderFK] = headerID
	tendModel.ColumnValues[event.LogFK] = tendLog.ID
	tendModel.ColumnValues[constants.BidIDColumn] = bidID
	tendModel.ColumnValues[event.AddressFK] = contractAddressID
	test_data.AssignMessageSenderID(tendLog, tendModel, db)
	tendErr := event.PersistModels([]event.InsertionModel{tendModel}, db)
	Expect(tendErr).NotTo(HaveOccurred())
}

func insertDent(blockNumber int64, bidID string, contractAddressID int64, db *postgres.DB) {
	headerID := insertHeader(db, blockNumber)
	dentLog := test_data.CreateTestLog(headerID, db)

	msgSender := shared.GetChecksumAddressString(test_data.DentEventLog.Log.Topics[1].Hex())
	msgSenderID, msgSenderErr := repository.GetOrCreateAddress(db, msgSender)
	Expect(msgSenderErr).NotTo(HaveOccurred())

	_, err := db.Exec(`INSERT into maker.dent (header_id, bid_id, lot, bid, msg_sender, address_id, log_id)
		VALUES($1, $2::NUMERIC, $3::NUMERIC, $4::NUMERIC, $5, $6, $7)`,
		headerID, bidID, 0, 0, msgSenderID, contractAddressID, dentLog.ID,
	)
	Expect(err).NotTo(HaveOccurred())
}

func insertDeal(blockNumber int64, bidID string, contractAddressID int64, db *postgres.DB) {
	headerID := insertHeader(db, blockNumber)
	dealLog := test_data.CreateTestLog(headerID, db)
	msgSenderID, msgSenderErr := repository.GetOrCreateAddress(db, test_data.DealEventLog.Log.Topics[1].Hex())
	Expect(msgSenderErr).NotTo(HaveOccurred())
	_, err := db.Exec(`INSERT into maker.deal (header_id, bid_id, address_id, log_id, msg_sender)
		VALUES($1, $2::NUMERIC, $3, $4, $5)`,
		headerID, bidID, contractAddressID, dealLog.ID, msgSenderID,
	)
	Expect(err).NotTo(HaveOccurred())
}

func insertYank(blockNumber int64, bidID string, contractAddressID int64, db *postgres.DB) {
	headerID := insertHeader(db, blockNumber)
	yankLog := test_data.CreateTestLog(headerID, db)
	msgSenderID, msgSenderErr := repository.GetOrCreateAddress(db, test_data.YankEventLog.Log.Topics[1].Hex())
	Expect(msgSenderErr).NotTo(HaveOccurred())
	_, err := db.Exec(`INSERT into maker.yank (header_id, bid_id, address_id, log_id, msg_sender)
		VALUES($1, $2::NUMERIC, $3, $4, $5)`,
		headerID, bidID, contractAddressID, yankLog.ID, msgSenderID,
	)
	Expect(err).NotTo(HaveOccurred())
}

func insertMedianKissSingle(a string, contractAddressID int64, db *postgres.DB) {
	headerID := insertHeader(db, rand.Int63n(1000))
	kissLog := test_data.CreateTestLog(headerID, db)
	addressID, addressErr := repository.GetOrCreateAddress(db, a)
	Expect(addressErr).NotTo(HaveOccurred())
	_, err := db.Exec(`INSERT into maker.median_kiss_single (header_id, address_id, log_id, msg_sender, a)
		VALUES($1, $2::NUMERIC, $3, $4, $4)`,
		headerID, contractAddressID, kissLog.ID, addressID)
	Expect(err).NotTo(HaveOccurred())
}

func insertMedianDissSingle(a string, contractAddressID int64, db *postgres.DB) {
	headerID := insertHeader(db, rand.Int63n(1000))
	dissLog := test_data.CreateTestLog(headerID, db)
	addressID, addressErr := repository.GetOrCreateAddress(db, a)
	Expect(addressErr).NotTo(HaveOccurred())
	_, err := db.Exec(`INSERT into maker.median_diss_single (header_id, address_id, log_id, msg_sender, a)
		VALUES($1, $2::NUMERIC, $3, $4, $4)`,
		headerID, contractAddressID, dissLog.ID, addressID)
	Expect(err).NotTo(HaveOccurred())
}

func insertMedianKissBatch(a []string, contractAddressID int64, db *postgres.DB) {
	headerID := insertHeader(db, rand.Int63n(1000))
	kissLog := test_data.CreateTestLog(headerID, db)
	msgSenderID, addressErr := repository.GetOrCreateAddress(db, a[0])
	Expect(addressErr).NotTo(HaveOccurred())
	_, err := db.Exec(`INSERT into maker.median_kiss_batch (header_id, address_id, log_id, msg_sender, a_length, a)
		VALUES($1, $2::NUMERIC, $3, $4, $5, $6)`,
		headerID, contractAddressID, kissLog.ID, msgSenderID, len(a), pq.Array(a))
	Expect(err).NotTo(HaveOccurred())
}

func insertMedianDissBatch(a []string, contractAddressID int64, db *postgres.DB) {
	headerID := insertHeader(db, rand.Int63n(1000))
	dissLog := test_data.CreateTestLog(headerID, db)
	msgSenderID, addressErr := repository.GetOrCreateAddress(db, a[0])
	Expect(addressErr).NotTo(HaveOccurred())
	_, err := db.Exec(`INSERT into maker.median_diss_batch (header_id, address_id, log_id, msg_sender, a_length, a)
		VALUES($1, $2::NUMERIC, $3, $4, $5, $6)`,
		headerID, contractAddressID, dissLog.ID, msgSenderID, len(a), pq.Array(a))
	Expect(err).NotTo(HaveOccurred())
}

func insertMedianLiftAddresses(a []string, contractAddressID int64, db *postgres.DB) {
	headerID := insertHeader(db, rand.Int63n(1000))
	liftLog := test_data.CreateTestLog(headerID, db)
	msgSenderID, addressErr := repository.GetOrCreateAddress(db, a[0])
	Expect(addressErr).NotTo(HaveOccurred())
	_, err := db.Exec(`INSERT INTO maker.median_lift (header_id, address_id, log_id, msg_sender, a_length, a)
		VALUES($1, $2::NUMERIC, $3, $4, $5, $6)`,
		headerID, contractAddressID, liftLog.ID, msgSenderID, len(a), pq.Array(a))
	Expect(err).NotTo(HaveOccurred())
}

func insertMedianDropAddresses(a []string, contractAddressID int64, db *postgres.DB) {
	headerID := insertHeader(db, rand.Int63n(1000))
	dropLog := test_data.CreateTestLog(headerID, db)
	msgSenderID, addressErr := repository.GetOrCreateAddress(db, a[0])
	Expect(addressErr).NotTo(HaveOccurred())
	_, err := db.Exec(`INSERT INTO maker.median_drop (header_id, address_id, log_id, msg_sender, a_length, a)
		VALUES($1, $2::NUMERIC, $3, $4, $5, $6)`,
		headerID, contractAddressID, dropLog.ID, msgSenderID, len(a), pq.Array(a))
	Expect(err).NotTo(HaveOccurred())
}

func insertCdpManagerCdpi(blockNumber int64, cdpi int, db *postgres.DB) {
	headerID := insertHeader(db, blockNumber)
	diffID := test_helpers.CreateFakeDiffRecordWithHeader(db, fakes.GetFakeHeader(blockNumber))
	_, err := db.Exec(`INSERT INTO maker.cdp_manager_cdpi (diff_id, header_id, cdpi)
		VALUES($1, $2, $3::NUMERIC)`,
		diffID, headerID, cdpi)
	Expect(err).NotTo(HaveOccurred())
}

func insertVatFold(urn string, blockNumber int64, db *postgres.DB) {
	headerID := insertHeader(db, blockNumber)
	vatFoldLog := test_data.CreateTestLog(headerID, db)
	ilkID, ilkErr := mcdShared.GetOrCreateIlk(query_helper.FakeIlk.Hex, db)
	Expect(ilkErr).NotTo(HaveOccurred())

	_, execErr := db.Exec(
		`INSERT INTO maker.vat_fold (header_id, log_id, ilk_id, u)
			VALUES($1, $2, $3, $4)`,
		headerID, vatFoldLog.ID, ilkID, urn,
	)
	Expect(execErr).NotTo(HaveOccurred())
}

func insertVowFlog(era string, blockNumber int64, db *postgres.DB) {
	headerID := insertHeader(db, blockNumber)

	msgSenderID, msgSenderErr := repository.GetOrCreateAddress(db, test_data.VowFlogEventLog.Log.Topics[1].Hex())
	Expect(msgSenderErr).NotTo(HaveOccurred())

	vowFlogLog := test_data.CreateTestLog(headerID, db)
	_, execErr := db.Exec(
		`INSERT INTO maker.vow_flog (header_id, era, log_id, msg_sender)
			VALUES($1, $2, $3, $4)`,
		headerID, era, vowFlogLog.ID, msgSenderID,
	)
	Expect(execErr).NotTo(HaveOccurred())
}

func insertVowFess(tab string, timestamp, blockNumber int64, db *postgres.DB) {
	headerRepository := repositories.NewHeaderRepository(db)
	fakeHeader := fakes.GetFakeHeaderWithTimestamp(timestamp, blockNumber)
	headerID, err := headerRepository.CreateOrUpdateHeader(fakeHeader)
	vowFessLog := test_data.CreateTestLog(headerID, db)
	Expect(err).NotTo(HaveOccurred())

	msgSenderID, msgSenderErr := repository.GetOrCreateAddress(db, test_data.VowFessEventLog.Log.Topics[1].Hex())
	Expect(msgSenderErr).NotTo(HaveOccurred())

	_, execErr := db.Exec(
		`INSERT INTO maker.vow_fess (header_id, msg_sender, tab, log_id)
			VALUES($1, $2, $3, $4)`,
		headerID, msgSenderID, tab, vowFessLog.ID,
	)
	Expect(execErr).NotTo(HaveOccurred())
}

func insertVatInit(ilk string, blockNumber int64, db *postgres.DB) {
	headerID := insertHeader(db, blockNumber)
	vatInitLog := test_data.CreateTestLog(headerID, db)
	ilkID, err := mcdShared.GetOrCreateIlk(ilk, db)
	Expect(err).NotTo(HaveOccurred())
	_, execErr := db.Exec(
		`INSERT INTO maker.vat_init (header_id, ilk_id, log_id)
			VALUES($1, $2, $3)`,
		headerID, ilkID, vatInitLog.ID,
	)
	Expect(execErr).NotTo(HaveOccurred())
}

func insertVatFlux(ilk, src, dst string, blockNumber int64, db *postgres.DB) {
	headerID := insertHeader(db, blockNumber)
	vatFluxLog := test_data.CreateTestLog(headerID, db)
	ilkID, err := mcdShared.GetOrCreateIlk(ilk, db)
	Expect(err).NotTo(HaveOccurred())

	_, execErr := db.Exec(
		`INSERT INTO maker.vat_flux (header_id, ilk_id, src, dst, log_id)
			VALUES($1, $2, $3, $4, $5)`,
		headerID, ilkID, src, dst, vatFluxLog.ID,
	)
	Expect(execErr).NotTo(HaveOccurred())
}

func insertVatFork(ilk, src, dst string, blockNumber int64, db *postgres.DB) {
	headerID := insertHeader(db, blockNumber)
	vatForkLog := test_data.CreateTestLog(headerID, db)
	ilkID, err := mcdShared.GetOrCreateIlk(ilk, db)
	Expect(err).NotTo(HaveOccurred())

	_, execErr := db.Exec(
		`INSERT INTO maker.vat_fork (header_id, ilk_id, src, dst, log_id)
			VALUES($1, $2, $3, $4, $5)`,
		headerID, ilkID, src, dst, vatForkLog.ID,
	)
	Expect(execErr).NotTo(HaveOccurred())
}

func insertVatFrob(ilk, urn, v, w string, blockNumber int64, db *postgres.DB) {
	headerID := insertHeader(db, blockNumber)
	vatFrobLog := test_data.CreateTestLog(headerID, db)
	urnID, err := mcdShared.GetOrCreateUrn(urn, ilk, db)
	Expect(err).NotTo(HaveOccurred())

	_, execErr := db.Exec(
		`INSERT INTO maker.vat_frob (header_id, urn_id, v, w, log_id)
			VALUES($1, $2, $3, $4, $5)`,
		headerID, urnID, v, w, vatFrobLog.ID,
	)
	Expect(execErr).NotTo(HaveOccurred())
}

func insertVatGrab(ilk, urn, v, w string, blockNumber int64, db *postgres.DB) {
	headerID := insertHeader(db, blockNumber)
	vatGrabLog := test_data.CreateTestLog(headerID, db)
	urnID, err := mcdShared.GetOrCreateUrn(urn, ilk, db)
	Expect(err).NotTo(HaveOccurred())

	_, execErr := db.Exec(
		`INSERT INTO maker.vat_grab (header_id, urn_id, v, w, log_id)
			VALUES($1, $2, $3, $4, $5)`,
		headerID, urnID, v, w, vatGrabLog.ID,
	)
	Expect(execErr).NotTo(HaveOccurred())
}

func insertVatSuck(u, v string, rad int, blockNumber int64, db *postgres.DB) {
	headerID := insertHeader(db, blockNumber)
	vatSuckLog := test_data.CreateTestLog(headerID, db)
	_, execErr := db.Exec(
		`INSERT INTO maker.vat_suck (header_id, u, v, rad, log_id)
			VALUES($1, $2, $3, $4, $5)`,
		headerID, u, v, rad, vatSuckLog.ID,
	)
	Expect(execErr).NotTo(HaveOccurred())
}

func insertVatHeal(blockNumber int64, transaction core.TransactionModel, db *postgres.DB) {
	// TODO: abstract to not init a new repo on every call
	headerRespository := repositories.NewHeaderRepository(db)
	headerID, insertHeaderErr := headerRespository.CreateOrUpdateHeader(fakes.GetFakeHeader(blockNumber))
	Expect(insertHeaderErr).NotTo(HaveOccurred())
	log := types.Log{TxIndex: uint(transaction.TxIndex), BlockNumber: uint64(blockNumber)}
	vatHealLogs := test_data.CreateLogs(headerID, []types.Log{log}, db)
	Expect(len(vatHealLogs)).To(Equal(1))
	insertTransaction(blockNumber, transaction, db)
	_, execErr := db.Exec(
		`INSERT INTO maker.vat_heal (header_id, log_id)
			VALUES($1, $2)`,
		headerID, vatHealLogs[0].ID,
	)
	Expect(execErr).NotTo(HaveOccurred())
}

func insertVatMove(src, dst string, blockNumber int64, db *postgres.DB) {
	headerID := insertHeader(db, blockNumber)
	vatMoveLog := test_data.CreateTestLog(headerID, db)
	_, execErr := db.Exec(
		`INSERT INTO maker.vat_move (header_id, src, dst, rad, log_id)
			VALUES($1, $2, $3, $4, $5)`,
		headerID, src, dst, 0, vatMoveLog.ID,
	)
	Expect(execErr).NotTo(HaveOccurred())
}

func insertVatSlip(ilk, usr string, blockNumber int64, db *postgres.DB) {
	headerID := insertHeader(db, blockNumber)
	vatSlipLog := test_data.CreateTestLog(headerID, db)
	ilkID, err := mcdShared.GetOrCreateIlk(ilk, db)
	Expect(err).NotTo(HaveOccurred())
	_, execErr := db.Exec(
		`INSERT INTO maker.vat_slip (header_id, ilk_id, usr, log_id)
				VALUES($1, $2, $3, $4)`,
		headerID, ilkID, usr, vatSlipLog.ID,
	)
	Expect(execErr).NotTo(HaveOccurred())
}

func insertVatAuthEvent(blockNumber int64, userAddress, tableName string, db *postgres.DB) {
	headerID := insertHeader(db, blockNumber)
	log := test_data.CreateTestLog(headerID, db)

	userAddressID, userAddressErr := repository.GetOrCreateAddress(db, userAddress)
	Expect(userAddressErr).NotTo(HaveOccurred())

	insertAuthEventQuery := fmt.Sprintf(`INSERT INTO %s (header_id, log_id, usr) VALUES ($1, $2, $3)`, tableName)
	_, insertErr := db.Exec(insertAuthEventQuery, headerID, log.ID, userAddressID)
	Expect(insertErr).NotTo(HaveOccurred())
}

func insertAuthEvent(blockNumber int64, contractAddress, msgSenderAddress, userAddress, tableName string, db *postgres.DB) {
	headerID := insertHeader(db, blockNumber)
	log := test_data.CreateTestLog(headerID, db)
	contractAddressID, contractAddressErr := repository.GetOrCreateAddress(db, contractAddress)
	Expect(contractAddressErr).NotTo(HaveOccurred())

	msgSenderAddressID, msgSenderAddressErr := repository.GetOrCreateAddress(db, msgSenderAddress)
	Expect(msgSenderAddressErr).NotTo(HaveOccurred())
	userAddressID, userAddressErr := repository.GetOrCreateAddress(db, userAddress)
	Expect(userAddressErr).NotTo(HaveOccurred())

	insertAuthEventQuery := fmt.Sprintf(`INSERT INTO %s (header_id, log_id, address_id, msg_sender, usr) VALUES ($1, $2, $3, $4, $5)`, tableName)
	_, insertErr := db.Exec(insertAuthEventQuery, headerID, log.ID, contractAddressID, msgSenderAddressID, userAddressID)
	Expect(insertErr).NotTo(HaveOccurred())
}

func insertPotPieUser(blockNumber int64, userAddress, tableName string, db *postgres.DB) {
	headerID := insertHeader(db, blockNumber)
	log := test_data.CreateTestLog(headerID, db)
	userAddressID, addressErr := repository.GetOrCreateAddress(db, userAddress)
	Expect(addressErr).NotTo(HaveOccurred())

	insertMsgSenderQuery := fmt.Sprintf(`INSERT INTO %s (header_id, log_id, msg_sender, wad) VALUES ($1, $2, $3, $4)`, tableName)
	_, insertErr := db.Exec(insertMsgSenderQuery, headerID, log.ID, userAddressID, rand.Int())
	Expect(insertErr).NotTo(HaveOccurred())
}

func insertHeader(db *postgres.DB, blockNumber int64) int64 {
	headerRepository := repositories.NewHeaderRepository(db)
	headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(blockNumber))
	Expect(err).NotTo(HaveOccurred())
	return headerID
}

func insertTransaction(blockNumber int64, transaction core.TransactionModel, db *postgres.DB) {
	var headerID int64
	err := db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, blockNumber)
	Expect(err).NotTo(HaveOccurred())

	headerRepository := repositories.NewHeaderRepository(db)
	err = headerRepository.CreateTransactions(headerID, []core.TransactionModel{transaction})
	Expect(err).NotTo(HaveOccurred())
}
