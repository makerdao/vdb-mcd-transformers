package queries

import (
	"math/rand"
	"time"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	ilkHelpers "github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	headerHelpers "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Time Ilk Snapshots query", func() {
	var (
		headerOne, headerTwo       core.Header
		timestampOne, timestampTwo int64
		datetimeOne, datetimeTwo   string

		ethIlkValuesOne, ethIlkValuesTwo, batIlkValuesOne, batIlkValuesTwo map[string]interface{}
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)

		blockOne := rand.Int()
		blockTwo := blockOne + 1
		timestampOne = int64(rand.Int31())
		timestampTwo = timestampOne + 86400
		datetimeOne = time.Unix(timestampOne, 0).UTC().Format(time.RFC3339)
		datetimeTwo = time.Unix(timestampTwo, 0).UTC().Format(time.RFC3339)
		headerOne = headerHelpers.CreateHeader(timestampOne, blockOne, db)
		headerTwo = headerHelpers.CreateHeader(timestampTwo, blockTwo, db)

		// set up 2 different ilks in 2 different blocks
		ethIlkValuesOne = ilkHelpers.GetIlkValues(0)
		ethIlkValuesTwo = ilkHelpers.GetIlkValues(1)
		batIlkValuesOne = ilkHelpers.GetIlkValues(2)
		batIlkValuesTwo = ilkHelpers.GetIlkValues(3)
		ilkHelpers.CreateIlk(db, headerOne, ethIlkValuesOne, ilkHelpers.FakeIlkVatMetadatas, ilkHelpers.FakeIlkCatMetadatas, ilkHelpers.FakeIlkJugMetadatas, ilkHelpers.FakeIlkSpotMetadatas)
		ilkHelpers.CreateIlk(db, headerTwo, ethIlkValuesTwo, ilkHelpers.FakeIlkVatMetadatas, ilkHelpers.FakeIlkCatMetadatas, ilkHelpers.FakeIlkJugMetadatas, ilkHelpers.FakeIlkSpotMetadatas)
		ilkHelpers.CreateIlk(db, headerOne, batIlkValuesOne, ilkHelpers.AnotherFakeIlkVatMetadatas, ilkHelpers.AnotherFakeIlkCatMetadatas, ilkHelpers.AnotherFakeIlkJugMetadatas, ilkHelpers.AnotherFakeIlkSpotMetadatas)
		ilkHelpers.CreateIlk(db, headerTwo, batIlkValuesTwo, ilkHelpers.AnotherFakeIlkVatMetadatas, ilkHelpers.AnotherFakeIlkCatMetadatas, ilkHelpers.AnotherFakeIlkJugMetadatas, ilkHelpers.AnotherFakeIlkSpotMetadatas)
	})

	Context("when called with first ilk identifier and a range of exactly the two blocks on a 1 day interval", func() {
		It("returns the first block's state of the first ilk", func() {
			expectedEthIlk := ilkHelpers.BucketedIlkSnapshotFromValues(ilkHelpers.FakeIlk.Hex, headerOne.Timestamp,
				headerOne.Timestamp, "1 day", timestampOne, timestampOne+86400, ethIlkValuesOne)

			var actualIlks []ilkHelpers.BucketedIlkSnapshot
			err := db.Select(&actualIlks, `SELECT ilk_identifier, bucket_start, bucket_end, bucket_interval, rate, art, spot, line, dust, chop, lump, dunk, flip, rho, duty, pip, mat, created, updated FROM api.time_ilk_snapshots($1, $2, $3)`,
				ilkHelpers.FakeIlk.Identifier, datetimeOne, datetimeTwo)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualIlks).To(ConsistOf(expectedEthIlk))
		})
	})

	Context("when called with second ilk identifier and a range of first block + 6 hours to first block + 8 hours on a 1 hour interval", func() {
		It("returns 2 results of the second ilk identifier both matching the first block", func() {
			sixHours := timestampOne + 21600
			sevenHours := timestampOne + 25200
			eightHours := timestampOne + 28800

			dateOne := time.Unix(sixHours, 0).UTC().Format(time.RFC3339)
			dateTwo := time.Unix(eightHours, 0).UTC().Format(time.RFC3339)

			expectedBatIlk1 := ilkHelpers.BucketedIlkSnapshotFromValues(ilkHelpers.AnotherFakeIlk.Hex, headerOne.Timestamp,
				headerOne.Timestamp, "01:00:00", sixHours, sevenHours, batIlkValuesOne)
			expectedBatIlk2 := ilkHelpers.BucketedIlkSnapshotFromValues(ilkHelpers.AnotherFakeIlk.Hex, headerOne.Timestamp,
				headerOne.Timestamp, "01:00:00", sevenHours, eightHours, batIlkValuesOne)

			var actualIlks []ilkHelpers.BucketedIlkSnapshot
			err := db.Select(&actualIlks, `SELECT ilk_identifier, bucket_start, bucket_end, bucket_interval, rate, art, spot, line, dust, chop, lump, dunk, flip, rho, duty, pip, mat, created, updated FROM api.time_ilk_snapshots($1, $2, $3, '1 hour'::INTERVAL)`,
				ilkHelpers.AnotherFakeIlk.Identifier, dateOne, dateTwo)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualIlks).To(ConsistOf(expectedBatIlk1, expectedBatIlk2))
		})
	})
})
