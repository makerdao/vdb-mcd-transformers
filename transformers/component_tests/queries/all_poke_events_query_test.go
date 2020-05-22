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
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("all poke events query", func() {
	var (
		headerRepo         datastore.HeaderRepository
		beginningTimeRange int
		endingTimeRange    int
		blockOne           int
		headerOne          core.Header
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)

		beginningTimeRange = test_helpers.GetRandomInt(1558710000, 1558720000)
		endingTimeRange = test_helpers.GetRandomInt(1558720001, 1558730000)
		headerRepo = repositories.NewHeaderRepository(db)

		blockOne = rand.Int()
		headerOne = createHeader(blockOne, beginningTimeRange, headerRepo)
	})

	It("returns poke events in different blocks between a time range", func() {
		spotPokeLog := test_data.CreateTestLog(headerOne.Id, db)

		spotPoke := generateSpotPoke(test_helpers.FakeIlk.Hex, 1, headerOne.Id, spotPokeLog.ID, db)
		spotPokeErr := event.PersistModels([]event.InsertionModel{spotPoke}, db)
		Expect(spotPokeErr).NotTo(HaveOccurred())

		headerTwo := createHeader(blockOne+1, endingTimeRange, headerRepo)
		anotherSpotPokeLog := test_data.CreateTestLog(headerTwo.Id, db)

		anotherSpotPoke := generateSpotPoke(test_helpers.AnotherFakeIlk.Hex, 1, headerTwo.Id, anotherSpotPokeLog.ID, db)
		err := event.PersistModels([]event.InsertionModel{anotherSpotPoke}, db)
		Expect(err).NotTo(HaveOccurred())

		expectedValues := []test_helpers.PokeEvent{
			{
				IlkId: strconv.FormatInt(anotherSpotPoke.ColumnValues[constants.IlkColumn].(int64), 10),
				Val:   anotherSpotPoke.ColumnValues[constants.ValueColumn].(string),
				Spot:  anotherSpotPoke.ColumnValues[constants.SpotColumn].(string),
			},
			{
				IlkId: strconv.FormatInt(spotPoke.ColumnValues[constants.IlkColumn].(int64), 10),
				Val:   spotPoke.ColumnValues[constants.ValueColumn].(string),
				Spot:  spotPoke.ColumnValues[constants.SpotColumn].(string),
			},
		}

		var dbPokeEvents []test_helpers.PokeEvent
		err = db.Select(&dbPokeEvents, `SELECT ilk_id, val, spot FROM api.all_poke_events($1, $2)`, beginningTimeRange, endingTimeRange)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbPokeEvents).To(ConsistOf(expectedValues))
	})

	It("returns poke events with transactions in the same block", func() {
		spotPokeLog := test_data.CreateTestLog(headerOne.Id, db)

		spotPoke := generateSpotPoke(test_helpers.FakeIlk.Hex, 1, headerOne.Id, spotPokeLog.ID, db)
		spotPokeErr := event.PersistModels([]event.InsertionModel{spotPoke}, db)
		Expect(spotPokeErr).NotTo(HaveOccurred())
		anotherSpotPokeLog := test_data.CreateTestLog(headerOne.Id, db)

		anotherSpotPoke := generateSpotPoke(test_helpers.AnotherFakeIlk.Hex, 1, headerOne.Id, anotherSpotPokeLog.ID, db)
		anotherErr := event.PersistModels([]event.InsertionModel{anotherSpotPoke}, db)
		Expect(anotherErr).NotTo(HaveOccurred())

		expectedValues := []test_helpers.PokeEvent{
			{
				IlkId: strconv.FormatInt(spotPoke.ColumnValues[constants.IlkColumn].(int64), 10),
				Val:   spotPoke.ColumnValues[constants.ValueColumn].(string),
				Spot:  spotPoke.ColumnValues[constants.SpotColumn].(string),
			},
			{
				IlkId: strconv.FormatInt(anotherSpotPoke.ColumnValues[constants.IlkColumn].(int64), 10),
				Val:   anotherSpotPoke.ColumnValues[constants.ValueColumn].(string),
				Spot:  anotherSpotPoke.ColumnValues[constants.SpotColumn].(string),
			},
		}

		var dbPokeEvents []test_helpers.PokeEvent
		anotherErr = db.Select(&dbPokeEvents, `SELECT ilk_id, val, spot FROM api.all_poke_events($1, $2)`, beginningTimeRange, endingTimeRange)
		Expect(anotherErr).NotTo(HaveOccurred())
		Expect(dbPokeEvents).To(ConsistOf(expectedValues))
	})

	It("ignores poke events not in time range", func() {
		spotPokeLog := test_data.CreateTestLog(headerOne.Id, db)
		spotPoke := generateSpotPoke(test_helpers.FakeIlk.Hex, 1, headerOne.Id, spotPokeLog.ID, db)
		spotPokeErr := event.PersistModels([]event.InsertionModel{spotPoke}, db)
		Expect(spotPokeErr).NotTo(HaveOccurred())

		headerTwo := createHeader(blockOne+1, endingTimeRange+1, headerRepo)
		anotherSpotPokeLog := test_data.CreateTestLog(headerTwo.Id, db)
		anotherSpotPoke := generateSpotPoke(test_helpers.AnotherFakeIlk.Hex, 1, headerTwo.Id, anotherSpotPokeLog.ID, db)
		anotherSpotPokeErr := event.PersistModels([]event.InsertionModel{anotherSpotPoke}, db)
		Expect(anotherSpotPokeErr).NotTo(HaveOccurred())

		expectedValues := []test_helpers.PokeEvent{
			{
				IlkId: strconv.FormatInt(spotPoke.ColumnValues[constants.IlkColumn].(int64), 10),
				Val:   spotPoke.ColumnValues[constants.ValueColumn].(string),
				Spot:  spotPoke.ColumnValues[constants.SpotColumn].(string),
			},
		}

		var dbPokeEvents []test_helpers.PokeEvent
		selectErr := db.Select(&dbPokeEvents, `SELECT ilk_id, val, spot FROM api.all_poke_events($1, $2)`, beginningTimeRange, endingTimeRange)
		Expect(selectErr).NotTo(HaveOccurred())
		Expect(dbPokeEvents).To(Equal(expectedValues))
	})

	Describe("result pagination", func() {
		var oldSpotPoke, recentSpotPoke event.InsertionModel

		BeforeEach(func() {
			logID := test_data.CreateTestLog(headerOne.Id, db).ID

			oldSpotPoke = generateSpotPoke(test_helpers.FakeIlk.Hex, 1, headerOne.Id, logID, db)
			oldSpotPokeErr := event.PersistModels([]event.InsertionModel{oldSpotPoke}, db)
			Expect(oldSpotPokeErr).NotTo(HaveOccurred())

			headerTwo := createHeader(blockOne+1, endingTimeRange, headerRepo)
			anotherLogID := test_data.CreateTestLog(headerTwo.Id, db).ID

			recentSpotPoke = generateSpotPoke(test_helpers.FakeIlk.Hex, 2, headerTwo.Id, anotherLogID, db)
			recentSpotPokeErr := event.PersistModels([]event.InsertionModel{recentSpotPoke}, db)
			Expect(recentSpotPokeErr).NotTo(HaveOccurred())
		})

		It("limits results to latest blocks if max_results argument is provided", func() {
			maxResults := 1
			var dbPokeEvents []test_helpers.PokeEvent
			selectErr := db.Select(&dbPokeEvents, `SELECT ilk_id, val, spot FROM api.all_poke_events($1, $2, $3)`,
				beginningTimeRange, endingTimeRange, maxResults)
			Expect(selectErr).NotTo(HaveOccurred())

			Expect(dbPokeEvents).To(ConsistOf(test_helpers.PokeEvent{
				IlkId: strconv.FormatInt(recentSpotPoke.ColumnValues[constants.IlkColumn].(int64), 10),
				Val:   recentSpotPoke.ColumnValues[constants.ValueColumn].(string),
				Spot:  recentSpotPoke.ColumnValues[constants.SpotColumn].(string),
			}))
		})

		It("offsets results if offset is provided", func() {
			maxResults := 1
			resultOffset := 1
			var dbPokeEvents []test_helpers.PokeEvent
			selectErr := db.Select(&dbPokeEvents, `SELECT ilk_id, val, spot FROM api.all_poke_events($1, $2, $3, $4)`,
				beginningTimeRange, endingTimeRange, maxResults, resultOffset)
			Expect(selectErr).NotTo(HaveOccurred())

			Expect(dbPokeEvents).To(ConsistOf(test_helpers.PokeEvent{
				IlkId: strconv.FormatInt(oldSpotPoke.ColumnValues[constants.IlkColumn].(int64), 10),
				Val:   oldSpotPoke.ColumnValues[constants.ValueColumn].(string),
				Spot:  oldSpotPoke.ColumnValues[constants.SpotColumn].(string),
			}))
		})
	})

	It("uses default arguments when none are passed in", func() {
		_, err := db.Exec(`SELECT * FROM api.all_poke_events()`)
		Expect(err).NotTo(HaveOccurred())
	})
})

func generateSpotPoke(ilk string, seed int, headerID, logID int64, db *postgres.DB) event.InsertionModel {
	ilkID, ilkErr := shared.GetOrCreateIlk(ilk, db)
	Expect(ilkErr).NotTo(HaveOccurred())

	spotPoke := test_data.SpotPokeModel()
	spotPoke.ColumnValues[event.HeaderFK] = headerID
	spotPoke.ColumnValues[event.LogFK] = logID
	spotPoke.ColumnValues[constants.IlkColumn] = ilkID
	spotPoke.ColumnValues[constants.ValueColumn] = strconv.Itoa(1 + seed)
	spotPoke.ColumnValues[constants.SpotColumn] = strconv.Itoa(2 + seed)
	return spotPoke
}
