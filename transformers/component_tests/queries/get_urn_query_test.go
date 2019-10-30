package queries

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	helper "github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vat"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
	"strconv"
	"time"
)

var _ = Describe("Single urn view", func() {
	var (
		db         *postgres.DB
		vatRepo    vat.VatStorageRepository
		headerRepo repositories.HeaderRepository
		urnOne     string
		urnTwo     string
		err        error
	)

	const getUrnQuery = `SELECT urn_identifier, ilk_identifier, ink, art, created, updated FROM api.get_urn($1, $2, $3)`

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		vatRepo = vat.VatStorageRepository{}
		vatRepo.SetDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		rand.Seed(time.Now().UnixNano())

		urnOne = test_data.RandomString(5)
		urnTwo = test_data.RandomString(5)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	It("gets only the specified urn", func() {
		blockOne := rand.Int()
		timestampOne := int(rand.Int31())
		blockTwo := blockOne + 1
		timestampTwo := timestampOne + 1

		urnOneMetadata := helper.GetUrnMetadata(helper.FakeIlk.Hex, urnOne)
		urnOneSetupData := helper.GetUrnSetupData(blockOne, timestampOne)
		helper.CreateUrn(urnOneSetupData, urnOneMetadata, vatRepo, headerRepo)

		expectedTimestampOne := helper.GetExpectedTimestamp(timestampOne)
		expectedUrn := helper.UrnState{
			UrnIdentifier: urnOne,
			IlkIdentifier: helper.FakeIlk.Identifier,
			Ink:           strconv.Itoa(urnOneSetupData.Ink),
			Art:           strconv.Itoa(urnOneSetupData.Art),
			Created:       helper.GetValidNullString(expectedTimestampOne),
			Updated:       helper.GetValidNullString(expectedTimestampOne),
		}

		urnTwoMetadata := helper.GetUrnMetadata(helper.AnotherFakeIlk.Hex, urnTwo)
		urnTwoSetupData := helper.GetUrnSetupData(blockTwo, timestampTwo)
		helper.CreateUrn(urnTwoSetupData, urnTwoMetadata, vatRepo, headerRepo)

		var actualUrn helper.UrnState
		err = db.Get(&actualUrn, getUrnQuery, helper.FakeIlk.Identifier, urnOne, blockTwo)
		Expect(err).NotTo(HaveOccurred())

		helper.AssertUrn(actualUrn, expectedUrn)
	})

	It("returns urn state without timestamps if corresponding headers aren't synced", func() {
		block := rand.Int()
		timestamp := int(rand.Int31())
		metadata := helper.GetUrnMetadata(helper.FakeIlk.Hex, urnOne)
		setupData := helper.GetUrnSetupData(block, timestamp)

		helper.CreateUrn(setupData, metadata, vatRepo, headerRepo)
		_, err = db.Exec(`DELETE FROM headers`)
		Expect(err).NotTo(HaveOccurred())

		var result helper.UrnState
		err = db.Get(&result, getUrnQuery, helper.FakeIlk.Identifier, urnOne, block)

		Expect(err).NotTo(HaveOccurred())
		Expect(result.Created.String).To(BeEmpty())
		Expect(result.Updated.String).To(BeEmpty())
	})

	It("fails if no argument is supplied (STRICT)", func() {
		_, err := db.Exec(`SELECT * FROM api.get_urn()`)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("function api.get_urn() does not exist"))
	})

	It("fails if only one argument is supplied (STRICT)", func() {
		_, err := db.Exec(`SELECT * FROM api.get_urn($1::text)`, helper.FakeIlk.Identifier)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("function api.get_urn(text) does not exist"))
	})

	It("allows blockHeight argument to be omitted", func() {
		_, err := db.Exec(`SELECT * FROM api.get_urn($1, $2)`, helper.FakeIlk.Identifier, urnOne)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("it includes diffs only up to given block height", func() {
		var (
			actualUrn    helper.UrnState
			blockOne     int
			timestampOne int
			setupDataOne helper.UrnSetupData
			metadata     helper.UrnMetadata
		)

		BeforeEach(func() {
			blockOne = rand.Int()
			timestampOne = int(rand.Int31())
			setupDataOne = helper.GetUrnSetupData(blockOne, timestampOne)
			metadata = helper.GetUrnMetadata(helper.FakeIlk.Hex, urnOne)
			helper.CreateUrn(setupDataOne, metadata, vatRepo, headerRepo)
		})

		It("gets urn state as of block one", func() {
			blockTwo := blockOne + 1
			hashTwo := test_data.RandomString(5)
			updatedInk := rand.Int()

			err = vatRepo.Create(blockTwo, hashTwo, metadata.UrnInk, strconv.Itoa(updatedInk))
			Expect(err).NotTo(HaveOccurred())

			expectedTimestampOne := helper.GetExpectedTimestamp(timestampOne)
			expectedUrn := helper.UrnState{
				UrnIdentifier: urnOne,
				IlkIdentifier: helper.FakeIlk.Identifier,
				Ink:           strconv.Itoa(setupDataOne.Ink),
				Art:           strconv.Itoa(setupDataOne.Art),
				Created:       helper.GetValidNullString(expectedTimestampOne),
				Updated:       helper.GetValidNullString(expectedTimestampOne),
			}

			err = db.Get(&actualUrn, getUrnQuery, helper.FakeIlk.Identifier, urnOne, blockOne)
			Expect(err).NotTo(HaveOccurred())

			helper.AssertUrn(actualUrn, expectedUrn)
		})

		It("gets urn state with updated values", func() {
			blockTwo := blockOne + 1
			timestampTwo := timestampOne + 1
			hashTwo := test_data.RandomString(5)
			updatedInk := rand.Int()

			err = vatRepo.Create(blockTwo, hashTwo, metadata.UrnInk, strconv.Itoa(updatedInk))
			Expect(err).NotTo(HaveOccurred())

			expectedTimestampOne := helper.GetExpectedTimestamp(timestampOne)
			expectedTimestampTwo := helper.GetExpectedTimestamp(timestampTwo)
			expectedUrn := helper.UrnState{
				UrnIdentifier: urnOne,
				IlkIdentifier: helper.FakeIlk.Identifier,
				Ink:           strconv.Itoa(updatedInk),
				Art:           strconv.Itoa(setupDataOne.Art), // Not changed
				Created:       helper.GetValidNullString(expectedTimestampOne),
				Updated:       helper.GetValidNullString(expectedTimestampTwo),
			}

			fakeHeaderTwo := fakes.GetFakeHeader(int64(blockTwo))
			fakeHeaderTwo.Timestamp = strconv.Itoa(timestampTwo)
			fakeHeaderTwo.Hash = hashTwo

			_, err = headerRepo.CreateOrUpdateHeader(fakeHeaderTwo)
			Expect(err).NotTo(HaveOccurred())

			err = db.Get(&actualUrn, getUrnQuery, helper.FakeIlk.Identifier, urnOne, blockTwo)
			Expect(err).NotTo(HaveOccurred())

			helper.AssertUrn(actualUrn, expectedUrn)
		})
	})
})
