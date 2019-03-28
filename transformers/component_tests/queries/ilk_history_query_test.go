package queries

import (
	"database/sql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/storage/cat"
	"github.com/vulcanize/mcd_transformers/transformers/storage/jug"
	"github.com/vulcanize/mcd_transformers/transformers/storage/pit"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vat"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
)

var _ = Describe("Ilk State History Query", func() {
	var (
		db                       *postgres.DB
		blockOne                 = rand.Int()
		blockTwo                 = blockOne + 1
		vatRepository            vat.VatStorageRepository
		pitRepository            pit.PitStorageRepository
		catRepository            cat.CatStorageRepository
		dripRepository           jug.JugStorageRepository
		fakeIlk                  = test_helpers.FakeIlk
		blockOneIlkState         map[string]string
		blockTwoIlkState         map[string]string
		expectedBlockOneIlkState test_helpers.IlkState
		expectedBlockTwoIlkState test_helpers.IlkState
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		vatRepository.SetDB(db)
		pitRepository.SetDB(db)
		catRepository.SetDB(db)
		dripRepository.SetDB(db)
		headerRepository = repositories.NewHeaderRepository(db)

		blockOneHeader = fakes.GetFakeHeader(int64(blockOne))
		_, err := headerRepository.CreateOrUpdateHeader(blockOneHeader)
		Expect(err).NotTo(HaveOccurred())

		blockTwoHeader = fakes.GetFakeHeader(int64(blockTwo))
		blockTwoHeader.Timestamp = blockTwoHeader.Timestamp + "2"
		blockTwoHeader.Hash = "block2Hash"
		_, err = headerRepository.CreateOrUpdateHeader(blockTwoHeader)
		Expect(err).NotTo(HaveOccurred())

		blockOneIlkState = test_helpers.GetIlkState(0)

		test_helpers.CreateVatRecords(blockOneHeader, blockOneIlkState, test_helpers.FakeIlkVatMetadatas, vatRepository)
		test_helpers.CreatePitRecords(blockOneHeader, blockOneIlkState, test_helpers.FakeIlkPitMetadatas, pitRepository)
		test_helpers.CreateCatRecords(blockOneHeader, blockOneIlkState, test_helpers.FakeIlkCatMetadatas, catRepository)
		test_helpers.CreateJugRecords(blockOneHeader, blockOneIlkState, test_helpers.FakeIlkJugMetadatas, dripRepository)

		expectedBlockOneIlkState = test_helpers.IlkState{
			Ilk:  fakeIlk,
			Take: blockOneIlkState[vat.IlkTake],
			Rate: blockOneIlkState[vat.IlkRate],
			Ink:  blockOneIlkState[vat.IlkInk],
			Art:  blockOneIlkState[vat.IlkArt],
			Spot: blockOneIlkState[pit.IlkSpot],
			Line: blockOneIlkState[pit.IlkLine],
			Chop: blockOneIlkState[cat.IlkChop],
			Lump: blockOneIlkState[cat.IlkLump],
			Flip: blockOneIlkState[cat.IlkFlip],
			Rho:  blockOneIlkState[jug.IlkRho],
			Tax:  blockOneIlkState[jug.IlkTax],
			Created: sql.NullString{
				String: blockOneHeader.Timestamp,
				Valid:  true,
			},
			Updated: sql.NullString{
				String: blockOneHeader.Timestamp,
				Valid:  true,
			},
		}

		blockTwoIlkState = test_helpers.GetIlkState(1)
		test_helpers.CreateVatRecords(blockTwoHeader, blockTwoIlkState, test_helpers.FakeIlkVatMetadatas, vatRepository)
		test_helpers.CreatePitRecords(blockTwoHeader, blockTwoIlkState, test_helpers.FakeIlkPitMetadatas, pitRepository)
		test_helpers.CreateCatRecords(blockTwoHeader, blockTwoIlkState, test_helpers.FakeIlkCatMetadatas, catRepository)
		test_helpers.CreateJugRecords(blockTwoHeader, blockTwoIlkState, test_helpers.FakeIlkJugMetadatas, dripRepository)
		expectedBlockTwoIlkState = test_helpers.IlkState{
			Ilk:  fakeIlk,
			Take: blockTwoIlkState[vat.IlkTake],
			Rate: blockTwoIlkState[vat.IlkRate],
			Ink:  blockTwoIlkState[vat.IlkInk],
			Art:  blockTwoIlkState[vat.IlkArt],
			Spot: blockTwoIlkState[pit.IlkSpot],
			Line: blockTwoIlkState[pit.IlkLine],
			Chop: blockTwoIlkState[cat.IlkChop],
			Lump: blockTwoIlkState[cat.IlkLump],
			Flip: blockTwoIlkState[cat.IlkFlip],
			Rho:  blockTwoIlkState[jug.IlkRho],
			Tax:  blockTwoIlkState[jug.IlkTax],
			Created: sql.NullString{
				String: blockOneHeader.Timestamp,
				Valid:  true,
			},
			Updated: sql.NullString{
				String: blockTwoHeader.Timestamp,
				Valid:  true,
			},
		}
	})

	It("returns the history of an ilk from the given block number", func() {
		var ilkId int
		err := db.Get(&ilkId, `SELECT id FROM maker.ilks WHERE ilk = $1`, fakeIlk)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []test_helpers.IlkState
		err = db.Select(&dbResult,
			`SELECT ilk, take, rate, ink, art, spot, line, chop, lump, flip, rho, tax, created, updated from maker.get_ilk_history_before_block($1, $2)`,
			blockTwo,
			ilkId)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(2))
		Expect(dbResult).To(ConsistOf([]test_helpers.IlkState{
			expectedBlockOneIlkState,
			expectedBlockTwoIlkState,
		}))
	})

	It("can handle multiple ilks in the db", func() {
		blockOneAnotherFakeIlkState := test_helpers.GetIlkState(3)

		test_helpers.CreateVatRecords(blockOneHeader, blockOneAnotherFakeIlkState, test_helpers.AnotherFakeIlkVatMetadatas, vatRepository)
		test_helpers.CreatePitRecords(blockOneHeader, blockOneAnotherFakeIlkState, test_helpers.AnotherFakeIlkPitMetadatas, pitRepository)
		test_helpers.CreateCatRecords(blockOneHeader, blockOneAnotherFakeIlkState, test_helpers.AnotherFakeIlkCatMetadatas, catRepository)
		test_helpers.CreateJugRecords(blockOneHeader, blockOneAnotherFakeIlkState, test_helpers.AnotherFakeIlkJugMetadatas, dripRepository)

		var ilkId int
		err := db.Get(&ilkId, `SELECT id FROM maker.ilks WHERE ilk = $1`, test_helpers.AnotherFakeIlk)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []test_helpers.IlkState
		err = db.Select(&dbResult,
			`SELECT ilk, take, rate, ink, art, spot, line, chop, lump, flip, rho, tax, created, updated from maker.get_ilk_history_before_block($1, $2)`,
			blockTwo,
			ilkId)
		Expect(err).NotTo(HaveOccurred())
		expectedBlockOneAnotherIlkState := test_helpers.IlkState{
			Ilk:  test_helpers.AnotherFakeIlk,
			Take: blockOneAnotherFakeIlkState[vat.IlkTake],
			Rate: blockOneAnotherFakeIlkState[vat.IlkRate],
			Ink:  blockOneAnotherFakeIlkState[vat.IlkInk],
			Art:  blockOneAnotherFakeIlkState[vat.IlkArt],
			Spot: blockOneAnotherFakeIlkState[pit.IlkSpot],
			Line: blockOneAnotherFakeIlkState[pit.IlkLine],
			Chop: blockOneAnotherFakeIlkState[cat.IlkChop],
			Lump: blockOneAnotherFakeIlkState[cat.IlkLump],
			Flip: blockOneAnotherFakeIlkState[cat.IlkFlip],
			Rho:  blockOneAnotherFakeIlkState[jug.IlkRho],
			Tax:  blockOneAnotherFakeIlkState[jug.IlkTax],
			Created: sql.NullString{
				String: blockOneHeader.Timestamp,
				Valid:  true,
			},
			Updated: sql.NullString{
				String: blockOneHeader.Timestamp,
				Valid:  true,
			},
		}

		//does not include fake ilk's results
		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult).To(ConsistOf([]test_helpers.IlkState{
			expectedBlockOneAnotherIlkState,
		}))
	})

	It("handles a query with a block number before the ilk is in the db", func() {
		blockZero := blockOne - 1

		blockZeroHeader := fakes.GetFakeHeader(int64(blockZero))
		_, err := headerRepository.CreateOrUpdateHeader(blockZeroHeader)
		Expect(err).NotTo(HaveOccurred())

		//the ilk is created in block 1 in the BeforeEach
		var ilkId int
		err = db.Get(&ilkId, `SELECT id FROM maker.ilks WHERE ilk = $1`, fakeIlk)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []test_helpers.IlkState
		err = db.Select(&dbResult,
			`SELECT ilk, take, rate, ink, art, spot, line, chop, lump, flip, rho, tax, created, updated from maker.get_ilk_history_before_block($1, $2)`,
			blockZero,
			ilkId)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult).To(BeEmpty())
	})

	It("handles when there have been no recent updates to the ilk", func() {
		blockOneHundred := int64(blockOne + 100)
		blockOneHundredHeader := fakes.GetFakeHeader(blockOneHundred)
		_, err := headerRepository.CreateOrUpdateHeader(blockOneHundredHeader)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.CreatePitRecords(blockOneHundredHeader, test_helpers.GetIlkState(100), test_helpers.AnotherFakeIlkPitMetadatas, pitRepository)

		var ilkId int
		err = db.Get(&ilkId, `SELECT id FROM maker.ilks WHERE ilk = $1`, fakeIlk)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []test_helpers.IlkState
		err = db.Select(&dbResult,
			`SELECT ilk, take, rate, ink, art, spot, line, chop, lump, flip, rho, tax, created, updated from maker.get_ilk_history_before_block($1, $2)`,
			blockOneHundred,
			ilkId)
		Expect(err).NotTo(HaveOccurred())

		Expect(dbResult).To(ConsistOf(expectedBlockOneIlkState, expectedBlockTwoIlkState))
	})
})
