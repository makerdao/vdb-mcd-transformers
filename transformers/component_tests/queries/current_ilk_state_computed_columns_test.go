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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/bite"
	"github.com/vulcanize/mcd_transformers/transformers/events/spot_file/mat"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_file/ilk"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_frob"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("current ilk state computed columns", func() {
	var (
		db               *postgres.DB
		fakeBlock        int
		fakeGuy          = "fakeAddress"
		fakeHeader       core.Header
		headerId, logId  int64
		headerRepository repositories.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)

		headerRepository = repositories.NewHeaderRepository(db)
		fakeBlock = rand.Int()
		fakeHeader = fakes.GetFakeHeader(int64(fakeBlock))
		var insertHeaderErr error
		headerId, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
		fakeHeaderSyncLog := test_data.CreateTestLog(headerId, db)
		logId = fakeHeaderSyncLog.ID

		ilkValues := test_helpers.GetIlkValues(0)
		test_helpers.CreateIlk(db, fakeHeader, ilkValues, test_helpers.FakeIlkVatMetadatas,
			test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas, test_helpers.FakeIlkSpotMetadatas)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("current_ilk_state_frobs", func() {
		It("returns relevant frobs for a current_ilk_state", func() {
			frobRepo := vat_frob.VatFrobRepository{}
			frobRepo.SetDB(db)
			frobEvent := test_data.VatFrobModelWithPositiveDart()
			frobEvent.ForeignKeyValues[constants.UrnFK] = fakeGuy
			frobEvent.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
			frobEvent.ColumnValues[constants.HeaderFK] = headerId
			frobEvent.ColumnValues[constants.LogFK] = logId
			insertFrobErr := frobRepo.Create([]shared.InsertionModel{frobEvent})
			Expect(insertFrobErr).NotTo(HaveOccurred())

			var actualFrobs []test_helpers.FrobEvent
			getFrobsErr := db.Select(&actualFrobs,
				`SELECT ilk_identifier, urn_identifier, dink, dart FROM api.current_ilk_state_frobs(
					(SELECT (ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated)::api.current_ilk_state
					 FROM api.current_ilk_state))`)
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
				newBlock         int
				oldFrob, newFrob shared.InsertionModel
			)

			BeforeEach(func() {
				frobRepo := vat_frob.VatFrobRepository{}
				frobRepo.SetDB(db)
				oldFrob = test_data.VatFrobModelWithPositiveDart()
				oldFrob.ForeignKeyValues[constants.UrnFK] = fakeGuy
				oldFrob.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
				oldFrob.ColumnValues[constants.HeaderFK] = headerId
				oldFrob.ColumnValues[constants.LogFK] = logId
				insertOldFrobErr := frobRepo.Create([]shared.InsertionModel{oldFrob})
				Expect(insertOldFrobErr).NotTo(HaveOccurred())

				newBlock = fakeBlock + 1
				newHeader := fakes.GetFakeHeader(int64(newBlock))
				newHeaderId, newHeaderErr := headerRepository.CreateOrUpdateHeader(newHeader)
				Expect(newHeaderErr).NotTo(HaveOccurred())
				newLogId := test_data.CreateTestLog(newHeaderId, db).ID

				newFrob = test_data.VatFrobModelWithNegativeDink()
				newFrob.ForeignKeyValues[constants.UrnFK] = fakeGuy
				newFrob.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
				newFrob.ColumnValues[constants.HeaderFK] = newHeaderId
				newFrob.ColumnValues[constants.LogFK] = newLogId
				insertNewFrobErr := frobRepo.Create([]shared.InsertionModel{newFrob})
				Expect(insertNewFrobErr).NotTo(HaveOccurred())
			})

			It("limits results if max_results argument is provided", func() {
				maxResults := 1
				var actualFrobs []test_helpers.FrobEvent
				getFrobsErr := db.Select(&actualFrobs,
					`SELECT ilk_identifier, urn_identifier, dink, dart FROM api.current_ilk_state_frobs(
						(SELECT (ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated)::api.current_ilk_state
						 FROM api.current_ilk_state), $1)`,
					maxResults)
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
					`SELECT ilk_identifier, urn_identifier, dink, dart FROM api.current_ilk_state_frobs(
						(SELECT (ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated)::api.current_ilk_state
						 FROM api.current_ilk_state), $1, $2)`,
					maxResults, resultOffset)
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

	Describe("current_ilk_state_ilk_file_events", func() {
		It("returns ilk file events for a current_ilk_state", func() {
			fileRepo := ilk.VatFileIlkRepository{}
			fileRepo.SetDB(db)
			fileEvent := test_data.VatFileIlkDustModel()
			fileEvent.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
			fileEvent.ColumnValues[constants.HeaderFK] = headerId
			fileEvent.ColumnValues[constants.LogFK] = logId
			insertFileErr := fileRepo.Create([]shared.InsertionModel{fileEvent})
			Expect(insertFileErr).NotTo(HaveOccurred())

			var actualFiles []test_helpers.IlkFileEvent
			getFilesErr := db.Select(&actualFiles,
				`SELECT ilk_identifier, what, data FROM api.current_ilk_state_ilk_file_events(
					(SELECT (ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated)::api.current_ilk_state
					 FROM api.current_ilk_state))`)
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
				newBlock               int
				fileEvent, spotFileMat shared.InsertionModel
			)

			BeforeEach(func() {
				fileRepo := ilk.VatFileIlkRepository{}
				fileRepo.SetDB(db)
				fileEvent = test_data.VatFileIlkDustModel()
				fileEvent.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
				fileEvent.ColumnValues[constants.HeaderFK] = headerId
				fileEvent.ColumnValues[constants.LogFK] = logId
				insertFileErr := fileRepo.Create([]shared.InsertionModel{fileEvent})
				Expect(insertFileErr).NotTo(HaveOccurred())

				newBlock = fakeBlock + 1
				newHeader := fakes.GetFakeHeader(int64(newBlock))
				newHeaderId, insertNewHeaderErr := headerRepository.CreateOrUpdateHeader(newHeader)
				Expect(insertNewHeaderErr).NotTo(HaveOccurred())
				newLogId := test_data.CreateTestLog(newHeaderId, db).ID

				spotFileMatRepo := mat.SpotFileMatRepository{}
				spotFileMatRepo.SetDB(db)
				spotFileMat = test_data.SpotFileMatModel()
				spotFileMat.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
				spotFileMat.ColumnValues[constants.HeaderFK] = newHeaderId
				spotFileMat.ColumnValues[constants.LogFK] = newLogId
				spotFileMatErr := spotFileMatRepo.Create([]shared.InsertionModel{spotFileMat})
				Expect(spotFileMatErr).NotTo(HaveOccurred())
			})

			It("limits results to latest blocks if max_results argument is provided", func() {
				maxResults := 1
				var actualFiles []test_helpers.IlkFileEvent
				getFilesErr := db.Select(&actualFiles,
					`SELECT ilk_identifier, what, data FROM api.current_ilk_state_ilk_file_events(
						(SELECT (ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated)::api.current_ilk_state
						 FROM api.current_ilk_state), $1)`,
					maxResults)
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
					`SELECT ilk_identifier, what, data FROM api.current_ilk_state_ilk_file_events(
						(SELECT (ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated)::api.current_ilk_state
						 FROM api.current_ilk_state), $1, $2)`,
					maxResults, resultOffset)
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

	Describe("current_ilk_state_bites", func() {
		It("returns bite event for a current ilk state", func() {
			biteRepo := bite.BiteRepository{}
			biteRepo.SetDB(db)
			biteEvent := test_data.BiteModel()
			biteEvent.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
			biteEvent.ColumnValues[constants.HeaderFK] = headerId
			biteEvent.ColumnValues[constants.LogFK] = logId
			insertBiteErr := biteRepo.Create([]shared.InsertionModel{biteEvent})
			Expect(insertBiteErr).NotTo(HaveOccurred())

			var actualBites []test_helpers.BiteEvent
			getBitesErr := db.Select(&actualBites, `
				SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.current_ilk_state_bites(
					(SELECT (ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated)::api.current_ilk_state
					 FROM api.current_ilk_state))`)
			Expect(getBitesErr).NotTo(HaveOccurred())

			expectedBites := []test_helpers.BiteEvent{{
				IlkIdentifier: test_helpers.FakeIlk.Identifier,
				UrnIdentifier: biteEvent.ForeignKeyValues[constants.UrnFK],
				Ink:           biteEvent.ColumnValues["ink"].(string),
				Art:           biteEvent.ColumnValues["art"].(string),
				Tab:           biteEvent.ColumnValues["tab"].(string),
			}}

			Expect(actualBites).To(Equal(expectedBites))
		})

		Describe("result pagination", func() {
			var (
				newBlock         int
				oldBite, newBite shared.InsertionModel
			)

			BeforeEach(func() {
				biteRepo := bite.BiteRepository{}
				biteRepo.SetDB(db)
				oldBite = test_data.BiteModel()
				oldBite.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
				oldBite.ColumnValues[constants.HeaderFK] = headerId
				oldBite.ColumnValues[constants.LogFK] = logId
				insertOldBiteErr := biteRepo.Create([]shared.InsertionModel{oldBite})
				Expect(insertOldBiteErr).NotTo(HaveOccurred())

				newBlock = fakeBlock + 1
				newHeader := fakes.GetFakeHeader(int64(newBlock))
				newHeaderId, insertNewHeaderErr := headerRepository.CreateOrUpdateHeader(newHeader)
				Expect(insertNewHeaderErr).NotTo(HaveOccurred())
				newLogId := test_data.CreateTestLog(newHeaderId, db).ID

				newBite = test_data.BiteModel()
				newBite.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
				newBite.ForeignKeyValues[constants.UrnFK] = test_data.FakeUrn
				newBite.ColumnValues[constants.HeaderFK] = newHeaderId
				newBite.ColumnValues[constants.LogFK] = newLogId
				insertNewBiteErr := biteRepo.Create([]shared.InsertionModel{newBite})
				Expect(insertNewBiteErr).NotTo(HaveOccurred())
			})

			It("limits results to recent blocks if max_results argument is provided", func() {
				maxResults := 1
				var actualBites []test_helpers.BiteEvent
				getBitesErr := db.Select(&actualBites, `
					SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.current_ilk_state_bites(
						(SELECT (ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated)::api.current_ilk_state
						 FROM api.current_ilk_state), $1)`,
					maxResults)
				Expect(getBitesErr).NotTo(HaveOccurred())

				expectedBite := test_helpers.BiteEvent{
					IlkIdentifier: test_helpers.FakeIlk.Identifier,
					UrnIdentifier: newBite.ForeignKeyValues[constants.UrnFK],
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
					SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.current_ilk_state_bites(
						(SELECT (ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated)::api.current_ilk_state
						 FROM api.current_ilk_state), $1, $2)`,
					maxResults, resultOffset)
				Expect(getBitesErr).NotTo(HaveOccurred())

				expectedBite := test_helpers.BiteEvent{
					IlkIdentifier: test_helpers.FakeIlk.Identifier,
					UrnIdentifier: oldBite.ForeignKeyValues[constants.UrnFK],
					Ink:           oldBite.ColumnValues["ink"].(string),
					Art:           oldBite.ColumnValues["art"].(string),
					Tab:           oldBite.ColumnValues["tab"].(string),
				}
				Expect(actualBites).To(ConsistOf(expectedBite))
			})
		})
	})
})
