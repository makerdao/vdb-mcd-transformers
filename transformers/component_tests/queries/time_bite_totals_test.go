package queries

import (
	"math/rand"
	"time"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Time Bite Totals query", func() {
	var (
		headerRepo   datastore.HeaderRepository
		blockOne     int
		timestampOne int64
		fakeUrn      = test_data.RandomString(5)
		headerOne    core.Header
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		blockOne = rand.Int()
		timestampOne = int64(rand.Int31())
		headerOne = createHeader(blockOne, int(timestampOne), headerRepo)
	})

	Context("when called with an hourly 2 hour range with the range start on the first block", func() {
		It("returns the all the bite results under the first hour and 0 for the second hour", func() {
			biteLog := test_data.CreateTestLog(headerOne.Id, db)

			oneHour := timestampOne + 3600
			twoHours := timestampOne + 7200

			dateStart := time.Unix(timestampOne, 0).UTC().Format(time.RFC3339)
			dateMiddle := time.Unix(oneHour, 0).UTC().Format(time.RFC3339)
			dateEnd := time.Unix(twoHours, 0).UTC().Format(time.RFC3339)

			biteOne := generateBite(test_helpers.FakeIlk.Hex, fakeUrn, headerOne.Id, biteLog.ID, db)
			createErr := event.PersistModels([]event.InsertionModel{biteOne}, db)
			Expect(createErr).NotTo(HaveOccurred())

			var actualBiteTotals []test_helpers.BucketedBiteTotals
			queryErr := db.Select(&actualBiteTotals, `SELECT bucket_start, bucket_end, bucket_interval, ink, art, tab FROM api.time_bite_totals($1, $2, $3, '1 hour'::INTERVAL)`, test_helpers.FakeIlk.Identifier, dateStart, dateEnd)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBiteTotals).To(ConsistOf(
				test_helpers.BucketedBiteTotals{BucketStart: dateStart, BucketEnd: dateMiddle, BucketInterval: "01:00:00", Ink: biteOne.ColumnValues["ink"].(string), Art: biteOne.ColumnValues["art"].(string), Tab: biteOne.ColumnValues["tab"].(string)},
				test_helpers.BucketedBiteTotals{BucketStart: dateMiddle, BucketEnd: dateEnd, BucketInterval: "01:00:00", Ink: "0", Art: "0", Tab: "0"}))
		})
	})
})
