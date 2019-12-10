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
	storage_helper "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"math/rand"
	"strconv"

	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/vow"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/utils"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Queued sin computed columns", func() {
	Describe("queued_sin_sin_queue_events", func() {
		var (
			blockOne, timestampOne int
			fakeEra                string
			headerOne              core.Header
			fakeTab                = strconv.Itoa(rand.Int())
			sinMappingMetadata     utils.StorageValueMetadata
			vowRepository          vow.VowStorageRepository
			headerRepository       repositories.HeaderRepository
			diffID                 int64
		)

		BeforeEach(func() {
			test_config.CleanTestDB(db)

			headerRepository = repositories.NewHeaderRepository(db)
			blockOne = rand.Int()
			timestampOne = int(rand.Int31())
			fakeEra = strconv.Itoa(timestampOne)
			headerOne = createHeader(blockOne, timestampOne, headerRepository)

			diffID = storage_helper.CreateFakeDiffRecord(db)

			vowRepository = vow.VowStorageRepository{}
			vowRepository.SetDB(db)
			sinMappingKeys := map[utils.Key]string{constants.Timestamp: fakeEra}
			sinMappingMetadata = utils.GetStorageValueMetadata(vow.SinMapping, sinMappingKeys, utils.Uint256)
			insertSinMappingErr := vowRepository.Create(diffID, headerOne.Id, sinMappingMetadata, fakeTab)
			Expect(insertSinMappingErr).NotTo(HaveOccurred())

			vowFlogLog := test_data.CreateTestLog(headerOne.Id, db)
			vowFlogEvent := test_data.VowFlogModel
			vowFlogEvent.ColumnValues[constants.EraColumn] = fakeEra
			vowFlogEvent.ColumnValues[constants.HeaderFK] = headerOne.Id
			vowFlogEvent.ColumnValues[constants.LogFK] = vowFlogLog.ID
			vowFlogErr := event.PersistModels([]event.InsertionModel{vowFlogEvent}, db)
			Expect(vowFlogErr).NotTo(HaveOccurred())
		})

		It("returns sin queue events for queued sin", func() {
			vowFessLog := test_data.CreateTestLog(headerOne.Id, db)
			vowFessEvent := test_data.VowFessModel
			vowFessEvent.ColumnValues[constants.HeaderFK] = headerOne.Id
			vowFessEvent.ColumnValues[constants.LogFK] = vowFessLog.ID
			vowFessErr := event.PersistModels([]event.InsertionModel{vowFessEvent}, db)
			Expect(vowFessErr).NotTo(HaveOccurred())

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
		})

		Describe("result pagination", func() {
			var headerTwo core.Header

			BeforeEach(func() {
				headerTwo = createHeader(blockOne+1, timestampOne+1, headerRepository)

				// add flog event for same sin in later block
				vowFlogLogTwo := test_data.CreateTestLog(headerTwo.Id, db)
				vowFlogEventTwo := test_data.VowFlogModel
				vowFlogEventTwo.ColumnValues[constants.EraColumn] = fakeEra
				vowFlogEventTwo.ColumnValues[constants.HeaderFK] = headerTwo.Id
				vowFlogEventTwo.ColumnValues[constants.LogFK] = vowFlogLogTwo.ID
				vowFlogErr := event.PersistModels([]event.InsertionModel{vowFlogEventTwo}, db)
				Expect(vowFlogErr).NotTo(HaveOccurred())
			})

			It("limits results to latest blocks if max_results argument is provided", func() {
				maxResults := 1
				var actualEvents []test_helpers.SinQueueEvent
				err := db.Select(&actualEvents,
					`SELECT era, act, block_height
					FROM api.queued_sin_sin_queue_events(
						(SELECT (era, tab, flogged, created, updated)::api.queued_sin FROM api.get_queued_sin($1)), $2)`,
					fakeEra, maxResults)

				Expect(err).NotTo(HaveOccurred())

				Expect(actualEvents).To(ConsistOf(
					test_helpers.SinQueueEvent{
						Era:         fakeEra,
						Act:         "flog",
						BlockHeight: strconv.FormatInt(headerTwo.BlockNumber, 10),
					},
				))
			})

			It("offsets results if offset is provided", func() {
				maxResults := 1
				resultOffset := 1
				var actualEvents []test_helpers.SinQueueEvent
				err := db.Select(&actualEvents,
					`SELECT era, act, block_height
					FROM api.queued_sin_sin_queue_events(
						(SELECT (era, tab, flogged, created, updated)::api.queued_sin FROM api.get_queued_sin($1)), $2, $3)`,
					fakeEra, maxResults, resultOffset)

				Expect(err).NotTo(HaveOccurred())

				Expect(actualEvents).To(ConsistOf(
					test_helpers.SinQueueEvent{
						Era:         fakeEra,
						Act:         "flog",
						BlockHeight: strconv.FormatInt(headerOne.BlockNumber, 10),
					},
				))
			})
		})
	})
})
