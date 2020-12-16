// VulcanizeDB
// Copyright © 2019 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package queries

import (
	"math/rand"
	"time"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	storage_helper "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Time Frob Totals query", func() {
	var (
		headerRepo                 datastore.HeaderRepository
		fakeIlkHex                 = test_helpers.FakeIlk.Hex
		fakeIlkIdentifier          = test_helpers.FakeIlk.Identifier
		fakeUrn                    = test_data.RandomString(40)
		blockOne                   int
		timestampOne, timestampTwo int64
		datetimeOne, datetimeTwo   string
		headerOne                  core.Header
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)

		blockOne = rand.Int()
		timestampOne = int64(rand.Int31())
		timestampTwo = timestampOne + 86400
		datetimeOne = time.Unix(timestampOne, 0).UTC().Format(time.RFC3339)
		datetimeTwo = time.Unix(timestampTwo, 0).UTC().Format(time.RFC3339)
		headerOne = createHeader(blockOne, int(timestampOne), headerRepo)

		storage_helper.CreateFakeDiffRecord(db)
	})

	Describe("frob_totals", func() {
		It("returns frob totals for single frob", func() {
			vatFrobLog := test_data.CreateTestLog(headerOne.Id, db)
			vatFrobEvent := getFakeVatFrobEvent(fakeIlkHex, fakeUrn, headerOne.Id, vatFrobLog.ID)
			insertFrobErr := event.PersistModels([]event.InsertionModel{vatFrobEvent}, db)
			Expect(insertFrobErr).NotTo(HaveOccurred())

			var actualFrobs []test_helpers.BucketedFrobTotals
			getFrobsErr := db.Select(&actualFrobs, `SELECT bucket_start, bucket_end, bucket_interval, count, dink, dart, lock, free, draw, wipe FROM api.time_frob_totals($1, $2, $3)`, fakeIlkIdentifier, datetimeOne, datetimeTwo)
			Expect(getFrobsErr).NotTo(HaveOccurred())

			Expect(actualFrobs).To(ConsistOf(
				test_helpers.BucketedFrobTotals{
					BucketStart:    datetimeOne,
					BucketEnd:      datetimeTwo,
					BucketInterval: "1 day",
					Count:          "1",
					Dink:           vatFrobEvent.ColumnValues[constants.DinkColumn].(string),
					Dart:           vatFrobEvent.ColumnValues[constants.DartColumn].(string),
					Lock:           vatFrobEvent.ColumnValues[constants.DinkColumn].(string),
					Free:           "0",
					Draw:           vatFrobEvent.ColumnValues[constants.DartColumn].(string),
					Wipe:           "0",
				},
			))
		})
	})
})
