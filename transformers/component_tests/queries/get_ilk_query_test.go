package queries

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/storage/cat"
	"github.com/vulcanize/mcd_transformers/transformers/storage/jug"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vat"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
)

var _ = Describe("Ilk State Query", func() {
	var (
		db               *postgres.DB
		blockOne         = rand.Int()
		blockTwo         = blockOne + 1
		blockThree       = blockOne + 2
		blockOneHeader   core.Header
		blockTwoHeader   core.Header
		blockThreeHeader core.Header
		vatRepository    vat.VatStorageRepository
		catRepository    cat.CatStorageRepository
		jugRepository    jug.JugStorageRepository
		headerRepository repositories.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		vatRepository.SetDB(db)
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
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	It("gets an ilk", func() {
		ilkValues := test_helpers.GetIlkValues(0)
		test_helpers.CreateIlk(db, blockOneHeader, ilkValues, test_helpers.FakeIlkVatMetadatas,
			test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas)

		var dbResult test_helpers.IlkState
		err := db.Get(&dbResult,
			`SELECT ilk_name, rate, art, spot, line, dust, chop, lump, flip, rho, duty, created, updated from api.get_ilk($1, $2)`,
			test_helpers.FakeIlk.Name, blockOne)
		Expect(err).NotTo(HaveOccurred())

		expectedIlk := test_helpers.IlkStateFromValues(
			test_helpers.FakeIlk.Hex, blockOneHeader.Timestamp, blockOneHeader.Timestamp, ilkValues)
		Expect(dbResult).To(Equal(expectedIlk))
	})

	It("returns the correct data if there are several ilks", func() {
		ilkValues := test_helpers.GetIlkValues(1)
		anotherIlkValues := test_helpers.GetIlkValues(2)
		test_helpers.CreateIlk(db, blockOneHeader, ilkValues, test_helpers.FakeIlkVatMetadatas,
			test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas)
		test_helpers.CreateIlk(db, blockOneHeader, anotherIlkValues, test_helpers.AnotherFakeIlkVatMetadatas,
			test_helpers.AnotherFakeIlkCatMetadatas, test_helpers.AnotherFakeIlkJugMetadatas)

		var fakeIlkResult test_helpers.IlkState
		err := db.Get(&fakeIlkResult,
			`SELECT ilk_name, rate, art, spot, line, dust, chop, lump, flip, rho, duty, created, updated from api.get_ilk($1, $2)`,
			test_helpers.FakeIlk.Name, blockOne)
		Expect(err).NotTo(HaveOccurred())

		var anotherFakeIlkResult test_helpers.IlkState
		err = db.Get(&anotherFakeIlkResult,
			`SELECT ilk_name, rate, art, spot, line, dust, chop, lump, flip, rho, duty, created, updated from api.get_ilk($1, $2)`,
			test_helpers.AnotherFakeIlk.Name, blockOne)
		Expect(err).NotTo(HaveOccurred())

		expectedFakeIlk := test_helpers.IlkStateFromValues(
			test_helpers.FakeIlk.Hex, blockOneHeader.Timestamp, blockOneHeader.Timestamp, ilkValues)
		expectedAnotherFakeIlk := test_helpers.IlkStateFromValues(
			test_helpers.AnotherFakeIlk.Hex, blockOneHeader.Timestamp, blockOneHeader.Timestamp, anotherIlkValues)

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
			test_helpers.CreateIlk(db, blockOneHeader, fakeIlkvalues, test_helpers.FakeIlkVatMetadatas,
				test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas)

			var blockOneDbResult test_helpers.IlkState
			err := db.Get(&blockOneDbResult, `SELECT ilk_name, rate, art, spot, line, dust, chop, lump, flip, rho, duty, created, updated from api.get_ilk($1, $2)`,
				test_helpers.FakeIlk.Name, blockOne)
			Expect(err).NotTo(HaveOccurred())

			expectedIlk := test_helpers.IlkStateFromValues(test_helpers.FakeIlk.Hex, blockOneHeader.Timestamp,
				blockOneHeader.Timestamp, fakeIlkvalues)
			Expect(blockOneDbResult).To(Equal(expectedIlk))
		})

		It("gets the Ilk for block two", func() {
			blockOneFakeIlkValues := test_helpers.GetIlkValues(1)
			test_helpers.CreateIlk(db, blockOneHeader, blockOneFakeIlkValues, test_helpers.FakeIlkVatMetadatas,
				test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas)
			// block two doesn't update rate or art
			blockTwoFakeIlkValues := test_helpers.GetIlkValues(2)
			vatMetadatasWithoutRateOrArt := []utils.StorageValueMetadata{test_helpers.FakeIlkSpotMetadata, test_helpers.FakeIlkLineMetadata}
			test_helpers.CreateIlk(db, blockTwoHeader, blockTwoFakeIlkValues, vatMetadatasWithoutRateOrArt,
				test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas)

			var blockTwoDbResult test_helpers.IlkState
			err := db.Get(&blockTwoDbResult, `SELECT rate, art, spot, line from api.get_ilk($1, $2)`,
				test_helpers.FakeIlk.Name, blockTwo)
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
			test_helpers.CreateIlk(db, blockOneHeader, blockOneFakeIlkValues, test_helpers.FakeIlkVatMetadatas,
				test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas)
			// block two doesn't update rate or art
			blockTwoFakeIlkValues := test_helpers.GetIlkValues(1)
			vatMetadatasWithoutRateOrArt := []utils.StorageValueMetadata{test_helpers.FakeIlkSpotMetadata, test_helpers.FakeIlkLineMetadata}
			test_helpers.CreateIlk(db, blockTwoHeader, blockTwoFakeIlkValues, vatMetadatasWithoutRateOrArt,
				test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas)
			// block three doesn't update art
			blockThreeFakeIlkValues := test_helpers.GetIlkValues(3)
			vatMetadatasWithoutArt := []utils.StorageValueMetadata{test_helpers.FakeIlkRateMetadata, test_helpers.FakeIlkSpotMetadata, test_helpers.FakeIlkLineMetadata}
			test_helpers.CreateIlk(db, blockThreeHeader, blockThreeFakeIlkValues, vatMetadatasWithoutArt,
				test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas)

			var blockThreeDbResult test_helpers.IlkState
			err := db.Get(&blockThreeDbResult, `SELECT rate, art, spot, line from api.get_ilk($1, $2)`,
				test_helpers.FakeIlk.Name, blockThree)
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
			test_helpers.CreateIlk(db, blockOneHeader, blockOneFakeIlkValues, test_helpers.FakeIlkVatMetadatas,
				test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas)
			blockOneAnotherFakeIlkState := test_helpers.GetIlkValues(2)
			test_helpers.CreateIlk(db, blockOneHeader, blockOneAnotherFakeIlkState, test_helpers.AnotherFakeIlkVatMetadatas,
				test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas)
			// block two doesn't update rate or art for fakeIlk
			// and doesn't update rate, art or line for anotherFakeIlk
			blockTwoFakeIlkValues := test_helpers.GetIlkValues(1)
			vatMetadatasWithoutRateOrArt := []utils.StorageValueMetadata{test_helpers.FakeIlkSpotMetadata, test_helpers.FakeIlkLineMetadata}
			test_helpers.CreateIlk(db, blockTwoHeader, blockTwoFakeIlkValues, vatMetadatasWithoutRateOrArt,
				test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas)
			// block three doesn't update ink
			// and doesn't update take, rate or ink for anotherFakeIlk
			blockThreeFakeIlkValues := test_helpers.GetIlkValues(3)
			vatMetadatasWithoutRate := []utils.StorageValueMetadata{test_helpers.FakeIlkArtMetadata, test_helpers.FakeIlkSpotMetadata, test_helpers.FakeIlkLineMetadata}
			test_helpers.CreateIlk(db, blockThreeHeader, blockThreeFakeIlkValues, vatMetadatasWithoutRate,
				test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas)

			var fakeIlkResult test_helpers.IlkState
			err := db.Get(&fakeIlkResult, `SELECT ilk_name, rate, art, spot, line from api.get_ilk($1, $2)`,
				test_helpers.FakeIlk.Name, blockThree)
			Expect(err).NotTo(HaveOccurred())

			var anotherFakeIlkResult test_helpers.IlkState
			err = db.Get(&anotherFakeIlkResult, `SELECT ilk_name, rate, art, spot, line from api.get_ilk($1, $2)`,
				test_helpers.AnotherFakeIlk.Name, blockThree)
			Expect(err).NotTo(HaveOccurred())

			blockThreeExpectedFakeIlk := test_helpers.IlkState{
				IlkName: test_helpers.FakeIlk.Name,
				Rate:    blockOneFakeIlkValues[vat.IlkRate], // value hasn't changed since block 1
				Art:     blockThreeFakeIlkValues[vat.IlkArt],
				Spot:    blockThreeFakeIlkValues[vat.IlkSpot],
				Line:    blockThreeFakeIlkValues[vat.IlkLine],
			}
			blockThreeExpectedAnotherFakeIlk := test_helpers.IlkState{
				IlkName: test_helpers.AnotherFakeIlk.Name,
				Rate:    blockOneAnotherFakeIlkState[vat.IlkRate], // value hasn't changed since block 1
				Art:     blockOneAnotherFakeIlkState[vat.IlkArt],  // value hasn't changed since block 1
				Spot:    blockOneAnotherFakeIlkState[vat.IlkSpot], // value hasn't changed since block 1
				Line:    blockOneAnotherFakeIlkState[vat.IlkLine], // value hasn't changed since block 1
			}
			Expect(fakeIlkResult).To(Equal(blockThreeExpectedFakeIlk))
			Expect(anotherFakeIlkResult).To(Equal(blockThreeExpectedAnotherFakeIlk))
		})
	})

	It("returns the rest of the Ilk data if a field has missing data", func() {
		vatMetadatasWithoutRate := test_helpers.FakeIlkVatMetadatas[1:]
		blockOneFakeIlkValues := test_helpers.GetIlkValues(1)
		blockOneFakeIlkValues[vat.IlkRate] = ""
		test_helpers.CreateIlk(db, blockOneHeader, blockOneFakeIlkValues, vatMetadatasWithoutRate,
			test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas)

		var blockOneDbResult test_helpers.IlkState
		err := db.Get(&blockOneDbResult, `SELECT ilk_name, art, spot, line, dust, chop, lump, flip, rho, duty, created, updated from api.get_ilk($1, $2)`,
			test_helpers.FakeIlk.Name, blockOne)

		Expect(err).NotTo(HaveOccurred())

		expectedIlk := test_helpers.IlkStateFromValues(test_helpers.FakeIlk.Hex, blockOneHeader.Timestamp, blockOneHeader.Timestamp, blockOneFakeIlkValues)
		Expect(blockOneDbResult).To(Equal(expectedIlk))
	})

	It("gets the correct timestamps", func() {
		// fakeIlk created at block1
		// fakeIlk updated at block2
		// anotherFakeIlk created at block2
		blockOneFakeIlkValues := test_helpers.GetIlkValues(0)
		test_helpers.CreateIlk(db, blockOneHeader, blockOneFakeIlkValues, test_helpers.FakeIlkVatMetadatas,
			test_helpers.EmptyMetadatas, test_helpers.EmptyMetadatas)

		blockTwoFakeIlkValues := test_helpers.GetIlkValues(1)
		test_helpers.CreateIlk(db, blockTwoHeader, blockTwoFakeIlkValues, test_helpers.EmptyMetadatas,
			test_helpers.EmptyMetadatas, test_helpers.FakeIlkJugMetadatas)

		blockTwoAnotherFakeIlkValues := test_helpers.GetIlkValues(2)
		test_helpers.CreateIlk(db, blockTwoHeader, blockTwoAnotherFakeIlkValues, test_helpers.EmptyMetadatas,
			test_helpers.EmptyMetadatas, test_helpers.AnotherFakeIlkJugMetadatas)

		var fakeIlkBlockOneDbResult test_helpers.IlkState
		err := db.Get(&fakeIlkBlockOneDbResult, `SELECT ilk_name, created, updated from api.get_ilk($1, $2)`,
			test_helpers.FakeIlk.Name, blockOne)
		Expect(err).NotTo(HaveOccurred())
		expectedBlockOneFakeIlkState := test_helpers.IlkState{
			IlkName: test_helpers.FakeIlk.Name,
			Created: test_helpers.GetValidNullString("1973-07-10T00:11:51Z"),
			Updated: test_helpers.GetValidNullString("1973-07-10T00:11:51Z"),
		}
		Expect(fakeIlkBlockOneDbResult).To(Equal(expectedBlockOneFakeIlkState))

		var fakeIlkBlockTwoDbResult test_helpers.IlkState
		err = db.Get(&fakeIlkBlockTwoDbResult, `SELECT ilk_name, created, updated from api.get_ilk($1, $2)`,
			test_helpers.FakeIlk.Name, blockTwo)
		Expect(err).NotTo(HaveOccurred())
		expectedBlockTwoFakeIlkState := test_helpers.IlkState{
			IlkName: test_helpers.FakeIlk.Name,
			Created: test_helpers.GetValidNullString("1973-07-10T00:11:51Z"),
			Updated: test_helpers.GetValidNullString("2005-03-18T01:58:32Z"),
		}
		Expect(fakeIlkBlockTwoDbResult).To(Equal(expectedBlockTwoFakeIlkState))

		var anotherFakeIlkDbResult test_helpers.IlkState
		err = db.Get(&anotherFakeIlkDbResult, `SELECT ilk_name, created, updated from api.get_ilk($1, $2)`,
			test_helpers.AnotherFakeIlk.Name, blockTwo)
		Expect(err).NotTo(HaveOccurred())
		expectedBlockTwoAnotherFakeIlkState := test_helpers.IlkState{
			IlkName: "FKE2",
			Created: test_helpers.GetValidNullString("2005-03-18T01:58:32Z"),
			Updated: test_helpers.GetValidNullString("2005-03-18T01:58:32Z"),
		}

		Expect(anotherFakeIlkDbResult).To(Equal(expectedBlockTwoAnotherFakeIlkState))
	})
})
