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
	"strconv"

	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vow_flog"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sin queue events query", func() {
	var (
		db                     *postgres.DB
		headerRepo             repositories.HeaderRepository
		blockOne, timestampOne int
		headerOne              core.Header
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)

		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepo)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("all_sin_queue_events", func() {
		It("returns vow fess events", func() {
			fakeEra := strconv.Itoa(timestampOne)
			vowFessLog := test_data.CreateTestLog(headerOne.Id, db)

			vowFessEvent := test_data.VowFessModel
			vowFessEvent.ColumnValues[constants.HeaderFK] = headerOne.Id
			vowFessEvent.ColumnValues[constants.LogFK] = vowFessLog.ID
			vowFessErr := event.PersistModels([]event.InsertionModel{vowFessEvent}, db)
			Expect(vowFessErr).NotTo(HaveOccurred())

			var actualEvents []test_helpers.SinQueueEvent
			getErr := db.Select(&actualEvents, `SELECT era, act FROM api.all_sin_queue_events($1)`, fakeEra)
			Expect(getErr).NotTo(HaveOccurred())

			Expect(actualEvents).To(ConsistOf(
				test_helpers.SinQueueEvent{Era: fakeEra, Act: "fess"},
			))
		})

		It("returns vow flog events", func() {
			vowFlogLog := test_data.CreateTestLog(headerOne.Id, db)

			fakeEra := strconv.Itoa(int(rand.Int31()))
			vowFlogRepo := vow_flog.VowFlogRepository{}
			vowFlogRepo.SetDB(db)
			vowFlogEvent := test_data.VowFlogModel
			vowFlogEvent.ColumnValues["era"] = fakeEra
			vowFlogEvent.ColumnValues[constants.HeaderFK] = headerOne.Id
			vowFlogEvent.ColumnValues[constants.LogFK] = vowFlogLog.ID
			createErr := vowFlogRepo.Create([]shared.InsertionModel{vowFlogEvent})
			Expect(createErr).NotTo(HaveOccurred())

			var actualEvents []test_helpers.SinQueueEvent
			getErr := db.Select(&actualEvents, `SELECT era, act FROM api.all_sin_queue_events($1)`, fakeEra)
			Expect(getErr).NotTo(HaveOccurred())

			Expect(actualEvents).To(ConsistOf(
				test_helpers.SinQueueEvent{Era: fakeEra, Act: "flog"},
			))
		})

		It("returns events from multiple blocks", func() {
			fakeEra := strconv.Itoa(timestampOne)

			vowFessLog := test_data.CreateTestLog(headerOne.Id, db)
			vowFessEvent := test_data.VowFessModel
			vowFessEvent.ColumnValues[constants.HeaderFK] = headerOne.Id
			vowFessEvent.ColumnValues[constants.LogFK] = vowFessLog.ID
			vowFessErr := event.PersistModels([]event.InsertionModel{vowFessEvent}, db)
			Expect(vowFessErr).NotTo(HaveOccurred())

			// New block
			timestampTwo := timestampOne + 1
			headerTwo := createHeader(blockOne+1, timestampTwo, headerRepo)

			vowFlogLog := test_data.CreateTestLog(headerTwo.Id, db)
			vowFlogRepo := vow_flog.VowFlogRepository{}
			vowFlogRepo.SetDB(db)
			vowFlogEvent := test_data.VowFlogModel
			vowFlogEvent.ColumnValues["era"] = fakeEra
			vowFlogEvent.ColumnValues[constants.HeaderFK] = headerTwo.Id
			vowFlogEvent.ColumnValues[constants.LogFK] = vowFlogLog.ID
			createFlogErr := vowFlogRepo.Create([]shared.InsertionModel{vowFlogEvent})
			Expect(createFlogErr).NotTo(HaveOccurred())

			var actualEvents []test_helpers.SinQueueEvent
			getErr := db.Select(&actualEvents, `SELECT era, act FROM api.all_sin_queue_events($1)`, fakeEra)
			Expect(getErr).NotTo(HaveOccurred())

			Expect(actualEvents).To(ConsistOf(
				test_helpers.SinQueueEvent{Era: fakeEra, Act: "fess"},
				test_helpers.SinQueueEvent{Era: fakeEra, Act: "flog"},
			))
		})

		It("ignores sin queue events with irrelevant eras", func() {
			vowFlogLog := test_data.CreateTestLog(headerOne.Id, db)

			rawEra := int(rand.Int31())
			fakeEra := strconv.Itoa(rawEra)
			irrelevantEra := strconv.Itoa(rawEra + 1)

			vowFlogRepo := vow_flog.VowFlogRepository{}
			vowFlogRepo.SetDB(db)
			vowFlogEvent := test_data.VowFlogModel
			vowFlogEvent.ColumnValues["era"] = fakeEra
			vowFlogEvent.ColumnValues[constants.HeaderFK] = headerOne.Id
			vowFlogEvent.ColumnValues[constants.LogFK] = vowFlogLog.ID
			createErr := vowFlogRepo.Create([]shared.InsertionModel{vowFlogEvent})
			Expect(createErr).NotTo(HaveOccurred())

			var actualEvents []test_helpers.SinQueueEvent
			getErr := db.Select(&actualEvents, `SELECT era, act FROM api.all_sin_queue_events($1)`, irrelevantEra)
			Expect(getErr).NotTo(HaveOccurred())

			Expect(actualEvents).To(BeEmpty())
		})

		Describe("result pagination", func() {
			var fakeEra string

			BeforeEach(func() {
				fakeEra = strconv.Itoa(timestampOne)
				logId := test_data.CreateTestLog(headerOne.Id, db).ID

				vowFessEvent := test_data.VowFessModel
				vowFessEvent.ColumnValues[constants.HeaderFK] = headerOne.Id
				vowFessEvent.ColumnValues[constants.LogFK] = logId
				vowFessErr := event.PersistModels([]event.InsertionModel{vowFessEvent}, db)
				Expect(vowFessErr).NotTo(HaveOccurred())

				// New block
				timestampTwo := timestampOne + 1
				headerTwo := createHeader(blockOne+1, timestampTwo, headerRepo)
				logTwoId := test_data.CreateTestLog(headerTwo.Id, db).ID

				vowFlogRepo := vow_flog.VowFlogRepository{}
				vowFlogRepo.SetDB(db)
				vowFlogEvent := test_data.VowFlogModel
				vowFlogEvent.ColumnValues["era"] = fakeEra
				vowFlogEvent.ColumnValues[constants.HeaderFK] = headerTwo.Id
				vowFlogEvent.ColumnValues[constants.LogFK] = logTwoId
				vowFlogErr := vowFlogRepo.Create([]shared.InsertionModel{vowFlogEvent})
				Expect(vowFlogErr).NotTo(HaveOccurred())
			})

			It("limits results to latest block if max_results argument is provided", func() {
				maxResults := 1
				var actualEvents []test_helpers.SinQueueEvent
				err := db.Select(&actualEvents, `SELECT era, act FROM api.all_sin_queue_events($1, $2)`,
					fakeEra, maxResults)
				Expect(err).NotTo(HaveOccurred())

				Expect(actualEvents).To(ConsistOf(
					test_helpers.SinQueueEvent{Era: fakeEra, Act: "flog"},
				))
			})

			It("offsets results if offset is provided", func() {
				maxResults := 1
				resultOffset := 1
				var actualEvents []test_helpers.SinQueueEvent
				err := db.Select(&actualEvents, `SELECT era, act FROM api.all_sin_queue_events($1, $2, $3)`,
					fakeEra, maxResults, resultOffset)
				Expect(err).NotTo(HaveOccurred())

				Expect(actualEvents).To(ConsistOf(
					test_helpers.SinQueueEvent{Era: fakeEra, Act: "fess"},
				))
			})
		})
	})
})
