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
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_file/ilk"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ilk File Events Query", func() {
	var (
		logOneId               int64
		blockOne, timestampOne int
		headerOne              core.Header
		headerRepo             datastore.HeaderRepository
		relevantIlkIdentifier  = test_helpers.GetValidNullString(test_helpers.FakeIlk.Identifier)
		vatFileRepo            ilk.VatFileIlkRepository
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepo)
		logOneId = test_data.CreateTestLog(headerOne.Id, db).ID
		vatFileRepo = ilk.VatFileIlkRepository{}
		vatFileRepo.SetDB(db)
	})

	It("returns all ilk file events for ilk", func() {
		catFileChopLumpLog := test_data.CreateTestLog(headerOne.Id, db)
		catFileChopLump := test_data.CatFileChopModel()
		ilkID, createIlkError := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
		Expect(createIlkError).NotTo(HaveOccurred())

		catFileChopLump.ColumnValues[constants.IlkColumn] = ilkID
		catFileChopLump.ColumnValues[constants.HeaderFK] = headerOne.Id
		catFileChopLump.ColumnValues[constants.LogFK] = catFileChopLumpLog.ID
		chopLumpErr := event.PersistModels([]event.InsertionModel{catFileChopLump}, db)
		Expect(chopLumpErr).NotTo(HaveOccurred())

		catFileFlipLog := test_data.CreateTestLog(headerOne.Id, db)
		catFileFlip := test_data.CatFileFlipModel()
		catFileFlip.ColumnValues[constants.IlkColumn] = ilkID
		catFileFlip.ColumnValues[constants.HeaderFK] = headerOne.Id
		catFileFlip.ColumnValues[constants.LogFK] = catFileFlipLog.ID
		flipErr := event.PersistModels([]event.InsertionModel{catFileFlip}, db)
		Expect(flipErr).NotTo(HaveOccurred())

		jugFileLog := test_data.CreateTestLog(headerOne.Id, db)
		jugFile := test_data.JugFileIlkModel()
		jugFile.ColumnValues[constants.IlkColumn] = ilkID
		jugFile.ColumnValues[constants.HeaderFK] = headerOne.Id
		jugFile.ColumnValues[constants.LogFK] = jugFileLog.ID
		jugErr := event.PersistModels([]event.InsertionModel{jugFile}, db)
		Expect(jugErr).NotTo(HaveOccurred())

		spotFileMatLog := test_data.CreateTestLog(headerOne.Id, db)
		spotFileMat := test_data.SpotFileMatModel()
		spotFileMat.ColumnValues[constants.HeaderFK] = headerOne.Id
		spotFileMat.ColumnValues[constants.LogFK] = spotFileMatLog.ID
		spotFileMat.ColumnValues[constants.IlkColumn] = ilkID
		spotFileMatErr := event.PersistModels([]event.InsertionModel{spotFileMat}, db)
		Expect(spotFileMatErr).NotTo(HaveOccurred())

		spotFilePipLog := test_data.CreateTestLog(headerOne.Id, db)
		spotFilePip := test_data.SpotFilePipModel()
		spotFilePip.ColumnValues[constants.HeaderFK] = headerOne.Id
		spotFilePip.ColumnValues[constants.LogFK] = spotFilePipLog.ID
		spotFilePip.ColumnValues[constants.IlkColumn] = ilkID
		spotFilePipErr := event.PersistModels([]event.InsertionModel{spotFilePip}, db)
		Expect(spotFilePipErr).NotTo(HaveOccurred())

		vatFileLog := test_data.CreateTestLog(headerOne.Id, db)
		vatFile := test_data.VatFileIlkDustModel()
		vatFile.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
		vatFile.ColumnValues[constants.HeaderFK] = headerOne.Id
		vatFile.ColumnValues[constants.LogFK] = vatFileLog.ID
		vatErr := vatFileRepo.Create([]shared.InsertionModel{vatFile})
		Expect(vatErr).NotTo(HaveOccurred())

		var actualFiles []test_helpers.IlkFileEvent
		filesErr := db.Select(&actualFiles, `SELECT ilk_identifier, what, data FROM api.all_ilk_file_events($1)`, test_helpers.FakeIlk.Identifier)
		Expect(filesErr).NotTo(HaveOccurred())

		Expect(actualFiles).To(ConsistOf(
			test_helpers.IlkFileEvent{
				IlkIdentifier: relevantIlkIdentifier,
				What:          catFileChopLump.ColumnValues["what"].(string),
				Data:          catFileChopLump.ColumnValues["data"].(string),
			},
			test_helpers.IlkFileEvent{
				IlkIdentifier: relevantIlkIdentifier,
				What:          catFileFlip.ColumnValues["what"].(string),
				Data:          catFileFlip.ColumnValues["flip"].(string),
			},
			test_helpers.IlkFileEvent{
				IlkIdentifier: relevantIlkIdentifier,
				What:          jugFile.ColumnValues["what"].(string),
				Data:          jugFile.ColumnValues["data"].(string),
			},
			test_helpers.IlkFileEvent{
				IlkIdentifier: relevantIlkIdentifier,
				What:          spotFileMat.ColumnValues[constants.WhatColumn].(string),
				Data:          spotFileMat.ColumnValues[constants.DataColumn].(string),
			},
			test_helpers.IlkFileEvent{
				IlkIdentifier: relevantIlkIdentifier,
				What:          spotFilePip.ColumnValues[constants.WhatColumn].(string),
				Data:          spotFilePip.ColumnValues[constants.PipColumn].(string),
			},
			test_helpers.IlkFileEvent{
				IlkIdentifier: relevantIlkIdentifier,
				What:          vatFile.ColumnValues["what"].(string),
				Data:          vatFile.ColumnValues["data"].(string),
			},
		))
	})

	It("includes results across blocks", func() {
		fileBlockOne := test_data.VatFileIlkDustModel()
		fileBlockOne.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
		fileBlockOne.ColumnValues["data"] = strconv.Itoa(rand.Int())
		fileBlockOne.ColumnValues[constants.HeaderFK] = headerOne.Id
		fileBlockOne.ColumnValues[constants.LogFK] = logOneId
		fileBlockOneErr := vatFileRepo.Create([]shared.InsertionModel{fileBlockOne})
		Expect(fileBlockOneErr).NotTo(HaveOccurred())

		headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)

		logBlockTwo := test_data.CreateTestLog(headerTwo.Id, db)
		fileBlockTwo := test_data.VatFileIlkDustModel()
		fileBlockTwo.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
		fileBlockTwo.ColumnValues["data"] = strconv.Itoa(rand.Int())
		fileBlockTwo.ColumnValues[constants.HeaderFK] = headerTwo.Id
		fileBlockTwo.ColumnValues[constants.LogFK] = logBlockTwo.ID
		fileBlockTwoErr := vatFileRepo.Create([]shared.InsertionModel{fileBlockTwo})
		Expect(fileBlockTwoErr).NotTo(HaveOccurred())

		var actualFiles []test_helpers.IlkFileEvent
		filesErr := db.Select(&actualFiles, `SELECT ilk_identifier, what, data FROM api.all_ilk_file_events($1)`, test_helpers.FakeIlk.Identifier)
		Expect(filesErr).NotTo(HaveOccurred())

		Expect(actualFiles).To(ConsistOf(
			test_helpers.IlkFileEvent{
				IlkIdentifier: relevantIlkIdentifier,
				What:          fileBlockOne.ColumnValues["what"].(string),
				Data:          fileBlockOne.ColumnValues["data"].(string),
			},
			test_helpers.IlkFileEvent{
				IlkIdentifier: relevantIlkIdentifier,
				What:          fileBlockTwo.ColumnValues["what"].(string),
				Data:          fileBlockTwo.ColumnValues["data"].(string),
			},
		))
	})

	Describe("result pagination", func() {
		var fileBlockOne, fileBlockTwo shared.InsertionModel

		BeforeEach(func() {
			fileBlockOne = test_data.VatFileIlkDustModel()
			fileBlockOne.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
			fileBlockOne.ColumnValues["data"] = strconv.Itoa(rand.Int())
			fileBlockOne.ColumnValues[constants.HeaderFK] = headerOne.Id
			fileBlockOne.ColumnValues[constants.LogFK] = logOneId
			fileBlockOneErr := vatFileRepo.Create([]shared.InsertionModel{fileBlockOne})
			Expect(fileBlockOneErr).NotTo(HaveOccurred())

			headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
			logTwoID := test_data.CreateTestLog(headerOne.Id, db).ID

			fileBlockTwo = test_data.VatFileIlkDustModel()
			fileBlockTwo.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
			fileBlockTwo.ColumnValues["data"] = strconv.Itoa(rand.Int())
			fileBlockTwo.ColumnValues[constants.HeaderFK] = headerTwo.Id
			fileBlockTwo.ColumnValues[constants.LogFK] = logTwoID
			fileBlockTwoErr := vatFileRepo.Create([]shared.InsertionModel{fileBlockTwo})
			Expect(fileBlockTwoErr).NotTo(HaveOccurred())
		})

		It("limits results to most recent blocks if max_results argument is provided", func() {
			maxResults := 1
			var actualFiles []test_helpers.IlkFileEvent
			filesErr := db.Select(&actualFiles, `SELECT ilk_identifier, what, data FROM api.all_ilk_file_events($1, $2)`,
				test_helpers.FakeIlk.Identifier, maxResults)
			Expect(filesErr).NotTo(HaveOccurred())

			Expect(actualFiles).To(ConsistOf(
				test_helpers.IlkFileEvent{
					IlkIdentifier: relevantIlkIdentifier,
					What:          fileBlockTwo.ColumnValues["what"].(string),
					Data:          fileBlockTwo.ColumnValues["data"].(string),
				},
			))
		})

		It("offsets results if offset is provided", func() {
			maxResults := 1
			resultOffset := 1
			var actualFiles []test_helpers.IlkFileEvent
			filesErr := db.Select(&actualFiles, `SELECT ilk_identifier, what, data FROM api.all_ilk_file_events($1, $2, $3)`,
				test_helpers.FakeIlk.Identifier, maxResults, resultOffset)
			Expect(filesErr).NotTo(HaveOccurred())

			Expect(actualFiles).To(ConsistOf(
				test_helpers.IlkFileEvent{
					IlkIdentifier: relevantIlkIdentifier,
					What:          fileBlockOne.ColumnValues["what"].(string),
					Data:          fileBlockOne.ColumnValues["data"].(string),
				},
			))
		})
	})

	It("does not include ilk file events for a different ilk", func() {
		relevantFile := test_data.VatFileIlkDustModel()
		relevantFile.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
		relevantFile.ColumnValues["data"] = strconv.Itoa(rand.Int())
		relevantFile.ColumnValues[constants.HeaderFK] = headerOne.Id
		relevantFile.ColumnValues[constants.LogFK] = logOneId

		irrelevantLog := test_data.CreateTestLog(headerOne.Id, db)
		irrelevantFile := test_data.VatFileIlkDustModel()
		irrelevantFile.ForeignKeyValues[constants.IlkFK] = test_helpers.AnotherFakeIlk.Hex
		irrelevantFile.ColumnValues["data"] = strconv.Itoa(rand.Int())
		irrelevantFile.ColumnValues[constants.HeaderFK] = headerOne.Id
		irrelevantFile.ColumnValues[constants.LogFK] = irrelevantLog.ID

		models := []shared.InsertionModel{relevantFile, irrelevantFile}
		vatBlockOneErr := vatFileRepo.Create(models)
		Expect(vatBlockOneErr).NotTo(HaveOccurred())

		var actualFiles []test_helpers.IlkFileEvent
		filesErr := db.Select(&actualFiles, `SELECT ilk_identifier, what, data FROM api.all_ilk_file_events($1)`, test_helpers.FakeIlk.Identifier)
		Expect(filesErr).NotTo(HaveOccurred())

		Expect(actualFiles).To(ConsistOf(
			test_helpers.IlkFileEvent{
				IlkIdentifier: relevantIlkIdentifier,
				What:          relevantFile.ColumnValues["what"].(string),
				Data:          relevantFile.ColumnValues["data"].(string),
			},
		))
	})

	It("fails if no argument is supplied (STRICT)", func() {
		_, err := db.Exec(`SELECT * FROM api.all_ilk_file_events()`)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("function api.all_ilk_file_events() does not exist"))
	})
})
