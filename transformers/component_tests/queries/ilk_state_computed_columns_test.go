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

var _ = Describe("Ilk state computed columns", func() {
	var (
		db               *postgres.DB
		fakeBlock        int
		fakeGuy          = "fakeAddress"
		fakeHeader       core.Header
		headerId         int64
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

		ilkValues := test_helpers.GetIlkValues(0)
		test_helpers.CreateIlk(db, fakeHeader, ilkValues, test_helpers.FakeIlkVatMetadatas,
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
			frobEvent := test_data.CopyModel(test_data.VatFrobModelWithPositiveDart)
			frobEvent.ForeignKeyValues[constants.UrnFK] = fakeGuy
			frobEvent.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
			insertFrobErr := frobRepo.Create(headerId, []shared.InsertionModel{frobEvent})
			Expect(insertFrobErr).NotTo(HaveOccurred())

			var actualFrobs []test_helpers.FrobEvent
			getFrobsErr := db.Select(&actualFrobs,
				`SELECT ilk_identifier, urn_identifier, dink, dart FROM api.ilk_state_frobs(
					(SELECT (ilk_identifier, block_height, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated)::api.ilk_state
					 FROM api.get_ilk($1, $2)))`,
				test_helpers.FakeIlk.Identifier, fakeBlock)
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
				oldFrob = test_data.CopyModel(test_data.VatFrobModelWithPositiveDart)
				oldFrob.ForeignKeyValues[constants.UrnFK] = fakeGuy
				oldFrob.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
				insertOldFrobErr := frobRepo.Create(headerId, []shared.InsertionModel{oldFrob})
				Expect(insertOldFrobErr).NotTo(HaveOccurred())

				newBlock = fakeBlock + 1
				newHeader := fakes.GetFakeHeader(int64(newBlock))
				newHeaderId, newHeaderErr := headerRepository.CreateOrUpdateHeader(newHeader)
				Expect(newHeaderErr).NotTo(HaveOccurred())

				newFrob = test_data.CopyModel(test_data.VatFrobModelWithNegativeDink)
				newFrob.ForeignKeyValues[constants.UrnFK] = fakeGuy
				newFrob.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
				insertNewFrobErr := frobRepo.Create(newHeaderId, []shared.InsertionModel{newFrob})
				Expect(insertNewFrobErr).NotTo(HaveOccurred())
			})

			It("limits results if max_results argument is provided", func() {
				maxResults := 1
				var actualFrobs []test_helpers.FrobEvent
				getFrobsErr := db.Select(&actualFrobs,
					`SELECT ilk_identifier, urn_identifier, dink, dart FROM api.ilk_state_frobs(
						(SELECT (ilk_identifier, block_height, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated)::api.ilk_state
						 FROM api.get_ilk($1, $2)), $3)`,
					test_helpers.FakeIlk.Identifier, newBlock, maxResults)
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
					test_helpers.FakeIlk.Identifier, newBlock, maxResults, resultOffset)
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
			fileEvent := test_data.VatFileIlkDustModel
			fileEvent.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
			insertFileErr := fileRepo.Create(headerId, []shared.InsertionModel{fileEvent})
			Expect(insertFileErr).NotTo(HaveOccurred())

			var actualFiles []test_helpers.IlkFileEvent
			getFilesErr := db.Select(&actualFiles,
				`SELECT ilk_identifier, what, data FROM api.ilk_state_ilk_file_events(
					(SELECT (ilk_identifier, block_height, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated)::api.ilk_state
					 FROM api.get_ilk($1, $2)))`,
				test_helpers.FakeIlk.Identifier, fakeBlock)
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
				fileEvent = test_data.VatFileIlkDustModel
				fileEvent.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
				insertFileErr := fileRepo.Create(headerId, []shared.InsertionModel{fileEvent})
				Expect(insertFileErr).NotTo(HaveOccurred())

				newBlock = fakeBlock + 1
				newHeader := fakes.GetFakeHeader(int64(newBlock))
				newHeaderId, insertNewHeaderErr := headerRepository.CreateOrUpdateHeader(newHeader)
				Expect(insertNewHeaderErr).NotTo(HaveOccurred())

				spotFileMatRepo := mat.SpotFileMatRepository{}
				spotFileMatRepo.SetDB(db)
				spotFileMat = test_data.SpotFileMatModel
				spotFileMat.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
				spotFileMatErr := spotFileMatRepo.Create(newHeaderId, []shared.InsertionModel{spotFileMat})
				Expect(spotFileMatErr).NotTo(HaveOccurred())
			})

			It("limits results to latest blocks if max_results argument is provided", func() {
				maxResults := 1
				var actualFiles []test_helpers.IlkFileEvent
				getFilesErr := db.Select(&actualFiles,
					`SELECT ilk_identifier, what, data FROM api.ilk_state_ilk_file_events(
						(SELECT (ilk_identifier, block_height, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated)::api.ilk_state
						 FROM api.get_ilk($1, $2)), $3)`,
					test_helpers.FakeIlk.Identifier, newBlock, maxResults)
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
					test_helpers.FakeIlk.Identifier, newBlock, maxResults, resultOffset)
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
			biteRepo := bite.BiteRepository{}
			biteRepo.SetDB(db)
			biteEvent := test_data.BiteModel
			biteEvent.Ilk = test_helpers.FakeIlk.Hex
			insertBiteErr := biteRepo.Create(headerId, []interface{}{biteEvent})
			Expect(insertBiteErr).NotTo(HaveOccurred())

			var actualBites []test_helpers.BiteEvent
			getBitesErr := db.Select(&actualBites, `
				SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.ilk_state_bites(
					(SELECT (ilk_identifier, block_height, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated)::api.ilk_state
					FROM api.get_ilk($1, $2))
				)`,
				test_helpers.FakeIlk.Identifier,
				fakeBlock)
			Expect(getBitesErr).NotTo(HaveOccurred())

			expectedBites := []test_helpers.BiteEvent{{
				IlkIdentifier: test_helpers.FakeIlk.Identifier,
				UrnIdentifier: biteEvent.Urn,
				Ink:           biteEvent.Ink,
				Art:           biteEvent.Art,
				Tab:           biteEvent.Tab,
			}}

			Expect(actualBites).To(Equal(expectedBites))
		})

		Describe("result pagination", func() {
			var (
				newBlock         int
				oldBite, newBite bite.BiteModel
			)

			BeforeEach(func() {
				biteRepo := bite.BiteRepository{}
				biteRepo.SetDB(db)
				oldBite = test_data.BiteModel
				oldBite.Ilk = test_helpers.FakeIlk.Hex
				insertOldBiteErr := biteRepo.Create(headerId, []interface{}{oldBite})
				Expect(insertOldBiteErr).NotTo(HaveOccurred())

				newBlock = fakeBlock + 1
				newHeader := fakes.GetFakeHeader(int64(newBlock))
				newHeaderId, insertNewHeaderErr := headerRepository.CreateOrUpdateHeader(newHeader)
				Expect(insertNewHeaderErr).NotTo(HaveOccurred())

				newBite = test_data.BiteModel
				newBite.Ilk = test_helpers.FakeIlk.Hex
				newBite.Urn = test_data.FakeUrn
				insertNewBiteErr := biteRepo.Create(newHeaderId, []interface{}{newBite})
				Expect(insertNewBiteErr).NotTo(HaveOccurred())
			})

			It("limits results to recent blocks if max_results argument is provided", func() {
				maxResults := 1
				var actualBites []test_helpers.BiteEvent
				getBitesErr := db.Select(&actualBites, `
				SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.ilk_state_bites(
					(SELECT (ilk_identifier, block_height, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated)::api.ilk_state
					FROM api.get_ilk($1, $2)), $3)`, test_helpers.FakeIlk.Identifier, newBlock, maxResults)
				Expect(getBitesErr).NotTo(HaveOccurred())

				expectedBite := test_helpers.BiteEvent{
					IlkIdentifier: test_helpers.FakeIlk.Identifier,
					UrnIdentifier: newBite.Urn,
					Ink:           newBite.Ink,
					Art:           newBite.Art,
					Tab:           newBite.Tab,
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
					test_helpers.FakeIlk.Identifier, newBlock, maxResults, resultOffset)
				Expect(getBitesErr).NotTo(HaveOccurred())

				expectedBite := test_helpers.BiteEvent{
					IlkIdentifier: test_helpers.FakeIlk.Identifier,
					UrnIdentifier: oldBite.Urn,
					Ink:           oldBite.Ink,
					Art:           oldBite.Art,
					Tab:           oldBite.Tab,
				}
				Expect(actualBites).To(ConsistOf(expectedBite))
			})
		})
	})
})
