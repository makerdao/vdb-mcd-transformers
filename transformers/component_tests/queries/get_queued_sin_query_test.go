package queries

import (
	"database/sql"
	"math/rand"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/vow_flog"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vow"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("QueuedSin", func() {
	var (
		db                 *postgres.DB
		fakeBlock          int
		fakeEra            = "1557920248"
		fakeHeader         core.Header
		fakeTab            = "123"
		headerID           int64
		sinMappingMetadata utils.StorageValueMetadata
		vowRepository      vow.VowStorageRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)

		headerRepository = repositories.NewHeaderRepository(db)
		fakeBlock = rand.Int()
		fakeHeader = fakes.GetFakeHeader(int64(fakeBlock))
		var insertHeaderErr error
		headerID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())

		vowRepository = vow.VowStorageRepository{}
		vowRepository.SetDB(db)
		sinMappingKeys := map[utils.Key]string{constants.Timestamp: fakeEra}
		sinMappingMetadata = utils.GetStorageValueMetadata(vow.SinMapping, sinMappingKeys, utils.Uint256)
		insertSinMappingErr := vowRepository.Create(int(fakeHeader.BlockNumber), fakeHeader.Hash, sinMappingMetadata, fakeTab)
		Expect(insertSinMappingErr).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("getting a single queued sin for an era", func() {
		It("gets queued sin for an era", func() {
			var result QueuedSin
			err := db.Get(&result, `SELECT era, tab, flogged, created, updated from api.get_queued_sin($1)`, fakeEra)
			Expect(err).NotTo(HaveOccurred())

			Expect(result.Era).To(Equal(fakeEra))
			Expect(result.Tab).To(Equal(fakeTab))
			Expect(result.Flogged).To(BeFalse())
			timestampAsInt, convertErr := strconv.ParseInt(fakeHeader.Timestamp, 10, 64)
			Expect(convertErr).NotTo(HaveOccurred())
			expectedTimestamp := time.Unix(timestampAsInt, 0).UTC().Format(time.RFC3339)
			Expect(result.Created).To(Equal(expectedTimestamp))
			Expect(result.Updated).To(Equal(expectedTimestamp))
		})

		It("returns flogged as true if era has been flogged", func() {
			vowFlogRepository := vow_flog.VowFlogRepository{}
			vowFlogRepository.SetDB(db)
			vowFlogEvent := test_data.VowFlogModel
			vowFlogEvent.Era = fakeEra
			insertVowFlogErr := vowFlogRepository.Create(headerID, []interface{}{vowFlogEvent})
			Expect(insertVowFlogErr).NotTo(HaveOccurred())

			var result QueuedSin
			err := db.Get(&result, `SELECT era, tab, flogged, created, updated from api.get_queued_sin($1)`, fakeEra)
			Expect(err).NotTo(HaveOccurred())

			Expect(result.Flogged).To(BeTrue())
		})

		It("does not return queued sin for another era", func() {
			anotherFakeEra := "1111111111"
			var result QueuedSin
			err := db.Get(&result, `SELECT era, tab, flogged, created, updated from api.get_queued_sin($1)`, anotherFakeEra)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(sql.ErrNoRows))
		})

		Context("when tab for an era has been updated", func() {
			var (
				anotherFakeTab = "321"
				laterTimestamp string
			)

			BeforeEach(func() {
				laterHeader := fakeHeader
				laterHeader.BlockNumber = int64(fakeBlock + 1000)
				laterHeader.Hash = test_data.RandomString(5)
				fakeHeaderTimestampAsInt, convertCreatedErr := strconv.ParseInt(fakeHeader.Timestamp, 10, 64)
				Expect(convertCreatedErr).NotTo(HaveOccurred())
				laterTimestamp = strconv.Itoa(int(fakeHeaderTimestampAsInt + 1))
				laterHeader.Timestamp = laterTimestamp
				_, insertHeaderErr := headerRepository.CreateOrUpdateHeader(laterHeader)
				Expect(insertHeaderErr).NotTo(HaveOccurred())

				insertVowMappingErr := vowRepository.Create(int(laterHeader.BlockNumber), laterHeader.Hash, sinMappingMetadata, anotherFakeTab)
				Expect(insertVowMappingErr).NotTo(HaveOccurred())
			})

			It("returns most recent 'updated' timestamp", func() {
				var result QueuedSin
				err := db.Get(&result, `SELECT era, tab, flogged, created, updated from api.get_queued_sin($1)`, fakeEra)
				Expect(err).NotTo(HaveOccurred())

				createdTimestampAsInt, convertCreatedErr := strconv.ParseInt(fakeHeader.Timestamp, 10, 64)
				Expect(convertCreatedErr).NotTo(HaveOccurred())
				expectedCreatedTimestamp := time.Unix(createdTimestampAsInt, 0).UTC().Format(time.RFC3339)
				Expect(result.Created).To(Equal(expectedCreatedTimestamp))

				updatedTimestampAsInt, convertUpdatedErr := strconv.ParseInt(laterTimestamp, 10, 64)
				Expect(convertUpdatedErr).NotTo(HaveOccurred())
				expectedUpdatedTimestamp := time.Unix(updatedTimestampAsInt, 0).UTC().Format(time.RFC3339)
				Expect(result.Updated).To(Equal(expectedUpdatedTimestamp))
			})

			It("returns most recent tab value", func() {
				var result QueuedSin
				err := db.Get(&result, `SELECT era, tab, flogged, created, updated from api.get_queued_sin($1)`, fakeEra)
				Expect(err).NotTo(HaveOccurred())

				Expect(result.Tab).To(Equal(anotherFakeTab))
			})
		})
	})

	Describe("getting all queued sins", func() {
		It("returns queued sin for every era", func() {
			anotherFakeEra := "1111111111"
			anotherFakeTab := "321"
			anotherSinMappingKeys := map[utils.Key]string{constants.Timestamp: anotherFakeEra}
			anotherSinMappingMetadata := utils.GetStorageValueMetadata(vow.SinMapping, anotherSinMappingKeys, utils.Uint256)
			insertSinMappingErr := vowRepository.Create(int(fakeHeader.BlockNumber), fakeHeader.Hash, anotherSinMappingMetadata, anotherFakeTab)
			Expect(insertSinMappingErr).NotTo(HaveOccurred())

			var results []QueuedSin
			err := db.Select(&results, `SELECT era, tab, flogged, created, updated from api.all_queued_sin()`)
			Expect(err).NotTo(HaveOccurred())

			Expect(len(results)).To(Equal(2))
			Expect(results[0].Era).To(Or(Equal(fakeEra), Equal(anotherFakeEra)))
			Expect(results[0].Tab).To(Or(Equal(fakeTab), Equal(anotherFakeTab)))
		})
	})
})

type QueuedSin struct {
	Era     string
	Tab     string
	Flogged bool
	Created string
	Updated string
}
