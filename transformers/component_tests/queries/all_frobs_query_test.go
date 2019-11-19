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

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_frob"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Frobs query", func() {
	var (
		db                *postgres.DB
		frobRepo          vat_frob.VatFrobRepository
		headerRepo        repositories.HeaderRepository
		fakeIlkHex        = test_helpers.FakeIlk.Hex
		fakeIlkIdentifier = test_helpers.FakeIlk.Identifier
		fakeUrn           = test_data.RandomString(40)
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		frobRepo = vat_frob.VatFrobRepository{}
		frobRepo.SetDB(db)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("urn_frobs", func() {
		It("returns frob for ilk/urn combination", func() {
			headerID, blockNumber := insertHeader(rand.Int63(), headerRepo)
			vatFrobLog := test_data.CreateTestLog(headerID, db)
			ilkRate := insertIlkRate(fakeIlkHex, blockNumber, db)
			vatFrobEvent := getFakeVatFrobEvent(fakeIlkHex, fakeUrn, headerID, vatFrobLog.ID)
			insertFrobErr := frobRepo.Create([]shared.InsertionModel{vatFrobEvent})
			Expect(insertFrobErr).NotTo(HaveOccurred())

			var actualFrobs []test_helpers.FrobEvent
			getFrobsErr := db.Select(&actualFrobs, `SELECT ilk_identifier, urn_identifier, dink, dart, ilk_rate FROM api.urn_frobs($1, $2)`, fakeIlkIdentifier, fakeUrn)
			Expect(getFrobsErr).NotTo(HaveOccurred())

			Expect(actualFrobs).To(ConsistOf(
				test_helpers.FrobEvent{
					IlkIdentifier: fakeIlkIdentifier,
					UrnIdentifier: fakeUrn,
					Dink:          vatFrobEvent.ColumnValues["dink"].(string),
					Dart:          vatFrobEvent.ColumnValues["dart"].(string),
					Rate:          strconv.Itoa(ilkRate),
				},
			))
		})

		It("returns matching frobs from multiple blocks", func() {
			headerOneId, headerOneBlockNumber := insertHeader(rand.Int63(), headerRepo)
			vatFrobLog := test_data.CreateTestLog(headerOneId, db)
			frobBlockOne := getFakeVatFrobEvent(fakeIlkHex, fakeUrn, headerOneId, vatFrobLog.ID)
			insertFrobErr := frobRepo.Create([]shared.InsertionModel{frobBlockOne})
			Expect(insertFrobErr).NotTo(HaveOccurred())

			headerTwoId, _ := insertHeader(headerOneBlockNumber+1, headerRepo)
			vatFrobLogTwo := test_data.CreateTestLog(headerTwoId, db)
			frobBlockTwo := getFakeVatFrobEvent(fakeIlkHex, fakeUrn, headerTwoId, vatFrobLogTwo.ID)
			insertFrobTwoErr := frobRepo.Create([]shared.InsertionModel{frobBlockTwo})
			Expect(insertFrobTwoErr).NotTo(HaveOccurred())

			ilkRate := insertIlkRate(fakeIlkHex, headerOneBlockNumber, db)

			var actualFrobs []test_helpers.FrobEvent
			getFrobsErr := db.Select(&actualFrobs, `SELECT ilk_identifier, urn_identifier, dink, dart, ilk_rate FROM api.urn_frobs($1, $2)`, fakeIlkIdentifier, fakeUrn)
			Expect(getFrobsErr).NotTo(HaveOccurred())

			Expect(actualFrobs).To(ConsistOf(
				test_helpers.FrobEvent{
					IlkIdentifier: fakeIlkIdentifier,
					UrnIdentifier: fakeUrn,
					Dink:          frobBlockOne.ColumnValues["dink"].(string),
					Dart:          frobBlockOne.ColumnValues["dart"].(string),
					Rate:          strconv.Itoa(ilkRate),
				},
				test_helpers.FrobEvent{
					IlkIdentifier: fakeIlkIdentifier,
					UrnIdentifier: fakeUrn,
					Dink:          frobBlockTwo.ColumnValues["dink"].(string),
					Dart:          frobBlockTwo.ColumnValues["dart"].(string),
					Rate:          strconv.Itoa(ilkRate),
				},
			))
		})

		Describe("result pagination", func() {
			var (
				frobBlockOne, frobBlockTwo shared.InsertionModel
				ilkRate                    int
			)

			BeforeEach(func() {
				headerOneId, headerOneBlockNumber := insertHeader(rand.Int63(), headerRepo)
				logId := test_data.CreateTestLog(headerOneId, db).ID
				frobBlockOne = getFakeVatFrobEvent(fakeIlkHex, fakeUrn, headerOneId, logId)
				insertFrobErr := frobRepo.Create([]shared.InsertionModel{frobBlockOne})
				Expect(insertFrobErr).NotTo(HaveOccurred())

				headerTwoId, _ := insertHeader(headerOneBlockNumber+1, headerRepo)
				logTwoId := test_data.CreateTestLog(headerTwoId, db).ID
				frobBlockTwo = getFakeVatFrobEvent(fakeIlkHex, fakeUrn, headerTwoId, logTwoId)
				insertFrobTwoErr := frobRepo.Create([]shared.InsertionModel{frobBlockTwo})
				Expect(insertFrobTwoErr).NotTo(HaveOccurred())

				ilkRate = insertIlkRate(fakeIlkHex, headerOneBlockNumber, db)
			})

			It("limits results to latest blocks if max_results argument is provided", func() {
				maxResults := 1
				var actualFrobs []test_helpers.FrobEvent
				getFrobsErr := db.Select(&actualFrobs, `SELECT ilk_identifier, urn_identifier, dink, dart, ilk_rate FROM api.urn_frobs($1, $2, $3)`,
					fakeIlkIdentifier, fakeUrn, maxResults)
				Expect(getFrobsErr).NotTo(HaveOccurred())

				Expect(actualFrobs).To(ConsistOf(
					test_helpers.FrobEvent{
						IlkIdentifier: fakeIlkIdentifier,
						UrnIdentifier: fakeUrn,
						Dink:          frobBlockTwo.ColumnValues["dink"].(string),
						Dart:          frobBlockTwo.ColumnValues["dart"].(string),
						Rate:          strconv.Itoa(ilkRate),
					},
				))
			})

			It("offsets results if offset is provided", func() {
				maxResults := 1
				resultOffset := 1
				var actualFrobs []test_helpers.FrobEvent
				getFrobsErr := db.Select(&actualFrobs, `SELECT ilk_identifier, urn_identifier, dink, dart, ilk_rate FROM api.urn_frobs($1, $2, $3, $4)`,
					fakeIlkIdentifier, fakeUrn, maxResults, resultOffset)
				Expect(getFrobsErr).NotTo(HaveOccurred())

				Expect(actualFrobs).To(ConsistOf(
					test_helpers.FrobEvent{
						IlkIdentifier: fakeIlkIdentifier,
						UrnIdentifier: fakeUrn,
						Dink:          frobBlockOne.ColumnValues["dink"].(string),
						Dart:          frobBlockOne.ColumnValues["dart"].(string),
						Rate:          strconv.Itoa(ilkRate),
					},
				))
			})
		})

		It("does not include frobs for a different urn", func() {
			headerID, blockNumber := insertHeader(rand.Int63(), headerRepo)
			vatFrobLog := test_data.CreateTestLog(headerID, db)
			vatFrobEvent := getFakeVatFrobEvent(fakeIlkHex, fakeUrn, headerID, vatFrobLog.ID)
			vatFrobLogTwo := test_data.CreateTestLog(headerID, db)
			irrelevantVatFrobEvent := getFakeVatFrobEvent(fakeIlkHex, test_data.RandomString(40), headerID, vatFrobLogTwo.ID)
			createFrobsErr := frobRepo.Create([]shared.InsertionModel{vatFrobEvent, irrelevantVatFrobEvent})
			Expect(createFrobsErr).NotTo(HaveOccurred())

			ilkRate := insertIlkRate(fakeIlkHex, blockNumber, db)

			var actualFrobs []test_helpers.FrobEvent
			getFrobsErr := db.Select(&actualFrobs, `SELECT ilk_identifier, urn_identifier, dink, dart, ilk_rate FROM api.urn_frobs($1, $2)`, fakeIlkIdentifier, fakeUrn)
			Expect(getFrobsErr).NotTo(HaveOccurred())

			Expect(actualFrobs).To(ConsistOf(
				test_helpers.FrobEvent{
					IlkIdentifier: fakeIlkIdentifier,
					UrnIdentifier: fakeUrn,
					Dink:          vatFrobEvent.ColumnValues["dink"].(string),
					Dart:          vatFrobEvent.ColumnValues["dart"].(string),
					Rate:          strconv.Itoa(ilkRate),
				},
			))
		})

		It("provides most recent rate for each frob across blocks", func() {
			blockNumber := rand.Int63()
			headerOneId, headerOneBlockNumber := insertHeader(blockNumber, headerRepo)
			vatFrobLog := test_data.CreateTestLog(headerOneId, db)
			ilkRateBlockOne := insertIlkRate(fakeIlkHex, headerOneBlockNumber, db)
			frobBlockOne := getFakeVatFrobEvent(fakeIlkHex, fakeUrn, headerOneId, vatFrobLog.ID)
			insertFrobOneErr := frobRepo.Create([]shared.InsertionModel{frobBlockOne})
			Expect(insertFrobOneErr).NotTo(HaveOccurred())

			headerTwoId, headerTwoBlockNumber := insertHeader(blockNumber+1, headerRepo)
			vatFrobLogTwo := test_data.CreateTestLog(headerTwoId, db)
			ilkRateBlockTwo := insertIlkRate(fakeIlkHex, headerTwoBlockNumber, db)
			frobBlockTwo := getFakeVatFrobEvent(fakeIlkHex, fakeUrn, headerTwoId, vatFrobLogTwo.ID)
			insertFrobTwoErr := frobRepo.Create([]shared.InsertionModel{frobBlockTwo})
			Expect(insertFrobTwoErr).NotTo(HaveOccurred())

			var actualFrobs []test_helpers.FrobEvent
			getFrobsErr := db.Select(&actualFrobs, `SELECT ilk_identifier, urn_identifier, dink, dart, ilk_rate FROM api.urn_frobs($1, $2)`, fakeIlkIdentifier, fakeUrn)
			Expect(getFrobsErr).NotTo(HaveOccurred())

			Expect(actualFrobs).To(ConsistOf(
				test_helpers.FrobEvent{
					IlkIdentifier: fakeIlkIdentifier,
					UrnIdentifier: fakeUrn,
					Dink:          frobBlockOne.ColumnValues["dink"].(string),
					Dart:          frobBlockOne.ColumnValues["dart"].(string),
					Rate:          strconv.Itoa(ilkRateBlockOne),
				},
				test_helpers.FrobEvent{
					IlkIdentifier: fakeIlkIdentifier,
					UrnIdentifier: fakeUrn,
					Dink:          frobBlockTwo.ColumnValues["dink"].(string),
					Dart:          frobBlockTwo.ColumnValues["dart"].(string),
					Rate:          strconv.Itoa(ilkRateBlockTwo),
				},
			))
		})

		It("gets the most recent ilk rate when not updated in the same block as frob", func() {
			_, headerOneBlockNumber := insertHeader(rand.Int63(), headerRepo)
			ilkRateBlockOne := insertIlkRate(fakeIlkHex, headerOneBlockNumber, db)

			headerTwoID, _ := insertHeader(headerOneBlockNumber+1, headerRepo)
			vatFrobLog := test_data.CreateTestLog(headerTwoID, db)
			frobBlockTwo := getFakeVatFrobEvent(fakeIlkHex, fakeUrn, headerTwoID, vatFrobLog.ID)
			insertFrobErr := frobRepo.Create([]shared.InsertionModel{frobBlockTwo})
			Expect(insertFrobErr).NotTo(HaveOccurred())

			var actualFrobs []test_helpers.FrobEvent
			getFrobsErr := db.Select(&actualFrobs, `SELECT ilk_identifier, urn_identifier, dink, dart, ilk_rate FROM api.urn_frobs($1, $2)`, fakeIlkIdentifier, fakeUrn)
			Expect(getFrobsErr).NotTo(HaveOccurred())

			Expect(actualFrobs).To(ConsistOf(
				test_helpers.FrobEvent{
					IlkIdentifier: fakeIlkIdentifier,
					UrnIdentifier: fakeUrn,
					Dink:          frobBlockTwo.ColumnValues["dink"].(string),
					Dart:          frobBlockTwo.ColumnValues["dart"].(string),
					Rate:          strconv.Itoa(ilkRateBlockOne),
				},
			))
		})

		It("fails if no argument is supplied (STRICT)", func() {
			_, err := db.Exec(`SELECT * FROM api.urn_frobs()`)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("function api.urn_frobs() does not exist"))
		})

		It("fails if only one argument is supplied (STRICT)", func() {
			_, err := db.Exec(`SELECT * FROM api.urn_frobs($1::text)`, fakeIlkIdentifier)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("function api.urn_frobs(text) does not exist"))
		})
	})

	Describe("all_frobs", func() {
		It("returns frobs for all urns on an ilk", func() {
			headerID, blockNumber := insertHeader(rand.Int63(), headerRepo)
			vatFrobLog := test_data.CreateTestLog(headerID, db)
			frobOne := getFakeVatFrobEvent(fakeIlkHex, fakeUrn, headerID, vatFrobLog.ID)
			anotherFakeUrn := test_data.RandomString(40)
			anotherVatFrobLog := test_data.CreateTestLog(headerID, db)
			frobTwo := getFakeVatFrobEvent(fakeIlkHex, anotherFakeUrn, headerID, anotherVatFrobLog.ID)
			insertFrobsErr := frobRepo.Create([]shared.InsertionModel{frobOne, frobTwo})
			Expect(insertFrobsErr).NotTo(HaveOccurred())
			ilkRate := insertIlkRate(fakeIlkHex, blockNumber, db)

			var actualFrobs []test_helpers.FrobEvent
			getFrobsErr := db.Select(&actualFrobs, `SELECT ilk_identifier, urn_identifier, dink, dart, ilk_rate FROM api.all_frobs($1)`, fakeIlkIdentifier)
			Expect(getFrobsErr).NotTo(HaveOccurred())

			Expect(actualFrobs).To(ConsistOf(
				test_helpers.FrobEvent{
					IlkIdentifier: fakeIlkIdentifier,
					UrnIdentifier: fakeUrn,
					Dink:          frobOne.ColumnValues["dink"].(string),
					Dart:          frobOne.ColumnValues["dart"].(string),
					Rate:          strconv.Itoa(ilkRate),
				},
				test_helpers.FrobEvent{
					IlkIdentifier: fakeIlkIdentifier,
					UrnIdentifier: anotherFakeUrn,
					Dink:          frobTwo.ColumnValues["dink"].(string),
					Dart:          frobTwo.ColumnValues["dart"].(string),
					Rate:          strconv.Itoa(ilkRate),
				},
			))
		})

		Describe("result pagination", func() {
			var (
				blockOne, blockTwo int64
				frobOne, frobTwo   shared.InsertionModel
				anotherFakeUrn     string
			)

			BeforeEach(func() {
				var headerOneId int64
				headerOneId, blockOne = insertHeader(rand.Int63(), headerRepo)
				logId := test_data.CreateTestLog(headerOneId, db).ID
				frobOne = getFakeVatFrobEvent(fakeIlkHex, fakeUrn, headerOneId, logId)
				insertFrobOneErr := frobRepo.Create([]shared.InsertionModel{frobOne})
				Expect(insertFrobOneErr).NotTo(HaveOccurred())

				anotherFakeUrn = test_data.RandomString(40)
				var headerTwoId int64
				headerTwoId, blockTwo = insertHeader(blockOne+1, headerRepo)
				logTwoId := test_data.CreateTestLog(headerTwoId, db).ID
				frobTwo = getFakeVatFrobEvent(fakeIlkHex, anotherFakeUrn, headerTwoId, logTwoId)
				insertFrobTwoErr := frobRepo.Create([]shared.InsertionModel{frobTwo})
				Expect(insertFrobTwoErr).NotTo(HaveOccurred())
			})

			It("limits results if max_results argument is provided", func() {
				ilkRateTwo := insertIlkRate(fakeIlkHex, blockTwo, db)

				maxResults := 1
				var actualFrobs []test_helpers.FrobEvent
				getFrobsErr := db.Select(&actualFrobs, `SELECT ilk_identifier, urn_identifier, dink, dart, ilk_rate FROM api.all_frobs($1, $2)`,
					fakeIlkIdentifier, maxResults)
				Expect(getFrobsErr).NotTo(HaveOccurred())

				Expect(actualFrobs).To(Equal(
					[]test_helpers.FrobEvent{
						{
							IlkIdentifier: fakeIlkIdentifier,
							UrnIdentifier: anotherFakeUrn,
							Dink:          frobTwo.ColumnValues["dink"].(string),
							Dart:          frobTwo.ColumnValues["dart"].(string),
							Rate:          strconv.Itoa(ilkRateTwo),
						},
					},
				))
			})

			It("offsets results if offset is provided", func() {
				ilkRateOne := insertIlkRate(fakeIlkHex, blockOne, db)

				maxResults := 1
				resultOffset := 1
				var actualFrobs []test_helpers.FrobEvent
				getFrobsErr := db.Select(&actualFrobs, `SELECT ilk_identifier, urn_identifier, dink, dart, ilk_rate FROM api.all_frobs($1, $2, $3)`,
					fakeIlkIdentifier, maxResults, resultOffset)
				Expect(getFrobsErr).NotTo(HaveOccurred())

				Expect(actualFrobs).To(Equal(
					[]test_helpers.FrobEvent{
						{
							IlkIdentifier: fakeIlkIdentifier,
							UrnIdentifier: fakeUrn,
							Dink:          frobOne.ColumnValues["dink"].(string),
							Dart:          frobOne.ColumnValues["dart"].(string),
							Rate:          strconv.Itoa(ilkRateOne),
						},
					},
				))
			})
		})

		It("returns frobs across multiple blocks", func() {
			headerOneID, headerOneBlockNumber := insertHeader(rand.Int63(), headerRepo)
			vatFrobLog := test_data.CreateTestLog(headerOneID, db)
			frobOne := getFakeVatFrobEvent(fakeIlkHex, fakeUrn, headerOneID, vatFrobLog.ID)
			insertFrobOneErr := frobRepo.Create([]shared.InsertionModel{frobOne})
			Expect(insertFrobOneErr).NotTo(HaveOccurred())

			ilkRate := insertIlkRate(fakeIlkHex, headerOneBlockNumber, db)

			headerTwoID, _ := insertHeader(headerOneBlockNumber+1, headerRepo)
			anotherVatFrobLog := test_data.CreateTestLog(headerTwoID, db)
			anotherFakeUrn := test_data.RandomString(40)
			frobTwo := getFakeVatFrobEvent(fakeIlkHex, anotherFakeUrn, headerTwoID, anotherVatFrobLog.ID)
			insertFrobTwoErr := frobRepo.Create([]shared.InsertionModel{frobTwo})
			Expect(insertFrobTwoErr).NotTo(HaveOccurred())

			var actualFrobs []test_helpers.FrobEvent
			getFrobsErr := db.Select(&actualFrobs, `SELECT ilk_identifier, urn_identifier, dink, dart, ilk_rate FROM api.all_frobs($1)`, fakeIlkIdentifier)
			Expect(getFrobsErr).NotTo(HaveOccurred())

			Expect(actualFrobs).To(ConsistOf(
				test_helpers.FrobEvent{
					IlkIdentifier: fakeIlkIdentifier,
					UrnIdentifier: fakeUrn,
					Dink:          frobOne.ColumnValues["dink"].(string),
					Dart:          frobOne.ColumnValues["dart"].(string),
					Rate:          strconv.Itoa(ilkRate),
				},
				test_helpers.FrobEvent{
					IlkIdentifier: fakeIlkIdentifier,
					UrnIdentifier: anotherFakeUrn,
					Dink:          frobTwo.ColumnValues["dink"].(string),
					Dart:          frobTwo.ColumnValues["dart"].(string),
					Rate:          strconv.Itoa(ilkRate),
				},
			))
		})

		It("provides most recent rate for each block", func() {
			headerOneID, headerOneBlockNumber := insertHeader(rand.Int63(), headerRepo)
			vatFrobLog := test_data.CreateTestLog(headerOneID, db)
			frobOne := getFakeVatFrobEvent(fakeIlkHex, fakeUrn, headerOneID, vatFrobLog.ID)
			insertFrobOneErr := frobRepo.Create([]shared.InsertionModel{frobOne})
			Expect(insertFrobOneErr).NotTo(HaveOccurred())

			ilkRateBlockOne := insertIlkRate(fakeIlkHex, headerOneBlockNumber, db)

			headerTwoID, headerTwoBlockNumber := insertHeader(headerOneBlockNumber+1, headerRepo)
			anotherVatFrobLog := test_data.CreateTestLog(headerTwoID, db)
			anotherFakeUrn := test_data.RandomString(40)
			frobTwo := getFakeVatFrobEvent(fakeIlkHex, anotherFakeUrn, headerTwoID, anotherVatFrobLog.ID)
			insertFrobTwoErr := frobRepo.Create([]shared.InsertionModel{frobTwo})
			Expect(insertFrobTwoErr).NotTo(HaveOccurred())

			ilkRateBlockTwo := insertIlkRate(fakeIlkHex, headerTwoBlockNumber, db)

			var actualFrobs []test_helpers.FrobEvent
			getFrobsErr := db.Select(&actualFrobs, `SELECT ilk_identifier, urn_identifier, dink, dart, ilk_rate FROM api.all_frobs($1)`, fakeIlkIdentifier)
			Expect(getFrobsErr).NotTo(HaveOccurred())

			Expect(actualFrobs).To(ConsistOf(
				test_helpers.FrobEvent{
					IlkIdentifier: fakeIlkIdentifier,
					UrnIdentifier: fakeUrn,
					Dink:          frobOne.ColumnValues["dink"].(string),
					Dart:          frobOne.ColumnValues["dart"].(string),
					Rate:          strconv.Itoa(ilkRateBlockOne),
				},
				test_helpers.FrobEvent{
					IlkIdentifier: fakeIlkIdentifier,
					UrnIdentifier: anotherFakeUrn,
					Dink:          frobTwo.ColumnValues["dink"].(string),
					Dart:          frobTwo.ColumnValues["dart"].(string),
					Rate:          strconv.Itoa(ilkRateBlockTwo),
				},
			))
		})

		It("gets most recent rate when rate not updated in same block", func() {
			_, headerOneBlockNumber := insertHeader(rand.Int63(), headerRepo)
			ilkRate := insertIlkRate(fakeIlkHex, headerOneBlockNumber, db)

			headerTwoID, _ := insertHeader(headerOneBlockNumber+1, headerRepo)
			vatFrobLog := test_data.CreateTestLog(headerTwoID, db)
			vatFrobEvent := getFakeVatFrobEvent(fakeIlkHex, fakeUrn, headerTwoID, vatFrobLog.ID)

			insertFrobsErr := frobRepo.Create([]shared.InsertionModel{vatFrobEvent})
			Expect(insertFrobsErr).NotTo(HaveOccurred())

			var actualFrobs []test_helpers.FrobEvent
			err := db.Select(&actualFrobs, `SELECT ilk_identifier, urn_identifier, dink, dart, ilk_rate FROM api.all_frobs($1)`, fakeIlkIdentifier)
			Expect(err).NotTo(HaveOccurred())

			Expect(actualFrobs).To(ConsistOf(
				test_helpers.FrobEvent{
					IlkIdentifier: fakeIlkIdentifier,
					UrnIdentifier: fakeUrn,
					Dink:          vatFrobEvent.ColumnValues["dink"].(string),
					Dart:          vatFrobEvent.ColumnValues["dart"].(string),
					Rate:          strconv.Itoa(ilkRate),
				},
			))
		})

		It("fails if no argument is supplied (STRICT)", func() {
			_, err := db.Exec(`SELECT * FROM api.all_frobs()`)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("function api.all_frobs() does not exist"))
		})
	})
})

func getFakeVatFrobEvent(ilk, urn string, headerID, logID int64) shared.InsertionModel {
	vatFrobEvent := test_data.VatFrobModelWithPositiveDart()
	vatFrobEvent.ForeignKeyValues[constants.IlkFK] = ilk
	vatFrobEvent.ForeignKeyValues[constants.UrnFK] = urn
	vatFrobEvent.ColumnValues["dink"] = strconv.Itoa(rand.Int())
	vatFrobEvent.ColumnValues["dart"] = strconv.Itoa(rand.Int())
	vatFrobEvent.ColumnValues[constants.HeaderFK] = headerID
	vatFrobEvent.ColumnValues[constants.LogFK] = logID
	return vatFrobEvent
}

func insertIlkRate(ilk string, blockNumber int64, db *postgres.DB) int {
	ilkID, ilkErr := shared.GetOrCreateIlk(ilk, db)
	Expect(ilkErr).NotTo(HaveOccurred())
	ilkRate := rand.Int()
	_, insertIlkRateErr := db.Exec(`INSERT INTO maker.vat_ilk_rate (block_number, ilk_id, rate) VALUES ($1, $2, $3)`, blockNumber, ilkID, ilkRate)
	Expect(insertIlkRateErr).NotTo(HaveOccurred())
	return ilkRate
}

func insertHeader(blockNumber int64, headerRepo repositories.HeaderRepository) (int64, int64) {
	header := fakes.GetFakeHeader(blockNumber)
	headerID, insertHeaderErr := headerRepo.CreateOrUpdateHeader(header)
	Expect(insertHeaderErr).NotTo(HaveOccurred())
	return headerID, blockNumber
}
