package queries

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/vow_fess"
	"github.com/vulcanize/mcd_transformers/transformers/events/vow_flog"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vow"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
)

var _ = Describe("Queued sin computed columns", func() {
	Describe("queued_sin_sin_queue_events", func() {
		var (
			db                 *postgres.DB
			fakeBlock          int
			fakeEra            = "1557920248"
			fakeHeader         core.Header
			fakeTab            = "123"
			headerID           int64
			sinMappingMetadata utils.StorageValueMetadata
			vowRepository      vow.VowStorageRepository
			headerRepository   repositories.HeaderRepository
		)

		BeforeEach(func() {
			db = test_config.NewTestDB(test_config.NewTestNode())
			test_config.CleanTestDB(db)

			headerRepository = repositories.NewHeaderRepository(db)
			fakeBlock = rand.Int()
			fakeHeader = fakes.GetFakeHeader(int64(fakeBlock))
			fakeHeader.Timestamp = fakeEra
			var insertHeaderErr error
			headerID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakeHeader)
			Expect(insertHeaderErr).NotTo(HaveOccurred())

			vowRepository = vow.VowStorageRepository{}
			vowRepository.SetDB(db)
			sinMappingKeys := map[utils.Key]string{constants.Timestamp: fakeEra}
			sinMappingMetadata = utils.GetStorageValueMetadata(vow.SinMapping, sinMappingKeys, utils.Uint256)
			insertSinMappingErr := vowRepository.Create(int(fakeHeader.BlockNumber), fakeHeader.Hash, sinMappingMetadata, fakeTab)
			Expect(insertSinMappingErr).NotTo(HaveOccurred())

			vowFessRepo := vow_fess.VowFessRepository{}
			vowFessRepo.SetDB(db)
			vowFessEvent := test_data.VowFessModel
			insertVowFessErr := vowFessRepo.Create(headerID, []interface{}{vowFessEvent})
			Expect(insertVowFessErr).NotTo(HaveOccurred())

			vowFlogRepo := vow_flog.VowFlogRepository{}
			vowFlogRepo.SetDB(db)
			vowFlogEvent := test_data.VowFlogModel
			vowFlogEvent.Era = fakeEra
			insertVowFlogErr := vowFlogRepo.Create(headerID, []interface{}{vowFlogEvent})
			Expect(insertVowFlogErr).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			closeErr := db.Close()
			Expect(closeErr).NotTo(HaveOccurred())
		})

		It("returns sin queue events for queued sin", func() {
			var actualEvents []test_helpers.SinQueueEvent
			err := db.Select(&actualEvents,
				`SELECT era, act
                    FROM api.queued_sin_sin_queue_events(
                        (SELECT (era, tab, flogged, created, updated)::api.queued_sin FROM api.get_queued_sin($1))
                    )`, fakeEra)

			Expect(err).NotTo(HaveOccurred())

			Expect(actualEvents).To(ConsistOf(
				test_helpers.SinQueueEvent{Era: fakeEra, Act: "fess"},
				test_helpers.SinQueueEvent{Era: fakeEra, Act: "flog"},
			))

			Expect(0)
		})
	})
})
