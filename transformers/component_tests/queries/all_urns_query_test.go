package queries

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	helper "github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/vat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("All Urns function", func() {
	var (
		vatRepo                vat.StorageRepository
		headerRepo             datastore.HeaderRepository
		headerOne              core.Header
		blockOne, timestampOne int
		urnOne                 string
		urnTwo                 string
		err                    error
		diffID                 int64
	)

	const allUrnsQuery = `SELECT urn_identifier, ilk_identifier, block_height, ink, art, created, updated FROM api.all_urns($1)`

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

		diffID = test_helpers.CreateFakeDiffRecord(db)
	})

	It("returns one urn", func() {
		setupData := helper.GetUrnSetupData()
		metadata := helper.GetUrnMetadata(helper.FakeIlk.Hex, urnOne)
		helper.CreateUrn(db, setupData, headerOne, metadata, vatRepo)

		var actualUrn helper.UrnState
		err = db.Get(&actualUrn, allUrnsQuery, blockOne)
		Expect(err).NotTo(HaveOccurred())

		expectedTimestamp := helper.GetExpectedTimestamp(timestampOne)
		expectedUrn := helper.UrnState{
			UrnIdentifier: urnOne,
			IlkIdentifier: helper.FakeIlk.Identifier,
			BlockHeight:   blockOne,
			Ink:           strconv.Itoa(setupData[vat.UrnInk].(int)),
			Art:           strconv.Itoa(setupData[vat.UrnArt].(int)),
			Created:       helper.GetValidNullString(expectedTimestamp),
			Updated:       helper.GetValidNullString(expectedTimestamp),
		}

		helper.AssertUrn(actualUrn, expectedUrn)
	})

	It("returns the correct data for multiple urns", func() {
		urnOneMetadata := helper.GetUrnMetadata(helper.FakeIlk.Hex, urnOne)
		urnOneSetupData := helper.GetUrnSetupData()
		helper.CreateUrn(db, urnOneSetupData, headerOne, urnOneMetadata, vatRepo)

		blockTwo := blockOne + 1
		timestampTwo := timestampOne + 1
		headerTwo := createHeader(blockTwo, timestampTwo, headerRepo)

		expectedTimestamp := time.Unix(int64(timestampOne), 0).UTC().Format(time.RFC3339)
		expectedUrnOne := helper.UrnState{
			UrnIdentifier: urnOne,
			IlkIdentifier: helper.FakeIlk.Identifier,
			BlockHeight:   blockOne,
			Ink:           strconv.Itoa(urnOneSetupData[vat.UrnInk].(int)),
			Art:           strconv.Itoa(urnOneSetupData[vat.UrnArt].(int)),
			Created:       helper.GetValidNullString(expectedTimestamp),
			Updated:       helper.GetValidNullString(expectedTimestamp),
		}

		urnTwoMetadata := helper.GetUrnMetadata(helper.AnotherFakeIlk.Hex, urnTwo)
		urnTwoSetupData := helper.GetUrnSetupData()
		helper.CreateUrn(db, urnTwoSetupData, headerTwo, urnTwoMetadata, vatRepo)

		expectedTimestampTwo := helper.GetExpectedTimestamp(timestampTwo)
		expectedUrnTwo := helper.UrnState{
			UrnIdentifier: urnTwo,
			IlkIdentifier: helper.AnotherFakeIlk.Identifier,
			BlockHeight:   blockTwo,
			Ink:           strconv.Itoa(urnTwoSetupData[vat.UrnInk].(int)),
			Art:           strconv.Itoa(urnTwoSetupData[vat.UrnArt].(int)),
			Created:       helper.GetValidNullString(expectedTimestampTwo),
			Updated:       helper.GetValidNullString(expectedTimestampTwo),
		}

		var result []helper.UrnState
		err = db.Select(&result, allUrnsQuery+` ORDER BY created`, blockTwo)
		Expect(err).NotTo(HaveOccurred())

		helper.AssertUrn(result[0], expectedUrnOne)
		helper.AssertUrn(result[1], expectedUrnTwo)
	})

	It("returns only the most recent urn for multiple snapshots of the same urn", func() {
		oldUrnMetadata := helper.GetUrnMetadata(helper.FakeIlk.Hex, urnOne)
		oldUrnSetupData := helper.GetUrnSetupData()
		helper.CreateUrn(db, oldUrnSetupData, headerOne, oldUrnMetadata, vatRepo)

		newUrnBlock := blockOne + 1
		newUrnTimestamp := timestampOne + 1
		newUrnHeader := createHeader(newUrnBlock, newUrnTimestamp, headerRepo)

		newUrnMetadata := helper.GetUrnMetadata(helper.FakeIlk.Hex, urnOne)
		newUrnSetupData := helper.GetUrnSetupData()
		helper.CreateUrn(db, newUrnSetupData, newUrnHeader, newUrnMetadata, vatRepo)

		createdTimestamp := helper.GetExpectedTimestamp(timestampOne)
		expectedTimestamp := helper.GetExpectedTimestamp(newUrnTimestamp)
		expectedUrn := helper.UrnState{
			UrnIdentifier: urnOne,
			IlkIdentifier: helper.FakeIlk.Identifier,
			BlockHeight:   newUrnBlock,
			Ink:           strconv.Itoa(newUrnSetupData[vat.UrnInk].(int)),
			Art:           strconv.Itoa(newUrnSetupData[vat.UrnArt].(int)),
			Created:       helper.GetValidNullString(createdTimestamp),
			Updated:       helper.GetValidNullString(expectedTimestamp),
		}

		var result []helper.UrnState
		err = db.Select(&result, allUrnsQuery, newUrnBlock)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(result)).To(Equal(1))
		helper.AssertUrn(result[0], expectedUrn)
	})

	It("returns available data if urn has ink but no art", func() {
		fakeInk := rand.Int()
		urnInkMetadata := types.GetValueMetadata(vat.UrnInk, map[types.Key]string{constants.Ilk: helper.FakeIlk.Hex, constants.Guy: urnOne}, types.Uint256)
		insertInkErr := vatRepo.Create(diffID, headerOne.Id, urnInkMetadata, strconv.Itoa(fakeInk))
		Expect(insertInkErr).NotTo(HaveOccurred())

		var result []helper.UrnState
		err = db.Select(&result, allUrnsQuery+` ORDER BY created`, headerOne.BlockNumber)
		Expect(err).NotTo(HaveOccurred())

		expectedTimestamp := helper.GetExpectedTimestamp(timestampOne)
		expectedUrn := helper.UrnState{
			UrnIdentifier: urnOne,
			IlkIdentifier: helper.FakeIlk.Identifier,
			BlockHeight:   blockOne,
			Ink:           strconv.Itoa(fakeInk),
			Art:           "0",
			Created:       helper.GetValidNullString(expectedTimestamp),
			Updated:       helper.GetValidNullString(expectedTimestamp),
		}

		Expect(len(result)).To(Equal(1))
		helper.AssertUrn(result[0], expectedUrn)
	})

	Describe("result pagination", func() {
		var (
			urnOneSetupData, urnTwoSetupData map[string]interface{}
			timestampTwo                     int
			blockTwo                         int
		)

		BeforeEach(func() {
			urnOneMetadata := helper.GetUrnMetadata(helper.FakeIlk.Hex, urnOne)
			urnOneSetupData = helper.GetUrnSetupData()
			helper.CreateUrn(db, urnOneSetupData, headerOne, urnOneMetadata, vatRepo)

			// New block
			blockTwo = blockOne + 1
			timestampTwo = timestampOne + 1
			headerTwo := createHeader(blockTwo, timestampTwo, headerRepo)

			urnTwoMetadata := helper.GetUrnMetadata(helper.AnotherFakeIlk.Hex, urnTwo)
			urnTwoSetupData = helper.GetUrnSetupData()
			helper.CreateUrn(db, urnTwoSetupData, headerTwo, urnTwoMetadata, vatRepo)
		})

		It("limits results if max_results argument is provided", func() {
			maxResults := 1
			var result []helper.UrnState

			err = db.Select(&result, `SELECT urn_identifier, ilk_identifier, ink, art, created, updated
			FROM api.all_urns($1, $2)`, blockTwo, maxResults)

			Expect(err).NotTo(HaveOccurred())

		})

		It("offsets results if offset is provided", func() {
			maxResults := 2 // We'll only get 1 because of the offset of 1, and a total of 2
			resultOffset := 1

			var result []helper.UrnState
			err = db.Select(&result, `SELECT urn_identifier, ilk_identifier, ink, art, created, updated
			FROM api.all_urns($1, $2, $3)`, blockTwo, maxResults, resultOffset)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(result)).To(Equal(1))
		})
	})

	It("allows blockHeight argument to be omitted", func() {
		_, err := db.Exec(`SELECT * FROM api.all_urns()`)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("it includes diffs only up to given block height", func() {
		var (
			actualUrn    helper.UrnState
			setupDataOne map[string]interface{}
			metadata     helper.UrnMetadata
		)

		BeforeEach(func() {
			setupDataOne = helper.GetUrnSetupData()
			metadata = helper.GetUrnMetadata(helper.FakeIlk.Hex, urnOne)
			helper.CreateUrn(db, setupDataOne, headerOne, metadata, vatRepo)
		})

		It("gets urn state as of block one", func() {
			err = db.Get(&actualUrn, allUrnsQuery, blockOne)
			Expect(err).NotTo(HaveOccurred())

			expectedTimestamp := helper.GetExpectedTimestamp(timestampOne)
			expectedUrn := helper.UrnState{
				UrnIdentifier: urnOne,
				IlkIdentifier: helper.FakeIlk.Identifier,
				BlockHeight:   blockOne,
				Ink:           strconv.Itoa(setupDataOne[vat.UrnInk].(int)),
				Art:           strconv.Itoa(setupDataOne[vat.UrnArt].(int)),
				Created:       helper.GetValidNullString(expectedTimestamp),
				Updated:       helper.GetValidNullString(expectedTimestamp),
			}

			helper.AssertUrn(actualUrn, expectedUrn)
		})

		It("gets urn state with updated values", func() {
			blockTwo := blockOne + 1
			timestampTwo := timestampOne + 1
			hashTwo := test_data.RandomString(5)
			updatedInk := rand.Int()
			fakeHeaderTwo := fakes.GetFakeHeader(int64(blockTwo))
			fakeHeaderTwo.Timestamp = strconv.Itoa(timestampTwo)
			fakeHeaderTwo.Hash = hashTwo

			fakeHeaderTwoID, err := headerRepo.CreateOrUpdateHeader(fakeHeaderTwo)
			Expect(err).NotTo(HaveOccurred())

			diffIdTwo := test_helpers.CreateFakeDiffRecordWithHeader(db, fakeHeaderTwo)
			err = vatRepo.Create(diffIdTwo, fakeHeaderTwoID, metadata.UrnInk, strconv.Itoa(updatedInk))
			Expect(err).NotTo(HaveOccurred())

			expectedTimestampOne := helper.GetExpectedTimestamp(timestampOne)
			expectedTimestampTwo := helper.GetExpectedTimestamp(timestampTwo)
			expectedUrn := helper.UrnState{
				UrnIdentifier: urnOne,
				IlkIdentifier: helper.FakeIlk.Identifier,
				BlockHeight:   blockTwo,
				Ink:           strconv.Itoa(updatedInk),
				Art:           strconv.Itoa(setupDataOne[vat.UrnArt].(int)), // Not changed
				Created:       helper.GetValidNullString(expectedTimestampOne),
				Updated:       helper.GetValidNullString(expectedTimestampTwo),
			}

			err = db.Get(&actualUrn, allUrnsQuery, blockTwo)
			Expect(err).NotTo(HaveOccurred())

			helper.AssertUrn(actualUrn, expectedUrn)
		})
	})
})
