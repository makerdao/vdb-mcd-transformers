package queries

import (
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bites query", func() {
	var (
		headerRepo             datastore.HeaderRepository
		blockOne, timestampOne int
		fakeUrn                = test_data.RandomString(5)
		fakeFlipAddress        = fakes.FakeAddress.Hex()
		headerOne              core.Header
	)

	const allBitesQuery = `SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.all_bites($1)`
	const urnBitesQuery = `SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.urn_bites($1, $2)`

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepo)
	})

	Describe("all_bites", func() {
		It("returns bites for an ilk", func() {
			biteLog := test_data.CreateTestLog(headerOne.Id, db)

			biteOne := generateBite(test_helpers.FakeIlk.Hex, fakeUrn, fakeFlipAddress, headerOne.Id, biteLog.ID, db)
			createErr := event.PersistModels([]event.InsertionModel{biteOne}, db)
			Expect(createErr).NotTo(HaveOccurred())

			var actualBites []test_helpers.BiteEvent
			getErr := db.Select(&actualBites, allBitesQuery,
				test_helpers.FakeIlk.Identifier)
			Expect(getErr).NotTo(HaveOccurred())

			Expect(actualBites).To(ConsistOf(
				test_helpers.BiteEvent{
					IlkIdentifier: test_helpers.FakeIlk.Identifier,
					UrnIdentifier: fakeUrn,
					Ink:           biteOne.ColumnValues["ink"].(string),
					Art:           biteOne.ColumnValues["art"].(string),
					Tab:           biteOne.ColumnValues["tab"].(string),
				},
			))
		})

		It("returns bites from multiple blocks", func() {
			biteBlockOneLog := test_data.CreateTestLog(headerOne.Id, db)

			biteBlockOne := generateBite(test_helpers.FakeIlk.Hex, fakeUrn, fakeFlipAddress, headerOne.Id, biteBlockOneLog.ID, db)
			createErr := event.PersistModels([]event.InsertionModel{biteBlockOne}, db)
			Expect(createErr).NotTo(HaveOccurred())

			// New block
			headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
			biteBlockTwoLog := test_data.CreateTestLog(headerTwo.Id, db)

			biteBlockTwo := generateBite(test_helpers.FakeIlk.Hex, fakeUrn, fakeFlipAddress, headerTwo.Id, biteBlockTwoLog.ID, db)
			createErrTwo := event.PersistModels([]event.InsertionModel{biteBlockTwo}, db)
			Expect(createErrTwo).NotTo(HaveOccurred())

			var actualBites []test_helpers.BiteEvent
			getErr := db.Select(&actualBites, allBitesQuery, test_helpers.FakeIlk.Identifier)
			Expect(getErr).NotTo(HaveOccurred())

			Expect(actualBites).To(ConsistOf(
				test_helpers.BiteEvent{
					IlkIdentifier: test_helpers.FakeIlk.Identifier,
					UrnIdentifier: fakeUrn,
					Ink:           biteBlockTwo.ColumnValues["ink"].(string),
					Art:           biteBlockTwo.ColumnValues["art"].(string),
					Tab:           biteBlockTwo.ColumnValues["tab"].(string),
				},
				test_helpers.BiteEvent{
					IlkIdentifier: test_helpers.FakeIlk.Identifier,
					UrnIdentifier: fakeUrn,
					Ink:           biteBlockOne.ColumnValues["ink"].(string),
					Art:           biteBlockOne.ColumnValues["art"].(string),
					Tab:           biteBlockOne.ColumnValues["tab"].(string),
				},
			))
		})

		It("ignores bites from irrelevant ilks", func() {
			biteLog := test_data.CreateTestLog(headerOne.Id, db)
			irrelevantBiteLog := test_data.CreateTestLog(headerOne.Id, db)

			bite := generateBite(test_helpers.FakeIlk.Hex, fakeUrn, fakeFlipAddress, headerOne.Id, biteLog.ID, db)
			irrelevantBite := generateBite(test_helpers.AnotherFakeIlk.Hex, fakeUrn, fakeFlipAddress, headerOne.Id, irrelevantBiteLog.ID, db)

			createErr := event.PersistModels([]event.InsertionModel{bite, irrelevantBite}, db)
			Expect(createErr).NotTo(HaveOccurred())

			var actualBites []test_helpers.BiteEvent
			getErr := db.Select(&actualBites, allBitesQuery, test_helpers.FakeIlk.Identifier)
			Expect(getErr).NotTo(HaveOccurred())

			Expect(actualBites).To(ConsistOf(
				test_helpers.BiteEvent{
					IlkIdentifier: test_helpers.FakeIlk.Identifier,
					UrnIdentifier: fakeUrn,
					Ink:           bite.ColumnValues["ink"].(string),
					Art:           bite.ColumnValues["art"].(string),
					Tab:           bite.ColumnValues["tab"].(string),
				},
			))
		})

		It("fails if no argument is supplied (STRICT)", func() {
			_, err := db.Exec(`SELECT * FROM api.all_bites()`)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("function api.all_bites() does not exist"))
		})

		Describe("result pagination", func() {
			var biteBlockOne, biteBlockTwo event.InsertionModel

			BeforeEach(func() {
				biteBlockOneLog := test_data.CreateTestLog(headerOne.Id, db)

				biteBlockOne = generateBite(test_helpers.FakeIlk.Hex, fakeUrn, fakeFlipAddress, headerOne.Id, biteBlockOneLog.ID, db)
				biteOneErr := event.PersistModels([]event.InsertionModel{biteBlockOne}, db)
				Expect(biteOneErr).NotTo(HaveOccurred())

				// New block
				headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
				biteBlockTwoLog := test_data.CreateTestLog(headerTwo.Id, db)

				biteBlockTwo = generateBite(test_helpers.FakeIlk.Hex, fakeUrn, fakeFlipAddress, headerTwo.Id, biteBlockTwoLog.ID, db)
				biteTwoErr := event.PersistModels([]event.InsertionModel{biteBlockTwo}, db)
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
						Ink:           biteBlockTwo.ColumnValues["ink"].(string),
						Art:           biteBlockTwo.ColumnValues["art"].(string),
						Tab:           biteBlockTwo.ColumnValues["tab"].(string),
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
						Ink:           biteBlockOne.ColumnValues["ink"].(string),
						Art:           biteBlockOne.ColumnValues["art"].(string),
						Tab:           biteBlockOne.ColumnValues["tab"].(string),
					},
				))
			})
		})
	})

	Describe("urn_bites", func() {
		It("returns bites for relevant ilk/urn", func() {
			biteOneLog := test_data.CreateTestLog(headerOne.Id, db)

			biteOne := generateBite(test_helpers.FakeIlk.Hex, fakeUrn, fakeFlipAddress, headerOne.Id, biteOneLog.ID, db)
			createErr := event.PersistModels([]event.InsertionModel{biteOne}, db)
			Expect(createErr).NotTo(HaveOccurred())

			var actualBites []test_helpers.BiteEvent
			getErr := db.Select(&actualBites, urnBitesQuery, test_helpers.FakeIlk.Identifier, fakeUrn)
			Expect(getErr).NotTo(HaveOccurred())

			Expect(actualBites).To(ConsistOf(
				test_helpers.BiteEvent{
					IlkIdentifier: test_helpers.FakeIlk.Identifier,
					UrnIdentifier: fakeUrn,
					Ink:           biteOne.ColumnValues["ink"].(string),
					Art:           biteOne.ColumnValues["art"].(string),
					Tab:           biteOne.ColumnValues["tab"].(string),
				},
			))
		})

		It("returns bites from multiple blocks", func() {
			biteOneLog := test_data.CreateTestLog(headerOne.Id, db)

			biteBlockOne := generateBite(test_helpers.FakeIlk.Hex, fakeUrn, fakeFlipAddress, headerOne.Id, biteOneLog.ID, db)
			createErr := event.PersistModels([]event.InsertionModel{biteBlockOne}, db)
			Expect(createErr).NotTo(HaveOccurred())

			// New block
			headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
			biteBlockTwoLog := test_data.CreateTestLog(headerTwo.Id, db)

			biteBlockTwo := generateBite(test_helpers.FakeIlk.Hex, fakeUrn, fakeFlipAddress, headerTwo.Id, biteBlockTwoLog.ID, db)
			createErrTwo := event.PersistModels([]event.InsertionModel{biteBlockTwo}, db)
			Expect(createErrTwo).NotTo(HaveOccurred())

			var actualBites []test_helpers.BiteEvent
			getErr := db.Select(&actualBites, urnBitesQuery, test_helpers.FakeIlk.Identifier, fakeUrn)
			Expect(getErr).NotTo(HaveOccurred())

			Expect(actualBites).To(ConsistOf(
				test_helpers.BiteEvent{
					IlkIdentifier: test_helpers.FakeIlk.Identifier,
					UrnIdentifier: fakeUrn,
					Ink:           biteBlockTwo.ColumnValues["ink"].(string),
					Art:           biteBlockTwo.ColumnValues["art"].(string),
					Tab:           biteBlockTwo.ColumnValues["tab"].(string),
				},
				test_helpers.BiteEvent{
					IlkIdentifier: test_helpers.FakeIlk.Identifier,
					UrnIdentifier: fakeUrn,
					Ink:           biteBlockOne.ColumnValues["ink"].(string),
					Art:           biteBlockOne.ColumnValues["art"].(string),
					Tab:           biteBlockOne.ColumnValues["tab"].(string),
				},
			))
		})

		Describe("result pagination", func() {
			var biteBlockOne, biteBlockTwo event.InsertionModel

			BeforeEach(func() {
				biteBlockOneLog := test_data.CreateTestLog(headerOne.Id, db)

				biteBlockOne = generateBite(test_helpers.FakeIlk.Hex, fakeUrn, fakeFlipAddress, headerOne.Id, biteBlockOneLog.ID, db)
				createErr := event.PersistModels([]event.InsertionModel{biteBlockOne}, db)
				Expect(createErr).NotTo(HaveOccurred())

				// New block
				headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
				biteBlockTwoLog := test_data.CreateTestLog(headerTwo.Id, db)

				biteBlockTwo = generateBite(test_helpers.FakeIlk.Hex, fakeUrn, fakeFlipAddress, headerTwo.Id, biteBlockTwoLog.ID, db)
				createErrTwo := event.PersistModels([]event.InsertionModel{biteBlockTwo}, db)
				Expect(createErrTwo).NotTo(HaveOccurred())
			})

			It("limits results to latest block if max_results argument is provided", func() {
				maxResults := 1
				var actualBites []test_helpers.BiteEvent
				err := db.Select(&actualBites, `SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.urn_bites($1, $2, $3)`,
					test_helpers.FakeIlk.Identifier, fakeUrn, maxResults)
				Expect(err).NotTo(HaveOccurred())

				Expect(actualBites).To(ConsistOf(
					test_helpers.BiteEvent{
						IlkIdentifier: test_helpers.FakeIlk.Identifier,
						UrnIdentifier: fakeUrn,
						Ink:           biteBlockTwo.ColumnValues["ink"].(string),
						Art:           biteBlockTwo.ColumnValues["art"].(string),
						Tab:           biteBlockTwo.ColumnValues["tab"].(string),
					},
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
					test_helpers.BiteEvent{
						IlkIdentifier: test_helpers.FakeIlk.Identifier,
						UrnIdentifier: fakeUrn,
						Ink:           biteBlockOne.ColumnValues["ink"].(string),
						Art:           biteBlockOne.ColumnValues["art"].(string),
						Tab:           biteBlockOne.ColumnValues["tab"].(string),
					},
				))
			})
		})

		It("ignores bites from irrelevant urns", func() {
			biteLog := test_data.CreateTestLog(headerOne.Id, db)
			irrelevantBiteLog := test_data.CreateTestLog(headerOne.Id, db)

			bite := generateBite(test_helpers.FakeIlk.Hex, fakeUrn, fakeFlipAddress, headerOne.Id, biteLog.ID, db)
			irrelevantBite := generateBite(test_helpers.FakeIlk.Hex, "irrelevantUrn", fakeFlipAddress, headerOne.Id, irrelevantBiteLog.ID, db)

			createErr := event.PersistModels([]event.InsertionModel{bite, irrelevantBite}, db)
			Expect(createErr).NotTo(HaveOccurred())

			var actualBites []test_helpers.BiteEvent
			getErr := db.Select(&actualBites, urnBitesQuery, test_helpers.FakeIlk.Identifier, fakeUrn)
			Expect(getErr).NotTo(HaveOccurred())

			Expect(actualBites).To(ConsistOf(
				test_helpers.BiteEvent{
					IlkIdentifier: test_helpers.FakeIlk.Identifier,
					UrnIdentifier: fakeUrn,
					Ink:           bite.ColumnValues["ink"].(string),
					Art:           bite.ColumnValues["art"].(string),
					Tab:           bite.ColumnValues["tab"].(string),
				},
			))
		})
	})
})

func generateBite(ilk, urn, flipAddress string, headerID, logID int64, db *postgres.DB) event.InsertionModel {
	urnID, urnErr := shared.GetOrCreateUrn(urn, ilk, db)
	Expect(urnErr).NotTo(HaveOccurred())
	addressID, addressErr := repository.GetOrCreateAddress(db, test_data.Cat100Address())
	Expect(addressErr).NotTo(HaveOccurred())
	biteEvent := test_data.BiteModel()
	test_data.AssignAddressID(test_data.BiteEventLog, biteEvent, db)
	biteEvent.ColumnValues["ink"] = strconv.Itoa(rand.Int())
	biteEvent.ColumnValues["art"] = strconv.Itoa(rand.Int())
	biteEvent.ColumnValues["tab"] = strconv.Itoa(rand.Int())
	biteEvent.ColumnValues["bid_id"] = strconv.Itoa(rand.Int())
	biteEvent.ColumnValues[constants.UrnColumn] = urnID
	biteEvent.ColumnValues[event.HeaderFK] = headerID
	biteEvent.ColumnValues[event.LogFK] = logID
	biteEvent.ColumnValues[event.AddressFK] = addressID
	biteEvent.ColumnValues[constants.FlipColumn] = flipAddress
	return biteEvent
}
