package queries

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
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

var _ = Describe("All Ilks query", func() {

	var (
		vatRepository    vat.VatStorageRepository
		catRepository    cat.CatStorageRepository
		jugRepository    jug.JugStorageRepository
		spotRepository   spot.SpotStorageRepository
		headerRepository repositories.HeaderRepository

		blockOne, blockTwo         int
		timestampOne, timestampTwo int
		headerOne, headerTwo       core.Header

		fakeIlk                   = test_helpers.FakeIlk
		fakeIlkStateBlock1        = test_helpers.GetIlkValues(1)
		anotherFakeIlk            = test_helpers.AnotherFakeIlk
		anotherFakeIlkStateBlock2 = test_helpers.GetIlkValues(2)
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		vatRepository.SetDB(db)
		catRepository.SetDB(db)
		jugRepository.SetDB(db)
		spotRepository.SetDB(db)

		blockOne = rand.Int()
		blockTwo = blockOne + 1
		timestampOne = int(rand.Int31())
		timestampTwo = timestampOne + 1
		headerOne = createHeader(blockOne, timestampOne, headerRepository)
		headerTwo = createHeader(blockTwo, timestampTwo, headerRepository)

		//creating fakeIlk at block 1
		test_helpers.CreateVatRecords(db, headerOne, fakeIlkStateBlock1, test_helpers.FakeIlkVatMetadatas, vatRepository)
		test_helpers.CreateCatRecords(db, headerOne, fakeIlkStateBlock1, test_helpers.FakeIlkCatMetadatas, catRepository)
		test_helpers.CreateJugRecords(db, headerOne, fakeIlkStateBlock1, test_helpers.FakeIlkJugMetadatas, jugRepository)
		test_helpers.CreateSpotRecords(db, headerOne, fakeIlkStateBlock1, test_helpers.FakeIlkSpotMetadatas, spotRepository)
		//creating anotherFakeIlk at block 2
		test_helpers.CreateVatRecords(db, headerTwo, anotherFakeIlkStateBlock2, test_helpers.AnotherFakeIlkVatMetadatas, vatRepository)
		test_helpers.CreateCatRecords(db, headerTwo, anotherFakeIlkStateBlock2, test_helpers.AnotherFakeIlkCatMetadatas, catRepository)
		test_helpers.CreateJugRecords(db, headerTwo, anotherFakeIlkStateBlock2, test_helpers.AnotherFakeIlkJugMetadatas, jugRepository)
		test_helpers.CreateSpotRecords(db, headerTwo, anotherFakeIlkStateBlock2, test_helpers.AnotherFakeIlkSpotMetadatas, spotRepository)
	})

	Context("When the headerSync is complete", func() {
		It("returns ilks as of block 1", func() {
			var dbResult []test_helpers.IlkState
			expectedResult := test_helpers.IlkState{
				IlkIdentifier: fakeIlk.Identifier,
				Rate:          fakeIlkStateBlock1[vat.IlkRate].(string),
				Art:           fakeIlkStateBlock1[vat.IlkArt].(string),
				Spot:          fakeIlkStateBlock1[vat.IlkSpot].(string),
				Line:          fakeIlkStateBlock1[vat.IlkLine].(string),
				Dust:          fakeIlkStateBlock1[vat.IlkDust].(string),
				Chop:          fakeIlkStateBlock1[cat.IlkChop].(string),
				Lump:          fakeIlkStateBlock1[cat.IlkLump].(string),
				Flip:          fakeIlkStateBlock1[cat.IlkFlip].(string),
				Rho:           fakeIlkStateBlock1[jug.IlkRho].(string),
				Duty:          fakeIlkStateBlock1[jug.IlkDuty].(string),
				Pip:           fakeIlkStateBlock1[spot.IlkPip].(string),
				Mat:           fakeIlkStateBlock1[spot.IlkMat].(string),
				Created:       test_helpers.GetValidNullString(getFormattedTimestamp(headerOne.Timestamp)),
				Updated:       test_helpers.GetValidNullString(getFormattedTimestamp(headerOne.Timestamp)),
			}
			err := db.Select(&dbResult,
				`SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated FROM api.all_ilks($1)`,
				headerOne.BlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbResult)).To(Equal(1))
			Expect(dbResult[0]).To(Equal(expectedResult))
		})

		It("returns ilks as of block 2", func() {
			var dbResult []test_helpers.IlkState
			//fakeIlk was created at block 1
			fakeIlkExpectedResult := test_helpers.IlkState{
				IlkIdentifier: fakeIlk.Identifier,
				Rate:          fakeIlkStateBlock1[vat.IlkRate].(string),
				Art:           fakeIlkStateBlock1[vat.IlkArt].(string),
				Spot:          fakeIlkStateBlock1[vat.IlkSpot].(string),
				Line:          fakeIlkStateBlock1[vat.IlkLine].(string),
				Dust:          fakeIlkStateBlock1[vat.IlkDust].(string),
				Chop:          fakeIlkStateBlock1[cat.IlkChop].(string),
				Lump:          fakeIlkStateBlock1[cat.IlkLump].(string),
				Flip:          fakeIlkStateBlock1[cat.IlkFlip].(string),
				Rho:           fakeIlkStateBlock1[jug.IlkRho].(string),
				Duty:          fakeIlkStateBlock1[jug.IlkDuty].(string),
				Pip:           fakeIlkStateBlock1[spot.IlkPip].(string),
				Mat:           fakeIlkStateBlock1[spot.IlkMat].(string),
				Created:       test_helpers.GetValidNullString(getFormattedTimestamp(headerOne.Timestamp)),
				Updated:       test_helpers.GetValidNullString(getFormattedTimestamp(headerOne.Timestamp)),
			}
			//anotherFakeIlk was created at block 2
			anotherFakeIlkExpectedResult := test_helpers.IlkState{
				IlkIdentifier: anotherFakeIlk.Identifier,
				Rate:          anotherFakeIlkStateBlock2[vat.IlkRate].(string),
				Art:           anotherFakeIlkStateBlock2[vat.IlkArt].(string),
				Spot:          anotherFakeIlkStateBlock2[vat.IlkSpot].(string),
				Line:          anotherFakeIlkStateBlock2[vat.IlkLine].(string),
				Dust:          anotherFakeIlkStateBlock2[vat.IlkDust].(string),
				Chop:          anotherFakeIlkStateBlock2[cat.IlkChop].(string),
				Lump:          anotherFakeIlkStateBlock2[cat.IlkLump].(string),
				Flip:          anotherFakeIlkStateBlock2[cat.IlkFlip].(string),
				Rho:           anotherFakeIlkStateBlock2[jug.IlkRho].(string),
				Duty:          anotherFakeIlkStateBlock2[jug.IlkDuty].(string),
				Pip:           anotherFakeIlkStateBlock2[spot.IlkPip].(string),
				Mat:           anotherFakeIlkStateBlock2[spot.IlkMat].(string),
				Created:       test_helpers.GetValidNullString(getFormattedTimestamp(headerTwo.Timestamp)),
				Updated:       test_helpers.GetValidNullString(getFormattedTimestamp(headerTwo.Timestamp)),
			}
			err := db.Select(&dbResult,
				`SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated FROM api.all_ilks($1)`,
				headerTwo.BlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbResult)).To(Equal(2))
			Expect(dbResult).To(ConsistOf(fakeIlkExpectedResult, anotherFakeIlkExpectedResult))
		})

		It("returns updated values as of block 3", func() {
			headerThree := createHeader(blockTwo+1, timestampTwo+1, headerRepository)

			//updating fakeIlk spot value at block 3
			fakeIlkStateBlock3 := test_helpers.GetIlkValues(3)
			spotMetadata := []utils.StorageValueMetadata{test_helpers.FakeIlkSpotMetadata}
			test_helpers.CreateVatRecords(db, headerThree, fakeIlkStateBlock3, spotMetadata, vatRepository)

			//updating all anotherFakeIlk values at block 3
			anotherFakeIlkStateBlock3 := test_helpers.GetIlkValues(4)
			test_helpers.CreateVatRecords(db, headerThree, anotherFakeIlkStateBlock3, test_helpers.AnotherFakeIlkVatMetadatas, vatRepository)
			test_helpers.CreateCatRecords(db, headerThree, anotherFakeIlkStateBlock3, test_helpers.AnotherFakeIlkCatMetadatas, catRepository)
			test_helpers.CreateJugRecords(db, headerThree, anotherFakeIlkStateBlock3, test_helpers.AnotherFakeIlkJugMetadatas, jugRepository)
			test_helpers.CreateSpotRecords(db, headerThree, anotherFakeIlkStateBlock3, test_helpers.AnotherFakeIlkSpotMetadatas, spotRepository)

			var dbResult []test_helpers.IlkState
			fakeIlkExpectedResult := test_helpers.IlkState{
				IlkIdentifier: fakeIlk.Identifier,
				Rate:          fakeIlkStateBlock1[vat.IlkRate].(string),
				Art:           fakeIlkStateBlock1[vat.IlkArt].(string),
				Spot:          fakeIlkStateBlock3[vat.IlkSpot].(string),
				Line:          fakeIlkStateBlock1[vat.IlkLine].(string),
				Dust:          fakeIlkStateBlock1[vat.IlkDust].(string),
				Chop:          fakeIlkStateBlock1[cat.IlkChop].(string),
				Lump:          fakeIlkStateBlock1[cat.IlkLump].(string),
				Flip:          fakeIlkStateBlock1[cat.IlkFlip].(string),
				Rho:           fakeIlkStateBlock1[jug.IlkRho].(string),
				Duty:          fakeIlkStateBlock1[jug.IlkDuty].(string),
				Pip:           fakeIlkStateBlock1[spot.IlkPip].(string),
				Mat:           fakeIlkStateBlock1[spot.IlkMat].(string),
				Created:       test_helpers.GetValidNullString(getFormattedTimestamp(headerOne.Timestamp)),
				Updated:       test_helpers.GetValidNullString(getFormattedTimestamp(headerThree.Timestamp)),
			}
			anotherFakeIlkExpectedResult := test_helpers.IlkState{
				IlkIdentifier: anotherFakeIlk.Identifier,
				Rate:          anotherFakeIlkStateBlock3[vat.IlkRate].(string),
				Art:           anotherFakeIlkStateBlock3[vat.IlkArt].(string),
				Spot:          anotherFakeIlkStateBlock3[vat.IlkSpot].(string),
				Line:          anotherFakeIlkStateBlock3[vat.IlkLine].(string),
				Dust:          anotherFakeIlkStateBlock3[vat.IlkDust].(string),
				Chop:          anotherFakeIlkStateBlock3[cat.IlkChop].(string),
				Lump:          anotherFakeIlkStateBlock3[cat.IlkLump].(string),
				Flip:          anotherFakeIlkStateBlock3[cat.IlkFlip].(string),
				Rho:           anotherFakeIlkStateBlock3[jug.IlkRho].(string),
				Duty:          anotherFakeIlkStateBlock3[jug.IlkDuty].(string),
				Pip:           anotherFakeIlkStateBlock3[spot.IlkPip].(string),
				Mat:           anotherFakeIlkStateBlock3[spot.IlkMat].(string),
				Created:       test_helpers.GetValidNullString(getFormattedTimestamp(headerTwo.Timestamp)),
				Updated:       test_helpers.GetValidNullString(getFormattedTimestamp(headerThree.Timestamp)),
			}
			err := db.Select(&dbResult,
				`SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated FROM api.all_ilks($1)`,
				headerThree.BlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbResult)).To(Equal(2))
			Expect(dbResult).To(ConsistOf(fakeIlkExpectedResult, anotherFakeIlkExpectedResult))
		})

		Describe("result pagination", func() {
			It("limits results to latest block numbers if max_results argument is provided", func() {
				//anotherFakeIlk was created at block 2
				anotherFakeIlkExpectedResult := test_helpers.IlkState{
					IlkIdentifier: anotherFakeIlk.Identifier,
					Rate:          anotherFakeIlkStateBlock2[vat.IlkRate].(string),
					Art:           anotherFakeIlkStateBlock2[vat.IlkArt].(string),
					Spot:          anotherFakeIlkStateBlock2[vat.IlkSpot].(string),
					Line:          anotherFakeIlkStateBlock2[vat.IlkLine].(string),
					Dust:          anotherFakeIlkStateBlock2[vat.IlkDust].(string),
					Chop:          anotherFakeIlkStateBlock2[cat.IlkChop].(string),
					Lump:          anotherFakeIlkStateBlock2[cat.IlkLump].(string),
					Flip:          anotherFakeIlkStateBlock2[cat.IlkFlip].(string),
					Rho:           anotherFakeIlkStateBlock2[jug.IlkRho].(string),
					Duty:          anotherFakeIlkStateBlock2[jug.IlkDuty].(string),
					Pip:           anotherFakeIlkStateBlock2[spot.IlkPip].(string),
					Mat:           anotherFakeIlkStateBlock2[spot.IlkMat].(string),
					Created:       test_helpers.GetValidNullString(getFormattedTimestamp(headerTwo.Timestamp)),
					Updated:       test_helpers.GetValidNullString(getFormattedTimestamp(headerTwo.Timestamp)),
				}

				maxResults := 1
				var dbResult []test_helpers.IlkState
				err := db.Select(&dbResult,
					`SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated FROM api.all_ilks($1, $2)`,
					headerTwo.BlockNumber, maxResults)
				Expect(err).NotTo(HaveOccurred())

				Expect(dbResult).To(Equal([]test_helpers.IlkState{anotherFakeIlkExpectedResult}))
			})

			It("offsets results if offset is provided", func() {
				fakeIlkExpectedResult := test_helpers.IlkState{
					IlkIdentifier: fakeIlk.Identifier,
					Rate:          fakeIlkStateBlock1[vat.IlkRate].(string),
					Art:           fakeIlkStateBlock1[vat.IlkArt].(string),
					Spot:          fakeIlkStateBlock1[vat.IlkSpot].(string),
					Line:          fakeIlkStateBlock1[vat.IlkLine].(string),
					Dust:          fakeIlkStateBlock1[vat.IlkDust].(string),
					Chop:          fakeIlkStateBlock1[cat.IlkChop].(string),
					Lump:          fakeIlkStateBlock1[cat.IlkLump].(string),
					Flip:          fakeIlkStateBlock1[cat.IlkFlip].(string),
					Rho:           fakeIlkStateBlock1[jug.IlkRho].(string),
					Duty:          fakeIlkStateBlock1[jug.IlkDuty].(string),
					Pip:           fakeIlkStateBlock1[spot.IlkPip].(string),
					Mat:           fakeIlkStateBlock1[spot.IlkMat].(string),
					Created:       test_helpers.GetValidNullString(getFormattedTimestamp(headerOne.Timestamp)),
					Updated:       test_helpers.GetValidNullString(getFormattedTimestamp(headerOne.Timestamp)),
				}

				maxResults := 1
				resultOffset := 1
				var dbResult []test_helpers.IlkState
				err := db.Select(&dbResult,
					`SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated FROM api.all_ilks($1, $2, $3)`,
					headerTwo.BlockNumber, maxResults, resultOffset)
				Expect(err).NotTo(HaveOccurred())

				Expect(dbResult).To(ConsistOf(fakeIlkExpectedResult))
			})
		})

		It("uses default value for blockHeight if not supplied", func() {
			_, err := db.Exec(`SELECT * FROM api.all_ilks()`)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	It("handles cases where some of the data is null", func() {
		headerFour := createHeader(blockOne+3, timestampOne+3, headerRepository)

		//updating fakeIlk spot value at block 3
		newIlk := test_helpers.TestIlk{Identifier: "newIlk", Hex: "0x6e6577496c6b0000000000000000000000000000000000000000000000000000"}
		newIlkStateBlock4 := test_helpers.GetIlkValues(4)
		metadata := []utils.StorageValueMetadata{{
			Name: vat.IlkRate,
			Keys: map[utils.Key]string{constants.Ilk: newIlk.Hex},
			Type: utils.Uint256,
		}}
		//only creating a vat_ilk_rate record
		test_helpers.CreateVatRecords(db, headerFour, newIlkStateBlock4, metadata, vatRepository)

		var dbResult []test_helpers.IlkState
		newIlkExpectedResult := test_helpers.IlkState{
			IlkIdentifier: newIlk.Identifier,
			Rate:          newIlkStateBlock4[vat.IlkRate].(string),
			Created:       test_helpers.GetValidNullString(getFormattedTimestamp(headerFour.Timestamp)),
			Updated:       test_helpers.GetValidNullString(getFormattedTimestamp(headerFour.Timestamp)),
		}
		err := db.Select(&dbResult,
			`SELECT ilk_identifier, rate, created, updated FROM api.all_ilks($1)`,
			headerFour.BlockNumber)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(dbResult)).To(Equal(3))
		Expect(dbResult).To(ContainElement(newIlkExpectedResult))
	})
})

func getFormattedTimestamp(timestampString string) string {
	parsed, _ := strconv.ParseInt(timestampString, 10, 64)
	return time.Unix(parsed, 0).UTC().Format(time.RFC3339)
}
