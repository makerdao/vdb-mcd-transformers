package queries

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	helper "github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/vat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
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
	)

	const urnsByIlkQuery = `SELECT urn_identifier, ilk_identifier, block_height, ink, art, created, updated FROM api.get_urns_by_ilk($1, $2)`

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
	})

	It("returns multiple urns for same ilk", func() {
		urnOneMetadata := helper.GetUrnMetadata(helper.FakeIlk.Hex, urnOne)
		urnOneSetupData := helper.GetUrnSetupData()
		helper.CreateUrn(db, urnOneSetupData, headerOne, urnOneMetadata, vatRepo)

		blockTwo := blockOne + 1
		timestampTwo := timestampOne + 1
		headerTwo := createHeader(blockTwo, timestampTwo, headerRepo)
		urnTwoMetadata := helper.GetUrnMetadata(helper.FakeIlk.Hex, urnTwo)
		urnTwoSetupData := helper.GetUrnSetupData()
		helper.CreateUrn(db, urnTwoSetupData, headerTwo, urnTwoMetadata, vatRepo)

		var result []helper.UrnState
		err = db.Select(&result, urnsByIlkQuery+` ORDER BY created`, helper.FakeIlk.Identifier, blockTwo)
		Expect(err).NotTo(HaveOccurred())

		urnOneTimestamp := time.Unix(int64(timestampOne), 0).UTC().Format(time.RFC3339)
		expectedUrnOne := helper.UrnState{
			UrnIdentifier: urnOne,
			IlkIdentifier: helper.FakeIlk.Identifier,
			BlockHeight:   blockOne,
			Ink:           strconv.Itoa(urnOneSetupData[vat.UrnInk].(int)),
			Art:           strconv.Itoa(urnOneSetupData[vat.UrnArt].(int)),
			Created:       helper.GetValidNullString(urnOneTimestamp),
			Updated:       helper.GetValidNullString(urnOneTimestamp),
		}
		urnTwoTimestamp := time.Unix(int64(timestampTwo), 0).UTC().Format(time.RFC3339)
		expectedUrnTwo := helper.UrnState{
			UrnIdentifier: urnTwo,
			IlkIdentifier: helper.FakeIlk.Identifier,
			BlockHeight:   blockTwo,
			Ink:           strconv.Itoa(urnTwoSetupData[vat.UrnInk].(int)),
			Art:           strconv.Itoa(urnTwoSetupData[vat.UrnArt].(int)),
			Created:       helper.GetValidNullString(urnTwoTimestamp),
			Updated:       helper.GetValidNullString(urnTwoTimestamp),
		}

		Expect(len(result)).To(Equal(2))
		helper.AssertUrn(result[0], expectedUrnOne)
		helper.AssertUrn(result[1], expectedUrnTwo)
	})

	It("does not return urns for other ilks", func() {
		urnOneMetadata := helper.GetUrnMetadata(helper.FakeIlk.Hex, urnOne)
		urnOneSetupData := helper.GetUrnSetupData()
		helper.CreateUrn(db, urnOneSetupData, headerOne, urnOneMetadata, vatRepo)

		urnTwoMetadata := helper.GetUrnMetadata(helper.AnotherFakeIlk.Hex, urnTwo)
		urnTwoSetupData := helper.GetUrnSetupData()
		helper.CreateUrn(db, urnTwoSetupData, headerOne, urnTwoMetadata, vatRepo)

		var result []helper.UrnState
		err = db.Select(&result, urnsByIlkQuery+` ORDER BY created`, helper.FakeIlk.Identifier, blockOne)
		Expect(err).NotTo(HaveOccurred())

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

		Expect(len(result)).To(Equal(1))
		helper.AssertUrn(result[0], expectedUrnOne)
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
		updatedTimestamp := helper.GetExpectedTimestamp(newUrnTimestamp)
		expectedUrn := helper.UrnState{
			UrnIdentifier: urnOne,
			IlkIdentifier: helper.FakeIlk.Identifier,
			BlockHeight:   newUrnBlock,
			Ink:           strconv.Itoa(newUrnSetupData[vat.UrnInk].(int)),
			Art:           strconv.Itoa(newUrnSetupData[vat.UrnArt].(int)),
			Created:       helper.GetValidNullString(createdTimestamp),
			Updated:       helper.GetValidNullString(updatedTimestamp),
		}

		var result []helper.UrnState
		err = db.Select(&result, urnsByIlkQuery, helper.FakeIlk.Identifier, newUrnBlock)
		Expect(err).NotTo(HaveOccurred())

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

			urnTwoMetadata := helper.GetUrnMetadata(helper.FakeIlk.Hex, urnTwo)
			urnTwoSetupData = helper.GetUrnSetupData()
			helper.CreateUrn(db, urnTwoSetupData, headerTwo, urnTwoMetadata, vatRepo)
		})

		It("limits results if max_results argument is provided", func() {
			maxResults := 1

			var result []helper.UrnState
			err = db.Select(&result, `SELECT urn_identifier, ilk_identifier, ink, art, created, updated
			FROM api.get_urns_by_ilk($1, $2, $3, $4)`, helper.FakeIlk.Identifier, blockTwo, blockOne, maxResults)
			Expect(err).NotTo(HaveOccurred())

			Expect(len(result)).To(Equal(maxResults))
			Expect(result[0].UrnIdentifier).To(Equal(urnTwo))
		})

		It("offsets results if offset is provided", func() {
			maxResults := 2 // We'll only get 1 because of the offset of 1, and a total of 2
			resultOffset := 1

			var result []helper.UrnState
			err = db.Select(&result, `SELECT urn_identifier, ilk_identifier, ink, art, created, updated
			FROM api.get_urns_by_ilk($1, $2, $3, $4, $5)`, helper.FakeIlk.Identifier, blockTwo, blockOne, maxResults, resultOffset)
			Expect(err).NotTo(HaveOccurred())

			Expect(len(result)).To(Equal(1))
			Expect(result[0].UrnIdentifier).To(Equal(urnOne))
		})
	})

	It("allows blockHeight argument to be omitted", func() {
		_, err := db.Exec(`SELECT * FROM api.get_urns_by_ilk($1)`, helper.FakeIlk.Identifier)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("selecting diffs by block height", func() {
		var (
			actualUrn                                 helper.UrnState
			blockTwo, timestampTwo, blockTwoInk       int
			blockThree, timestampThree, blockThreeInk int
			setupDataOne                              map[string]interface{}
			metadata                                  helper.UrnMetadata
		)

		BeforeEach(func() {
			setupDataOne = helper.GetUrnSetupData()
			metadata = helper.GetUrnMetadata(helper.FakeIlk.Hex, urnOne)
			helper.CreateUrn(db, setupDataOne, headerOne, metadata, vatRepo)

			blockTwo = blockOne + 1
			fakeHeaderTwo := fakes.GetFakeHeader(int64(blockTwo))
			timestampTwo = timestampOne + 1
			fakeHeaderTwo.Timestamp = strconv.Itoa(timestampTwo)
			hashTwo := test_data.RandomString(5)
			fakeHeaderTwo.Hash = hashTwo
			fakeHeaderTwoID, err := headerRepo.CreateOrUpdateHeader(fakeHeaderTwo)
			Expect(err).NotTo(HaveOccurred())

			diffIdTwo := test_helpers.CreateFakeDiffRecordWithHeader(db, fakeHeaderTwo)
			blockTwoInk = rand.Int()
			err = vatRepo.Create(diffIdTwo, fakeHeaderTwoID, metadata.UrnInk, strconv.Itoa(blockTwoInk))
			Expect(err).NotTo(HaveOccurred())

			blockThree = blockTwo + 1
			fakeHeaderThree := fakes.GetFakeHeader(int64(blockThree))
			timestampThree = timestampTwo + 1
			fakeHeaderThree.Timestamp = strconv.Itoa(timestampThree)
			hashThree := test_data.RandomString(5)
			fakeHeaderThree.Hash = hashThree
			fakeHeaderThreeID, err := headerRepo.CreateOrUpdateHeader(fakeHeaderThree)
			Expect(err).NotTo(HaveOccurred())

			diffIdThree := test_helpers.CreateFakeDiffRecordWithHeader(db, fakeHeaderThree)
			blockThreeInk = rand.Int()
			err = vatRepo.Create(diffIdThree, fakeHeaderThreeID, metadata.UrnInk, strconv.Itoa(blockThreeInk))
			Expect(err).NotTo(HaveOccurred())
		})

		It("gets past urn state as of block one", func() {
			err = db.Get(&actualUrn, urnsByIlkQuery, helper.FakeIlk.Identifier, blockOne)
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

		It("gets latest urn state with updated values as of block three", func() {
			createdTimestamp := helper.GetExpectedTimestamp(timestampOne)
			updatedTimestampThree := helper.GetExpectedTimestamp(timestampThree)
			expectedUrn := helper.UrnState{
				UrnIdentifier: urnOne,
				IlkIdentifier: helper.FakeIlk.Identifier,
				BlockHeight:   blockThree,
				Ink:           strconv.Itoa(blockThreeInk),
				Art:           strconv.Itoa(setupDataOne[vat.UrnArt].(int)), // Not changed
				Created:       helper.GetValidNullString(createdTimestamp),
				Updated:       helper.GetValidNullString(updatedTimestampThree),
			}

			err = db.Get(&actualUrn, urnsByIlkQuery, helper.FakeIlk.Identifier, blockThree)
			Expect(err).NotTo(HaveOccurred())

			helper.AssertUrn(actualUrn, expectedUrn)
		})

		It("gets urn state as of block two by filtering by both min and max block height", func() {
			createdTimestamp := helper.GetExpectedTimestamp(timestampOne)
			updatedTimestampTwo := helper.GetExpectedTimestamp(timestampTwo)
			expectedUrn := helper.UrnState{
				UrnIdentifier: urnOne,
				IlkIdentifier: helper.FakeIlk.Identifier,
				BlockHeight:   blockTwo,
				Ink:           strconv.Itoa(blockTwoInk),
				Art:           strconv.Itoa(setupDataOne[vat.UrnArt].(int)), // Not changed
				Created:       helper.GetValidNullString(createdTimestamp),
				Updated:       helper.GetValidNullString(updatedTimestampTwo),
			}

			urnsByIlkMinHeightQuery := `SELECT urn_identifier, ilk_identifier, block_height, ink, art, created, updated FROM api.get_urns_by_ilk($1, $2, $3)`
			err = db.Get(&actualUrn, urnsByIlkMinHeightQuery, helper.FakeIlk.Identifier, blockTwo, blockOne)
			Expect(err).NotTo(HaveOccurred())

			helper.AssertUrn(actualUrn, expectedUrn)
		})
	})
})
