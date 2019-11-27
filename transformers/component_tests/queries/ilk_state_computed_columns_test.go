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

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/spot_file/mat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_file/ilk"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_frob"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ilk state computed columns", func() {
	var (
		blockOne, timestampOne int
		db                     *postgres.DB
		fakeGuy                = fakes.RandomString(42)
		headerOne              core.Header
		headerRepository       repositories.HeaderRepository
		logId                  int64
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)

		headerRepository = repositories.NewHeaderRepository(db)
		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepository)
		fakeHeaderSyncLog := test_data.CreateTestLog(headerOne.Id, db)
		logId = fakeHeaderSyncLog.ID

		ilkValues := test_helpers.GetIlkValues(0)
		test_helpers.CreateIlk(db, headerOne, ilkValues, test_helpers.FakeIlkVatMetadatas,
			test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas, test_helpers.FakeIlkSpotMetadatas)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("ilk_state_frobs", func() {
		It("returns relevant frobs for an ilk_state", func() {
			frobRepo := vat_frob.VatFrobRepository{}
			frobRepo.SetDB(db)
			frobEvent := test_data.VatFrobModelWithPositiveDart()
			frobEvent.ForeignKeyValues[constants.UrnFK] = fakeGuy
			frobEvent.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
			frobEvent.ColumnValues[constants.HeaderFK] = headerOne.Id
			frobEvent.ColumnValues[constants.LogFK] = logId
			insertFrobErr := frobRepo.Create([]shared.InsertionModel{frobEvent})
			Expect(insertFrobErr).NotTo(HaveOccurred())

			var actualFrobs []test_helpers.FrobEvent
			getFrobsErr := db.Select(&actualFrobs,
				`SELECT ilk_identifier, urn_identifier, dink, dart FROM api.ilk_state_frobs(
					(SELECT (ilk_identifier, block_height, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated)::api.ilk_state
					 FROM api.get_ilk($1, $2)))`,
				test_helpers.FakeIlk.Identifier, blockOne)
			Expect(getFrobsErr).NotTo(HaveOccurred())

			expectedFrobs := []test_helpers.FrobEvent{{
				IlkIdentifier: test_helpers.FakeIlk.Identifier,
				UrnIdentifier: fakeGuy,
				Dink:          frobEvent.ColumnValues["dink"].(string),
				Dart:          frobEvent.ColumnValues["dart"].(string),
			}}

			Expect(actualFrobs).To(Equal(expectedFrobs))
		})

		Describe("result pagination", func() {
			var (
				headerTwo        core.Header
				oldFrob, newFrob shared.InsertionModel
			)

			BeforeEach(func() {
				frobRepo := vat_frob.VatFrobRepository{}
				frobRepo.SetDB(db)
				oldFrob = test_data.VatFrobModelWithPositiveDart()
				oldFrob.ForeignKeyValues[constants.UrnFK] = fakeGuy
				oldFrob.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
				oldFrob.ColumnValues[constants.HeaderFK] = headerOne.Id
				oldFrob.ColumnValues[constants.LogFK] = logId
				insertOldFrobErr := frobRepo.Create([]shared.InsertionModel{oldFrob})
				Expect(insertOldFrobErr).NotTo(HaveOccurred())

				headerTwo = createHeader(blockOne+1, timestampOne+1, headerRepository)
				newLogId := test_data.CreateTestLog(headerTwo.Id, db).ID

				newFrob = test_data.VatFrobModelWithNegativeDink()
				newFrob.ForeignKeyValues[constants.UrnFK] = fakeGuy
				newFrob.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
				newFrob.ColumnValues[constants.HeaderFK] = headerTwo.Id
				newFrob.ColumnValues[constants.LogFK] = newLogId
				insertNewFrobErr := frobRepo.Create([]shared.InsertionModel{newFrob})
				Expect(insertNewFrobErr).NotTo(HaveOccurred())
			})

			It("limits results if max_results argument is provided", func() {
				maxResults := 1
				var actualFrobs []test_helpers.FrobEvent
				getFrobsErr := db.Select(&actualFrobs,
					`SELECT ilk_identifier, urn_identifier, dink, dart FROM api.ilk_state_frobs(
						(SELECT (ilk_identifier, block_height, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated)::api.ilk_state
						 FROM api.get_ilk($1, $2)), $3)`,
					test_helpers.FakeIlk.Identifier, headerTwo.BlockNumber, maxResults)
				Expect(getFrobsErr).NotTo(HaveOccurred())

				expectedFrobs := []test_helpers.FrobEvent{{
					IlkIdentifier: test_helpers.FakeIlk.Identifier,
					UrnIdentifier: fakeGuy,
					Dink:          newFrob.ColumnValues["dink"].(string),
					Dart:          newFrob.ColumnValues["dart"].(string),
				}}
				Expect(actualFrobs).To(Equal(expectedFrobs))
			})

			It("offsets results if offset is provided", func() {
				maxResults := 1
				resultOffset := 1
				var actualFrobs []test_helpers.FrobEvent
				getFrobsErr := db.Select(&actualFrobs,
					`SELECT ilk_identifier, urn_identifier, dink, dart FROM api.ilk_state_frobs(
						(SELECT (ilk_identifier, block_height, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated)::api.ilk_state
						 FROM api.get_ilk($1, $2)), $3, $4)`,
					test_helpers.FakeIlk.Identifier, headerTwo.BlockNumber, maxResults, resultOffset)
				Expect(getFrobsErr).NotTo(HaveOccurred())

				expectedFrobs := []test_helpers.FrobEvent{{
					IlkIdentifier: test_helpers.FakeIlk.Identifier,
					UrnIdentifier: fakeGuy,
					Dink:          oldFrob.ColumnValues["dink"].(string),
					Dart:          oldFrob.ColumnValues["dart"].(string),
				}}
				Expect(actualFrobs).To(Equal(expectedFrobs))
			})
		})
	})

	Describe("ilks_state_ilk_file_events", func() {
		It("returns ilk file events for an ilk state", func() {
			fileRepo := ilk.VatFileIlkRepository{}
			fileRepo.SetDB(db)
			fileEvent := test_data.VatFileIlkDustModel()
			fileEvent.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
			fileEvent.ColumnValues[constants.HeaderFK] = headerOne.Id
			fileEvent.ColumnValues[constants.LogFK] = logId
			insertFileErr := fileRepo.Create([]shared.InsertionModel{fileEvent})
			Expect(insertFileErr).NotTo(HaveOccurred())

			var actualFiles []test_helpers.IlkFileEvent
			getFilesErr := db.Select(&actualFiles,
				`SELECT ilk_identifier, what, data FROM api.ilk_state_ilk_file_events(
					(SELECT (ilk_identifier, block_height, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated)::api.ilk_state
					 FROM api.get_ilk($1, $2)))`,
				test_helpers.FakeIlk.Identifier, blockOne)
			Expect(getFilesErr).NotTo(HaveOccurred())

			expectedFiles := []test_helpers.IlkFileEvent{{
				IlkIdentifier: test_helpers.GetValidNullString(test_helpers.FakeIlk.Identifier),
				What:          fileEvent.ColumnValues["what"].(string),
				Data:          fileEvent.ColumnValues["data"].(string),
			}}

			Expect(actualFiles).To(Equal(expectedFiles))
		})

		Describe("result pagination", func() {
			var (
				headerTwo              core.Header
				fileEvent, spotFileMat shared.InsertionModel
			)

			BeforeEach(func() {
				fileRepo := ilk.VatFileIlkRepository{}
				fileRepo.SetDB(db)
				fileEvent = test_data.VatFileIlkDustModel()
				fileEvent.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
				fileEvent.ColumnValues[constants.HeaderFK] = headerOne.Id
				fileEvent.ColumnValues[constants.LogFK] = logId
				insertFileErr := fileRepo.Create([]shared.InsertionModel{fileEvent})
				Expect(insertFileErr).NotTo(HaveOccurred())

				headerTwo = createHeader(blockOne+1, timestampOne+1, headerRepository)
				newLogId := test_data.CreateTestLog(headerTwo.Id, db).ID

				spotFileMatRepo := mat.SpotFileMatRepository{}
				spotFileMatRepo.SetDB(db)
				spotFileMat = test_data.SpotFileMatModel()
				spotFileMat.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
				spotFileMat.ColumnValues[constants.HeaderFK] = headerTwo.Id
				spotFileMat.ColumnValues[constants.LogFK] = newLogId
				spotFileMatErr := spotFileMatRepo.Create([]shared.InsertionModel{spotFileMat})
				Expect(spotFileMatErr).NotTo(HaveOccurred())
			})

			It("limits results to latest blocks if max_results argument is provided", func() {
				maxResults := 1
				var actualFiles []test_helpers.IlkFileEvent
				getFilesErr := db.Select(&actualFiles,
					`SELECT ilk_identifier, what, data FROM api.ilk_state_ilk_file_events(
						(SELECT (ilk_identifier, block_height, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated)::api.ilk_state
						 FROM api.get_ilk($1, $2)), $3)`,
					test_helpers.FakeIlk.Identifier, headerTwo.BlockNumber, maxResults)
				Expect(getFilesErr).NotTo(HaveOccurred())

				expectedFile := test_helpers.IlkFileEvent{
					IlkIdentifier: test_helpers.GetValidNullString(test_helpers.FakeIlk.Identifier),
					What:          spotFileMat.ColumnValues["what"].(string),
					Data:          spotFileMat.ColumnValues["data"].(string),
				}
				Expect(actualFiles).To(ConsistOf(expectedFile))
			})

			It("offsets results if offset is provided", func() {
				maxResults := 1
				resultOffset := 1
				var actualFiles []test_helpers.IlkFileEvent
				getFilesErr := db.Select(&actualFiles,
					`SELECT ilk_identifier, what, data FROM api.ilk_state_ilk_file_events(
                        (SELECT (ilk_identifier, block_height, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated)::api.ilk_state
                         FROM api.get_ilk($1, $2)), $3, $4)`,
					test_helpers.FakeIlk.Identifier, headerTwo.BlockNumber, maxResults, resultOffset)
				Expect(getFilesErr).NotTo(HaveOccurred())

				expectedFile := test_helpers.IlkFileEvent{
					IlkIdentifier: test_helpers.GetValidNullString(test_helpers.FakeIlk.Identifier),
					What:          fileEvent.ColumnValues["what"].(string),
					Data:          fileEvent.ColumnValues["data"].(string),
				}
				Expect(actualFiles).To(ConsistOf(expectedFile))
			})
		})
	})

	Describe("ilk_state_bites", func() {
		It("returns bite event for an ilk state", func() {
			biteEvent := generateBite(test_helpers.FakeIlk.Hex, fakeGuy, headerOne.Id, logId, db)
			insertBiteErr := event.PersistModels([]event.InsertionModel{biteEvent}, db)
			Expect(insertBiteErr).NotTo(HaveOccurred())

			var actualBites []test_helpers.BiteEvent
			getBitesErr := db.Select(&actualBites, `
				SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.ilk_state_bites(
					(SELECT (ilk_identifier, block_height, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated)::api.ilk_state
					FROM api.get_ilk($1, $2))
				)`,
				test_helpers.FakeIlk.Identifier,
				blockOne)
			Expect(getBitesErr).NotTo(HaveOccurred())

			expectedBites := []test_helpers.BiteEvent{{
				IlkIdentifier: test_helpers.FakeIlk.Identifier,
				UrnIdentifier: fakeGuy,
				Ink:           biteEvent.ColumnValues["ink"].(string),
				Art:           biteEvent.ColumnValues["art"].(string),
				Tab:           biteEvent.ColumnValues["tab"].(string),
			}}

			Expect(actualBites).To(Equal(expectedBites))
		})

		Describe("result pagination", func() {
			var (
				headerTwo        core.Header
				oldBite, newBite event.InsertionModel
				oldGuy           = "0x7d7bEe5fCfD8028cf7b00876C5b1421c800561A6"
			)

			BeforeEach(func() {
				oldBite = generateBite(test_helpers.FakeIlk.Hex, oldGuy, headerOne.Id, logId, db)
				insertOldBiteErr := event.PersistModels([]event.InsertionModel{oldBite}, db)
				Expect(insertOldBiteErr).NotTo(HaveOccurred())

				headerTwo = createHeader(blockOne+1, timestampOne+1, headerRepository)
				newLogId := test_data.CreateTestLog(headerTwo.Id, db).ID

				newBite = generateBite(test_helpers.FakeIlk.Hex, fakeGuy, headerTwo.Id, newLogId, db)
				insertNewBiteErr := event.PersistModels([]event.InsertionModel{newBite}, db)
				Expect(insertNewBiteErr).NotTo(HaveOccurred())
			})

			It("limits results to recent blocks if max_results argument is provided", func() {
				maxResults := 1
				var actualBites []test_helpers.BiteEvent
				getBitesErr := db.Select(&actualBites, `
				SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.ilk_state_bites(
					(SELECT (ilk_identifier, block_height, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated)::api.ilk_state
					FROM api.get_ilk($1, $2)), $3)`, test_helpers.FakeIlk.Identifier, headerTwo.BlockNumber, maxResults)
				Expect(getBitesErr).NotTo(HaveOccurred())

				expectedBite := test_helpers.BiteEvent{
					IlkIdentifier: test_helpers.FakeIlk.Identifier,
					UrnIdentifier: fakeGuy,
					Ink:           newBite.ColumnValues["ink"].(string),
					Art:           newBite.ColumnValues["art"].(string),
					Tab:           newBite.ColumnValues["tab"].(string),
				}
				Expect(actualBites).To(ConsistOf(expectedBite))
			})

			It("offsets results if offset is provided", func() {
				maxResults := 1
				resultOffset := 1
				var actualBites []test_helpers.BiteEvent
				getBitesErr := db.Select(&actualBites, `
				SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.ilk_state_bites(
					(SELECT (ilk_identifier, block_height, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated)::api.ilk_state
					FROM api.get_ilk($1, $2)), $3, $4)`,
					test_helpers.FakeIlk.Identifier, headerTwo.BlockNumber, maxResults, resultOffset)
				Expect(getBitesErr).NotTo(HaveOccurred())

				expectedBite := test_helpers.BiteEvent{
					IlkIdentifier: test_helpers.FakeIlk.Identifier,
					UrnIdentifier: oldGuy,
					Ink:           oldBite.ColumnValues["ink"].(string),
					Art:           oldBite.ColumnValues["art"].(string),
					Tab:           oldBite.ColumnValues["tab"].(string),
				}
				Expect(actualBites).To(ConsistOf(expectedBite))
			})
		})
	})
})
