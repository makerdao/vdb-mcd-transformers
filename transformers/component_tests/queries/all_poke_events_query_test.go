package queries

import (
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
		persistedLog := test_data.CreateTestLog(headerID, db)

		spotPoke := generateSpotPoke(test_helpers.FakeIlk.Hex, 1, headerID, persistedLog.ID)
		ilkIdBlockOne, err := shared.GetOrCreateIlk(spotPoke.Ilk, db)
		err = spotPokeRepo.Create([]interface{}{spotPoke})
		Expect(err).NotTo(HaveOccurred())

		fakeHeaderTwo := fakes.GetFakeHeaderWithTimestamp(endingTimeRange, fakeHeaderOne.BlockNumber+1)
		anotherHeaderID, err := headerRepo.CreateOrUpdateHeader(fakeHeaderTwo)
		Expect(err).NotTo(HaveOccurred())
		anotherPersistedLog := test_data.CreateTestLog(anotherHeaderID, db)

		anotherSpotPoke := generateSpotPoke(test_helpers.AnotherFakeIlk.Hex, 1, anotherHeaderID, anotherPersistedLog.ID)
		anotherIlkId, err := shared.GetOrCreateIlk(anotherSpotPoke.Ilk, db)
		Expect(err).NotTo(HaveOccurred())
		err = spotPokeRepo.Create([]interface{}{anotherSpotPoke})
		Expect(err).NotTo(HaveOccurred())

		expectedValues := []test_helpers.PokeEvent{
			{
				IlkId: strconv.Itoa(int(anotherIlkId)),
				Val:   anotherSpotPoke.Value,
				Spot:  anotherSpotPoke.Spot,
			},
			{
				IlkId: strconv.Itoa(int(ilkIdBlockOne)),
				Val:   spotPoke.Value,
				Spot:  spotPoke.Spot,
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
		persistedLog := test_data.CreateTestLog(headerID, db)
		anotherPersistedLog := test_data.CreateTestLog(headerID, db)

		spotPoke := generateSpotPoke(test_helpers.FakeIlk.Hex, 1, headerID, persistedLog.ID)
		ilkIdBlockOne, err := shared.GetOrCreateIlk(spotPoke.Ilk, db)
		err = spotPokeRepo.Create([]interface{}{spotPoke})
		Expect(err).NotTo(HaveOccurred())

		anotherSpotPoke := generateSpotPoke(test_helpers.AnotherFakeIlk.Hex, 1, headerID, anotherPersistedLog.ID)
		anotherIlkId, err := shared.GetOrCreateIlk(anotherSpotPoke.Ilk, db)
		Expect(err).NotTo(HaveOccurred())
		err = spotPokeRepo.Create([]interface{}{anotherSpotPoke})
		Expect(err).NotTo(HaveOccurred())

		expectedValues := []test_helpers.PokeEvent{
			{
				IlkId: strconv.Itoa(int(ilkIdBlockOne)),
				Val:   spotPoke.Value,
				Spot:  spotPoke.Spot,
			},
			{
				IlkId: strconv.Itoa(int(anotherIlkId)),
				Val:   anotherSpotPoke.Value,
				Spot:  anotherSpotPoke.Spot,
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
		persistedLogID := test_data.CreateTestLog(headerID, db).ID

		spotPoke := generateSpotPoke(test_helpers.FakeIlk.Hex, 1, headerID, persistedLogID)
		ilkIdBlockOne, err := shared.GetOrCreateIlk(spotPoke.Ilk, db)
		err = spotPokeRepo.Create([]interface{}{spotPoke})
		Expect(err).NotTo(HaveOccurred())

		fakeHeaderTwo := fakes.GetFakeHeaderWithTimestamp(endingTimeRange+1, fakeHeaderOne.BlockNumber+1)
		anotherHeaderID, err := headerRepo.CreateOrUpdateHeader(fakeHeaderTwo)
		Expect(err).NotTo(HaveOccurred())
		anotherPersistedLogID := test_data.CreateTestLog(anotherHeaderID, db).ID

		anotherSpotPoke := generateSpotPoke(test_helpers.AnotherFakeIlk.Hex, 1, anotherHeaderID, anotherPersistedLogID)
		_, err = shared.GetOrCreateIlk(anotherSpotPoke.Ilk, db)
		Expect(err).NotTo(HaveOccurred())
		err = spotPokeRepo.Create([]interface{}{anotherSpotPoke})
		Expect(err).NotTo(HaveOccurred())

		expectedValues := []test_helpers.PokeEvent{
			{
				IlkId: strconv.Itoa(int(ilkIdBlockOne)),
				Val:   spotPoke.Value,
				Spot:  spotPoke.Spot,
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
			oldSpotPoke, recentSpotPoke spot_poke.SpotPokeModel
		)
		BeforeEach(func() {
			fakeHeaderOne := fakes.GetFakeHeaderWithTimestamp(beginningTimeRange, int64(test_data.SpotPokeHeaderSyncLog.Log.BlockNumber))
			headerID, headerOneErr := headerRepo.CreateOrUpdateHeader(fakeHeaderOne)
			Expect(headerOneErr).NotTo(HaveOccurred())
			logID := test_data.CreateTestLog(headerID, db).ID

			oldSpotPoke = generateSpotPoke(test_helpers.FakeIlk.Hex, 1, headerID, logID)
			var ilkErr error
			ilkId, ilkErr = shared.GetOrCreateIlk(oldSpotPoke.Ilk, db)
			Expect(ilkErr).NotTo(HaveOccurred())
			oldSpotPokeErr := spotPokeRepo.Create([]interface{}{oldSpotPoke})
			Expect(oldSpotPokeErr).NotTo(HaveOccurred())

			fakeHeaderTwo := fakes.GetFakeHeaderWithTimestamp(endingTimeRange, fakeHeaderOne.BlockNumber+1)
			anotherHeaderID, headerTwoErr := headerRepo.CreateOrUpdateHeader(fakeHeaderTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())
			anotherLogID := test_data.CreateTestLog(anotherHeaderID, db).ID

			recentSpotPoke = generateSpotPoke(test_helpers.FakeIlk.Hex, 2, anotherHeaderID, anotherLogID)
			recentSpotPokeErr := spotPokeRepo.Create([]interface{}{recentSpotPoke})
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
				Val:   recentSpotPoke.Value,
				Spot:  recentSpotPoke.Spot,
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
				Val:   oldSpotPoke.Value,
				Spot:  oldSpotPoke.Spot,
			}))
		})
	})

	It("uses default arguments when none are passed in", func() {
		_, err := db.Exec(`SELECT * FROM api.all_poke_events()`)
		Expect(err).NotTo(HaveOccurred())
	})
})

func generateSpotPoke(ilk string, seed int, headerID, logID int64) spot_poke.SpotPokeModel {
	spotPoke := test_data.SpotPokeModel
	spotPoke.Ilk = ilk
	spotPoke.Value = strconv.Itoa(1 + seed)
	spotPoke.Spot = strconv.Itoa(2 + seed)
	spotPoke.HeaderID = headerID
	spotPoke.LogID = logID
	return spotPoke
}
