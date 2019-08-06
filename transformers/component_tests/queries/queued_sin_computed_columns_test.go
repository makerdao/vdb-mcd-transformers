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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/vow_fess"
	"github.com/vulcanize/mcd_transformers/transformers/events/vow_flog"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vow"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
	"strconv"
)

var _ = Describe("Queued sin computed columns", func() {
	Describe("queued_sin_sin_queue_events", func() {
		var (
			db                 *postgres.DB
			fakeBlock          int
			fakeEra            = rand.Intn(1000000000)
			fakeHeader         core.Header
			fakeTab            = "123"
			headerID           int64
			sinMappingMetadata utils.StorageValueMetadata
			vowRepository      vow.VowStorageRepository
			vowFlogRepo        vow_flog.VowFlogRepository
			headerRepository   repositories.HeaderRepository
		)

		BeforeEach(func() {
			db = test_config.NewTestDB(test_config.NewTestNode())
			test_config.CleanTestDB(db)

			headerRepository = repositories.NewHeaderRepository(db)
			fakeBlock = rand.Int()
			fakeHeader = fakes.GetFakeHeader(int64(fakeBlock))
			fakeHeader.Timestamp = strconv.Itoa(fakeEra)
			var insertHeaderErr error
			headerID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakeHeader)
			Expect(insertHeaderErr).NotTo(HaveOccurred())
			persistedLogOne := test_data.CreateTestLog(headerID, db)
			persistedLogTwo := test_data.CreateTestLog(headerID, db)

			vowRepository = vow.VowStorageRepository{}
			vowRepository.SetDB(db)
			sinMappingKeys := map[utils.Key]string{constants.Timestamp: strconv.Itoa(fakeEra)}
			sinMappingMetadata = utils.GetStorageValueMetadata(vow.SinMapping, sinMappingKeys, utils.Uint256)
			insertSinMappingErr := vowRepository.Create(int(fakeHeader.BlockNumber), fakeHeader.Hash, sinMappingMetadata, fakeTab)
			Expect(insertSinMappingErr).NotTo(HaveOccurred())

			vowFessRepo := vow_fess.VowFessRepository{}
			vowFessRepo.SetDB(db)
			vowFessEvent := test_data.VowFessModel
			vowFessEvent.ColumnValues["header_id"] = headerID
			vowFessEvent.ColumnValues["log_id"] = persistedLogOne.ID
			insertVowFessErr := vowFessRepo.Create([]shared.InsertionModel{vowFessEvent})
			Expect(insertVowFessErr).NotTo(HaveOccurred())

			vowFlogRepo = vow_flog.VowFlogRepository{}
			vowFlogRepo.SetDB(db)
			vowFlogEvent := test_data.VowFlogModel
			vowFlogEvent.ColumnValues["era"] = fakeEra
			vowFlogEvent.ColumnValues["header_id"] = headerID
			vowFlogEvent.ColumnValues["log_id"] = persistedLogTwo.ID
			insertVowFlogErr := vowFlogRepo.Create([]shared.InsertionModel{vowFlogEvent})
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
				test_helpers.SinQueueEvent{Era: strconv.Itoa(fakeEra), Act: "fess"},
				test_helpers.SinQueueEvent{Era: strconv.Itoa(fakeEra), Act: "flog"},
			))
		})

		Describe("result pagination", func() {
			var headerOne, headerTwo core.Header

			BeforeEach(func() {
				blockOne := fakeBlock + 1
				headerOne = fakes.GetFakeHeader(int64(blockOne))
				headerOne.Timestamp = strconv.Itoa(fakeEra + 1)
				headerOneId, headerOneErr := headerRepository.CreateOrUpdateHeader(headerOne)
				Expect(headerOneErr).NotTo(HaveOccurred())

				blockTwo := blockOne + 1
				headerTwo = fakes.GetFakeHeader(int64(blockTwo))
				headerTwo.Timestamp = strconv.Itoa(fakeEra + 1)
				headerTwoId, headerTwoError := headerRepository.CreateOrUpdateHeader(headerTwo)
				Expect(headerTwoError).NotTo(HaveOccurred())

				// add flog event for same sin in later block
				vowFlogEventOne := test_data.VowFlogModel
				vowFlogEventOne.ColumnValues["era"] = fakeEra
				vowFlogEventOne.ColumnValues["header_id"] = headerOneId
				vowFlogOneErr := vowFlogRepo.Create([]shared.InsertionModel{vowFlogEventOne})
				Expect(vowFlogOneErr).NotTo(HaveOccurred())

				// add flog event for same sin in later block
				vowFlogEventTwo := test_data.VowFlogModel
				vowFlogEventTwo.ColumnValues["era"] = fakeEra
				vowFlogEventTwo.ColumnValues["header_id"] = headerTwoId
				vowFlogTwoErr := vowFlogRepo.Create([]shared.InsertionModel{vowFlogEventTwo})
				Expect(vowFlogTwoErr).NotTo(HaveOccurred())
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
						Era:         strconv.Itoa(fakeEra),
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
						Era:         strconv.Itoa(fakeEra),
						Act:         "flog",
						BlockHeight: strconv.FormatInt(headerOne.BlockNumber, 10),
					},
				))
			})
		})
	})
})
