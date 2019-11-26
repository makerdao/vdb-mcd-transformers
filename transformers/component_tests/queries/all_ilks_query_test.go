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
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("All Ilks query", func() {

	var (
		db               *postgres.DB
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
		db = test_config.NewTestDB(test_config.NewTestNode())
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
		test_helpers.CreateVatRecords(headerOne, fakeIlkStateBlock1, test_helpers.FakeIlkVatMetadatas, vatRepository)
		test_helpers.CreateCatRecords(headerOne, fakeIlkStateBlock1, test_helpers.FakeIlkCatMetadatas, catRepository)
		test_helpers.CreateJugRecords(headerOne, fakeIlkStateBlock1, test_helpers.FakeIlkJugMetadatas, jugRepository)
		test_helpers.CreateSpotRecords(headerOne, fakeIlkStateBlock1, test_helpers.FakeIlkSpotMetadatas, spotRepository)
		//creating anotherFakeIlk at block 2
		test_helpers.CreateVatRecords(headerTwo, anotherFakeIlkStateBlock2, test_helpers.AnotherFakeIlkVatMetadatas, vatRepository)
		test_helpers.CreateCatRecords(headerTwo, anotherFakeIlkStateBlock2, test_helpers.AnotherFakeIlkCatMetadatas, catRepository)
		test_helpers.CreateJugRecords(headerTwo, anotherFakeIlkStateBlock2, test_helpers.AnotherFakeIlkJugMetadatas, jugRepository)
		test_helpers.CreateSpotRecords(headerTwo, anotherFakeIlkStateBlock2, test_helpers.AnotherFakeIlkSpotMetadatas, spotRepository)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Context("When the headerSync is complete", func() {
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
				Pip:           fakeIlkStateBlock1[spot.IlkPip],
				Mat:           fakeIlkStateBlock1[spot.IlkMat],
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
				Pip:           fakeIlkStateBlock1[spot.IlkPip],
				Mat:           fakeIlkStateBlock1[spot.IlkMat],
				Created:       test_helpers.GetValidNullString(getFormattedTimestamp(headerOne.Timestamp)),
				Updated:       test_helpers.GetValidNullString(getFormattedTimestamp(headerOne.Timestamp)),
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
				Pip:           anotherFakeIlkStateBlock2[spot.IlkPip],
				Mat:           anotherFakeIlkStateBlock2[spot.IlkMat],
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
			test_helpers.CreateVatRecords(headerThree, fakeIlkStateBlock3, spotMetadata, vatRepository)

			//updating all anotherFakeIlk values at block 3
			anotherFakeIlkStateBlock3 := test_helpers.GetIlkValues(4)
			test_helpers.CreateVatRecords(headerThree, anotherFakeIlkStateBlock3, test_helpers.AnotherFakeIlkVatMetadatas, vatRepository)
			test_helpers.CreateCatRecords(headerThree, anotherFakeIlkStateBlock3, test_helpers.AnotherFakeIlkCatMetadatas, catRepository)
			test_helpers.CreateJugRecords(headerThree, anotherFakeIlkStateBlock3, test_helpers.AnotherFakeIlkJugMetadatas, jugRepository)
			test_helpers.CreateSpotRecords(headerThree, anotherFakeIlkStateBlock3, test_helpers.AnotherFakeIlkSpotMetadatas, spotRepository)

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
				Pip:           fakeIlkStateBlock1[spot.IlkPip],
				Mat:           fakeIlkStateBlock1[spot.IlkMat],
				Created:       test_helpers.GetValidNullString(getFormattedTimestamp(headerOne.Timestamp)),
				Updated:       test_helpers.GetValidNullString(getFormattedTimestamp(headerThree.Timestamp)),
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
				Pip:           anotherFakeIlkStateBlock3[spot.IlkPip],
				Mat:           anotherFakeIlkStateBlock3[spot.IlkMat],
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
					Pip:           anotherFakeIlkStateBlock2[spot.IlkPip],
					Mat:           anotherFakeIlkStateBlock2[spot.IlkMat],
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
					Pip:           fakeIlkStateBlock1[spot.IlkPip],
					Mat:           fakeIlkStateBlock1[spot.IlkMat],
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
		test_helpers.CreateVatRecords(headerFour, newIlkStateBlock4, metadata, vatRepository)

		var dbResult []test_helpers.IlkState
		newIlkExpectedResult := test_helpers.IlkState{
			IlkIdentifier: newIlk.Identifier,
			Rate:          newIlkStateBlock4[vat.IlkRate],
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
