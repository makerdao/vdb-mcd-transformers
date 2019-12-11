package queries

import (
	"math/rand"

	storage_helper "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/cat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/jug"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/spot"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/vat"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/utils"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ilk State Query", func() {
	var (
		headerOne, headerTwo, headerThree core.Header
		vatRepository                     vat.VatStorageRepository
		catRepository                     cat.CatStorageRepository
		jugRepository                     jug.JugStorageRepository
		spotRepository                    spot.SpotStorageRepository
		headerRepository                  repositories.HeaderRepository
		diffID                            int64
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		vatRepository.SetDB(db)
		catRepository.SetDB(db)
		jugRepository.SetDB(db)
		spotRepository.SetDB(db)
		headerRepository = repositories.NewHeaderRepository(db)

		blockOne := rand.Int()
		blockTwo := blockOne + 1
		blockThree := blockTwo + 1
		timestampOne := int(rand.Int31())
		timestampTwo := timestampOne + 1000000
		timestampThree := timestampTwo + 1000000
		headerOne = createHeader(blockOne, timestampOne, headerRepository)
		headerTwo = createHeader(blockTwo, timestampTwo, headerRepository)
		headerThree = createHeader(blockThree, timestampThree, headerRepository)

		diffID = storage_helper.CreateFakeDiffRecord(db)
	})

	It("gets an ilk", func() {
		ilkValues := test_helpers.GetIlkValues(0)
		test_helpers.CreateIlk(db, diffID, headerOne, ilkValues, test_helpers.FakeIlkVatMetadatas, test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas, test_helpers.FakeIlkSpotMetadatas)

		var dbResult test_helpers.IlkState
		err := db.Get(&dbResult,
			`SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated from api.get_ilk($1, $2)`,
			test_helpers.FakeIlk.Identifier, headerOne.BlockNumber)
		Expect(err).NotTo(HaveOccurred())

		expectedIlk := test_helpers.IlkStateFromValues(
			test_helpers.FakeIlk.Hex, headerOne.Timestamp, headerOne.Timestamp, ilkValues)
		Expect(dbResult).To(Equal(expectedIlk))
	})

	It("returns the correct data if there are several ilks", func() {
		ilkValues := test_helpers.GetIlkValues(1)
		anotherIlkValues := test_helpers.GetIlkValues(2)
		test_helpers.CreateIlk(db, diffID, headerOne, ilkValues, test_helpers.FakeIlkVatMetadatas, test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas, test_helpers.FakeIlkSpotMetadatas)
		test_helpers.CreateIlk(db, diffID, headerOne, anotherIlkValues, test_helpers.AnotherFakeIlkVatMetadatas, test_helpers.AnotherFakeIlkCatMetadatas, test_helpers.AnotherFakeIlkJugMetadatas, test_helpers.AnotherFakeIlkSpotMetadatas)

		var fakeIlkResult test_helpers.IlkState
		err := db.Get(&fakeIlkResult,
			`SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated from api.get_ilk($1, $2)`,
			test_helpers.FakeIlk.Identifier, headerOne.BlockNumber)
		Expect(err).NotTo(HaveOccurred())

		var anotherFakeIlkResult test_helpers.IlkState
		err = db.Get(&anotherFakeIlkResult,
			`SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated from api.get_ilk($1, $2)`,
			test_helpers.AnotherFakeIlk.Identifier, headerOne.BlockNumber)
		Expect(err).NotTo(HaveOccurred())

		expectedFakeIlk := test_helpers.IlkStateFromValues(
			test_helpers.FakeIlk.Hex, headerOne.Timestamp, headerOne.Timestamp, ilkValues)
		expectedAnotherFakeIlk := test_helpers.IlkStateFromValues(
			test_helpers.AnotherFakeIlk.Hex, headerOne.Timestamp, headerOne.Timestamp, anotherIlkValues)

		Expect(fakeIlkResult).To(Equal(expectedFakeIlk))
		Expect(anotherFakeIlkResult).To(Equal(expectedAnotherFakeIlk))
	})

	It("fails if no argument is supplied (STRICT)", func() {
		_, err := db.Exec(`SELECT * FROM api.get_ilk()`)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("function api.get_ilk() does not exist"))
	})

	It("allows blockHeight argument to be omitted", func() {
		_, err := db.Exec(`SELECT * FROM api.get_ilk($1)`, 0)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("handles getting the most recent Ilk values as of a given block", func() {
		It("gets the Ilk for block one", func() {
			fakeIlkvalues := test_helpers.GetIlkValues(0)
			test_helpers.CreateIlk(db, diffID, headerOne, fakeIlkvalues, test_helpers.FakeIlkVatMetadatas, test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas, test_helpers.FakeIlkSpotMetadatas)

			var blockOneDbResult test_helpers.IlkState
			err := db.Get(&blockOneDbResult, `SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated from api.get_ilk($1, $2)`,
				test_helpers.FakeIlk.Identifier, headerOne.BlockNumber)
			Expect(err).NotTo(HaveOccurred())

			expectedIlk := test_helpers.IlkStateFromValues(test_helpers.FakeIlk.Hex, headerOne.Timestamp,
				headerOne.Timestamp, fakeIlkvalues)
			Expect(blockOneDbResult).To(Equal(expectedIlk))
		})

		It("gets the Ilk for block two", func() {
			blockOneFakeIlkValues := test_helpers.GetIlkValues(1)
			test_helpers.CreateIlk(db, diffID, headerOne, blockOneFakeIlkValues, test_helpers.FakeIlkVatMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas)
			// block two doesn't update rate or art
			blockTwoFakeIlkValues := test_helpers.GetIlkValues(2)
			vatMetadatasWithoutRateOrArt := []utils.StorageValueMetadata{test_helpers.FakeIlkSpotMetadata, test_helpers.FakeIlkLineMetadata}
			test_helpers.CreateIlk(db, diffID, headerTwo, blockTwoFakeIlkValues, vatMetadatasWithoutRateOrArt, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas)

			var blockTwoDbResult test_helpers.IlkState
			err := db.Get(&blockTwoDbResult, `SELECT rate, art, spot, line from api.get_ilk($1, $2)`,
				test_helpers.FakeIlk.Identifier, headerTwo.BlockNumber)
			Expect(err).NotTo(HaveOccurred())

			blockTwoExpectedIlk := test_helpers.IlkState{
				Rate: blockOneFakeIlkValues[vat.IlkRate], // value hasn't changed since block 1
				Art:  blockOneFakeIlkValues[vat.IlkArt],  // value hasn't changed since block 1
				Spot: blockTwoFakeIlkValues[vat.IlkSpot],
				Line: blockTwoFakeIlkValues[vat.IlkLine],
			}
			Expect(blockTwoDbResult).To(Equal(blockTwoExpectedIlk))
		})

		It("gets the Ilk for block three", func() {
			//no updates to ink
			blockOneFakeIlkValues := test_helpers.GetIlkValues(1)
			test_helpers.CreateIlk(db, diffID, headerOne, blockOneFakeIlkValues, test_helpers.FakeIlkVatMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas)
			// block two doesn't update rate or art
			blockTwoFakeIlkValues := test_helpers.GetIlkValues(1)
			vatMetadatasWithoutRateOrArt := []utils.StorageValueMetadata{test_helpers.FakeIlkSpotMetadata, test_helpers.FakeIlkLineMetadata}
			test_helpers.CreateIlk(db, diffID, headerTwo, blockTwoFakeIlkValues, vatMetadatasWithoutRateOrArt, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas)
			// block three doesn't update art
			blockThreeFakeIlkValues := test_helpers.GetIlkValues(3)
			vatMetadatasWithoutArt := []utils.StorageValueMetadata{test_helpers.FakeIlkRateMetadata, test_helpers.FakeIlkSpotMetadata, test_helpers.FakeIlkLineMetadata}
			test_helpers.CreateIlk(db, diffID, headerThree, blockThreeFakeIlkValues, vatMetadatasWithoutArt, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas)

			var blockThreeDbResult test_helpers.IlkState
			err := db.Get(&blockThreeDbResult, `SELECT rate, art, spot, line from api.get_ilk($1, $2)`,
				test_helpers.FakeIlk.Identifier, headerThree.BlockNumber)
			Expect(err).NotTo(HaveOccurred())

			blockThreeExpectedIlk := test_helpers.IlkState{
				Rate: blockThreeFakeIlkValues[vat.IlkRate],
				Art:  blockOneFakeIlkValues[vat.IlkArt], // hasn't been updated since block 1
				Spot: blockThreeFakeIlkValues[vat.IlkSpot],
				Line: blockThreeFakeIlkValues[vat.IlkLine],
			}
			Expect(blockThreeDbResult).To(Equal(blockThreeExpectedIlk))
		})

		It("gets more than one ilk as of block three", func() {
			blockOneFakeIlkValues := test_helpers.GetIlkValues(1)
			test_helpers.CreateIlk(db, diffID, headerOne, blockOneFakeIlkValues, test_helpers.FakeIlkVatMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas)
			blockOneAnotherFakeIlkState := test_helpers.GetIlkValues(2)
			test_helpers.CreateIlk(db, diffID, headerOne, blockOneAnotherFakeIlkState, test_helpers.AnotherFakeIlkVatMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas)
			// block two doesn't update rate or art for fakeIlk
			// and doesn't update rate, art or line for anotherFakeIlk
			blockTwoFakeIlkValues := test_helpers.GetIlkValues(1)
			vatMetadatasWithoutRateOrArt := []utils.StorageValueMetadata{test_helpers.FakeIlkSpotMetadata, test_helpers.FakeIlkLineMetadata}
			test_helpers.CreateIlk(db, diffID, headerTwo, blockTwoFakeIlkValues, vatMetadatasWithoutRateOrArt, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas)
			// block three doesn't update ink
			// and doesn't update take, rate or ink for anotherFakeIlk
			blockThreeFakeIlkValues := test_helpers.GetIlkValues(3)
			vatMetadatasWithoutRate := []utils.StorageValueMetadata{test_helpers.FakeIlkArtMetadata, test_helpers.FakeIlkSpotMetadata, test_helpers.FakeIlkLineMetadata}
			test_helpers.CreateIlk(db, diffID, headerThree, blockThreeFakeIlkValues, vatMetadatasWithoutRate, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas)

			var fakeIlkResult test_helpers.IlkState
			err := db.Get(&fakeIlkResult, `SELECT ilk_identifier, rate, art, spot, line from api.get_ilk($1, $2)`,
				test_helpers.FakeIlk.Identifier, headerThree.BlockNumber)
			Expect(err).NotTo(HaveOccurred())

			var anotherFakeIlkResult test_helpers.IlkState
			err = db.Get(&anotherFakeIlkResult, `SELECT ilk_identifier, rate, art, spot, line from api.get_ilk($1, $2)`,
				test_helpers.AnotherFakeIlk.Identifier, headerThree.BlockNumber)
			Expect(err).NotTo(HaveOccurred())

			blockThreeExpectedFakeIlk := test_helpers.IlkState{
				IlkIdentifier: test_helpers.FakeIlk.Identifier,
				Rate:          blockOneFakeIlkValues[vat.IlkRate], // value hasn't changed since block 1
				Art:           blockThreeFakeIlkValues[vat.IlkArt],
				Spot:          blockThreeFakeIlkValues[vat.IlkSpot],
				Line:          blockThreeFakeIlkValues[vat.IlkLine],
			}
			blockThreeExpectedAnotherFakeIlk := test_helpers.IlkState{
				IlkIdentifier: test_helpers.AnotherFakeIlk.Identifier,
				Rate:          blockOneAnotherFakeIlkState[vat.IlkRate], // value hasn't changed since block 1
				Art:           blockOneAnotherFakeIlkState[vat.IlkArt],  // value hasn't changed since block 1
				Spot:          blockOneAnotherFakeIlkState[vat.IlkSpot], // value hasn't changed since block 1
				Line:          blockOneAnotherFakeIlkState[vat.IlkLine], // value hasn't changed since block 1
			}
			Expect(fakeIlkResult).To(Equal(blockThreeExpectedFakeIlk))
			Expect(anotherFakeIlkResult).To(Equal(blockThreeExpectedAnotherFakeIlk))
		})
	})

	It("returns the rest of the Ilk data if a field has missing data", func() {
		vatMetadatasWithoutRate := test_helpers.FakeIlkVatMetadatas[1:]
		blockOneFakeIlkValues := test_helpers.GetIlkValues(1)
		blockOneFakeIlkValues[vat.IlkRate] = ""
		test_helpers.CreateIlk(db, diffID, headerOne, blockOneFakeIlkValues, vatMetadatasWithoutRate, test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas, test_helpers.FakeIlkSpotMetadatas)

		var blockOneDbResult test_helpers.IlkState
		err := db.Get(&blockOneDbResult, `SELECT ilk_identifier, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated from api.get_ilk($1, $2)`,
			test_helpers.FakeIlk.Identifier, headerOne.BlockNumber)
		Expect(err).NotTo(HaveOccurred())

		expectedIlk := test_helpers.IlkStateFromValues(test_helpers.FakeIlk.Hex, headerOne.Timestamp, headerOne.Timestamp, blockOneFakeIlkValues)
		Expect(blockOneDbResult).To(Equal(expectedIlk))
	})

	It("gets the correct timestamps", func() {
		// fakeIlk created at block1
		// fakeIlk updated at block2
		// anotherFakeIlk created at block2
		blockOneFakeIlkValues := test_helpers.GetIlkValues(0)
		test_helpers.CreateIlk(db, diffID, headerOne, blockOneFakeIlkValues, test_helpers.FakeIlkVatMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas)

		blockTwoFakeIlkValues := test_helpers.GetIlkValues(1)
		test_helpers.CreateIlk(db, diffID, headerTwo, blockTwoFakeIlkValues, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas, test_helpers.FakeIlkJugMetadatas, test_helpers.EmptyMetadatas)

		blockTwoAnotherFakeIlkValues := test_helpers.GetIlkValues(2)
		test_helpers.CreateIlk(db, diffID, headerTwo, blockTwoAnotherFakeIlkValues, test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas, test_helpers.AnotherFakeIlkJugMetadatas, test_helpers.EmptyMetadatas)

		var fakeIlkBlockOneDbResult test_helpers.IlkState
		err := db.Get(&fakeIlkBlockOneDbResult, `SELECT ilk_identifier, created, updated from api.get_ilk($1, $2)`,
			test_helpers.FakeIlk.Identifier, headerOne.BlockNumber)
		Expect(err).NotTo(HaveOccurred())

		formattedTimestampOne := getFormattedTimestamp(headerOne.Timestamp)
		formattedTimestampTwo := getFormattedTimestamp(headerTwo.Timestamp)

		expectedBlockOneFakeIlkState := test_helpers.IlkState{
			IlkIdentifier: test_helpers.FakeIlk.Identifier,
			Created:       test_helpers.GetValidNullString(formattedTimestampOne),
			Updated:       test_helpers.GetValidNullString(formattedTimestampOne),
		}
		Expect(fakeIlkBlockOneDbResult).To(Equal(expectedBlockOneFakeIlkState))

		var fakeIlkBlockTwoDbResult test_helpers.IlkState
		err = db.Get(&fakeIlkBlockTwoDbResult, `SELECT ilk_identifier, created, updated from api.get_ilk($1, $2)`,
			test_helpers.FakeIlk.Identifier, headerTwo.BlockNumber)
		Expect(err).NotTo(HaveOccurred())
		expectedBlockTwoFakeIlkState := test_helpers.IlkState{
			IlkIdentifier: test_helpers.FakeIlk.Identifier,
			Created:       test_helpers.GetValidNullString(formattedTimestampOne),
			Updated:       test_helpers.GetValidNullString(formattedTimestampTwo),
		}
		Expect(fakeIlkBlockTwoDbResult).To(Equal(expectedBlockTwoFakeIlkState))

		var anotherFakeIlkDbResult test_helpers.IlkState
		err = db.Get(&anotherFakeIlkDbResult, `SELECT ilk_identifier, created, updated from api.get_ilk($1, $2)`,
			test_helpers.AnotherFakeIlk.Identifier, headerTwo.BlockNumber)
		Expect(err).NotTo(HaveOccurred())
		expectedBlockTwoAnotherFakeIlkState := test_helpers.IlkState{
			IlkIdentifier: test_helpers.AnotherFakeIlk.Identifier,
			Created:       test_helpers.GetValidNullString(formattedTimestampTwo),
			Updated:       test_helpers.GetValidNullString(formattedTimestampTwo),
		}

		Expect(anotherFakeIlkDbResult).To(Equal(expectedBlockTwoAnotherFakeIlkState))
	})
})
