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
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/storage"
)

var _ = Describe("Maker storage repository", func() {
	var (
		db         *postgres.DB
		repository storage.IMakerStorageRepository
		ilk1       = "494c4b31" // ILK1
		ilk2       = "494c4b32" // ILK2
		guy1       = "47555931" // GUY1
		guy2       = "47555932" // GUY2
		guy3       = "47555933" // GUY3
		era        = big.NewInt(0).SetBytes(common.FromHex("0x000000000000000000000000000000000000000000000000000000005bb48864")).String()
		tab        = big.NewInt(0).SetBytes(common.FromHex("0x0000000000000000000000000000000000000000000002544faa778090e00000")).String()
		timestamp  = int64(1538558053)
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

		It("fetches guy from v field on vat_heal", func() {
			insertVatHeal(guy2, guy1, 1, db)

			keys, err := repository.GetDaiKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(keys)).To(Equal(1))
			Expect(keys).To(ConsistOf(guy1))
		})

		It("fetches unique guys from vat_move + vat_frob + vat_heal + vat_fold", func() {
			guy4 := "47555934"
			guy5 := "47555935"
			guy6 := "47555936"
			insertVatMove(guy1, guy2, 1, db)
			insertVatFrob(ilk1, guy1, guy1, guy3, 2, db)
			insertVatHeal(guy6, guy4, 3, db)
			insertVatFold(guy5, 4, db)
			// duplicates
			insertVatMove(guy3, guy1, 5, db)
			insertVatFrob(ilk2, guy2, guy2, guy5, 6, db)
			insertVatHeal(guy6, guy2, 7, db)
			insertVatFold(guy4, 8, db)

			keys, err := repository.GetDaiKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(keys)).To(Equal(5))
			Expect(keys).To(ConsistOf(guy1, guy2, guy3, guy4, guy5))
		})

		It("does not return error if no matching rows", func() {
			daiKeys, err := repository.GetDaiKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(daiKeys)).To(BeZero())
		})
	})

	Describe("getting max flip", func() {
		It("fetches the max flip", func() {
			insertCatNFlip("1", db)
			insertCatNFlip("3", db)
			insertCatNFlip("2", db)

			maxFlip, err := repository.GetMaxFlip()

			Expect(err).NotTo(HaveOccurred())
			Expect(maxFlip).To(Equal(int64(3)))
		})

		It("returns ErrNoFlips if no flips from which to draw max", func() {
			_, err := repository.GetMaxFlip()

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(storage.ErrNoFlips))
		})
	})

	Describe("getting gem keys", func() {
		It("fetches guy from both src and dst field on vat_flux", func() {
			insertVatFlux(ilk1, guy1, guy2, 1, db)

			gems, err := repository.GetGemKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(gems)).To(Equal(2))
			Expect(gems).To(ConsistOf([]storage.Urn{{
				Ilk: ilk1,
				Guy: guy1,
			}, {
				Ilk: ilk1,
				Guy: guy2,
			}}))
		})

		It("fetches guy from v field on vat_frob + vat_grab", func() {
			insertVatFrob(ilk1, guy1, guy2, guy1, 1, db)
			insertVatGrab(ilk1, guy1, guy3, guy1, 2, db)

			gems, err := repository.GetGemKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(gems)).To(Equal(2))
			Expect(gems).To(ConsistOf([]storage.Urn{{
				Ilk: ilk1,
				Guy: guy2,
			}, {
				Ilk: ilk1,
				Guy: guy3,
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
				Ilk: ilk1,
				Guy: guy1,
			}, {
				Ilk: ilk1,
				Guy: guy2,
			}, {
				Ilk: ilk1,
				Guy: guy3,
			}, {
				Ilk: ilk2,
				Guy: guy1,
			}, {
				Ilk: ilk2,
				Guy: guy2,
			}, {
				Ilk: ilk2,
				Guy: guy3,
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

		It("fetches guy from u field of vat heal", func() {
			insertVatHeal(guy1, guy2, 1, db)

			sinKeys, err := repository.GetVatSinKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(sinKeys)).To(Equal(1))
			Expect(sinKeys).To(ConsistOf(guy1))
		})

		It("fetches unique sin keys from vat_grab + vat_heal", func() {
			insertVatGrab(guy3, guy3, guy3, guy1, 1, db)
			insertVatHeal(guy2, guy3, 2, db)
			// duplicates
			insertVatGrab(guy2, guy2, guy2, guy2, 3, db)
			insertVatHeal(guy1, guy2, 4, db)

			sinKeys, err := repository.GetVatSinKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(sinKeys)).To(Equal(2))
			Expect(sinKeys).To(ConsistOf(guy1, guy2))
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
		It("fetches unique urns from vat_frob + vat_grab events", func() {
			insertVatFrob(ilk1, guy1, guy1, guy1, 1, db)
			insertVatFrob(ilk1, guy2, guy1, guy1, 2, db)
			insertVatFrob(ilk2, guy1, guy1, guy1, 3, db)
			insertVatFrob(ilk1, guy1, guy1, guy1, 4, db)
			insertVatGrab(ilk1, guy1, guy1, guy1, 5, db)
			insertVatGrab(ilk1, guy3, guy1, guy1, 6, db)

			urns, err := repository.GetUrns()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(urns)).To(Equal(4))
			Expect(urns).To(ConsistOf([]storage.Urn{{
				Ilk: ilk1,
				Guy: guy1,
			}, {
				Ilk: ilk1,
				Guy: guy2,
			}, {
				Ilk: ilk2,
				Guy: guy1,
			}, {
				Ilk: ilk1,
				Guy: guy3,
			}}))
		})

		It("does not return error if no matching rows", func() {
			urns, err := repository.GetUrns()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(urns)).To(BeZero())
		})
	})
})

func insertCatNFlip(nflip string, db *postgres.DB) {
	_, err := db.Exec(`INSERT INTO maker.cat_nflip (nflip) VALUES ($1)`, nflip)
	Expect(err).NotTo(HaveOccurred())
}

func insertVatFold(urn string, blockNumber int64, db *postgres.DB) {
	headerRepository := repositories.NewHeaderRepository(db)
	headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(blockNumber))
	Expect(err).NotTo(HaveOccurred())
	ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk, db)
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

func insertVatHeal(urn, v string, blockNumber int64, db *postgres.DB) {
	headerRepository := repositories.NewHeaderRepository(db)
	headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(blockNumber))
	Expect(err).NotTo(HaveOccurred())
	_, execErr := db.Exec(
		`INSERT INTO maker.vat_heal (header_id, urn, v, log_idx, tx_idx)
			VALUES($1, $2, $3, $4, $5)`,
		headerID, urn, v, 0, 0,
	)
	Expect(execErr).NotTo(HaveOccurred())
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
