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
	"database/sql"
	storage_helper "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"math/rand"
	"strconv"
	"time"

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

var _ = Describe("QueuedSin", func() {
	var (
		blockOne, timestampOne int
		fakeEra                string
		fakeTab                = strconv.Itoa(int(rand.Int31()))
		headerOne              core.Header
		headerRepository       repositories.HeaderRepository
		logId                  int64
		rawEra                 int
		sinMappingMetadata     utils.StorageValueMetadata
		vowRepository          vow.VowStorageRepository
		diffID                 int64
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)

		headerRepository = repositories.NewHeaderRepository(db)
		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepository)
		fakeHeaderSyncLog := test_data.CreateTestLog(headerOne.Id, db)
		logId = fakeHeaderSyncLog.ID

		rawEra = int(rand.Int31())
		fakeEra = strconv.Itoa(rawEra)

		diffID = storage_helper.CreateFakeDiffRecord(db)

		vowRepository = vow.VowStorageRepository{}
		vowRepository.SetDB(db)
		sinMappingKeys := map[utils.Key]string{constants.Timestamp: fakeEra}
		sinMappingMetadata = utils.GetStorageValueMetadata(vow.SinMapping, sinMappingKeys, utils.Uint256)
		insertSinMappingErr := vowRepository.Create(diffID, headerOne.Id, sinMappingMetadata, fakeTab)
		Expect(insertSinMappingErr).NotTo(HaveOccurred())
	})

	Describe("getting a single queued sin for an era", func() {
		It("gets queued sin for an era", func() {
			var result QueuedSin
			err := db.Get(&result, `SELECT era, tab, flogged, created, updated from api.get_queued_sin($1)`, fakeEra)
			Expect(err).NotTo(HaveOccurred())

			Expect(result.Era).To(Equal(test_helpers.GetValidNullString(fakeEra)))
			Expect(result.Tab).To(Equal(test_helpers.GetValidNullString(fakeTab)))
			Expect(result.Flogged).To(Equal(sql.NullBool{Bool: false, Valid: true}))
			timestampAsInt, convertErr := strconv.ParseInt(headerOne.Timestamp, 10, 64)
			Expect(convertErr).NotTo(HaveOccurred())
			timestampAsStr := time.Unix(timestampAsInt, 0).UTC().Format(time.RFC3339)
			expectedTimestamp := test_helpers.GetValidNullString(timestampAsStr)
			Expect(result.Created).To(Equal(expectedTimestamp))
			Expect(result.Updated).To(Equal(expectedTimestamp))
		})

		It("returns flogged as true if era has been flogged", func() {
			vowFlogEvent := test_data.VowFlogModel
			vowFlogEvent.ColumnValues[constants.EraColumn] = fakeEra
			vowFlogEvent.ColumnValues[constants.HeaderFK] = headerOne.Id
			vowFlogEvent.ColumnValues[constants.LogFK] = logId
			vowFlogErr := event.PersistModels([]event.InsertionModel{vowFlogEvent}, db)
			Expect(vowFlogErr).NotTo(HaveOccurred())

			var result QueuedSin
			err := db.Get(&result, `SELECT era, tab, flogged, created, updated from api.get_queued_sin($1)`, fakeEra)
			Expect(err).NotTo(HaveOccurred())

			Expect(result.Flogged).To(Equal(sql.NullBool{Bool: true, Valid: true}))
		})

		It("does not return queued sin for another era", func() {
			anotherFakeEra := strconv.Itoa(int(rand.Int31()))
			var result QueuedSin
			err := db.Get(&result, `SELECT era, tab, flogged, created, updated from api.get_queued_sin($1)`, anotherFakeEra)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(BeZero())
		})

		Context("when tab for an era has been updated", func() {
			var (
				anotherFakeTab = "321"
				laterTimestamp string
			)

			BeforeEach(func() {
				timestampTwo := timestampOne + 1
				headerTwo := createHeader(blockOne+1, timestampTwo, headerRepository)
				laterTimestamp = strconv.Itoa(timestampTwo)
				insertVowMappingErr := vowRepository.Create(diffID, headerTwo.Id, sinMappingMetadata, anotherFakeTab)
				Expect(insertVowMappingErr).NotTo(HaveOccurred())
			})

			It("returns most recent 'updated' timestamp", func() {
				var result QueuedSin
				err := db.Get(&result, `SELECT era, tab, flogged, created, updated from api.get_queued_sin($1)`, fakeEra)
				Expect(err).NotTo(HaveOccurred())

				createdTimestampAsInt, convertCreatedErr := strconv.ParseInt(headerOne.Timestamp, 10, 64)
				Expect(convertCreatedErr).NotTo(HaveOccurred())
				expectedCreatedTimestamp := time.Unix(createdTimestampAsInt, 0).UTC().Format(time.RFC3339)
				Expect(result.Created).To(Equal(test_helpers.GetValidNullString(expectedCreatedTimestamp)))

				updatedTimestampAsInt, convertUpdatedErr := strconv.ParseInt(laterTimestamp, 10, 64)
				Expect(convertUpdatedErr).NotTo(HaveOccurred())
				expectedUpdatedTimestamp := time.Unix(updatedTimestampAsInt, 0).UTC().Format(time.RFC3339)
				Expect(result.Updated).To(Equal(test_helpers.GetValidNullString(expectedUpdatedTimestamp)))
			})

			It("returns most recent tab value", func() {
				var result QueuedSin
				err := db.Get(&result, `SELECT era, tab, flogged, created, updated from api.get_queued_sin($1)`, fakeEra)
				Expect(err).NotTo(HaveOccurred())

				Expect(result.Tab).To(Equal(test_helpers.GetValidNullString(anotherFakeTab)))
			})
		})
	})

	Describe("getting all queued sins", func() {
		It("returns queued sin for every era", func() {
			anotherFakeEra := strconv.Itoa(int(rand.Int31()))
			anotherFakeTab := strconv.Itoa(int(rand.Int31()))
			anotherSinMappingKeys := map[utils.Key]string{constants.Timestamp: anotherFakeEra}
			anotherSinMappingMetadata := utils.GetStorageValueMetadata(vow.SinMapping, anotherSinMappingKeys, utils.Uint256)
			insertSinMappingErr := vowRepository.Create(diffID, headerOne.Id, anotherSinMappingMetadata, anotherFakeTab)
			Expect(insertSinMappingErr).NotTo(HaveOccurred())

			var results []QueuedSin
			err := db.Select(&results, `SELECT era, tab, flogged, created, updated from api.all_queued_sin()`)
			Expect(err).NotTo(HaveOccurred())

			Expect(len(results)).To(Equal(2))
			fakeEraNullString := test_helpers.GetValidNullString(fakeEra)
			anotherFakeEraNullString := test_helpers.GetValidNullString(anotherFakeEra)
			Expect(results[0].Era).To(Or(Equal(fakeEraNullString), Equal(anotherFakeEraNullString)))
			fakeTabNullString := test_helpers.GetValidNullString(fakeTab)
			anotherFakeTabNullString := test_helpers.GetValidNullString(anotherFakeTab)
			Expect(results[0].Tab).To(Or(Equal(fakeTabNullString), Equal(anotherFakeTabNullString)))
		})

		Describe("result pagination", func() {
			var laterEra, anotherFakeTab string

			BeforeEach(func() {
				laterEra = strconv.Itoa(rawEra + 1)
				anotherFakeTab = strconv.Itoa(int(rand.Int31()))
				anotherSinMappingKeys := map[utils.Key]string{constants.Timestamp: laterEra}
				anotherSinMappingMetadata := utils.GetStorageValueMetadata(vow.SinMapping, anotherSinMappingKeys, utils.Uint256)

				insertSinMappingErr := vowRepository.Create(diffID, headerOne.Id, anotherSinMappingMetadata, anotherFakeTab)
				Expect(insertSinMappingErr).NotTo(HaveOccurred())
			})

			It("limits results to latest era if max_results argument is provided", func() {
				maxResults := 1
				var results []QueuedSin
				err := db.Select(&results, `SELECT era, tab, flogged, created, updated FROM api.all_queued_sin($1)`,
					maxResults)
				Expect(err).NotTo(HaveOccurred())

				Expect(len(results)).To(Equal(maxResults))
				Expect(results[0].Era).To(Equal(test_helpers.GetValidNullString(laterEra)))
				Expect(results[0].Tab).To(Equal(test_helpers.GetValidNullString(anotherFakeTab)))
			})

			It("offsets results if offset is provided", func() {
				maxResults := 1
				resultOffset := 1
				var results []QueuedSin
				err := db.Select(&results, `SELECT era, tab, flogged, created, updated FROM api.all_queued_sin($1, $2)`,
					maxResults, resultOffset)
				Expect(err).NotTo(HaveOccurred())

				Expect(len(results)).To(Equal(maxResults))
				Expect(results[0].Era).To(Equal(test_helpers.GetValidNullString(fakeEra)))
				Expect(results[0].Tab).To(Equal(test_helpers.GetValidNullString(fakeTab)))
			})
		})
	})
})

type QueuedSin struct {
	Era     sql.NullString
	Tab     sql.NullString
	Flogged sql.NullBool
	Created sql.NullString
	Updated sql.NullString
}
