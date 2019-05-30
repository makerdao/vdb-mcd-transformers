package queries

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage/cat"
	"github.com/vulcanize/mcd_transformers/transformers/storage/jug"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vat"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"strconv"
	"time"
)

var _ = Describe("All Ilks query", func() {

	var (
		db               *postgres.DB
		vatRepository    vat.VatStorageRepository
		catRepository    cat.CatStorageRepository
		jugRepository    jug.JugStorageRepository
		headerRepository repositories.HeaderRepository

		blockOneTimestamp  = int64(111111111)
		fakeIlk            = test_helpers.FakeIlk
		blockOneHeader     = fakes.GetFakeHeaderWithTimestamp(blockOneTimestamp, int64(1))
		fakeIlkStateBlock1 = test_helpers.GetIlkValues(1)

		blockTwoTimestamp         = int64(222222222)
		anotherFakeIlk            = test_helpers.AnotherFakeIlk
		blockTwoHeader            = fakes.GetFakeHeaderWithTimestamp(blockTwoTimestamp, int64(2))
		anotherFakeIlkStateBlock2 = test_helpers.GetIlkValues(2)
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		vatRepository.SetDB(db)
		catRepository.SetDB(db)
		jugRepository.SetDB(db)

		//creating fakeIlk at block 1
		test_helpers.CreateVatRecords(blockOneHeader, fakeIlkStateBlock1, test_helpers.FakeIlkVatMetadatas, vatRepository)
		test_helpers.CreateCatRecords(blockOneHeader, fakeIlkStateBlock1, test_helpers.FakeIlkCatMetadatas, catRepository)
		test_helpers.CreateJugRecords(blockOneHeader, fakeIlkStateBlock1, test_helpers.FakeIlkJugMetadatas, jugRepository)
		//creating anotherFakeIlk at block 2
		test_helpers.CreateVatRecords(blockTwoHeader, anotherFakeIlkStateBlock2, test_helpers.AnotherFakeIlkVatMetadatas, vatRepository)
		test_helpers.CreateCatRecords(blockTwoHeader, anotherFakeIlkStateBlock2, test_helpers.AnotherFakeIlkCatMetadatas, catRepository)
		test_helpers.CreateJugRecords(blockTwoHeader, anotherFakeIlkStateBlock2, test_helpers.AnotherFakeIlkJugMetadatas, jugRepository)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Context("When the headerSync is complete", func() {
		BeforeEach(func() {
			_, err := headerRepository.CreateOrUpdateHeader(blockOneHeader)
			Expect(err).NotTo(HaveOccurred())
			_, err = headerRepository.CreateOrUpdateHeader(blockTwoHeader)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns ilks as of block 1", func() {
			var dbResult []test_helpers.IlkState
			expectedResult := test_helpers.IlkState{
				IlkIdentifier: fakeIlk.Identifier,
				Rate:          fakeIlkStateBlock1[vat.IlkRate],
				Art:           fakeIlkStateBlock1[vat.IlkArt],
				Spot:          fakeIlkStateBlock1[vat.IlkSpot],
				Line:          fakeIlkStateBlock1[vat.IlkLine],
				Dust:          fakeIlkStateBlock1[vat.IlkDust],
				Chop:          fakeIlkStateBlock1[cat.IlkChop],
				Lump:          fakeIlkStateBlock1[cat.IlkLump],
				Flip:          fakeIlkStateBlock1[cat.IlkFlip],
				Rho:           fakeIlkStateBlock1[jug.IlkRho],
				Duty:          fakeIlkStateBlock1[jug.IlkDuty],
				Created:       test_helpers.GetValidNullString(getFormattedTimestamp(blockOneHeader.Timestamp)),
				Updated:       test_helpers.GetValidNullString(getFormattedTimestamp(blockOneHeader.Timestamp)),
			}
			err := db.Select(&dbResult,
				`SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, created, updated FROM api.all_ilks($1)`,
				blockOneHeader.BlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbResult)).To(Equal(1))
			Expect(dbResult[0]).To(Equal(expectedResult))
		})

		It("returns ilks as of block 2", func() {
			var dbResult []test_helpers.IlkState
			//fakeIlk was created at block 1
			fakeIlkExpectedResult := test_helpers.IlkState{
				IlkIdentifier: fakeIlk.Identifier,
				Rate:          fakeIlkStateBlock1[vat.IlkRate],
				Art:           fakeIlkStateBlock1[vat.IlkArt],
				Spot:          fakeIlkStateBlock1[vat.IlkSpot],
				Line:          fakeIlkStateBlock1[vat.IlkLine],
				Dust:          fakeIlkStateBlock1[vat.IlkDust],
				Chop:          fakeIlkStateBlock1[cat.IlkChop],
				Lump:          fakeIlkStateBlock1[cat.IlkLump],
				Flip:          fakeIlkStateBlock1[cat.IlkFlip],
				Rho:           fakeIlkStateBlock1[jug.IlkRho],
				Duty:          fakeIlkStateBlock1[jug.IlkDuty],
				Created:       test_helpers.GetValidNullString(getFormattedTimestamp(blockOneHeader.Timestamp)),
				Updated:       test_helpers.GetValidNullString(getFormattedTimestamp(blockOneHeader.Timestamp)),
			}
			//anotherFakeIlk was created at block 2
			anotherFakeIlkExpectedResult := test_helpers.IlkState{
				IlkIdentifier: anotherFakeIlk.Identifier,
				Rate:          anotherFakeIlkStateBlock2[vat.IlkRate],
				Art:           anotherFakeIlkStateBlock2[vat.IlkArt],
				Spot:          anotherFakeIlkStateBlock2[vat.IlkSpot],
				Line:          anotherFakeIlkStateBlock2[vat.IlkLine],
				Dust:          anotherFakeIlkStateBlock2[vat.IlkDust],
				Chop:          anotherFakeIlkStateBlock2[cat.IlkChop],
				Lump:          anotherFakeIlkStateBlock2[cat.IlkLump],
				Flip:          anotherFakeIlkStateBlock2[cat.IlkFlip],
				Rho:           anotherFakeIlkStateBlock2[jug.IlkRho],
				Duty:          anotherFakeIlkStateBlock2[jug.IlkDuty],
				Created:       test_helpers.GetValidNullString(getFormattedTimestamp(blockTwoHeader.Timestamp)),
				Updated:       test_helpers.GetValidNullString(getFormattedTimestamp(blockTwoHeader.Timestamp)),
			}
			err := db.Select(&dbResult,
				`SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, created, updated FROM api.all_ilks($1)`,
				blockTwoHeader.BlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbResult)).To(Equal(2))
			Expect(dbResult).To(ConsistOf(fakeIlkExpectedResult, anotherFakeIlkExpectedResult))
		})

		It("returns updated values as of block 3", func() {
			blockThreeTimestamp := int64(333333333)
			blockThreeHeader := fakes.GetFakeHeaderWithTimestamp(blockThreeTimestamp, int64(3))

			_, headerErr := headerRepository.CreateOrUpdateHeader(blockThreeHeader)
			Expect(headerErr).NotTo(HaveOccurred())

			//updating fakeIlk spot value at block 3
			fakeIlkStateBlock3 := test_helpers.GetIlkValues(3)
			spotMetadata := []utils.StorageValueMetadata{test_helpers.FakeIlkSpotMetadata}
			test_helpers.CreateVatRecords(blockThreeHeader, fakeIlkStateBlock3, spotMetadata, vatRepository)

			//updating all anotherFakeIlk values at block 3
			anotherFakeIlkStateBlock3 := test_helpers.GetIlkValues(4)
			test_helpers.CreateVatRecords(blockThreeHeader, anotherFakeIlkStateBlock3, test_helpers.AnotherFakeIlkVatMetadatas, vatRepository)
			test_helpers.CreateCatRecords(blockThreeHeader, anotherFakeIlkStateBlock3, test_helpers.AnotherFakeIlkCatMetadatas, catRepository)
			test_helpers.CreateJugRecords(blockThreeHeader, anotherFakeIlkStateBlock3, test_helpers.AnotherFakeIlkJugMetadatas, jugRepository)

			var dbResult []test_helpers.IlkState
			fakeIlkExpectedResult := test_helpers.IlkState{
				IlkIdentifier: fakeIlk.Identifier,
				Rate:          fakeIlkStateBlock1[vat.IlkRate],
				Art:           fakeIlkStateBlock1[vat.IlkArt],
				Spot:          fakeIlkStateBlock3[vat.IlkSpot],
				Line:          fakeIlkStateBlock1[vat.IlkLine],
				Dust:          fakeIlkStateBlock1[vat.IlkDust],
				Chop:          fakeIlkStateBlock1[cat.IlkChop],
				Lump:          fakeIlkStateBlock1[cat.IlkLump],
				Flip:          fakeIlkStateBlock1[cat.IlkFlip],
				Rho:           fakeIlkStateBlock1[jug.IlkRho],
				Duty:          fakeIlkStateBlock1[jug.IlkDuty],
				Created:       test_helpers.GetValidNullString(getFormattedTimestamp(blockOneHeader.Timestamp)),
				Updated:       test_helpers.GetValidNullString(getFormattedTimestamp(blockThreeHeader.Timestamp)),
			}
			anotherFakeIlkExpectedResult := test_helpers.IlkState{
				IlkIdentifier: anotherFakeIlk.Identifier,
				Rate:          anotherFakeIlkStateBlock3[vat.IlkRate],
				Art:           anotherFakeIlkStateBlock3[vat.IlkArt],
				Spot:          anotherFakeIlkStateBlock3[vat.IlkSpot],
				Line:          anotherFakeIlkStateBlock3[vat.IlkLine],
				Dust:          anotherFakeIlkStateBlock3[vat.IlkDust],
				Chop:          anotherFakeIlkStateBlock3[cat.IlkChop],
				Lump:          anotherFakeIlkStateBlock3[cat.IlkLump],
				Flip:          anotherFakeIlkStateBlock3[cat.IlkFlip],
				Rho:           anotherFakeIlkStateBlock3[jug.IlkRho],
				Duty:          anotherFakeIlkStateBlock3[jug.IlkDuty],
				Created:       test_helpers.GetValidNullString(getFormattedTimestamp(blockTwoHeader.Timestamp)),
				Updated:       test_helpers.GetValidNullString(getFormattedTimestamp(blockThreeHeader.Timestamp)),
			}
			err := db.Select(&dbResult,
				`SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, created, updated FROM api.all_ilks($1)`,
				blockThreeHeader.BlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbResult)).To(Equal(2))
			Expect(dbResult).To(ConsistOf(fakeIlkExpectedResult, anotherFakeIlkExpectedResult))
		})

		It("uses default value for blockHeight if not supplied", func() {
			_, err := db.Exec(`SELECT * FROM api.all_ilks()`)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	It("returns ilk states without timestamps if the corresponding header hasn't been synced yet", func() {
		var dbResult []test_helpers.IlkState
		expectedResult := test_helpers.IlkState{
			IlkIdentifier: fakeIlk.Identifier,
			Rate:          fakeIlkStateBlock1[vat.IlkRate],
			Art:           fakeIlkStateBlock1[vat.IlkArt],
			Spot:          fakeIlkStateBlock1[vat.IlkSpot],
			Line:          fakeIlkStateBlock1[vat.IlkLine],
			Dust:          fakeIlkStateBlock1[vat.IlkDust],
			Chop:          fakeIlkStateBlock1[cat.IlkChop],
			Lump:          fakeIlkStateBlock1[cat.IlkLump],
			Flip:          fakeIlkStateBlock1[cat.IlkFlip],
			Rho:           fakeIlkStateBlock1[jug.IlkRho],
			Duty:          fakeIlkStateBlock1[jug.IlkDuty],
			Created:       test_helpers.GetEmptyNullString(),
			Updated:       test_helpers.GetEmptyNullString(),
		}
		err := db.Select(&dbResult,
			`SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, created, updated FROM api.all_ilks($1)`,
			blockOneHeader.BlockNumber)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0]).To(Equal(expectedResult))
	})

	It("handles cases where some of the data is null", func() {
		blockFourHeader := fakes.GetFakeHeaderWithTimestamp(int64(444444444), int64(4))
		_, headerErr := headerRepository.CreateOrUpdateHeader(blockFourHeader)
		Expect(headerErr).NotTo(HaveOccurred())

		//updating fakeIlk spot value at block 3
		newIlk := test_helpers.TestIlk{Identifier: "newIlk", Hex: "6e6577496c6b0000000000000000000000000000000000000000000000000000"}
		newIlkStateBlock4 := test_helpers.GetIlkValues(4)
		metadata := []utils.StorageValueMetadata{{
			Name: vat.IlkRate,
			Keys: map[utils.Key]string{constants.Ilk: newIlk.Hex},
			Type: utils.Uint256,
		}}
		//only creating a vat_ilk_rate record
		test_helpers.CreateVatRecords(blockFourHeader, newIlkStateBlock4, metadata, vatRepository)

		var dbResult []test_helpers.IlkState
		newIlkExpectedResult := test_helpers.IlkState{
			IlkIdentifier: newIlk.Identifier,
			Rate:          newIlkStateBlock4[vat.IlkRate],
			Created:       test_helpers.GetValidNullString(getFormattedTimestamp(blockFourHeader.Timestamp)),
			Updated:       test_helpers.GetValidNullString(getFormattedTimestamp(blockFourHeader.Timestamp)),
		}
		err := db.Select(&dbResult,
			`SELECT ilk_identifier, rate, created, updated FROM api.all_ilks($1)`,
			blockFourHeader.BlockNumber)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(dbResult)).To(Equal(3))
		Expect(dbResult).To(ContainElement(newIlkExpectedResult))
	})
})

func getFormattedTimestamp(timestampString string) string {
	parsed, _ := strconv.ParseInt(timestampString, 10, 64)
	return time.Unix(parsed, 0).UTC().Format(time.RFC3339)
}
