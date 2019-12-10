package queries

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	helper "github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/vat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/utils"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Urn view", func() {
	var (
		vatRepo                vat.VatStorageRepository
		headerRepo             repositories.HeaderRepository
		headerOne              core.Header
		blockOne, timestampOne int
		urnOne                 string
		urnTwo                 string
		err                    error
	)

	const allUrnsQuery = `SELECT urn_identifier, ilk_identifier, block_height, ink, art, created, updated FROM api.all_urns($1)`

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		vatRepo = vat.VatStorageRepository{}
		vatRepo.SetDB(db)
		headerRepo = repositories.NewHeaderRepository(db)

		urnOne = test_data.RandomString(5)
		urnTwo = test_data.RandomString(5)

		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepo)
	})

	It("gets an urn", func() {
		setupData := helper.GetUrnSetupData()
		metadata := helper.GetUrnMetadata(helper.FakeIlk.Hex, urnOne)
		helper.CreateUrn(setupData, headerOne.Id, metadata, vatRepo)

		var actualUrn helper.UrnState
		err = db.Get(&actualUrn, allUrnsQuery, blockOne)
		Expect(err).NotTo(HaveOccurred())

		expectedTimestamp := helper.GetExpectedTimestamp(timestampOne)
		expectedUrn := helper.UrnState{
			UrnIdentifier: urnOne,
			IlkIdentifier: helper.FakeIlk.Identifier,
			BlockHeight:   blockOne,
			Ink:           strconv.Itoa(setupData[vat.UrnInk]),
			Art:           strconv.Itoa(setupData[vat.UrnArt]),
			Created:       helper.GetValidNullString(expectedTimestamp),
			Updated:       helper.GetValidNullString(expectedTimestamp),
		}

		helper.AssertUrn(actualUrn, expectedUrn)
	})

	It("returns the correct data for multiple urns", func() {
		urnOneMetadata := helper.GetUrnMetadata(helper.FakeIlk.Hex, urnOne)
		urnOneSetupData := helper.GetUrnSetupData()
		helper.CreateUrn(urnOneSetupData, headerOne.Id, urnOneMetadata, vatRepo)

		blockTwo := blockOne + 1
		timestampTwo := timestampOne + 1
		headerTwo := createHeader(blockTwo, timestampTwo, headerRepo)

		expectedTimestamp := time.Unix(int64(timestampOne), 0).UTC().Format(time.RFC3339)
		expectedUrnOne := helper.UrnState{
			UrnIdentifier: urnOne,
			IlkIdentifier: helper.FakeIlk.Identifier,
			BlockHeight:   blockTwo,
			Ink:           strconv.Itoa(urnOneSetupData[vat.UrnInk]),
			Art:           strconv.Itoa(urnOneSetupData[vat.UrnArt]),
			Created:       helper.GetValidNullString(expectedTimestamp),
			Updated:       helper.GetValidNullString(expectedTimestamp),
		}

		urnTwoMetadata := helper.GetUrnMetadata(helper.AnotherFakeIlk.Hex, urnTwo)
		urnTwoSetupData := helper.GetUrnSetupData()
		helper.CreateUrn(urnTwoSetupData, headerTwo.Id, urnTwoMetadata, vatRepo)

		expectedTimestampTwo := helper.GetExpectedTimestamp(timestampTwo)
		expectedUrnTwo := helper.UrnState{
			UrnIdentifier: urnTwo,
			IlkIdentifier: helper.AnotherFakeIlk.Identifier,
			BlockHeight:   blockTwo,
			Ink:           strconv.Itoa(urnTwoSetupData[vat.UrnInk]),
			Art:           strconv.Itoa(urnTwoSetupData[vat.UrnArt]),
			Created:       helper.GetValidNullString(expectedTimestampTwo),
			Updated:       helper.GetValidNullString(expectedTimestampTwo),
		}

		var result []helper.UrnState
		err = db.Select(&result, allUrnsQuery+` ORDER BY created`, blockTwo)
		Expect(err).NotTo(HaveOccurred())

		helper.AssertUrn(result[0], expectedUrnOne)
		helper.AssertUrn(result[1], expectedUrnTwo)
	})

	It("returns available data if urn has ink but no art", func() {
		fakeInk := rand.Int()
		urnInkMetadata := utils.GetStorageValueMetadata(vat.UrnInk, map[utils.Key]string{constants.Ilk: helper.FakeIlk.Hex, constants.Guy: urnOne}, utils.Uint256)
		insertInkErr := vatRepo.Create(headerOne.Id, urnInkMetadata, strconv.Itoa(fakeInk))
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
			urnOneSetupData, urnTwoSetupData map[string]int
			timestampTwo                     int
			blockTwo                         int
		)

		BeforeEach(func() {
			urnOneMetadata := helper.GetUrnMetadata(helper.FakeIlk.Hex, urnOne)
			urnOneSetupData = helper.GetUrnSetupData()
			helper.CreateUrn(urnOneSetupData, headerOne.Id, urnOneMetadata, vatRepo)

			// New block
			blockTwo = blockOne + 1
			timestampTwo = timestampOne + 1
			headerTwo := createHeader(blockTwo, timestampTwo, headerRepo)

			urnTwoMetadata := helper.GetUrnMetadata(helper.AnotherFakeIlk.Hex, urnTwo)
			urnTwoSetupData = helper.GetUrnSetupData()
			helper.CreateUrn(urnTwoSetupData, headerTwo.Id, urnTwoMetadata, vatRepo)
		})

		It("limits results if max_results argument is provided", func() {
			expectedTimestamp := helper.GetExpectedTimestamp(timestampTwo)
			expectedUrn := helper.UrnState{
				UrnIdentifier: urnTwo,
				IlkIdentifier: helper.AnotherFakeIlk.Identifier,
				Ink:           strconv.Itoa(urnTwoSetupData[vat.UrnInk]),
				Art:           strconv.Itoa(urnTwoSetupData[vat.UrnArt]),
				Created:       helper.GetValidNullString(expectedTimestamp),
				Updated:       helper.GetValidNullString(expectedTimestamp),
			}

			maxResults := 1
			var result []helper.UrnState
			err = db.Select(&result, `SELECT urn_identifier, ilk_identifier, ink, art, created, updated
			FROM api.all_urns($1, $2)`, blockTwo, maxResults)
			Expect(err).NotTo(HaveOccurred())

			Expect(len(result)).To(Equal(maxResults))
			helper.AssertUrn(result[0], expectedUrn)
		})

		It("offsets results if offset is provided", func() {
			expectedTimestamp := helper.GetExpectedTimestamp(timestampOne)
			expectedUrn := helper.UrnState{
				UrnIdentifier: urnOne,
				IlkIdentifier: helper.FakeIlk.Identifier,
				Ink:           strconv.Itoa(urnOneSetupData[vat.UrnInk]),
				Art:           strconv.Itoa(urnOneSetupData[vat.UrnArt]),
				Created:       helper.GetValidNullString(expectedTimestamp),
				Updated:       helper.GetValidNullString(expectedTimestamp),
			}

			maxResults := 1
			resultOffset := 1
			var result []helper.UrnState
			err = db.Select(&result, `SELECT urn_identifier, ilk_identifier, ink, art, created, updated
			FROM api.all_urns($1, $2, $3)`, blockTwo, maxResults, resultOffset)
			Expect(err).NotTo(HaveOccurred())

			Expect(len(result)).To(Equal(maxResults))
			helper.AssertUrn(result[0], expectedUrn)
		})
	})

	It("allows blockHeight argument to be omitted", func() {
		_, err := db.Exec(`SELECT * FROM api.all_urns()`)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("it includes diffs only up to given block height", func() {
		var (
			actualUrn    helper.UrnState
			setupDataOne map[string]int
			metadata     helper.UrnMetadata
		)

		BeforeEach(func() {
			setupDataOne = helper.GetUrnSetupData()
			metadata = helper.GetUrnMetadata(helper.FakeIlk.Hex, urnOne)
			helper.CreateUrn(setupDataOne, headerOne.Id, metadata, vatRepo)
		})

		It("gets urn state as of block one", func() {
			err = db.Get(&actualUrn, allUrnsQuery, blockOne)
			Expect(err).NotTo(HaveOccurred())

			expectedTimestamp := helper.GetExpectedTimestamp(timestampOne)
			expectedUrn := helper.UrnState{
				UrnIdentifier: urnOne,
				IlkIdentifier: helper.FakeIlk.Identifier,
				BlockHeight:   blockOne,
				Ink:           strconv.Itoa(setupDataOne[vat.UrnInk]),
				Art:           strconv.Itoa(setupDataOne[vat.UrnArt]),
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

			err = vatRepo.Create(fakeHeaderTwoID, metadata.UrnInk, strconv.Itoa(updatedInk))
			Expect(err).NotTo(HaveOccurred())

			expectedTimestampOne := helper.GetExpectedTimestamp(timestampOne)
			expectedTimestampTwo := helper.GetExpectedTimestamp(timestampTwo)
			expectedUrn := helper.UrnState{
				UrnIdentifier: urnOne,
				IlkIdentifier: helper.FakeIlk.Identifier,
				BlockHeight:   blockTwo,
				Ink:           strconv.Itoa(updatedInk),
				Art:           strconv.Itoa(setupDataOne[vat.UrnArt]), // Not changed
				Created:       helper.GetValidNullString(expectedTimestampOne),
				Updated:       helper.GetValidNullString(expectedTimestampTwo),
			}

			err = db.Get(&actualUrn, allUrnsQuery, blockTwo)
			Expect(err).NotTo(HaveOccurred())

			helper.AssertUrn(actualUrn, expectedUrn)
		})
	})
})
