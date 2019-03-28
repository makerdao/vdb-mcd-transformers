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
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
)

var (
	vatRepository                                    vat.VatStorageRepository
	pitRepository                                    pit.PitStorageRepository
	catRepository                                    cat.CatStorageRepository
	jugRepository                                    jug.JugStorageRepository
	headerRepository                                 repositories.HeaderRepository
	blockOneHeader, blockTwoHeader, blockThreeHeader core.Header
)

var _ = Describe("Ilk State Query", func() {
	var (
		db             *postgres.DB
		blockOne       = rand.Int()
		blockTwo       = blockOne + 1
		blockThree     = blockOne + 2
		fakeIlk        = test_helpers.FakeIlk
		anotherFakeIlk = test_helpers.AnotherFakeIlk
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		vatRepository.SetDB(db)
		pitRepository.SetDB(db)
		catRepository.SetDB(db)
		jugRepository.SetDB(db)
		headerRepository = repositories.NewHeaderRepository(db)

		blockOneHeader = fakes.GetFakeHeader(int64(blockOne))
		_, err := headerRepository.CreateOrUpdateHeader(blockOneHeader)
		Expect(err).NotTo(HaveOccurred())

		blockTwoHeader = fakes.GetFakeHeader(int64(blockTwo))
		blockTwoHeader.Timestamp = blockTwoHeader.Timestamp + "2"
		blockTwoHeader.Hash = "block2Hash"
		_, err = headerRepository.CreateOrUpdateHeader(blockTwoHeader)
		Expect(err).NotTo(HaveOccurred())

		blockThreeHeader = fakes.GetFakeHeader(int64(blockThree))
		blockThreeHeader.Timestamp = blockThreeHeader.Timestamp + "3"
		blockThreeHeader.Hash = "block3Hash"
		_, err = headerRepository.CreateOrUpdateHeader(blockThreeHeader)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		test_config.CleanTestDB(db)
	})

	It("gets an ilk", func() {
		ilkState := test_helpers.GetIlkState(0)
		createIlkAtBlock(blockOneHeader, ilkState, test_helpers.FakeIlkVatMetadatas, test_helpers.FakeIlkPitMetadatas, test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas)

		var ilkId int
		err := db.Get(&ilkId, `SELECT id FROM maker.ilks WHERE ilk = $1`, fakeIlk)
		Expect(err).NotTo(HaveOccurred())

		var dbResult test_helpers.IlkState
		err = db.Get(&dbResult,
			`SELECT ilk, take, rate, ink, art, spot, line, chop, lump, flip, rho, tax, created, updated from maker.get_ilk_at_block_number($1, $2)`,
			blockOne,
			ilkId)
		Expect(err).NotTo(HaveOccurred())

		expectedIlk := test_helpers.IlkState{
			Ilk:     fakeIlk,
			Take:    ilkState[vat.IlkTake],
			Rate:    ilkState[vat.IlkRate],
			Ink:     ilkState[vat.IlkInk],
			Art:     ilkState[vat.IlkArt],
			Spot:    ilkState[pit.IlkSpot],
			Line:    ilkState[pit.IlkLine],
			Chop:    ilkState[cat.IlkChop],
			Lump:    ilkState[cat.IlkLump],
			Flip:    ilkState[cat.IlkFlip],
			Rho:     ilkState[jug.IlkRho],
			Tax:     ilkState[jug.IlkTax],
			Created: sql.NullString{String: blockOneHeader.Timestamp, Valid: true},
			Updated: sql.NullString{String: blockOneHeader.Timestamp, Valid: true},
		}
		Expect(dbResult).To(Equal(expectedIlk))
	})

	It("returns the correct data if there are several ilks", func() {
		fakeIlkState := test_helpers.GetIlkState(1)
		anotherFakeIlkState := test_helpers.GetIlkState(2)
		createIlkAtBlock(blockOneHeader, fakeIlkState, test_helpers.FakeIlkVatMetadatas, test_helpers.FakeIlkPitMetadatas, test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas)
		createIlkAtBlock(blockOneHeader, anotherFakeIlkState, test_helpers.AnotherFakeIlkVatMetadatas, test_helpers.AnotherFakeIlkPitMetadatas, test_helpers.AnotherFakeIlkCatMetadatas, test_helpers.AnotherFakeIlkJugMetadatas)

		var fakeIlkId int
		err := db.Get(&fakeIlkId, `SELECT id FROM maker.ilks WHERE ilk = $1`, fakeIlk)
		Expect(err).NotTo(HaveOccurred())

		var anotherFakeIlkId int
		err = db.Get(&anotherFakeIlkId, `SELECT id FROM maker.ilks WHERE ilk = $1`, anotherFakeIlk)
		Expect(err).NotTo(HaveOccurred())

		var fakeIlkResult test_helpers.IlkState
		err = db.Get(&fakeIlkResult,
			`SELECT ilk, take, rate, ink, art, spot, line, chop, lump, flip, rho, tax, created, updated from maker.get_ilk_at_block_number($1, $2)`,
			blockOne,
			fakeIlkId)
		Expect(err).NotTo(HaveOccurred())

		var anotherFakeIlkResult test_helpers.IlkState
		err = db.Get(&anotherFakeIlkResult,
			`SELECT ilk, take, rate, ink, art, spot, line, chop, lump, flip, rho, tax, created, updated from maker.get_ilk_at_block_number($1, $2)`,
			blockOne,
			anotherFakeIlkId)
		Expect(err).NotTo(HaveOccurred())

		expectedFakeIlk := test_helpers.IlkState{
			Ilk:     fakeIlk,
			Take:    fakeIlkState[vat.IlkTake],
			Rate:    fakeIlkState[vat.IlkRate],
			Ink:     fakeIlkState[vat.IlkInk],
			Art:     fakeIlkState[vat.IlkArt],
			Spot:    fakeIlkState[pit.IlkSpot],
			Line:    fakeIlkState[pit.IlkLine],
			Chop:    fakeIlkState[cat.IlkChop],
			Lump:    fakeIlkState[cat.IlkLump],
			Flip:    fakeIlkState[cat.IlkFlip],
			Rho:     fakeIlkState[jug.IlkRho],
			Tax:     fakeIlkState[jug.IlkTax],
			Created: sql.NullString{String: blockOneHeader.Timestamp, Valid: true},
			Updated: sql.NullString{String: blockOneHeader.Timestamp, Valid: true},
		}
		expectedAnotherFakeIlk := test_helpers.IlkState{
			Ilk:     anotherFakeIlk,
			Take:    anotherFakeIlkState[vat.IlkTake],
			Rate:    anotherFakeIlkState[vat.IlkRate],
			Ink:     anotherFakeIlkState[vat.IlkInk],
			Art:     anotherFakeIlkState[vat.IlkArt],
			Spot:    anotherFakeIlkState[pit.IlkSpot],
			Line:    anotherFakeIlkState[pit.IlkLine],
			Chop:    anotherFakeIlkState[cat.IlkChop],
			Lump:    anotherFakeIlkState[cat.IlkLump],
			Flip:    anotherFakeIlkState[cat.IlkFlip],
			Rho:     anotherFakeIlkState[jug.IlkRho],
			Tax:     anotherFakeIlkState[jug.IlkTax],
			Created: sql.NullString{String: blockOneHeader.Timestamp, Valid: true},
			Updated: sql.NullString{String: blockOneHeader.Timestamp, Valid: true},
		}
		Expect(fakeIlkResult).To(Equal(expectedFakeIlk))
		Expect(anotherFakeIlkResult).To(Equal(expectedAnotherFakeIlk))
	})

	Describe("handles getting the most recent Ilk values as of a given block", func() {
		It("gets the Ilk for block one", func() {
			fakeIlkState := test_helpers.GetIlkState(0)
			createIlkAtBlock(blockOneHeader, fakeIlkState, test_helpers.FakeIlkVatMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas)

			var fakeIlkId int
			err := db.Get(&fakeIlkId, `SELECT id FROM maker.ilks WHERE ilk = $1`, fakeIlk)
			Expect(err).NotTo(HaveOccurred())

			var blockOneDbResult test_helpers.IlkState
			err = db.Get(&blockOneDbResult, `SELECT take, rate, ink, created from maker.get_ilk_at_block_number($1, $2)`, blockOne, fakeIlkId)
			Expect(err).NotTo(HaveOccurred())

			expectedIlk := test_helpers.IlkState{
				Take:    fakeIlkState[vat.IlkTake],
				Rate:    fakeIlkState[vat.IlkRate],
				Ink:     fakeIlkState[vat.IlkInk],
				Created: sql.NullString{String: blockOneHeader.Timestamp, Valid: true},
			}
			Expect(blockOneDbResult).To(Equal(expectedIlk))
		})

		It("gets the Ilk for block two", func() {
			blockOneFakeIlkState := test_helpers.GetIlkState(1)
			createIlkAtBlock(blockOneHeader, blockOneFakeIlkState, test_helpers.FakeIlkVatMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas)
			// block two doesn't update rate or ink
			blockTwoFakeIlkState := test_helpers.GetIlkState(2)
			createIlkAtBlock(blockTwoHeader, blockTwoFakeIlkState, []utils.StorageValueMetadata{test_helpers.FakeIlkTakeMetadata}, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas)

			var fakeIlkId int
			err := db.Get(&fakeIlkId, `SELECT id FROM maker.ilks WHERE ilk = $1`, fakeIlk)
			Expect(err).NotTo(HaveOccurred())

			var blockTwoDbResult test_helpers.IlkState
			err = db.Get(&blockTwoDbResult, `SELECT take, rate, ink from maker.get_ilk_at_block_number($1, $2)`, blockTwo, fakeIlkId)
			Expect(err).NotTo(HaveOccurred())

			blockTwoExpectedIlk := test_helpers.IlkState{
				Take: blockTwoFakeIlkState[vat.IlkTake],
				Rate: blockOneFakeIlkState[vat.IlkRate], // value hasn't changed since block 1
				Ink:  blockOneFakeIlkState[vat.IlkInk],  // value hasn't changed since block 1
			}
			Expect(blockTwoDbResult).To(Equal(blockTwoExpectedIlk))
		})

		It("gets the Ilk for block three", func() {
			//no updates to ink
			blockOneFakeIlkState := test_helpers.GetIlkState(1)
			createIlkAtBlock(blockOneHeader, blockOneFakeIlkState, test_helpers.FakeIlkVatMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas)
			// block two doesn't update rate or ink
			blockTwoFakeIlkState := test_helpers.GetIlkState(1)
			createIlkAtBlock(blockTwoHeader, blockTwoFakeIlkState, []utils.StorageValueMetadata{test_helpers.FakeIlkTakeMetadata}, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas)
			// block three doesn't update ink
			blockThreeFakeIlkState := test_helpers.GetIlkState(3)
			createIlkAtBlock(blockThreeHeader, blockThreeFakeIlkState, []utils.StorageValueMetadata{test_helpers.FakeIlkTakeMetadata, test_helpers.FakeIlkRateMetadata}, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas)

			var fakeIlkId int
			err := db.Get(&fakeIlkId, `SELECT id FROM maker.ilks WHERE ilk = $1`, fakeIlk)
			Expect(err).NotTo(HaveOccurred())

			var blockThreeDbResult test_helpers.IlkState
			err = db.Get(&blockThreeDbResult, `SELECT take, rate, ink from maker.get_ilk_at_block_number($1, $2)`, blockThree, fakeIlkId)
			Expect(err).NotTo(HaveOccurred())

			blockThreeExpectedIlk := test_helpers.IlkState{
				Take: blockThreeFakeIlkState[vat.IlkTake],
				Rate: blockThreeFakeIlkState[vat.IlkRate],
				Ink:  blockOneFakeIlkState[vat.IlkInk], // value hasn't changed since block 1
			}
			Expect(blockThreeDbResult).To(Equal(blockThreeExpectedIlk))
		})

		It("gets more than one ilk as of block three", func() {
			blockOneFakeIlkState := test_helpers.GetIlkState(1)
			createIlkAtBlock(blockOneHeader, blockOneFakeIlkState, test_helpers.FakeIlkVatMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas)
			blockOneAnotherFakeIlkState := test_helpers.GetIlkState(2)
			createIlkAtBlock(blockOneHeader, blockOneAnotherFakeIlkState, test_helpers.AnotherFakeIlkVatMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas)
			// block two doesn't update rate or ink for fakeIlk
			// and doesn't update take, rate or ink for anotherFakeIlk
			blockTwoFakeIlkState := test_helpers.GetIlkState(1)
			createIlkAtBlock(blockTwoHeader, blockTwoFakeIlkState, []utils.StorageValueMetadata{test_helpers.FakeIlkTakeMetadata}, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas)
			// block three doesn't update ink
			// and doesn't update take, rate or ink for anotherFakeIlk
			blockThreeFakeIlkState := test_helpers.GetIlkState(3)
			createIlkAtBlock(blockThreeHeader, blockThreeFakeIlkState, []utils.StorageValueMetadata{test_helpers.FakeIlkTakeMetadata, test_helpers.FakeIlkRateMetadata}, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas)

			var fakeIlkId int
			err := db.Get(&fakeIlkId, `SELECT id FROM maker.ilks WHERE ilk = $1`, fakeIlk)
			Expect(err).NotTo(HaveOccurred())

			var fakeIlkResult test_helpers.IlkState
			err = db.Get(&fakeIlkResult, `SELECT ilk, take, rate, ink from maker.get_ilk_at_block_number($1, $2)`, blockThree, fakeIlkId)
			Expect(err).NotTo(HaveOccurred())

			var anotherFakeIlkId int
			err = db.Get(&anotherFakeIlkId, `SELECT id FROM maker.ilks WHERE ilk = $1`, anotherFakeIlk)
			Expect(err).NotTo(HaveOccurred())

			var anotherFakeIlkResult test_helpers.IlkState
			err = db.Get(&anotherFakeIlkResult, `SELECT ilk, take, rate, ink from maker.get_ilk_at_block_number($1, $2)`, blockThree, anotherFakeIlkId)
			Expect(err).NotTo(HaveOccurred())

			blockThreeExpectedFakeIlk := test_helpers.IlkState{
				Ilk:  fakeIlk,
				Take: blockThreeFakeIlkState[vat.IlkTake],
				Rate: blockThreeFakeIlkState[vat.IlkRate],
				Ink:  blockOneFakeIlkState[vat.IlkInk], // value hasn't changed since block 1
			}
			blockThreeExpectedAnotherFakeIlk := test_helpers.IlkState{
				Ilk:  anotherFakeIlk,
				Take: blockOneAnotherFakeIlkState[vat.IlkTake], // value hasn't changed since block 1
				Rate: blockOneAnotherFakeIlkState[vat.IlkRate], // value hasn't changed since block 1
				Ink:  blockOneAnotherFakeIlkState[vat.IlkInk],  // value hasn't changed since block 1
			}
			Expect(fakeIlkResult).To(Equal(blockThreeExpectedFakeIlk))
			Expect(anotherFakeIlkResult).To(Equal(blockThreeExpectedAnotherFakeIlk))
		})
	})

	It("returns the rest of the Ilk data if a field has missing data", func() {
		vatMetadatasWithoutTake := test_helpers.FakeIlkVatMetadatas[1:]
		blockOneFakeIlkState := test_helpers.GetIlkState(1)
		createIlkAtBlock(blockOneHeader, blockOneFakeIlkState, vatMetadatasWithoutTake, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas)

		var fakeIlkId int
		err := db.Get(&fakeIlkId, `SELECT id FROM maker.ilks WHERE ilk = $1`, fakeIlk)
		Expect(err).NotTo(HaveOccurred())

		var blockOneDbResult test_helpers.IlkState
		err = db.Get(&blockOneDbResult, `SELECT rate, ink from maker.get_ilk_at_block_number($1, $2)`, blockOne, fakeIlkId)
		Expect(err).NotTo(HaveOccurred())

		expectedIlk := test_helpers.IlkState{
			Rate: blockOneFakeIlkState[vat.IlkRate],
			Ink:  blockOneFakeIlkState[vat.IlkInk],
		}
		Expect(blockOneDbResult).To(Equal(expectedIlk))
	})

	It("gets the correct timestamps", func() {
		// fakeIlk created at block1
		// fakeIlk updated at block2
		// anotherFakeIlk created at block2
		blockOneFakeIlkState := test_helpers.GetIlkState(0)
		createIlkAtBlock(blockOneHeader, blockOneFakeIlkState, test_helpers.FakeIlkVatMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas)

		blockTwoFakeIlkState := test_helpers.GetIlkState(1)
		createIlkAtBlock(blockTwoHeader, blockTwoFakeIlkState, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas, test_helpers.FakeIlkJugMetadatas)

		blockTwoAnotherFakeIlkState := test_helpers.GetIlkState(2)
		createIlkAtBlock(blockTwoHeader, blockTwoAnotherFakeIlkState, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas, test_helpers.AnotherFakeIlkJugMetadatas)

		var fakeIlkId int
		err := db.Get(&fakeIlkId, `SELECT id FROM maker.ilks WHERE ilk = $1`, fakeIlk)
		Expect(err).NotTo(HaveOccurred())

		var fakeIlkBlockOneDbResult test_helpers.IlkState
		err = db.Get(&fakeIlkBlockOneDbResult, `SELECT ilk, created, updated from maker.get_ilk_at_block_number($1, $2)`, blockOne, fakeIlkId)
		Expect(err).NotTo(HaveOccurred())
		expectedBlockOneFakeIlkState := test_helpers.IlkState{
			Ilk:     fakeIlk,
			Created: sql.NullString{String: blockOneHeader.Timestamp, Valid: true},
			Updated: sql.NullString{String: blockOneHeader.Timestamp, Valid: true},
		}
		Expect(fakeIlkBlockOneDbResult).To(Equal(expectedBlockOneFakeIlkState))

		var fakeIlkBlockTwoDbResult test_helpers.IlkState
		err = db.Get(&fakeIlkBlockTwoDbResult, `SELECT ilk, created, updated from maker.get_ilk_at_block_number($1, $2)`, blockTwo, fakeIlkId)
		Expect(err).NotTo(HaveOccurred())
		expectedBlockTwoFakeIlkState := test_helpers.IlkState{
			Ilk:     fakeIlk,
			Created: sql.NullString{String: blockOneHeader.Timestamp, Valid: true},
			Updated: sql.NullString{String: blockTwoHeader.Timestamp, Valid: true},
		}
		Expect(fakeIlkBlockTwoDbResult).To(Equal(expectedBlockTwoFakeIlkState))

		var anotherFakeIlkId int
		err = db.Get(&anotherFakeIlkId, `SELECT id FROM maker.ilks WHERE ilk = $1`, anotherFakeIlk)
		Expect(err).NotTo(HaveOccurred())

		var anotherFakeIlkDbResult test_helpers.IlkState
		err = db.Get(&anotherFakeIlkDbResult, `SELECT ilk, created, updated from maker.get_ilk_at_block_number($1, $2)`, blockTwo, anotherFakeIlkId)
		Expect(err).NotTo(HaveOccurred())
		expectedBlockTwoAnotherFakeIlkState := test_helpers.IlkState{
			Ilk:     anotherFakeIlk,
			Created: sql.NullString{String: blockTwoHeader.Timestamp, Valid: true},
			Updated: sql.NullString{String: blockTwoHeader.Timestamp, Valid: true},
		}

		Expect(anotherFakeIlkDbResult).To(Equal(expectedBlockTwoAnotherFakeIlkState))
	})
})

func createIlkAtBlock(header core.Header, valuesMap map[string]string, vatMetadatas, pitMetadatas, catMetadatas, jugMetadatas []utils.StorageValueMetadata) {
	test_helpers.CreateVatRecords(header, valuesMap, vatMetadatas, vatRepository)
	test_helpers.CreatePitRecords(header, valuesMap, pitMetadatas, pitRepository)
	test_helpers.CreateCatRecords(header, valuesMap, catMetadatas, catRepository)
	test_helpers.CreateJugRecords(header, valuesMap, jugMetadatas, jugRepository)
}
