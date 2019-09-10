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

			vowRepository = vow.VowStorageRepository{}
			vowRepository.SetDB(db)
			sinMappingKeys := map[utils.Key]string{constants.Timestamp: strconv.Itoa(fakeEra)}
			sinMappingMetadata = utils.GetStorageValueMetadata(vow.SinMapping, sinMappingKeys, utils.Uint256)
			insertSinMappingErr := vowRepository.Create(int(fakeHeader.BlockNumber), fakeHeader.Hash, sinMappingMetadata, fakeTab)
			Expect(insertSinMappingErr).NotTo(HaveOccurred())

			vowFessRepo := vow_fess.VowFessRepository{}
			vowFessRepo.SetDB(db)
			vowFessEvent := test_data.VowFessModel
			insertVowFessErr := vowFessRepo.Create(headerID, []shared.InsertionModel{vowFessEvent})
			Expect(insertVowFessErr).NotTo(HaveOccurred())

			vowFlogRepo = vow_flog.VowFlogRepository{}
			vowFlogRepo.SetDB(db)
			vowFlogEvent := test_data.VowFlogModel
			vowFlogEvent.ColumnValues["era"] = fakeEra
			insertVowFlogErr := vowFlogRepo.Create(headerID, []shared.InsertionModel{vowFlogEvent})
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

		It("limits results to latest blocks if max_results argument is provided", func() {
			newBlock := fakeBlock + 1
			newHeader := fakes.GetFakeHeader(int64(newBlock))
			newHeader.Timestamp = strconv.Itoa(fakeEra + 1)
			newHeaderId, newHeaderErr := headerRepository.CreateOrUpdateHeader(newHeader)
			Expect(newHeaderErr).NotTo(HaveOccurred())

			// add flog event for same sin in later block
			vowFlogEvent := test_data.VowFlogModel
			vowFlogEvent.ColumnValues["era"] = fakeEra
			insertVowFlogErr := vowFlogRepo.Create(newHeaderId, []shared.InsertionModel{vowFlogEvent})
			Expect(insertVowFlogErr).NotTo(HaveOccurred())

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
					BlockHeight: strconv.FormatInt(newHeader.BlockNumber, 10),
				},
			))
		})
	})
})
