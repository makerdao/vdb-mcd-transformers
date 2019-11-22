package queries

import (
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/spot_poke"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("all poke events query", func() {
	var (
		db                 *postgres.DB
		spotPokeRepo       spot_poke.SpotPokeRepository
		headerRepo         repositories.HeaderRepository
		beginningTimeRange int
		endingTimeRange    int
		blockOne           int
		headerOne          core.Header
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)

		beginningTimeRange = test_helpers.GetRandomInt(1558710000, 1558720000)
		endingTimeRange = test_helpers.GetRandomInt(1558720001, 1558730000)
		headerRepo = repositories.NewHeaderRepository(db)
		spotPokeRepo = spot_poke.SpotPokeRepository{}
		spotPokeRepo.SetDB(db)

		blockOne = rand.Int()
		headerOne = createHeader(blockOne, beginningTimeRange, headerRepo)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	It("returns poke events in different blocks between a time range", func() {
		spotPokeLog := test_data.CreateTestLog(headerOne.Id, db)

		spotPoke := generateSpotPoke(test_helpers.FakeIlk.Hex, 1, headerOne.Id, spotPokeLog.ID)
		ilkIdBlockOne, ilkErr := shared.GetOrCreateIlk(spotPoke.ForeignKeyValues[constants.IlkFK], db)
		Expect(ilkErr).NotTo(HaveOccurred())
		spotPokeErr := spotPokeRepo.Create([]shared.InsertionModel{spotPoke})
		Expect(spotPokeErr).NotTo(HaveOccurred())

		headerTwo := createHeader(blockOne+1, endingTimeRange, headerRepo)
		anotherSpotPokeLog := test_data.CreateTestLog(headerTwo.Id, db)

		anotherSpotPoke := generateSpotPoke(test_helpers.AnotherFakeIlk.Hex, 1, headerTwo.Id, anotherSpotPokeLog.ID)
		anotherIlkId, err := shared.GetOrCreateIlk(anotherSpotPoke.ForeignKeyValues[constants.IlkFK], db)
		Expect(err).NotTo(HaveOccurred())
		err = spotPokeRepo.Create([]shared.InsertionModel{anotherSpotPoke})
		Expect(err).NotTo(HaveOccurred())

		expectedValues := []test_helpers.PokeEvent{
			{
				IlkId: strconv.Itoa(int(anotherIlkId)),
				Val:   anotherSpotPoke.ColumnValues["value"].(string),
				Spot:  anotherSpotPoke.ColumnValues["spot"].(string),
			},
			{
				IlkId: strconv.Itoa(int(ilkIdBlockOne)),
				Val:   spotPoke.ColumnValues["value"].(string),
				Spot:  spotPoke.ColumnValues["spot"].(string),
			},
		}

		var dbPokeEvents []test_helpers.PokeEvent
		err = db.Select(&dbPokeEvents, `SELECT ilk_id, val, spot FROM api.all_poke_events($1, $2)`, beginningTimeRange, endingTimeRange)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbPokeEvents).To(ConsistOf(expectedValues))
	})

	It("returns poke events with transactions in the same block", func() {
		spotPokeLog := test_data.CreateTestLog(headerOne.Id, db)

		spotPoke := generateSpotPoke(test_helpers.FakeIlk.Hex, 1, headerOne.Id, spotPokeLog.ID)
		ilkIdBlockOne, ilkErr := shared.GetOrCreateIlk(spotPoke.ForeignKeyValues[constants.IlkFK], db)
		Expect(ilkErr).NotTo(HaveOccurred())
		spotPokeErr := spotPokeRepo.Create([]shared.InsertionModel{spotPoke})
		Expect(spotPokeErr).NotTo(HaveOccurred())
		anotherSpotPokeLog := test_data.CreateTestLog(headerOne.Id, db)

		anotherSpotPoke := generateSpotPoke(test_helpers.AnotherFakeIlk.Hex, 1, headerOne.Id, anotherSpotPokeLog.ID)
		anotherIlkId, err := shared.GetOrCreateIlk(anotherSpotPoke.ForeignKeyValues[constants.IlkFK], db)
		Expect(err).NotTo(HaveOccurred())
		err = spotPokeRepo.Create([]shared.InsertionModel{anotherSpotPoke})
		Expect(err).NotTo(HaveOccurred())

		expectedValues := []test_helpers.PokeEvent{
			{
				IlkId: strconv.Itoa(int(ilkIdBlockOne)),
				Val:   spotPoke.ColumnValues["value"].(string),
				Spot:  spotPoke.ColumnValues["spot"].(string),
			},
			{
				IlkId: strconv.Itoa(int(anotherIlkId)),
				Val:   anotherSpotPoke.ColumnValues["value"].(string),
				Spot:  anotherSpotPoke.ColumnValues["spot"].(string),
			},
		}

		var dbPokeEvents []test_helpers.PokeEvent
		err = db.Select(&dbPokeEvents, `SELECT ilk_id, val, spot FROM api.all_poke_events($1, $2)`, beginningTimeRange, endingTimeRange)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbPokeEvents).To(ConsistOf(expectedValues))
	})

	It("ignores poke events not in time range", func() {
		spotPokeLog := test_data.CreateTestLog(headerOne.Id, db)
		spotPoke := generateSpotPoke(test_helpers.FakeIlk.Hex, 1, headerOne.Id, spotPokeLog.ID)
		ilkIdBlockOne, ilkErr := shared.GetOrCreateIlk(spotPoke.ForeignKeyValues[constants.IlkFK], db)
		Expect(ilkErr).NotTo(HaveOccurred())
		spotPokeErr := spotPokeRepo.Create([]shared.InsertionModel{spotPoke})
		Expect(spotPokeErr).NotTo(HaveOccurred())

		headerTwo := createHeader(blockOne+1, endingTimeRange+1, headerRepo)
		anotherSpotPokeLog := test_data.CreateTestLog(headerTwo.Id, db)
		anotherSpotPoke := generateSpotPoke(test_helpers.AnotherFakeIlk.Hex, 1, headerTwo.Id, anotherSpotPokeLog.ID)

		_, anotherIlkErr := shared.GetOrCreateIlk(anotherSpotPoke.ForeignKeyValues[constants.IlkFK], db)
		Expect(anotherIlkErr).NotTo(HaveOccurred())
		anotherSpotPokeErr := spotPokeRepo.Create([]shared.InsertionModel{anotherSpotPoke})
		Expect(anotherSpotPokeErr).NotTo(HaveOccurred())

		expectedValues := []test_helpers.PokeEvent{
			{
				IlkId: strconv.Itoa(int(ilkIdBlockOne)),
				Val:   spotPoke.ColumnValues["value"].(string),
				Spot:  spotPoke.ColumnValues["spot"].(string),
			},
		}

		var dbPokeEvents []test_helpers.PokeEvent
		selectErr := db.Select(&dbPokeEvents, `SELECT ilk_id, val, spot FROM api.all_poke_events($1, $2)`, beginningTimeRange, endingTimeRange)
		Expect(selectErr).NotTo(HaveOccurred())
		Expect(dbPokeEvents).To(Equal(expectedValues))
	})

	Describe("result pagination", func() {
		var (
			ilkId                       int64
			oldSpotPoke, recentSpotPoke shared.InsertionModel
		)
		BeforeEach(func() {
			logID := test_data.CreateTestLog(headerOne.Id, db).ID

			oldSpotPoke = generateSpotPoke(test_helpers.FakeIlk.Hex, 1, headerOne.Id, logID)
			var ilkErr error
			ilkId, ilkErr = shared.GetOrCreateIlk(oldSpotPoke.ForeignKeyValues[constants.IlkFK], db)
			Expect(ilkErr).NotTo(HaveOccurred())
			oldSpotPokeErr := spotPokeRepo.Create([]shared.InsertionModel{oldSpotPoke})
			Expect(oldSpotPokeErr).NotTo(HaveOccurred())

			headerTwo := createHeader(blockOne+1, endingTimeRange, headerRepo)
			anotherLogID := test_data.CreateTestLog(headerTwo.Id, db).ID

			recentSpotPoke = generateSpotPoke(test_helpers.FakeIlk.Hex, 2, headerTwo.Id, anotherLogID)
			recentSpotPokeErr := spotPokeRepo.Create([]shared.InsertionModel{recentSpotPoke})
			Expect(recentSpotPokeErr).NotTo(HaveOccurred())
		})

		It("limits results to latest blocks if max_results argument is provided", func() {
			maxResults := 1
			var dbPokeEvents []test_helpers.PokeEvent
			selectErr := db.Select(&dbPokeEvents, `SELECT ilk_id, val, spot FROM api.all_poke_events($1, $2, $3)`,
				beginningTimeRange, endingTimeRange, maxResults)
			Expect(selectErr).NotTo(HaveOccurred())

			Expect(dbPokeEvents).To(ConsistOf(test_helpers.PokeEvent{
				IlkId: strconv.FormatInt(ilkId, 10),
				Val:   recentSpotPoke.ColumnValues["value"].(string),
				Spot:  recentSpotPoke.ColumnValues["spot"].(string),
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
				IlkId: strconv.FormatInt(ilkId, 10),
				Val:   oldSpotPoke.ColumnValues["value"].(string),
				Spot:  oldSpotPoke.ColumnValues["spot"].(string),
			}))
		})
	})

	It("uses default arguments when none are passed in", func() {
		_, err := db.Exec(`SELECT * FROM api.all_poke_events()`)
		Expect(err).NotTo(HaveOccurred())
	})
})

func generateSpotPoke(ilk string, seed int, headerID, logID int64) shared.InsertionModel {
	spotPoke := test_data.SpotPokeModel()
	spotPoke.ForeignKeyValues[constants.IlkFK] = ilk
	spotPoke.ColumnValues["value"] = strconv.Itoa(1 + seed)
	spotPoke.ColumnValues["spot"] = strconv.Itoa(2 + seed)
	spotPoke.ColumnValues[constants.HeaderFK] = headerID
	spotPoke.ColumnValues[constants.LogFK] = logID
	return spotPoke
}
