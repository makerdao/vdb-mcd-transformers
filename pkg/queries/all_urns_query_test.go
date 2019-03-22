package queries_test

import (
	"database/sql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	helper "github.com/vulcanize/mcd_transformers/pkg/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/storage/pit"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vat"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
	"strconv"
	"time"
)

var _ = Describe("Urn view", func() {
	var (
		db         *postgres.DB
		vatRepo    vat.VatStorageRepository
		pitRepo    pit.PitStorageRepository
		headerRepo repositories.HeaderRepository
		urnOne     string
		urnTwo     string
		ilkOne     string
		ilkTwo     string
		err        error
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		vatRepo = vat.VatStorageRepository{}
		vatRepo.SetDB(db)
		pitRepo = pit.PitStorageRepository{}
		pitRepo.SetDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		rand.Seed(time.Now().UnixNano())

		urnOne = test_data.RandomString(5)
		urnTwo = test_data.RandomString(5)
		ilkOne = test_data.RandomString(5)
		ilkTwo = test_data.RandomString(5)
	})

	It("gets an urn", func() {
		fakeBlockNo := rand.Int()
		fakeTimestamp := 12345

		setupData := helper.GetUrnSetupData(fakeBlockNo, fakeTimestamp)
		metadata := helper.GetUrnMetadata(ilkOne, urnOne)
		helper.CreateUrn(setupData, metadata, vatRepo, pitRepo, headerRepo)

		var actualUrn helper.UrnState
		err = db.Get(&actualUrn, `SELECT urnId, ilkId, blockHeight, ink, art, ratio, safe, created, updated
			FROM maker.get_all_urn_states_at_block($1)`, fakeBlockNo)
		Expect(err).NotTo(HaveOccurred())

		expectedRatio := helper.GetExpectedRatio(setupData.Ink, setupData.Spot, setupData.Art, setupData.Rate)

		expectedUrn := helper.UrnState{
			UrnId:       urnOne,
			IlkId:       ilkOne,
			BlockHeight: fakeBlockNo,
			Ink:         strconv.Itoa(setupData.Ink),
			Art:         strconv.Itoa(setupData.Art),
			Ratio:       sql.NullString{String: strconv.FormatFloat(expectedRatio, 'f', 8, 64), Valid: true},
			Safe:        expectedRatio >= 1,
			Created:     sql.NullString{String: "12345", Valid: true},
			Updated:     sql.NullString{String: "12345", Valid: true},
		}

		helper.AssertUrn(actualUrn, expectedUrn)
	})

	It("returns the correct data for multiple urns", func() {
		blockOne := rand.Int()
		timestampOne := rand.Int()

		urnOneMetadata := helper.GetUrnMetadata(ilkOne, urnOne)
		urnOneSetupData := helper.GetUrnSetupData(blockOne, timestampOne)
		helper.CreateUrn(urnOneSetupData, urnOneMetadata, vatRepo, pitRepo, headerRepo)
		expectedRatioOne := helper.GetExpectedRatio(urnOneSetupData.Ink, urnOneSetupData.Spot, urnOneSetupData.Art, urnOneSetupData.Rate)

		expectedUrnOne := helper.UrnState{
			UrnId:   urnOne,
			IlkId:   ilkOne,
			Ink:     strconv.Itoa(urnOneSetupData.Ink),
			Art:     strconv.Itoa(urnOneSetupData.Art),
			Ratio:   sql.NullString{String: strconv.FormatFloat(expectedRatioOne, 'f', 8, 64), Valid: true},
			Safe:    expectedRatioOne >= 1,
			Created: sql.NullString{String: strconv.Itoa(timestampOne), Valid: true},
			Updated: sql.NullString{String: strconv.Itoa(timestampOne), Valid: true},
		}

		// New block
		blockTwo := blockOne + 1
		timestampTwo := timestampOne + 1

		urnTwoMetadata := helper.GetUrnMetadata(ilkTwo, urnTwo)
		urnTwoSetupData := helper.GetUrnSetupData(blockTwo, timestampTwo)
		helper.CreateUrn(urnTwoSetupData, urnTwoMetadata, vatRepo, pitRepo, headerRepo)
		expectedRatioTwo := helper.GetExpectedRatio(urnTwoSetupData.Ink, urnTwoSetupData.Spot, urnTwoSetupData.Art, urnTwoSetupData.Rate)

		expectedUrnTwo := helper.UrnState{
			UrnId:   urnTwo,
			IlkId:   ilkTwo,
			Ink:     strconv.Itoa(urnTwoSetupData.Ink),
			Art:     strconv.Itoa(urnTwoSetupData.Art),
			Ratio:   sql.NullString{String: strconv.FormatFloat(expectedRatioTwo, 'f', 8, 64), Valid: true},
			Safe:    expectedRatioTwo >= 1,
			Created: sql.NullString{String: strconv.Itoa(timestampTwo), Valid: true},
			Updated: sql.NullString{String: strconv.Itoa(timestampTwo), Valid: true},
		}

		var result []helper.UrnState
		err = db.Select(&result, `SELECT urnId, ilkId, ink, art, ratio, safe, created, updated
			FROM maker.get_all_urn_states_at_block($1)`, blockTwo)
		Expect(err).NotTo(HaveOccurred())

		helper.AssertUrn(result[0], expectedUrnOne)
		helper.AssertUrn(result[1], expectedUrnTwo)
	})

	It("returns urn state without timestamps if corresponding headers aren't synced", func() {
		block := rand.Int()
		timestamp := rand.Int()
		metadata := helper.GetUrnMetadata(ilkOne, urnOne)
		setupData := helper.GetUrnSetupData(block, timestamp)

		helper.CreateUrn(setupData, metadata, vatRepo, pitRepo, headerRepo)
		_, err = db.Exec(`DELETE FROM headers`)
		Expect(err).NotTo(HaveOccurred())

		var result helper.UrnState
		err = db.Get(&result, `SELECT urnId, ilkId, ink, art, ratio, safe, created, updated
			FROM maker.get_all_urn_states_at_block($1)`, block)

		Expect(err).NotTo(HaveOccurred())
		Expect(result.Created.String).To(BeEmpty())
		Expect(result.Updated.String).To(BeEmpty())
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
			timestampOne = rand.Int()
			setupDataOne = helper.GetUrnSetupData(blockOne, timestampOne)
			metadata = helper.GetUrnMetadata(ilkOne, urnOne)
			helper.CreateUrn(setupDataOne, metadata, vatRepo, pitRepo, headerRepo)
		})

		It("gets urn state as of block one", func() {
			err = db.Get(&actualUrn, `SELECT urnId, ilkId, ink, art, ratio, safe, created, updated
				FROM maker.get_all_urn_states_at_block($1)`, blockOne)
			Expect(err).NotTo(HaveOccurred())

			expectedRatio := helper.GetExpectedRatio(setupDataOne.Ink, setupDataOne.Spot, setupDataOne.Art, setupDataOne.Rate)

			expectedUrn := helper.UrnState{
				UrnId:   urnOne,
				IlkId:   ilkOne,
				Ink:     strconv.Itoa(setupDataOne.Ink),
				Art:     strconv.Itoa(setupDataOne.Art),
				Ratio:   sql.NullString{String: strconv.FormatFloat(expectedRatio, 'f', 8, 64), Valid: true},
				Safe:    expectedRatio >= 1,
				Created: sql.NullString{String: strconv.Itoa(timestampOne), Valid: true},
				Updated: sql.NullString{String: strconv.Itoa(timestampOne), Valid: true},
			}

			helper.AssertUrn(actualUrn, expectedUrn)
		})

		It("gets urn state with updated values", func() {
			blockTwo := blockOne + 1
			timestampTwo := timestampOne + 1
			hashTwo := test_data.RandomString(5)
			updatedInk := rand.Int()

			err = vatRepo.Create(blockTwo, hashTwo, metadata.UrnInk, strconv.Itoa(updatedInk))
			Expect(err).NotTo(HaveOccurred())

			expectedRatio := helper.GetExpectedRatio(updatedInk, setupDataOne.Spot, setupDataOne.Art, setupDataOne.Rate)

			expectedUrn := helper.UrnState{
				UrnId:   urnOne,
				IlkId:   ilkOne,
				Ink:     strconv.Itoa(updatedInk),
				Art:     strconv.Itoa(setupDataOne.Art), // Not changed
				Ratio:   sql.NullString{String: strconv.FormatFloat(expectedRatio, 'f', 8, 64), Valid: true},
				Safe:    expectedRatio >= 1,
				Created: sql.NullString{String: strconv.Itoa(timestampOne), Valid: true},
				Updated: sql.NullString{String: strconv.Itoa(timestampTwo), Valid: true},
			}

			fakeHeaderTwo := fakes.GetFakeHeader(int64(blockTwo))
			fakeHeaderTwo.Timestamp = strconv.Itoa(timestampTwo)
			fakeHeaderTwo.Hash = hashTwo

			_, err = headerRepo.CreateOrUpdateHeader(fakeHeaderTwo)
			Expect(err).NotTo(HaveOccurred())

			err = db.Get(&actualUrn, `SELECT urnId, ilkId, ink, art, ratio, safe, created, updated
				FROM maker.get_all_urn_states_at_block($1)`, blockTwo)
			Expect(err).NotTo(HaveOccurred())

			helper.AssertUrn(actualUrn, expectedUrn)
		})
	})

	It("returns null ratio and urn being safe if there is no debt", func() {
		block := rand.Int()
		setupData := helper.GetUrnSetupData(block, 1)
		setupData.Art = 0
		metadata := helper.GetUrnMetadata(ilkOne, urnOne)
		helper.CreateUrn(setupData, metadata, vatRepo, pitRepo, headerRepo)

		fakeHeader := fakes.GetFakeHeader(int64(block))
		_, err = headerRepo.CreateOrUpdateHeader(fakeHeader)
		Expect(err).NotTo(HaveOccurred())

		var result helper.UrnState
		err = db.Get(&result, `SELECT urnId, ilkId, ink, art, ratio, safe, created, updated
			FROM maker.get_all_urn_states_at_block($1)`, block)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Ratio.String).To(BeEmpty())
		Expect(result.Safe).To(BeTrue())
	})
})
