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

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_frob"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Frobs query", func() {
	var (
		frobRepo               vat_frob.VatFrobRepository
		headerRepo             repositories.HeaderRepository
		fakeIlkHex             = test_helpers.FakeIlk.Hex
		fakeIlkIdentifier      = test_helpers.FakeIlk.Identifier
		fakeUrn                = test_data.RandomString(40)
		blockOne, timestampOne int
		headerOne              core.Header
		diffID                 int64
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		frobRepo = vat_frob.VatFrobRepository{}
		frobRepo.SetDB(db)

		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepo)

		diffID = storage_helper.CreateFakeDiffRecord(db)
	})

	Describe("urn_frobs", func() {
		It("returns frob for ilk/urn combination", func() {
			vatFrobLog := test_data.CreateTestLog(headerOne.Id, db)
			ilkRate := insertIlkRate(fakeIlkHex, diffID, headerOne.Id, db)
			vatFrobEvent := getFakeVatFrobEvent(fakeIlkHex, fakeUrn, headerOne.Id, vatFrobLog.ID)
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
			vatFrobLog := test_data.CreateTestLog(headerOne.Id, db)
			frobBlockOne := getFakeVatFrobEvent(fakeIlkHex, fakeUrn, headerOne.Id, vatFrobLog.ID)
			insertFrobErr := frobRepo.Create([]shared.InsertionModel{frobBlockOne})
			Expect(insertFrobErr).NotTo(HaveOccurred())

			headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
			vatFrobLogTwo := test_data.CreateTestLog(headerTwo.Id, db)
			frobBlockTwo := getFakeVatFrobEvent(fakeIlkHex, fakeUrn, headerTwo.Id, vatFrobLogTwo.ID)
			insertFrobTwoErr := frobRepo.Create([]shared.InsertionModel{frobBlockTwo})
			Expect(insertFrobTwoErr).NotTo(HaveOccurred())

			ilkRate := insertIlkRate(fakeIlkHex, diffID, headerOne.Id, db)

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
				logId := test_data.CreateTestLog(headerOne.Id, db).ID
				frobBlockOne = getFakeVatFrobEvent(fakeIlkHex, fakeUrn, headerOne.Id, logId)
				insertFrobErr := frobRepo.Create([]shared.InsertionModel{frobBlockOne})
				Expect(insertFrobErr).NotTo(HaveOccurred())

				headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
				logTwoId := test_data.CreateTestLog(headerTwo.Id, db).ID
				frobBlockTwo = getFakeVatFrobEvent(fakeIlkHex, fakeUrn, headerTwo.Id, logTwoId)
				insertFrobTwoErr := frobRepo.Create([]shared.InsertionModel{frobBlockTwo})
				Expect(insertFrobTwoErr).NotTo(HaveOccurred())

				ilkRate = insertIlkRate(fakeIlkHex, diffID, headerOne.Id, db)
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
			vatFrobLog := test_data.CreateTestLog(headerOne.Id, db)
			vatFrobEvent := getFakeVatFrobEvent(fakeIlkHex, fakeUrn, headerOne.Id, vatFrobLog.ID)
			vatFrobLogTwo := test_data.CreateTestLog(headerOne.Id, db)
			irrelevantVatFrobEvent := getFakeVatFrobEvent(fakeIlkHex, test_data.RandomString(40), headerOne.Id, vatFrobLogTwo.ID)
			createFrobsErr := frobRepo.Create([]shared.InsertionModel{vatFrobEvent, irrelevantVatFrobEvent})
			Expect(createFrobsErr).NotTo(HaveOccurred())

			ilkRate := insertIlkRate(fakeIlkHex, diffID, headerOne.Id, db)

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
			vatFrobLog := test_data.CreateTestLog(headerOne.Id, db)
			ilkRateBlockOne := insertIlkRate(fakeIlkHex, diffID, headerOne.Id, db)
			frobBlockOne := getFakeVatFrobEvent(fakeIlkHex, fakeUrn, headerOne.Id, vatFrobLog.ID)
			insertFrobOneErr := frobRepo.Create([]shared.InsertionModel{frobBlockOne})
			Expect(insertFrobOneErr).NotTo(HaveOccurred())

			headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
			vatFrobLogTwo := test_data.CreateTestLog(headerTwo.Id, db)
			ilkRateBlockTwo := insertIlkRate(fakeIlkHex, diffID, headerTwo.Id, db)
			frobBlockTwo := getFakeVatFrobEvent(fakeIlkHex, fakeUrn, headerTwo.Id, vatFrobLogTwo.ID)
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
			ilkRateBlockOne := insertIlkRate(fakeIlkHex, diffID, headerOne.Id, db)

			headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
			vatFrobLog := test_data.CreateTestLog(headerTwo.Id, db)
			frobBlockTwo := getFakeVatFrobEvent(fakeIlkHex, fakeUrn, headerTwo.Id, vatFrobLog.ID)
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
			vatFrobLog := test_data.CreateTestLog(headerOne.Id, db)
			frobOne := getFakeVatFrobEvent(fakeIlkHex, fakeUrn, headerOne.Id, vatFrobLog.ID)
			anotherFakeUrn := test_data.RandomString(40)
			anotherVatFrobLog := test_data.CreateTestLog(headerOne.Id, db)
			frobTwo := getFakeVatFrobEvent(fakeIlkHex, anotherFakeUrn, headerOne.Id, anotherVatFrobLog.ID)
			insertFrobsErr := frobRepo.Create([]shared.InsertionModel{frobOne, frobTwo})
			Expect(insertFrobsErr).NotTo(HaveOccurred())
			ilkRate := insertIlkRate(fakeIlkHex, diffID, headerOne.Id, db)

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
				frobOne, frobTwo shared.InsertionModel
				anotherFakeUrn   string
				headerTwo        core.Header
			)

			BeforeEach(func() {
				logId := test_data.CreateTestLog(headerOne.Id, db).ID
				frobOne = getFakeVatFrobEvent(fakeIlkHex, fakeUrn, headerOne.Id, logId)
				insertFrobOneErr := frobRepo.Create([]shared.InsertionModel{frobOne})
				Expect(insertFrobOneErr).NotTo(HaveOccurred())

				anotherFakeUrn = test_data.RandomString(40)
				headerTwo = createHeader(blockOne+1, timestampOne+1, headerRepo)
				logTwoId := test_data.CreateTestLog(headerTwo.Id, db).ID
				frobTwo = getFakeVatFrobEvent(fakeIlkHex, anotherFakeUrn, headerTwo.Id, logTwoId)
				insertFrobTwoErr := frobRepo.Create([]shared.InsertionModel{frobTwo})
				Expect(insertFrobTwoErr).NotTo(HaveOccurred())
			})

			It("limits results if max_results argument is provided", func() {
				ilkRateTwo := insertIlkRate(fakeIlkHex, diffID, headerTwo.Id, db)

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
				ilkRateOne := insertIlkRate(fakeIlkHex, diffID, headerOne.Id, db)

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
			vatFrobLog := test_data.CreateTestLog(headerOne.Id, db)
			frobOne := getFakeVatFrobEvent(fakeIlkHex, fakeUrn, headerOne.Id, vatFrobLog.ID)
			insertFrobOneErr := frobRepo.Create([]shared.InsertionModel{frobOne})
			Expect(insertFrobOneErr).NotTo(HaveOccurred())

			ilkRate := insertIlkRate(fakeIlkHex, diffID, headerOne.Id, db)

			headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
			anotherVatFrobLog := test_data.CreateTestLog(headerTwo.Id, db)
			anotherFakeUrn := test_data.RandomString(40)
			frobTwo := getFakeVatFrobEvent(fakeIlkHex, anotherFakeUrn, headerTwo.Id, anotherVatFrobLog.ID)
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
			vatFrobLog := test_data.CreateTestLog(headerOne.Id, db)
			frobOne := getFakeVatFrobEvent(fakeIlkHex, fakeUrn, headerOne.Id, vatFrobLog.ID)
			insertFrobOneErr := frobRepo.Create([]shared.InsertionModel{frobOne})
			Expect(insertFrobOneErr).NotTo(HaveOccurred())

			ilkRateBlockOne := insertIlkRate(fakeIlkHex, diffID, headerOne.Id, db)

			headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
			anotherVatFrobLog := test_data.CreateTestLog(headerTwo.Id, db)
			anotherFakeUrn := test_data.RandomString(40)
			frobTwo := getFakeVatFrobEvent(fakeIlkHex, anotherFakeUrn, headerTwo.Id, anotherVatFrobLog.ID)
			insertFrobTwoErr := frobRepo.Create([]shared.InsertionModel{frobTwo})
			Expect(insertFrobTwoErr).NotTo(HaveOccurred())

			ilkRateBlockTwo := insertIlkRate(fakeIlkHex, diffID, headerTwo.Id, db)

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
			ilkRate := insertIlkRate(fakeIlkHex, diffID, headerOne.Id, db)

			headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
			vatFrobLog := test_data.CreateTestLog(headerTwo.Id, db)
			vatFrobEvent := getFakeVatFrobEvent(fakeIlkHex, fakeUrn, headerTwo.Id, vatFrobLog.ID)

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

func insertIlkRate(ilk string, diffID, headerID int64, db *postgres.DB) int {
	ilkID, ilkErr := shared.GetOrCreateIlk(ilk, db)
	Expect(ilkErr).NotTo(HaveOccurred())
	ilkRate := rand.Int()
	_, insertIlkRateErr := db.Exec(`INSERT INTO maker.vat_ilk_rate (diff_id, header_id, ilk_id, rate) VALUES ($1, $2, $3, $4)`, diffID, headerID, ilkID, ilkRate)
	Expect(insertIlkRateErr).NotTo(HaveOccurred())
	return ilkRate
}
