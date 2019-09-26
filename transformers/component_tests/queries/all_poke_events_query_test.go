package queries

import (
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"math/rand"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/spot_poke"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("all poke events query", func() {
	var (
		db                 *postgres.DB
		spotPokeRepo       spot_poke.SpotPokeRepository
		headerRepo         repositories.HeaderRepository
		beginningTimeRange int64
		endingTimeRange    int64
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)

		beginningTimeRange = int64(test_helpers.GetRandomInt(1558710000, 1558720000))
		endingTimeRange = int64(test_helpers.GetRandomInt(1558720001, 1558730000))
		headerRepo = repositories.NewHeaderRepository(db)
		spotPokeRepo = spot_poke.SpotPokeRepository{}
		spotPokeRepo.SetDB(db)
		rand.Seed(GinkgoRandomSeed())
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	It("returns poke events in different blocks between a time range", func() {
		fakeHeaderOne := fakes.GetFakeHeaderWithTimestamp(beginningTimeRange, rand.Int63())
		headerID, err := headerRepo.CreateOrUpdateHeader(fakeHeaderOne)
		Expect(err).NotTo(HaveOccurred())
		spotPokeLog := test_data.CreateTestLog(headerID, db)

		spotPoke := generateSpotPoke(test_helpers.FakeIlk.Hex, 1, headerID, spotPokeLog.ID)
		ilkIdBlockOne, err := shared.GetOrCreateIlk(spotPoke.ForeignKeyValues[constants.IlkFK], db)
		err = spotPokeRepo.Create([]shared.InsertionModel{spotPoke})
		Expect(err).NotTo(HaveOccurred())

		fakeHeaderTwo := fakes.GetFakeHeaderWithTimestamp(endingTimeRange, fakeHeaderOne.BlockNumber+1)
		anotherHeaderID, err := headerRepo.CreateOrUpdateHeader(fakeHeaderTwo)
		Expect(err).NotTo(HaveOccurred())
		anotherSpotPokeLog := test_data.CreateTestLog(anotherHeaderID, db)

		anotherSpotPoke := generateSpotPoke(test_helpers.AnotherFakeIlk.Hex, 1, anotherHeaderID, anotherSpotPokeLog.ID)
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
		fakeHeaderOne := fakes.GetFakeHeaderWithTimestamp(beginningTimeRange, rand.Int63())
		headerID, err := headerRepo.CreateOrUpdateHeader(fakeHeaderOne)
		Expect(err).NotTo(HaveOccurred())
		spotPokeLog := test_data.CreateTestLog(headerID, db)

		spotPoke := generateSpotPoke(test_helpers.FakeIlk.Hex, 1, headerID, spotPokeLog.ID)
		ilkIdBlockOne, err := shared.GetOrCreateIlk(spotPoke.ForeignKeyValues[constants.IlkFK], db)
		err = spotPokeRepo.Create([]shared.InsertionModel{spotPoke})
		Expect(err).NotTo(HaveOccurred())
		anotherSpotPokeLog := test_data.CreateTestLog(headerID, db)

		anotherSpotPoke := generateSpotPoke(test_helpers.AnotherFakeIlk.Hex, 1, headerID, anotherSpotPokeLog.ID)
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
		fakeHeaderOne := fakes.GetFakeHeaderWithTimestamp(beginningTimeRange, rand.Int63())
		headerID, err := headerRepo.CreateOrUpdateHeader(fakeHeaderOne)
		Expect(err).NotTo(HaveOccurred())
		spotPokeLog := test_data.CreateTestLog(headerID, db)
		spotPoke := generateSpotPoke(test_helpers.FakeIlk.Hex, 1, headerID, spotPokeLog.ID)
		ilkIdBlockOne, err := shared.GetOrCreateIlk(spotPoke.ForeignKeyValues[constants.IlkFK], db)
		err = spotPokeRepo.Create([]shared.InsertionModel{spotPoke})
		Expect(err).NotTo(HaveOccurred())

		fakeHeaderTwo := fakes.GetFakeHeaderWithTimestamp(endingTimeRange+1, fakeHeaderOne.BlockNumber+1)
		anotherHeaderID, err := headerRepo.CreateOrUpdateHeader(fakeHeaderTwo)
		Expect(err).NotTo(HaveOccurred())
		anotherSpotPokeLog := test_data.CreateTestLog(anotherHeaderID, db)
		anotherSpotPoke := generateSpotPoke(test_helpers.AnotherFakeIlk.Hex, 1, anotherHeaderID, anotherSpotPokeLog.ID)

		_, err = shared.GetOrCreateIlk(anotherSpotPoke.ForeignKeyValues[constants.IlkFK], db)
		Expect(err).NotTo(HaveOccurred())
		err = spotPokeRepo.Create([]shared.InsertionModel{anotherSpotPoke})
		Expect(err).NotTo(HaveOccurred())

		expectedValues := []test_helpers.PokeEvent{
			{
				IlkId: strconv.Itoa(int(ilkIdBlockOne)),
				Val:   spotPoke.ColumnValues["value"].(string),
				Spot:  spotPoke.ColumnValues["spot"].(string),
			},
		}

		var dbPokeEvents []test_helpers.PokeEvent
		err = db.Select(&dbPokeEvents, `SELECT ilk_id, val, spot FROM api.all_poke_events($1, $2)`, beginningTimeRange, endingTimeRange)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbPokeEvents).To(Equal(expectedValues))
	})

	Describe("result pagination", func() {
		var (
			ilkId                       int64
			oldSpotPoke, recentSpotPoke shared.InsertionModel
		)
		BeforeEach(func() {
			fakeHeaderOne := fakes.GetFakeHeaderWithTimestamp(beginningTimeRange, int64(test_data.SpotPokeHeaderSyncLog.Log.BlockNumber))
			headerID, headerOneErr := headerRepo.CreateOrUpdateHeader(fakeHeaderOne)
			Expect(headerOneErr).NotTo(HaveOccurred())
			logID := test_data.CreateTestLog(headerID, db).ID

			oldSpotPoke = generateSpotPoke(test_helpers.FakeIlk.Hex, 1, headerID, logID)
			var ilkErr error
			ilkId, ilkErr = shared.GetOrCreateIlk(oldSpotPoke.ForeignKeyValues[constants.IlkFK], db)
			Expect(ilkErr).NotTo(HaveOccurred())
			oldSpotPokeErr := spotPokeRepo.Create([]shared.InsertionModel{oldSpotPoke})
			Expect(oldSpotPokeErr).NotTo(HaveOccurred())

			fakeHeaderTwo := fakes.GetFakeHeaderWithTimestamp(endingTimeRange, fakeHeaderOne.BlockNumber+1)
			anotherHeaderID, headerTwoErr := headerRepo.CreateOrUpdateHeader(fakeHeaderTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())
			anotherLogID := test_data.CreateTestLog(anotherHeaderID, db).ID

			recentSpotPoke = generateSpotPoke(test_helpers.FakeIlk.Hex, 2, anotherHeaderID, anotherLogID)
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
	spotPoke := test_data.CopyModel(test_data.SpotPokeModel)
	spotPoke.ForeignKeyValues[constants.IlkFK] = ilk
	spotPoke.ColumnValues["value"] = strconv.Itoa(1 + seed)
	spotPoke.ColumnValues["spot"] = strconv.Itoa(2 + seed)
	spotPoke.ColumnValues[constants.HeaderFK] = headerID
	spotPoke.ColumnValues[constants.LogFK] = logID
	return spotPoke
}
