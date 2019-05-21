package queries

import (
	"math/rand"
	"strconv"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/vow_fess"
	"github.com/vulcanize/mcd_transformers/transformers/events/vow_flog"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("Sin queue events query", func() {
	var (
		db         *postgres.DB
		headerRepo repositories.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("all_sin_queue_events", func() {
		It("returns vow fess events", func() {
			fakeEra := strconv.Itoa(int(rand.Int31()))
			headerOne := fakes.GetFakeHeader(1)
			headerOne.Timestamp = fakeEra
			headerOneId, err := headerRepo.CreateOrUpdateHeader(headerOne)
			Expect(err).NotTo(HaveOccurred())

			vowFessRepo := vow_fess.VowFessRepository{}
			vowFessRepo.SetDB(db)
			vowFessEvent := test_data.VowFessModel
			err = vowFessRepo.Create(headerOneId, []interface{}{vowFessEvent})
			Expect(err).NotTo(HaveOccurred())

			var actualEvents []test_helpers.SinQueueEvent
			err = db.Select(&actualEvents, `SELECT era, act FROM api.all_sin_queue_events($1)`, fakeEra)
			Expect(err).NotTo(HaveOccurred())

			Expect(actualEvents).To(ConsistOf(
				test_helpers.SinQueueEvent{Era: fakeEra, Act: "fess"},
			))
		})

		It("returns vow flog events", func() {
			headerOne := fakes.GetFakeHeader(1)
			headerOneId, err := headerRepo.CreateOrUpdateHeader(headerOne)
			Expect(err).NotTo(HaveOccurred())

			fakeEra := strconv.Itoa(int(rand.Int31()))
			vowFlogRepo := vow_flog.VowFlogRepository{}
			vowFlogRepo.SetDB(db)
			vowFlogEvent := test_data.VowFlogModel
			vowFlogEvent.Era = fakeEra
			err = vowFlogRepo.Create(headerOneId, []interface{}{vowFlogEvent})
			Expect(err).NotTo(HaveOccurred())

			var actualEvents []test_helpers.SinQueueEvent
			err = db.Select(&actualEvents, `SELECT era, act FROM api.all_sin_queue_events($1)`, fakeEra)
			Expect(err).NotTo(HaveOccurred())

			Expect(actualEvents).To(ConsistOf(
				test_helpers.SinQueueEvent{Era: fakeEra, Act: "flog"},
			))
		})

		It("returns events from multiple blocks", func() {
			headerOne := fakes.GetFakeHeader(1)
			headerOneId, err := headerRepo.CreateOrUpdateHeader(headerOne)
			Expect(err).NotTo(HaveOccurred())

			vowFlogRepo := vow_flog.VowFlogRepository{}
			vowFlogRepo.SetDB(db)
			fakeEra := strconv.Itoa(int(rand.Int31()))
			vowFlogEvent := test_data.VowFlogModel
			vowFlogEvent.Era = fakeEra
			err = vowFlogRepo.Create(headerOneId, []interface{}{vowFlogEvent})
			Expect(err).NotTo(HaveOccurred())

			// New block
			headerTwo := fakes.GetFakeHeader(2)
			headerTwo.Hash = "anotherHash"
			headerTwo.Timestamp = fakeEra
			headerTwoId, err := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(err).NotTo(HaveOccurred())

			vowFessRepo := vow_fess.VowFessRepository{}
			vowFessRepo.SetDB(db)
			vowFessEvent := test_data.VowFessModel
			err = vowFessRepo.Create(headerTwoId, []interface{}{vowFessEvent})
			Expect(err).NotTo(HaveOccurred())

			var actualEvents []test_helpers.SinQueueEvent
			err = db.Select(&actualEvents, `SELECT era, act FROM api.all_sin_queue_events($1)`, fakeEra)
			Expect(err).NotTo(HaveOccurred())

			Expect(actualEvents).To(ConsistOf(
				test_helpers.SinQueueEvent{Era: fakeEra, Act: "fess"},
				test_helpers.SinQueueEvent{Era: fakeEra, Act: "flog"},
			))
		})

		It("ignores sin queue events with irrelevant eras", func() {
			headerOne := fakes.GetFakeHeader(1)
			headerOneId, err := headerRepo.CreateOrUpdateHeader(headerOne)
			Expect(err).NotTo(HaveOccurred())

			rawEra := int(rand.Int31())
			fakeEra := strconv.Itoa(rawEra)
			irrelevantEra := strconv.Itoa(rawEra + 1)

			vowFlogRepo := vow_flog.VowFlogRepository{}
			vowFlogRepo.SetDB(db)
			vowFlogEvent := test_data.VowFlogModel
			vowFlogEvent.Era = fakeEra
			err = vowFlogRepo.Create(headerOneId, []interface{}{vowFlogEvent})
			Expect(err).NotTo(HaveOccurred())

			var actualEvents []test_helpers.SinQueueEvent
			err = db.Select(&actualEvents, `SELECT era, act FROM api.all_sin_queue_events($1)`, irrelevantEra)
			Expect(err).NotTo(HaveOccurred())

			Expect(actualEvents).To(BeEmpty())
		})
	})
})
