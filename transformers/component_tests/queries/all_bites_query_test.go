package queries

import (
	"fmt"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"math/rand"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/bite"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("Bites query", func() {
	var (
		db         *postgres.DB
		biteRepo   bite.BiteRepository
		headerRepo repositories.HeaderRepository
		fakeUrn    = test_data.RandomString(5)
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		biteRepo = bite.BiteRepository{}
		biteRepo.SetDB(db)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("all_bites", func() {
		It("returns bites for an ilk", func() {
			headerOne := fakes.GetFakeHeader(1)
			headerOneId, err := headerRepo.CreateOrUpdateHeader(headerOne)
			Expect(err).NotTo(HaveOccurred())
			biteLog := test_data.CreateTestLog(headerOneId, db)

			biteOne := generateBite(test_helpers.FakeIlk.Hex, fakeUrn, headerOneId, biteLog.ID)
			err = biteRepo.Create([]shared.InsertionModel{biteOne})
			Expect(err).NotTo(HaveOccurred())

			var actualBites []test_helpers.BiteEvent
			err = db.Select(&actualBites, `SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.all_bites($1)`,
				test_helpers.FakeIlk.Identifier)
			Expect(err).NotTo(HaveOccurred())

			Expect(actualBites).To(ConsistOf(
				test_helpers.BiteEvent{
					IlkIdentifier: test_helpers.FakeIlk.Identifier,
					UrnIdentifier: fakeUrn,
					Ink:           fmt.Sprint(biteOne.ColumnValues["ink"].(uint)),
					Art:           fmt.Sprint(biteOne.ColumnValues["art"].(uint)),
					Tab:           fmt.Sprint(biteOne.ColumnValues["tab"].(uint)),
				},
			))
		})

		It("returns bites from multiple blocks", func() {
			headerOne := fakes.GetFakeHeaderWithTimestamp(int64(111111111), 1)
			headerOneId, err := headerRepo.CreateOrUpdateHeader(headerOne)
			Expect(err).NotTo(HaveOccurred())
			biteBlockOneLog := test_data.CreateTestLog(headerOneId, db)

			biteBlockOne := generateBite(test_helpers.FakeIlk.Hex, fakeUrn, headerOneId, biteBlockOneLog.ID)
			err = biteRepo.Create([]shared.InsertionModel{biteBlockOne})
			Expect(err).NotTo(HaveOccurred())

			// New block
			headerTwo := fakes.GetFakeHeaderWithTimestamp(int64(222222222), 2)
			headerTwo.Hash = "anotherHash"
			headerTwoId, err := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(err).NotTo(HaveOccurred())
			biteBlockTwoLog := test_data.CreateTestLog(headerTwoId, db)

			biteBlockTwo := generateBite(test_helpers.FakeIlk.Hex, fakeUrn, headerTwoId, biteBlockTwoLog.ID)
			err = biteRepo.Create([]shared.InsertionModel{biteBlockTwo})
			Expect(err).NotTo(HaveOccurred())

			var actualBites []test_helpers.BiteEvent
			err = db.Select(&actualBites, `SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.all_bites($1)`, test_helpers.FakeIlk.Identifier)
			Expect(err).NotTo(HaveOccurred())

			Expect(actualBites).To(ConsistOf(
				test_helpers.BiteEvent{
					IlkIdentifier: test_helpers.FakeIlk.Identifier,
					UrnIdentifier: fakeUrn,
					Ink:           fmt.Sprint(biteBlockTwo.ColumnValues["ink"].(uint)),
					Art:           fmt.Sprint(biteBlockTwo.ColumnValues["art"].(uint)),
					Tab:           fmt.Sprint(biteBlockTwo.ColumnValues["tab"].(uint)),
				},
				test_helpers.BiteEvent{
					IlkIdentifier: test_helpers.FakeIlk.Identifier,
					UrnIdentifier: fakeUrn,
					Ink:           fmt.Sprint(biteBlockOne.ColumnValues["ink"].(uint)),
					Art:           fmt.Sprint(biteBlockOne.ColumnValues["art"].(uint)),
					Tab:           fmt.Sprint(biteBlockOne.ColumnValues["tab"].(uint)),
				},
			))
		})

		It("ignores bites from irrelevant ilks", func() {
			headerOne := fakes.GetFakeHeader(1)
			headerOneId, err := headerRepo.CreateOrUpdateHeader(headerOne)
			Expect(err).NotTo(HaveOccurred())
			biteLog := test_data.CreateTestLog(headerOneId, db)
			irrelevantBiteLog := test_data.CreateTestLog(headerOneId, db)

			bite := generateBite(test_helpers.FakeIlk.Hex, fakeUrn, headerOneId, biteLog.ID)
			irrelevantBite := generateBite(test_helpers.AnotherFakeIlk.Hex, fakeUrn, headerOneId, irrelevantBiteLog.ID)

			err = biteRepo.Create([]shared.InsertionModel{bite, irrelevantBite})
			Expect(err).NotTo(HaveOccurred())

			var actualBites []test_helpers.BiteEvent
			err = db.Select(&actualBites, `SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.all_bites($1)`, test_helpers.FakeIlk.Identifier)
			Expect(err).NotTo(HaveOccurred())

			Expect(actualBites).To(ConsistOf(
				test_helpers.BiteEvent{
					IlkIdentifier: test_helpers.FakeIlk.Identifier,
					UrnIdentifier: fakeUrn,
					Ink:           fmt.Sprint(bite.ColumnValues["ink"].(uint)),
					Art:           fmt.Sprint(bite.ColumnValues["art"].(uint)),
					Tab:           fmt.Sprint(bite.ColumnValues["tab"].(uint)),
				},
			))
		})

		It("fails if no argument is supplied (STRICT)", func() {
			_, err := db.Exec(`SELECT * FROM api.all_bites()`)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("function api.all_bites() does not exist"))
		})

		Describe("result pagination", func() {
			var biteBlockOne, biteBlockTwo shared.InsertionModel

			BeforeEach(func() {
				headerOne := fakes.GetFakeHeaderWithTimestamp(int64(111111111), 1)
				headerOneId, headerOneErr := headerRepo.CreateOrUpdateHeader(headerOne)
				Expect(headerOneErr).NotTo(HaveOccurred())
				biteBlockOneLog := test_data.CreateTestLog(headerOneId, db)

				biteBlockOne = generateBite(test_helpers.FakeIlk.Hex, fakeUrn, headerOneId, biteBlockOneLog.ID)
				biteOneErr := biteRepo.Create([]shared.InsertionModel{biteBlockOne})
				Expect(biteOneErr).NotTo(HaveOccurred())

				// New block
				headerTwo := fakes.GetFakeHeaderWithTimestamp(int64(222222222), 2)
				headerTwo.Hash = "anotherHash"
				headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
				Expect(headerTwoErr).NotTo(HaveOccurred())
				biteBlockTwoLog := test_data.CreateTestLog(headerTwoId, db)

				biteBlockTwo = generateBite(test_helpers.FakeIlk.Hex, fakeUrn, headerTwoId, biteBlockTwoLog.ID)
				biteTwoErr := biteRepo.Create([]shared.InsertionModel{biteBlockTwo})
				Expect(biteTwoErr).NotTo(HaveOccurred())
			})

			It("limits results to most recent blocks if max_results argument is provided", func() {
				maxResults := 1
				var actualBites []test_helpers.BiteEvent
				err := db.Select(&actualBites, `SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.all_bites($1, $2)`,
					test_helpers.FakeIlk.Identifier, maxResults)
				Expect(err).NotTo(HaveOccurred())

				Expect(actualBites).To(ConsistOf(
					test_helpers.BiteEvent{
						IlkIdentifier: test_helpers.FakeIlk.Identifier,
						UrnIdentifier: fakeUrn,
						Ink:           fmt.Sprint(biteBlockTwo.ColumnValues["ink"].(uint)),
						Art:           fmt.Sprint(biteBlockTwo.ColumnValues["art"].(uint)),
						Tab:           fmt.Sprint(biteBlockTwo.ColumnValues["tab"].(uint)),
					},
				))
			})

			It("offsets results if offset is provided", func() {
				maxResults := 1
				resultOffset := 1
				var actualBites []test_helpers.BiteEvent
				err := db.Select(&actualBites, `SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.all_bites($1, $2, $3)`,
					test_helpers.FakeIlk.Identifier, maxResults, resultOffset)
				Expect(err).NotTo(HaveOccurred())

				Expect(actualBites).To(ConsistOf(
					test_helpers.BiteEvent{
						IlkIdentifier: test_helpers.FakeIlk.Identifier,
						UrnIdentifier: fakeUrn,
						Ink:           fmt.Sprint(biteBlockOne.ColumnValues["ink"].(uint)),
						Art:           fmt.Sprint(biteBlockOne.ColumnValues["art"].(uint)),
						Tab:           fmt.Sprint(biteBlockOne.ColumnValues["tab"].(uint)),
					},
				))
			})
		})
	})

	Describe("urn_bites", func() {
		It("returns bites for relevant ilk/urn", func() {
			headerOne := fakes.GetFakeHeader(1)
			headerOneId, err := headerRepo.CreateOrUpdateHeader(headerOne)
			Expect(err).NotTo(HaveOccurred())
			biteOneLog := test_data.CreateTestLog(headerOneId, db)

			biteOne := generateBite(test_helpers.FakeIlk.Hex, fakeUrn, headerOneId, biteOneLog.ID)
			err = biteRepo.Create([]shared.InsertionModel{biteOne})
			Expect(err).NotTo(HaveOccurred())

			var actualBites []test_helpers.BiteEvent
			err = db.Select(&actualBites, `SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.urn_bites($1, $2)`, test_helpers.FakeIlk.Identifier, fakeUrn)
			Expect(err).NotTo(HaveOccurred())

			Expect(actualBites).To(ConsistOf(
				test_helpers.BiteEvent{
					IlkIdentifier: test_helpers.FakeIlk.Identifier,
					UrnIdentifier: fakeUrn,
					Ink:           fmt.Sprint(biteOne.ColumnValues["ink"].(uint)),
					Art:           fmt.Sprint(biteOne.ColumnValues["art"].(uint)),
					Tab:           fmt.Sprint(biteOne.ColumnValues["tab"].(uint)),
				},
			))
		})

		It("returns bites from multiple blocks", func() {
			headerOne := fakes.GetFakeHeaderWithTimestamp(int64(111111111), 1)
			headerOneId, err := headerRepo.CreateOrUpdateHeader(headerOne)
			Expect(err).NotTo(HaveOccurred())
			biteOneLog := test_data.CreateTestLog(headerOneId, db)

			biteBlockOne := generateBite(test_helpers.FakeIlk.Hex, fakeUrn, headerOneId, biteOneLog.ID)
			err = biteRepo.Create([]shared.InsertionModel{biteBlockOne})
			Expect(err).NotTo(HaveOccurred())

			// New block
			headerTwo := fakes.GetFakeHeaderWithTimestamp(int64(222222222), 2)
			headerTwo.Hash = "anotherHash"
			headerTwoId, err := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(err).NotTo(HaveOccurred())
			biteBlockTwoLog := test_data.CreateTestLog(headerTwoId, db)

			biteBlockTwo := generateBite(test_helpers.FakeIlk.Hex, fakeUrn, headerTwoId, biteBlockTwoLog.ID)
			err = biteRepo.Create([]shared.InsertionModel{biteBlockTwo})
			Expect(err).NotTo(HaveOccurred())

			var actualBites []test_helpers.BiteEvent
			err = db.Select(&actualBites, `SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.urn_bites($1, $2)`, test_helpers.FakeIlk.Identifier, fakeUrn)
			Expect(err).NotTo(HaveOccurred())

			Expect(actualBites).To(ConsistOf(
				test_helpers.BiteEvent{
					IlkIdentifier: test_helpers.FakeIlk.Identifier,
					UrnIdentifier: fakeUrn,
					Ink:           fmt.Sprint(biteBlockTwo.ColumnValues["ink"].(uint)),
					Art:           fmt.Sprint(biteBlockTwo.ColumnValues["art"].(uint)),
					Tab:           fmt.Sprint(biteBlockTwo.ColumnValues["tab"].(uint)),
				},
				test_helpers.BiteEvent{
					IlkIdentifier: test_helpers.FakeIlk.Identifier,
					UrnIdentifier: fakeUrn,
					Ink:           fmt.Sprint(biteBlockOne.ColumnValues["ink"].(uint)),
					Art:           fmt.Sprint(biteBlockOne.ColumnValues["art"].(uint)),
					Tab:           fmt.Sprint(biteBlockOne.ColumnValues["tab"].(uint)),
				},
			))
		})

		Describe("result pagination", func() {
			var biteBlockOne, biteBlockTwo shared.InsertionModel

			BeforeEach(func() {
				headerOne := fakes.GetFakeHeaderWithTimestamp(int64(111111111), 1)
				headerOneId, err := headerRepo.CreateOrUpdateHeader(headerOne)
				Expect(err).NotTo(HaveOccurred())
				biteBlockOneLog := test_data.CreateTestLog(headerOneId, db)

				biteBlockOne = generateBite(test_helpers.FakeIlk.Hex, fakeUrn, headerOneId, biteBlockOneLog.ID)
				err = biteRepo.Create([]shared.InsertionModel{biteBlockOne})
				Expect(err).NotTo(HaveOccurred())

				// New block
				headerTwo := fakes.GetFakeHeaderWithTimestamp(int64(222222222), 2)
				headerTwo.Hash = "anotherHash"
				headerTwoId, err := headerRepo.CreateOrUpdateHeader(headerTwo)
				Expect(err).NotTo(HaveOccurred())
				biteBlockTwoLog := test_data.CreateTestLog(headerTwoId, db)

				biteBlockTwo = generateBite(test_helpers.FakeIlk.Hex, fakeUrn, headerTwoId, biteBlockTwoLog.ID)
				err = biteRepo.Create([]shared.InsertionModel{biteBlockTwo})
				Expect(err).NotTo(HaveOccurred())
			})

			It("limits results to latest block if max_results argument is provided", func() {
				maxResults := 1
				var actualBites []test_helpers.BiteEvent
				err := db.Select(&actualBites, `SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.urn_bites($1, $2, $3)`,
					test_helpers.FakeIlk.Identifier, fakeUrn, maxResults)
				Expect(err).NotTo(HaveOccurred())

				Expect(actualBites).To(ConsistOf(
					test_helpers.BiteEvent{IlkIdentifier: test_helpers.FakeIlk.Identifier, UrnIdentifier: fakeUrn, Ink: biteBlockTwo.Ink, Art: biteBlockTwo.Art, Tab: biteBlockTwo.Tab},
				))
			})

			It("offsets results if offset is provided", func() {
				maxResults := 1
				resultOffset := 1
				var actualBites []test_helpers.BiteEvent
				err := db.Select(&actualBites, `SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.urn_bites($1, $2, $3, $4)`,
					test_helpers.FakeIlk.Identifier, fakeUrn, maxResults, resultOffset)
				Expect(err).NotTo(HaveOccurred())

				Expect(actualBites).To(ConsistOf(
					test_helpers.BiteEvent{IlkIdentifier: test_helpers.FakeIlk.Identifier, UrnIdentifier: fakeUrn, Ink: biteBlockOne.Ink, Art: biteBlockOne.Art, Tab: biteBlockOne.Tab},
				))
			})
		})

		It("ignores bites from irrelevant urns", func() {
			headerOne := fakes.GetFakeHeader(1)
			headerOneId, err := headerRepo.CreateOrUpdateHeader(headerOne)
			Expect(err).NotTo(HaveOccurred())
			biteLog := test_data.CreateTestLog(headerOneId, db)
			irrelevantBiteLog := test_data.CreateTestLog(headerOneId, db)

			bite := generateBite(test_helpers.FakeIlk.Hex, fakeUrn, headerOneId, biteLog.ID)
			irrelevantBite := generateBite(test_helpers.FakeIlk.Hex, "irrelevantUrn", headerOneId, irrelevantBiteLog.ID)

			err = biteRepo.Create([]shared.InsertionModel{bite, irrelevantBite})
			Expect(err).NotTo(HaveOccurred())

			var actualBites []test_helpers.BiteEvent
			err = db.Select(&actualBites, `SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.urn_bites($1, $2)`, test_helpers.FakeIlk.Identifier, fakeUrn)
			Expect(err).NotTo(HaveOccurred())

			Expect(actualBites).To(ConsistOf(
				test_helpers.BiteEvent{
					IlkIdentifier: test_helpers.FakeIlk.Identifier,
					UrnIdentifier: fakeUrn,
					Ink:           fmt.Sprint(bite.ColumnValues["ink"].(uint)),
					Art:           fmt.Sprint(bite.ColumnValues["art"].(uint)),
					Tab:           fmt.Sprint(bite.ColumnValues["tab"].(uint)),
				},
			))
		})
	})
})

func generateBite(ilk, urn string, headerID, logID int64) shared.InsertionModel {
	biteEvent := test_data.BiteModel
	biteEvent.ForeignKeyValues[constants.IlkFK] = ilk
	biteEvent.ForeignKeyValues[constants.UrnFK] = urn
	biteEvent.ColumnValues["ink"] = strconv.Itoa(rand.Int())
	biteEvent.ColumnValues["art"] = strconv.Itoa(rand.Int())
	biteEvent.ColumnValues["tab"] = strconv.Itoa(rand.Int())
	biteEvent.ColumnValues["bite_identifier"] = strconv.Itoa(rand.Int())
	biteEvent.ColumnValues[constants.HeaderFK] = headerID
	biteEvent.ColumnValues[constants.LogFK] = logID
	return biteEvent
}
