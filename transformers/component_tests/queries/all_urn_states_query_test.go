package queries

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	helper "github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vat"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
	"strconv"
)

var _ = Describe("Urn history query", func() {
	var (
		db         *postgres.DB
		vatRepo    vat.VatStorageRepository
		headerRepo repositories.HeaderRepository
		fakeUrn    string
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		headerRepo = repositories.NewHeaderRepository(db)
		vatRepo = vat.VatStorageRepository{}
		vatRepo.SetDB(db)

		fakeUrn = test_data.RandomString(5)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	It("returns a reverse chronological history for the given ilk and urn", func() {
		blockOne := rand.Int()
		timestampOne := int(rand.Int31())
		urnSetupData := helper.GetUrnSetupData(blockOne, timestampOne)
		urnMetadata := helper.GetUrnMetadata(helper.FakeIlk.Hex, fakeUrn)
		helper.CreateUrn(urnSetupData, urnMetadata, vatRepo, headerRepo)

		inkBlockOne := urnSetupData.Ink
		artBlockOne := urnSetupData.Art

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
		createFakeHeader(blockTwo, timestampTwo, headerRepo)

		// Relevant ink diff in block two
		inkBlockTwo := rand.Int()
		err := vatRepo.Create(blockTwo, fakes.FakeHash.String(), urnMetadata.UrnInk, strconv.Itoa(inkBlockTwo))
		Expect(err).NotTo(HaveOccurred())

		// Irrelevant art diff in block two
		wrongUrn := test_data.RandomString(5)
		wrongArt := strconv.Itoa(rand.Int())
		wrongMetadata := utils.GetStorageValueMetadata(vat.UrnArt,
			map[utils.Key]string{constants.Ilk: helper.FakeIlk.Hex, constants.Guy: wrongUrn}, utils.Uint256)
		err = vatRepo.Create(blockOne, fakes.FakeHash.String(), wrongMetadata, wrongArt)
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
		createFakeHeader(blockThree, timestampThree, headerRepo)

		// Relevant art diff in block three
		artBlockThree := 0
		err = vatRepo.Create(blockThree, fakes.FakeHash.String(), urnMetadata.UrnArt, strconv.Itoa(artBlockThree))
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
			urnCreatedTimestamp, urnUpdatedTimestamp int
			urnCreatedBlock, urnUpdatedBlock         int
			urnSetupData                             helper.UrnSetupData
		)

		BeforeEach(func() {
			urnCreatedBlock = rand.Int()
			urnCreatedTimestamp = int(rand.Int31())
			urnSetupData = helper.GetUrnSetupData(urnCreatedBlock, urnCreatedTimestamp)
			urnMetadata := helper.GetUrnMetadata(helper.FakeIlk.Hex, fakeUrn)
			helper.CreateUrn(urnSetupData, urnMetadata, vatRepo, headerRepo)

			// New block
			urnUpdatedBlock = urnCreatedBlock + 1
			urnUpdatedTimestamp = urnCreatedTimestamp + 1
			createFakeHeader(urnUpdatedBlock, urnUpdatedTimestamp, headerRepo)

			// diff in new block
			err := vatRepo.Create(urnUpdatedBlock, fakes.FakeHash.String(), urnMetadata.UrnInk, strconv.Itoa(urnSetupData.Ink))
			Expect(err).NotTo(HaveOccurred())
		})

		It("limits results to most recent blocks when limit argument is provided", func() {
			expectedTimeCreated := helper.GetExpectedTimestamp(urnCreatedTimestamp)
			expectedTimeUpdated := helper.GetExpectedTimestamp(urnUpdatedTimestamp)
			expectedUrn := helper.UrnState{
				UrnIdentifier: fakeUrn,
				IlkIdentifier: helper.FakeIlk.Identifier,
				BlockHeight:   urnUpdatedBlock,
				Ink:           strconv.Itoa(urnSetupData.Ink),
				Art:           strconv.Itoa(urnSetupData.Art),
				Created:       helper.GetValidNullString(expectedTimeCreated),
				Updated:       helper.GetValidNullString(expectedTimeUpdated),
			}

			maxResults := 1
			var result []helper.UrnState
			dbErr := db.Select(&result,
				`SELECT * FROM api.all_urn_states($1, $2, $3, $4)`,
				helper.FakeIlk.Identifier, fakeUrn, urnUpdatedBlock, maxResults)
			Expect(dbErr).NotTo(HaveOccurred())

			Expect(len(result)).To(Equal(maxResults))
			helper.AssertUrn(result[0], expectedUrn)
		})

		It("offsets results if offset is provided", func() {
			expectedTimeCreated := helper.GetExpectedTimestamp(urnCreatedTimestamp)
			expectedUrn := helper.UrnState{
				UrnIdentifier: fakeUrn,
				IlkIdentifier: helper.FakeIlk.Identifier,
				BlockHeight:   urnCreatedBlock,
				Ink:           strconv.Itoa(urnSetupData.Ink),
				Art:           strconv.Itoa(urnSetupData.Art),
				Created:       helper.GetValidNullString(expectedTimeCreated),
				Updated:       helper.GetValidNullString(expectedTimeCreated),
			}

			maxResults := 1
			resultOffset := 1
			var result []helper.UrnState
			dbErr := db.Select(&result,
				`SELECT * FROM api.all_urn_states($1, $2, $3, $4, $5)`,
				helper.FakeIlk.Identifier, fakeUrn, urnUpdatedBlock, maxResults, resultOffset)
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

func createFakeHeader(blockNumber, timestamp int, headerRepo repositories.HeaderRepository) {
	fakeHeaderOne := fakes.GetFakeHeader(int64(blockNumber))
	fakeHeaderOne.Timestamp = strconv.Itoa(timestamp)

	_, headerErr := headerRepo.CreateOrUpdateHeader(fakeHeaderOne)
	Expect(headerErr).NotTo(HaveOccurred())
}
