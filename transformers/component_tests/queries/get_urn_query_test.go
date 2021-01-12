package queries

import (
	"math/rand"
	"strconv"

	storage_helper "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/vat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Single urn view", func() {
	var (
		vatRepo                vat.StorageRepository
		headerRepo             datastore.HeaderRepository
		urnOne, urnTwo         string
		blockOne, timestampOne int
		headerOne              core.Header
		diffID                 int64
	)

	const getUrnQuery = `SELECT urn_identifier, ilk_identifier, ink, art, created, updated FROM api.get_urn($1, $2, $3)`

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		vatRepo = vat.StorageRepository{}
		vatRepo.SetDB(db)
		headerRepo = repositories.NewHeaderRepository(db)

		urnOne = test_data.RandomString(5)
		urnTwo = test_data.RandomString(5)

		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepo)

		diffID = storage_helper.CreateFakeDiffRecord(db)
	})

	It("gets only the specified urn", func() {
		blockTwo := blockOne + 1
		headerTwo := createHeader(blockTwo, timestampOne+1, headerRepo)

		urnOneMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, urnOne)
		urnOneSetupData := test_helpers.GetUrnSetupData()
		test_helpers.CreateUrn(db, urnOneSetupData, headerOne, urnOneMetadata, vatRepo)

		expectedTimestampOne := test_helpers.GetExpectedTimestamp(timestampOne)
		expectedUrn := test_helpers.UrnState{
			UrnIdentifier: urnOne,
			IlkIdentifier: test_helpers.FakeIlk.Identifier,
			Ink:           strconv.Itoa(urnOneSetupData[vat.UrnInk].(int)),
			Art:           strconv.Itoa(urnOneSetupData[vat.UrnArt].(int)),
			Created:       test_helpers.GetValidNullString(expectedTimestampOne),
			Updated:       test_helpers.GetValidNullString(expectedTimestampOne),
		}

		urnTwoMetadata := test_helpers.GetUrnMetadata(test_helpers.AnotherFakeIlk.Hex, urnTwo)
		urnTwoSetupData := test_helpers.GetUrnSetupData()
		test_helpers.CreateUrn(db, urnTwoSetupData, headerTwo, urnTwoMetadata, vatRepo)

		var actualUrn test_helpers.UrnState
		getErr := db.Get(&actualUrn, getUrnQuery, test_helpers.FakeIlk.Identifier, urnOne, blockTwo)
		Expect(getErr).NotTo(HaveOccurred())

		test_helpers.AssertUrn(actualUrn, expectedUrn)
	})

	It("fails if no argument is supplied (STRICT)", func() {
		_, err := db.Exec(`SELECT * FROM api.get_urn()`)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("function api.get_urn() does not exist"))
	})

	It("fails if only one argument is supplied (STRICT)", func() {
		_, err := db.Exec(`SELECT * FROM api.get_urn($1::text)`, test_helpers.FakeIlk.Identifier)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("function api.get_urn(text) does not exist"))
	})

	It("allows blockHeight argument to be omitted", func() {
		_, err := db.Exec(`SELECT * FROM api.get_urn($1, $2)`, test_helpers.FakeIlk.Identifier, urnOne)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("it includes diffs only up to given block height", func() {
		var (
			actualUrn              test_helpers.UrnState
			setupDataOne           map[string]interface{}
			metadata               test_helpers.UrnMetadata
			blockTwo, timestampTwo int
			headerTwo              core.Header
		)

		BeforeEach(func() {
			setupDataOne = test_helpers.GetUrnSetupData()
			metadata = test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, urnOne)
			test_helpers.CreateUrn(db, setupDataOne, headerOne, metadata, vatRepo)

			blockTwo = blockOne + 1
			timestampTwo = timestampOne + 1
			headerTwo = createHeader(blockTwo, timestampTwo, headerRepo)
		})

		It("gets urn state as of block one", func() {
			updatedInk := rand.Int()

			createErr := vatRepo.Create(diffID, headerTwo.Id, metadata.UrnInk, strconv.Itoa(updatedInk))
			Expect(createErr).NotTo(HaveOccurred())

			expectedTimestampOne := test_helpers.GetExpectedTimestamp(timestampOne)
			expectedUrn := test_helpers.UrnState{
				UrnIdentifier: urnOne,
				IlkIdentifier: test_helpers.FakeIlk.Identifier,
				Ink:           strconv.Itoa(setupDataOne[vat.UrnInk].(int)),
				Art:           strconv.Itoa(setupDataOne[vat.UrnArt].(int)),
				Created:       test_helpers.GetValidNullString(expectedTimestampOne),
				Updated:       test_helpers.GetValidNullString(expectedTimestampOne),
			}

			getErr := db.Get(&actualUrn, getUrnQuery, test_helpers.FakeIlk.Identifier, urnOne, blockOne)
			Expect(getErr).NotTo(HaveOccurred())

			test_helpers.AssertUrn(actualUrn, expectedUrn)
		})

		It("gets urn state with updated values", func() {
			updatedInk := rand.Int()

			diffIdTwo := storage_helper.CreateFakeDiffRecordWithHeader(db, headerTwo)
			createErr := vatRepo.Create(diffIdTwo, headerTwo.Id, metadata.UrnInk, strconv.Itoa(updatedInk))
			Expect(createErr).NotTo(HaveOccurred())

			expectedTimestampOne := test_helpers.GetExpectedTimestamp(timestampOne)
			expectedTimestampTwo := test_helpers.GetExpectedTimestamp(timestampTwo)
			expectedUrn := test_helpers.UrnState{
				UrnIdentifier: urnOne,
				IlkIdentifier: test_helpers.FakeIlk.Identifier,
				Ink:           strconv.Itoa(updatedInk),
				Art:           strconv.Itoa(setupDataOne[vat.UrnArt].(int)), // Not changed
				Created:       test_helpers.GetValidNullString(expectedTimestampOne),
				Updated:       test_helpers.GetValidNullString(expectedTimestampTwo),
			}

			getErr := db.Get(&actualUrn, getUrnQuery, test_helpers.FakeIlk.Identifier, urnOne, blockTwo)
			Expect(getErr).NotTo(HaveOccurred())

			test_helpers.AssertUrn(actualUrn, expectedUrn)
		})
	})
})
