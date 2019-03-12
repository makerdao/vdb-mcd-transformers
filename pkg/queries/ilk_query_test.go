package queries_test

import (
	"database/sql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage/cat"
	"github.com/vulcanize/mcd_transformers/transformers/storage/drip"
	"github.com/vulcanize/mcd_transformers/transformers/storage/pit"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vat"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
	"strconv"
)

var (
	vatRepository                                    vat.VatStorageRepository
	pitRepository                                    pit.PitStorageRepository
	catRepository                                    cat.CatStorageRepository
	dripRepository                                   drip.DripStorageRepository
	headerRepository                                 repositories.HeaderRepository
	blockOneHeader, blockTwoHeader, blockThreeHeader core.Header
)

var _ = Describe("Ilk State Query", func() {
	var (
		db             *postgres.DB
		blockOne       = rand.Int()
		blockTwo       = blockOne + 1
		blockThree     = blockOne + 2
		fakeIlk        = "fakeIlk"
		anotherFakeIlk = "anotherFakeIlk"

		emptyMetadatas      []utils.StorageValueMetadata
		fakeIlkTakeMetadata = getMetadata(vat.IlkTake, fakeIlk, utils.Uint256)
		fakeIlkRateMetadata = getMetadata(vat.IlkRate, fakeIlk, utils.Uint256)
		fakeIlkInkMetadata  = getMetadata(vat.IlkInk, fakeIlk, utils.Uint256)
		fakeIlkArtMetadata  = getMetadata(vat.IlkArt, fakeIlk, utils.Uint256)
		fakeIlkSpotMetadata = getMetadata(pit.IlkSpot, fakeIlk, utils.Uint256)
		fakeIlkLineMetadata = getMetadata(pit.IlkLine, fakeIlk, utils.Uint256)
		fakeIlkChopMetadata = getMetadata(cat.IlkChop, fakeIlk, utils.Uint256)
		fakeIlkLumpMetadata = getMetadata(cat.IlkLump, fakeIlk, utils.Uint256)
		fakeIlkFlipMetadata = getMetadata(cat.IlkFlip, fakeIlk, utils.Address)
		fakeIlkRhoMetadata  = getMetadata(drip.IlkRho, fakeIlk, utils.Uint256)
		fakeIlkTaxMetadata  = getMetadata(drip.IlkTax, fakeIlk, utils.Uint256)

		fakeIlkVatMetadatas = []utils.StorageValueMetadata{
			fakeIlkTakeMetadata,
			fakeIlkRateMetadata,
			fakeIlkInkMetadata,
			fakeIlkArtMetadata,
		}
		fakeIlkPitMetadatas = []utils.StorageValueMetadata{
			fakeIlkSpotMetadata,
			fakeIlkLineMetadata,
		}
		fakeIlkCatMetadatas = []utils.StorageValueMetadata{
			fakeIlkChopMetadata,
			fakeIlkLumpMetadata,
			fakeIlkFlipMetadata,
		}
		fakeIlkDripMetadatas = []utils.StorageValueMetadata{
			fakeIlkRhoMetadata,
			fakeIlkTaxMetadata,
		}

		anotherFakeIlkTakeMetadata = getMetadata(vat.IlkTake, anotherFakeIlk, utils.Uint256)
		anotherFakeIlkRateMetadata = getMetadata(vat.IlkRate, anotherFakeIlk, utils.Uint256)
		anotherFakeIlkInkMetadata  = getMetadata(vat.IlkInk, anotherFakeIlk, utils.Uint256)
		anotherFakeIlkArtMetadata  = getMetadata(vat.IlkArt, anotherFakeIlk, utils.Uint256)
		anotherFakeIlkSpotMetadata = getMetadata(pit.IlkSpot, anotherFakeIlk, utils.Uint256)
		anotherFakeIlkLineMetadata = getMetadata(pit.IlkLine, anotherFakeIlk, utils.Uint256)
		anotherFakeIlkChopMetadata = getMetadata(cat.IlkChop, anotherFakeIlk, utils.Uint256)
		anotherFakeIlkLumpMetadata = getMetadata(cat.IlkLump, anotherFakeIlk, utils.Uint256)
		anotherFakeIlkFlipMetadata = getMetadata(cat.IlkFlip, anotherFakeIlk, utils.Address)
		anotherFakeIlkRhoMetadata  = getMetadata(drip.IlkRho, anotherFakeIlk, utils.Uint256)
		anotherFakeIlkTaxMetadata  = getMetadata(drip.IlkTax, anotherFakeIlk, utils.Uint256)

		anotherFakeIlkVatMetadatas = []utils.StorageValueMetadata{
			anotherFakeIlkTakeMetadata,
			anotherFakeIlkRateMetadata,
			anotherFakeIlkInkMetadata,
			anotherFakeIlkArtMetadata,
		}
		anotherFakeIlkPitMetadatas = []utils.StorageValueMetadata{
			anotherFakeIlkSpotMetadata,
			anotherFakeIlkLineMetadata,
		}
		anotherFakeIlkCatMetadatas = []utils.StorageValueMetadata{
			anotherFakeIlkChopMetadata,
			anotherFakeIlkLumpMetadata,
			anotherFakeIlkFlipMetadata,
		}
		anotherFakeIlkDripMetadatas = []utils.StorageValueMetadata{
			anotherFakeIlkRhoMetadata,
			anotherFakeIlkTaxMetadata,
		}
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
		ilkState := getIlkState(0)
		createIlkAtBlock(blockOneHeader, ilkState, fakeIlkVatMetadatas, fakeIlkPitMetadatas, fakeIlkCatMetadatas, fakeIlkDripMetadatas)

		var ilkId int
		err := db.Get(&ilkId, `SELECT id FROM maker.ilks WHERE ilk = $1`, fakeIlk)
		Expect(err).NotTo(HaveOccurred())

		var dbResult IlkState
		err = db.Get(&dbResult,
			`SELECT ilk, take, rate, ink, art, spot, line, chop, lump, flip, rho, tax, created, updated from maker.get_ilk_at_block_number($1, $2)`,
			blockOne,
			ilkId)
		Expect(err).NotTo(HaveOccurred())

		expectedIlk := IlkState{
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
			Rho:     ilkState[drip.IlkRho],
			Tax:     ilkState[drip.IlkTax],
			Created: sql.NullString{String: blockOneHeader.Timestamp, Valid: true},
			Updated: sql.NullString{String: blockOneHeader.Timestamp, Valid: true},
		}
		Expect(dbResult).To(Equal(expectedIlk))
	})

	It("returns the correct data if there are several ilks", func() {
		fakeIlkState := getIlkState(1)
		anotherFakeIlkState := getIlkState(2)
		createIlkAtBlock(blockOneHeader, fakeIlkState, fakeIlkVatMetadatas, fakeIlkPitMetadatas, fakeIlkCatMetadatas, fakeIlkDripMetadatas)
		createIlkAtBlock(blockOneHeader, anotherFakeIlkState, anotherFakeIlkVatMetadatas, anotherFakeIlkPitMetadatas, anotherFakeIlkCatMetadatas, anotherFakeIlkDripMetadatas)

		var fakeIlkId int
		err := db.Get(&fakeIlkId, `SELECT id FROM maker.ilks WHERE ilk = $1`, fakeIlk)
		Expect(err).NotTo(HaveOccurred())

		var anotherFakeIlkId int
		err = db.Get(&anotherFakeIlkId, `SELECT id FROM maker.ilks WHERE ilk = $1`, anotherFakeIlk)
		Expect(err).NotTo(HaveOccurred())

		var fakeIlkResult IlkState
		err = db.Get(&fakeIlkResult,
			`SELECT ilk, take, rate, ink, art, spot, line, chop, lump, flip, rho, tax, created, updated from maker.get_ilk_at_block_number($1, $2)`,
			blockOne,
			fakeIlkId)
		Expect(err).NotTo(HaveOccurred())

		var anotherFakeIlkResult IlkState
		err = db.Get(&anotherFakeIlkResult,
			`SELECT ilk, take, rate, ink, art, spot, line, chop, lump, flip, rho, tax, created, updated from maker.get_ilk_at_block_number($1, $2)`,
			blockOne,
			anotherFakeIlkId)
		Expect(err).NotTo(HaveOccurred())

		expectedFakeIlk := IlkState{
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
			Rho:     fakeIlkState[drip.IlkRho],
			Tax:     fakeIlkState[drip.IlkTax],
			Created: sql.NullString{String: blockOneHeader.Timestamp, Valid: true},
			Updated: sql.NullString{String: blockOneHeader.Timestamp, Valid: true},
		}
		expectedAnotherFakeIlk := IlkState{
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
			Rho:     anotherFakeIlkState[drip.IlkRho],
			Tax:     anotherFakeIlkState[drip.IlkTax],
			Created: sql.NullString{String: blockOneHeader.Timestamp, Valid: true},
			Updated: sql.NullString{String: blockOneHeader.Timestamp, Valid: true},
		}
		Expect(fakeIlkResult).To(Equal(expectedFakeIlk))
		Expect(anotherFakeIlkResult).To(Equal(expectedAnotherFakeIlk))
	})

	Describe("handles getting the most recent Ilk values as of a given block", func() {
		It("gets the Ilk for block one", func() {
			fakeIlkState := getIlkState(0)
			createIlkAtBlock(blockOneHeader, fakeIlkState, fakeIlkVatMetadatas, emptyMetadatas, emptyMetadatas, emptyMetadatas)

			var fakeIlkId int
			err := db.Get(&fakeIlkId, `SELECT id FROM maker.ilks WHERE ilk = $1`, fakeIlk)
			Expect(err).NotTo(HaveOccurred())

			var blockOneDbResult IlkState
			err = db.Get(&blockOneDbResult, `SELECT take, rate, ink, created from maker.get_ilk_at_block_number($1, $2)`, blockOne, fakeIlkId)
			Expect(err).NotTo(HaveOccurred())

			expectedIlk := IlkState{
				Take:    fakeIlkState[vat.IlkTake],
				Rate:    fakeIlkState[vat.IlkRate],
				Ink:     fakeIlkState[vat.IlkInk],
				Created: sql.NullString{String: blockOneHeader.Timestamp, Valid: true},
			}
			Expect(blockOneDbResult).To(Equal(expectedIlk))
		})

		It("gets the Ilk for block two", func() {
			blockOneFakeIlkState := getIlkState(1)
			createIlkAtBlock(blockOneHeader, blockOneFakeIlkState, fakeIlkVatMetadatas, emptyMetadatas, emptyMetadatas, emptyMetadatas)
			// block two doesn't update rate or ink
			blockTwoFakeIlkState := getIlkState(2)
			createIlkAtBlock(blockTwoHeader, blockTwoFakeIlkState, []utils.StorageValueMetadata{fakeIlkTakeMetadata}, emptyMetadatas, emptyMetadatas, emptyMetadatas)

			var fakeIlkId int
			err := db.Get(&fakeIlkId, `SELECT id FROM maker.ilks WHERE ilk = $1`, fakeIlk)
			Expect(err).NotTo(HaveOccurred())

			var blockTwoDbResult IlkState
			err = db.Get(&blockTwoDbResult, `SELECT take, rate, ink from maker.get_ilk_at_block_number($1, $2)`, blockTwo, fakeIlkId)
			Expect(err).NotTo(HaveOccurred())

			blockTwoExpectedIlk := IlkState{
				Take: blockTwoFakeIlkState[vat.IlkTake],
				Rate: blockOneFakeIlkState[vat.IlkRate], // value hasn't changed since block 1
				Ink:  blockOneFakeIlkState[vat.IlkInk],  // value hasn't changed since block 1
			}
			Expect(blockTwoDbResult).To(Equal(blockTwoExpectedIlk))
		})

		It("gets the Ilk for block three", func() {
			//no updates to ink
			blockOneFakeIlkState := getIlkState(1)
			createIlkAtBlock(blockOneHeader, blockOneFakeIlkState, fakeIlkVatMetadatas, emptyMetadatas, emptyMetadatas, emptyMetadatas)
			// block two doesn't update rate or ink
			blockTwoFakeIlkState := getIlkState(1)
			createIlkAtBlock(blockTwoHeader, blockTwoFakeIlkState, []utils.StorageValueMetadata{fakeIlkTakeMetadata}, emptyMetadatas, emptyMetadatas, emptyMetadatas)
			// block three doesn't update ink
			blockThreeFakeIlkState := getIlkState(3)
			createIlkAtBlock(blockThreeHeader, blockThreeFakeIlkState, []utils.StorageValueMetadata{fakeIlkTakeMetadata, fakeIlkRateMetadata}, emptyMetadatas, emptyMetadatas, emptyMetadatas)

			var fakeIlkId int
			err := db.Get(&fakeIlkId, `SELECT id FROM maker.ilks WHERE ilk = $1`, fakeIlk)
			Expect(err).NotTo(HaveOccurred())

			var blockThreeDbResult IlkState
			err = db.Get(&blockThreeDbResult, `SELECT take, rate, ink from maker.get_ilk_at_block_number($1, $2)`, blockThree, fakeIlkId)
			Expect(err).NotTo(HaveOccurred())

			blockThreeExpectedIlk := IlkState{
				Take: blockThreeFakeIlkState[vat.IlkTake],
				Rate: blockThreeFakeIlkState[vat.IlkRate],
				Ink:  blockOneFakeIlkState[vat.IlkInk], // value hasn't changed since block 1
			}
			Expect(blockThreeDbResult).To(Equal(blockThreeExpectedIlk))
		})

		It("gets more than one ilk as of block three", func() {
			blockOneFakeIlkState := getIlkState(1)
			createIlkAtBlock(blockOneHeader, blockOneFakeIlkState, fakeIlkVatMetadatas, emptyMetadatas, emptyMetadatas, emptyMetadatas)
			blockOneAnotherFakeIlkState := getIlkState(2)
			createIlkAtBlock(blockOneHeader, blockOneAnotherFakeIlkState, anotherFakeIlkVatMetadatas, emptyMetadatas, emptyMetadatas, emptyMetadatas)
			// block two doesn't update rate or ink for fakeIlk
			// and doesn't update take, rate or ink for anotherFakeIlk
			blockTwoFakeIlkState := getIlkState(1)
			createIlkAtBlock(blockTwoHeader, blockTwoFakeIlkState, []utils.StorageValueMetadata{fakeIlkTakeMetadata}, emptyMetadatas, emptyMetadatas, emptyMetadatas)
			// block three doesn't update ink
			// and doesn't update take, rate or ink for anotherFakeIlk
			blockThreeFakeIlkState := getIlkState(3)
			createIlkAtBlock(blockThreeHeader, blockThreeFakeIlkState, []utils.StorageValueMetadata{fakeIlkTakeMetadata, fakeIlkRateMetadata}, emptyMetadatas, emptyMetadatas, emptyMetadatas)

			var fakeIlkId int
			err := db.Get(&fakeIlkId, `SELECT id FROM maker.ilks WHERE ilk = $1`, fakeIlk)
			Expect(err).NotTo(HaveOccurred())

			var fakeIlkResult IlkState
			err = db.Get(&fakeIlkResult, `SELECT ilk, take, rate, ink from maker.get_ilk_at_block_number($1, $2)`, blockThree, fakeIlkId)
			Expect(err).NotTo(HaveOccurred())

			var anotherFakeIlkId int
			err = db.Get(&anotherFakeIlkId, `SELECT id FROM maker.ilks WHERE ilk = $1`, anotherFakeIlk)
			Expect(err).NotTo(HaveOccurred())

			var anotherFakeIlkResult IlkState
			err = db.Get(&anotherFakeIlkResult, `SELECT ilk, take, rate, ink from maker.get_ilk_at_block_number($1, $2)`, blockThree, anotherFakeIlkId)
			Expect(err).NotTo(HaveOccurred())

			blockThreeExpectedFakeIlk := IlkState{
				Ilk:  fakeIlk,
				Take: blockThreeFakeIlkState[vat.IlkTake],
				Rate: blockThreeFakeIlkState[vat.IlkRate],
				Ink:  blockOneFakeIlkState[vat.IlkInk], // value hasn't changed since block 1
			}
			blockThreeExpectedAnotherFakeIlk := IlkState{
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
		vatMetadatasWithoutTake := fakeIlkVatMetadatas[1:]
		blockOneFakeIlkState := getIlkState(1)
		createIlkAtBlock(blockOneHeader, blockOneFakeIlkState, vatMetadatasWithoutTake, emptyMetadatas, emptyMetadatas, emptyMetadatas)

		var fakeIlkId int
		err := db.Get(&fakeIlkId, `SELECT id FROM maker.ilks WHERE ilk = $1`, fakeIlk)
		Expect(err).NotTo(HaveOccurred())

		var blockOneDbResult IlkState
		err = db.Get(&blockOneDbResult, `SELECT rate, ink from maker.get_ilk_at_block_number($1, $2)`, blockOne, fakeIlkId)
		Expect(err).NotTo(HaveOccurred())

		expectedIlk := IlkState{
			Rate: blockOneFakeIlkState[vat.IlkRate],
			Ink:  blockOneFakeIlkState[vat.IlkInk],
		}
		Expect(blockOneDbResult).To(Equal(expectedIlk))
	})

	It("gets the correct timestamps", func() {
		// fakeIlk created at block1
		// fakeIlk updated at block2
		// anotherFakeIlk created at block2
		blockOneFakeIlkState := getIlkState(0)
		createIlkAtBlock(blockOneHeader, blockOneFakeIlkState, fakeIlkVatMetadatas, emptyMetadatas, emptyMetadatas, emptyMetadatas)

		blockTwoFakeIlkState := getIlkState(1)
		createIlkAtBlock(blockTwoHeader, blockTwoFakeIlkState, emptyMetadatas, emptyMetadatas, emptyMetadatas, fakeIlkDripMetadatas)

		blockTwoAnotherFakeIlkState := getIlkState(2)
		createIlkAtBlock(blockTwoHeader, blockTwoAnotherFakeIlkState, emptyMetadatas, emptyMetadatas, emptyMetadatas, anotherFakeIlkDripMetadatas)

		var fakeIlkId int
		err := db.Get(&fakeIlkId, `SELECT id FROM maker.ilks WHERE ilk = $1`, fakeIlk)
		Expect(err).NotTo(HaveOccurred())

		var fakeIlkBlockOneDbResult IlkState
		err = db.Get(&fakeIlkBlockOneDbResult, `SELECT ilk, created, updated from maker.get_ilk_at_block_number($1, $2)`, blockOne, fakeIlkId)
		Expect(err).NotTo(HaveOccurred())
		expectedBlockOneFakeIlkState := IlkState{
			Ilk:     fakeIlk,
			Created: sql.NullString{String: blockOneHeader.Timestamp, Valid: true},
			Updated: sql.NullString{String: blockOneHeader.Timestamp, Valid: true},
		}
		Expect(fakeIlkBlockOneDbResult).To(Equal(expectedBlockOneFakeIlkState))

		var fakeIlkBlockTwoDbResult IlkState
		err = db.Get(&fakeIlkBlockTwoDbResult, `SELECT ilk, created, updated from maker.get_ilk_at_block_number($1, $2)`, blockTwo, fakeIlkId)
		Expect(err).NotTo(HaveOccurred())
		expectedBlockTwoFakeIlkState := IlkState{
			Ilk:     fakeIlk,
			Created: sql.NullString{String: blockOneHeader.Timestamp, Valid: true},
			Updated: sql.NullString{String: blockTwoHeader.Timestamp, Valid: true},
		}
		Expect(fakeIlkBlockTwoDbResult).To(Equal(expectedBlockTwoFakeIlkState))

		var anotherFakeIlkId int
		err = db.Get(&anotherFakeIlkId, `SELECT id FROM maker.ilks WHERE ilk = $1`, anotherFakeIlk)
		Expect(err).NotTo(HaveOccurred())

		var anotherFakeIlkDbResult IlkState
		err = db.Get(&anotherFakeIlkDbResult, `SELECT ilk, created, updated from maker.get_ilk_at_block_number($1, $2)`, blockTwo, anotherFakeIlkId)
		Expect(err).NotTo(HaveOccurred())
		expectedBlockTwoAnotherFakeIlkState := IlkState{
			Ilk:     anotherFakeIlk,
			Created: sql.NullString{String: blockTwoHeader.Timestamp, Valid: true},
			Updated: sql.NullString{String: blockTwoHeader.Timestamp, Valid: true},
		}

		Expect(anotherFakeIlkDbResult).To(Equal(expectedBlockTwoAnotherFakeIlkState))
	})
})

func getIlkState(seed int) map[string]string {
	valuesMap := make(map[string]string)
	valuesMap[vat.IlkTake] = strconv.Itoa(1 + seed)
	valuesMap[vat.IlkRate] = strconv.Itoa(2 + seed)
	valuesMap[vat.IlkInk] = strconv.Itoa(3 + seed)
	valuesMap[vat.IlkArt] = strconv.Itoa(4 + seed)
	valuesMap[pit.IlkSpot] = strconv.Itoa(5 + seed)
	valuesMap[pit.IlkLine] = strconv.Itoa(6 + seed)
	valuesMap[cat.IlkChop] = strconv.Itoa(7 + seed)
	valuesMap[cat.IlkLump] = strconv.Itoa(8 + seed)
	valuesMap[cat.IlkFlip] = "an address" + strconv.Itoa(seed)
	valuesMap[drip.IlkRho] = strconv.Itoa(9 + seed)
	valuesMap[drip.IlkTax] = strconv.Itoa(10 + seed)

	return valuesMap
}

func getMetadata(fieldType, ilk string, valueType utils.ValueType) utils.StorageValueMetadata {
	return utils.GetStorageValueMetadata(fieldType, map[utils.Key]string{constants.Ilk: ilk}, valueType)
}

func createIlkAtBlock(header core.Header, valuesMap map[string]string, vatMetadatas, pitMetadatas, catMetadatas, dripMetadatas []utils.StorageValueMetadata) {
	blockHash := header.Hash
	blockNumber := int(header.BlockNumber)
	for _, metadata := range vatMetadatas {
		value := valuesMap[metadata.Name]
		err := vatRepository.Create(blockNumber, blockHash, metadata, value)

		Expect(err).NotTo(HaveOccurred())
	}

	for _, metadata := range pitMetadatas {
		value := valuesMap[metadata.Name]
		err := pitRepository.Create(blockNumber, blockHash, metadata, value)

		Expect(err).NotTo(HaveOccurred())
	}

	for _, metadata := range catMetadatas {
		value := valuesMap[metadata.Name]
		err := catRepository.Create(blockNumber, blockHash, metadata, value)

		Expect(err).NotTo(HaveOccurred())
	}

	for _, metadata := range dripMetadatas {
		value := valuesMap[metadata.Name]
		err := dripRepository.Create(blockNumber, blockHash, metadata, value)

		Expect(err).NotTo(HaveOccurred())
	}
}

type IlkState struct {
	Ilk     string
	Take    string
	Rate    string
	Ink     string
	Art     string
	Spot    string
	Line    string
	Chop    string
	Lump    string
	Flip    string
	Rho     string
	Tax     string
	Created sql.NullString
	Updated sql.NullString
}
