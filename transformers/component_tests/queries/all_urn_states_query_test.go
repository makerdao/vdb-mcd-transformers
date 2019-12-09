package queries

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"math/rand"
	"strconv"

	helper "github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/vat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/utils"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Urn history query", func() {
	var (
		vatRepo                vat.VatStorageRepository
		headerRepo             repositories.HeaderRepository
		fakeUrn                string
		blockOne, timestampOne int
		headerOne              core.Header
		diffID int64
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		vatRepo = vat.VatStorageRepository{}
		vatRepo.SetDB(db)

		fakeUrn = test_data.RandomString(5)

		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepo)

		diffID = test_helpers.CreateDiffRecordWithHeader(db, headerOne)
	})

	It("returns a reverse chronological history for the given ilk and urn", func() {
		urnSetupData := helper.GetUrnSetupData()
		urnMetadata := helper.GetUrnMetadata(helper.FakeIlk.Hex, fakeUrn)
		helper.CreateUrn(urnSetupData, diffID, headerOne.Id, urnMetadata, vatRepo)

		inkBlockOne := urnSetupData[vat.UrnInk]
		artBlockOne := urnSetupData[vat.UrnArt]

		expectedTimestampOne := helper.GetExpectedTimestamp(timestampOne)
		expectedUrnBlockOne := helper.UrnState{
			UrnIdentifier: fakeUrn,
			IlkIdentifier: helper.FakeIlk.Identifier,
			BlockHeight:   blockOne,
			Ink:           strconv.Itoa(inkBlockOne),
			Art:           strconv.Itoa(artBlockOne),
			Created:       helper.GetValidNullString(expectedTimestampOne),
			Updated:       helper.GetValidNullString(expectedTimestampOne),
		}

		// New block
		blockTwo := blockOne + 1
		timestampTwo := timestampOne + 1
		headerTwo := createHeader(blockTwo, timestampTwo, headerRepo)

		// Diff for header2
		diffTwoID := test_helpers.CreateDiffRecordWithHeader(db, headerTwo)

		// Relevant ink diff in block two
		inkBlockTwo := rand.Int()
		err := vatRepo.Create(diffTwoID, headerTwo.Id, urnMetadata.UrnInk, strconv.Itoa(inkBlockTwo))
		Expect(err).NotTo(HaveOccurred())

		// Irrelevant art diff in block two
		wrongUrn := test_data.RandomString(5)
		wrongArt := strconv.Itoa(rand.Int())
		wrongMetadata := utils.GetStorageValueMetadata(vat.UrnArt,
			map[utils.Key]string{constants.Ilk: helper.FakeIlk.Hex, constants.Guy: wrongUrn}, utils.Uint256)
		err = vatRepo.Create(diffID, headerOne.Id, wrongMetadata, wrongArt)
		Expect(err).NotTo(HaveOccurred())

		expectedTimestampTwo := helper.GetExpectedTimestamp(timestampTwo)
		expectedUrnBlockTwo := helper.UrnState{
			UrnIdentifier: fakeUrn,
			IlkIdentifier: helper.FakeIlk.Identifier,
			BlockHeight:   blockTwo,
			Ink:           strconv.Itoa(inkBlockTwo),
			Art:           strconv.Itoa(artBlockOne),
			Created:       helper.GetValidNullString(expectedTimestampOne),
			Updated:       helper.GetValidNullString(expectedTimestampTwo),
		}

		// New block
		blockThree := blockTwo + 1
		timestampThree := timestampTwo + 1
		headerThree := createHeader(blockThree, timestampThree, headerRepo)

		// Diff for header3
		diffThreeID := test_helpers.CreateDiffRecordWithHeader(db, headerThree)

		// Relevant art diff in block three
		artBlockThree := 0
		err = vatRepo.Create(diffThreeID, headerThree.Id, urnMetadata.UrnArt, strconv.Itoa(artBlockThree))
		Expect(err).NotTo(HaveOccurred())

		expectedTimestampThree := helper.GetExpectedTimestamp(timestampThree)
		expectedUrnBlockThree := helper.UrnState{
			UrnIdentifier: fakeUrn,
			IlkIdentifier: helper.FakeIlk.Identifier,
			BlockHeight:   blockThree,
			Ink:           strconv.Itoa(inkBlockTwo),
			Art:           strconv.Itoa(artBlockThree),
			Created:       helper.GetValidNullString(expectedTimestampOne),
			Updated:       helper.GetValidNullString(expectedTimestampThree),
		}

		var result []helper.UrnState
		dbErr := db.Select(&result,
			`SELECT * FROM api.all_urn_states($1, $2, $3)`,
			helper.FakeIlk.Identifier, fakeUrn, blockThree)
		Expect(dbErr).NotTo(HaveOccurred())

		// Reverse chronological order
		helper.AssertUrn(result[0], expectedUrnBlockThree)
		helper.AssertUrn(result[1], expectedUrnBlockTwo)
		helper.AssertUrn(result[2], expectedUrnBlockOne)
	})

	Describe("result pagination", func() {
		var (
			blockTwo, timestampTwo int
			urnSetupData           map[string]int
		)

		BeforeEach(func() {
			urnSetupData = helper.GetUrnSetupData()
			urnMetadata := helper.GetUrnMetadata(helper.FakeIlk.Hex, fakeUrn)
			helper.CreateUrn(urnSetupData, diffID, headerOne.Id, urnMetadata, vatRepo)

			// New block
			blockTwo = blockOne + 1
			timestampTwo = timestampOne + 1
			headerTwo := createHeader(blockTwo, timestampTwo, headerRepo)

			// diff in new block
			err := vatRepo.Create(diffID, headerTwo.Id, urnMetadata.UrnInk, strconv.Itoa(urnSetupData[vat.UrnInk]))
			Expect(err).NotTo(HaveOccurred())
		})

		It("limits results to most recent blocks when limit argument is provided", func() {
			expectedTimeCreated := helper.GetExpectedTimestamp(timestampOne)
			expectedTimeUpdated := helper.GetExpectedTimestamp(timestampTwo)
			expectedUrn := helper.UrnState{
				UrnIdentifier: fakeUrn,
				IlkIdentifier: helper.FakeIlk.Identifier,
				BlockHeight:   blockTwo,
				Ink:           strconv.Itoa(urnSetupData[vat.UrnInk]),
				Art:           strconv.Itoa(urnSetupData[vat.UrnArt]),
				Created:       helper.GetValidNullString(expectedTimeCreated),
				Updated:       helper.GetValidNullString(expectedTimeUpdated),
			}

			maxResults := 1
			var result []helper.UrnState
			dbErr := db.Select(&result,
				`SELECT * FROM api.all_urn_states($1, $2, $3, $4)`,
				helper.FakeIlk.Identifier, fakeUrn, blockTwo, maxResults)
			Expect(dbErr).NotTo(HaveOccurred())

			Expect(len(result)).To(Equal(maxResults))
			helper.AssertUrn(result[0], expectedUrn)
		})

		It("offsets results if offset is provided", func() {
			expectedTimeCreated := helper.GetExpectedTimestamp(timestampOne)
			expectedUrn := helper.UrnState{
				UrnIdentifier: fakeUrn,
				IlkIdentifier: helper.FakeIlk.Identifier,
				BlockHeight:   blockOne,
				Ink:           strconv.Itoa(urnSetupData[vat.UrnInk]),
				Art:           strconv.Itoa(urnSetupData[vat.UrnArt]),
				Created:       helper.GetValidNullString(expectedTimeCreated),
				Updated:       helper.GetValidNullString(expectedTimeCreated),
			}

			maxResults := 1
			resultOffset := 1
			var result []helper.UrnState
			dbErr := db.Select(&result,
				`SELECT * FROM api.all_urn_states($1, $2, $3, $4, $5)`,
				helper.FakeIlk.Identifier, fakeUrn, blockTwo, maxResults, resultOffset)
			Expect(dbErr).NotTo(HaveOccurred())

			Expect(len(result)).To(Equal(maxResults))
			helper.AssertUrn(result[0], expectedUrn)
		})
	})

	It("fails if no argument is supplied", func() {
		_, err := db.Exec(`SELECT * FROM api.all_urn_states()`)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("function api.all_urn_states() does not exist"))
	})

	It("fails if only one argument is supplied", func() {
		_, err := db.Exec(`SELECT * FROM api.all_urn_states($1::text)`, helper.FakeIlk.Identifier)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("function api.all_urn_states(text) does not exist"))
	})

	It("allows blockHeight argument to be omitted", func() {
		_, err := db.Exec(`SELECT * FROM api.all_urn_states($1, $2)`, helper.FakeIlk.Identifier, fakeUrn)
		Expect(err).NotTo(HaveOccurred())
	})
})
