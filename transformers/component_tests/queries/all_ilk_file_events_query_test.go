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
		logOneID               int64
		blockOne, timestampOne int
		headerOne              core.Header
		headerRepo             datastore.HeaderRepository
		relevantIlkIdentifier  = test_helpers.GetValidNullString(test_helpers.FakeIlk.Identifier)
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepo)
		logOneID = test_data.CreateTestLog(headerOne.Id, db).ID
	})

	It("returns all ilk file events for ilk", func() {
		catFileChopLumpLog := test_data.CreateTestLogFromEventLog(headerOne.Id, test_data.CatFileChopEventLog.Log, db)
		catFileChopLump := test_data.CatFileChopModel()
		ilkID, createIlkError := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
		Expect(createIlkError).NotTo(HaveOccurred())

		catFileChopLump.ColumnValues[constants.IlkColumn] = ilkID
		catFileChopLump.ColumnValues[event.HeaderFK] = headerOne.Id
		catFileChopLump.ColumnValues[event.LogFK] = catFileChopLumpLog.ID
		test_data.AssignMessageSenderID(catFileChopLumpLog, catFileChopLump, db)
		chopLumpErr := event.PersistModels([]event.InsertionModel{catFileChopLump}, db)
		Expect(chopLumpErr).NotTo(HaveOccurred())

		catFileFlipLog := test_data.CreateTestLogFromEventLog(headerOne.Id, test_data.CatFileFlipEventLog.Log, db)
		catFileFlip := test_data.CatFileFlipModel()
		catFileFlip.ColumnValues[constants.IlkColumn] = ilkID
		catFileFlip.ColumnValues[event.HeaderFK] = headerOne.Id
		catFileFlip.ColumnValues[event.LogFK] = catFileFlipLog.ID
		test_data.AssignMessageSenderID(catFileFlipLog, catFileFlip, db)
		flipErr := event.PersistModels([]event.InsertionModel{catFileFlip}, db)
		Expect(flipErr).NotTo(HaveOccurred())

		jugFileLog := test_data.CreateTestLog(headerOne.Id, db)
		jugFile := test_data.JugFileIlkModel()
		jugFile.ColumnValues[constants.IlkColumn] = ilkID
		jugFile.ColumnValues[event.HeaderFK] = headerOne.Id
		jugFile.ColumnValues[event.LogFK] = jugFileLog.ID
		test_data.AssignMessageSenderID(test_data.JugFileIlkEventLog, jugFile, db)
		jugErr := event.PersistModels([]event.InsertionModel{jugFile}, db)
		Expect(jugErr).NotTo(HaveOccurred())

		spotFileMatLog := test_data.CreateTestLog(headerOne.Id, db)
		spotFileMat := test_data.SpotFileMatModel()
		spotFileMat.ColumnValues[event.HeaderFK] = headerOne.Id
		spotFileMat.ColumnValues[event.LogFK] = spotFileMatLog.ID
		spotFileMat.ColumnValues[constants.IlkColumn] = ilkID
		spotFileMatErr := event.PersistModels([]event.InsertionModel{spotFileMat}, db)
		Expect(spotFileMatErr).NotTo(HaveOccurred())

		spotFilePipLog := test_data.CreateTestLog(headerOne.Id, db)
		spotFilePip := test_data.SpotFilePipModel()
		spotFilePip.ColumnValues[event.HeaderFK] = headerOne.Id
		spotFilePip.ColumnValues[event.LogFK] = spotFilePipLog.ID
		spotFilePip.ColumnValues[constants.IlkColumn] = ilkID
		spotFilePipErr := event.PersistModels([]event.InsertionModel{spotFilePip}, db)
		Expect(spotFilePipErr).NotTo(HaveOccurred())

		vatFileLog := test_data.CreateTestLog(headerOne.Id, db)
		vatFile := test_data.VatFileIlkDustModel()
		vatFile.ColumnValues[event.HeaderFK] = headerOne.Id
		vatFile.ColumnValues[event.LogFK] = vatFileLog.ID
		vatFile.ColumnValues[constants.IlkColumn] = ilkID
		vatErr := event.PersistModels([]event.InsertionModel{vatFile}, db)
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
		ilkID, createIlkError := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
		Expect(createIlkError).NotTo(HaveOccurred())

		fileBlockOne.ColumnValues[constants.IlkColumn] = ilkID
		fileBlockOne.ColumnValues[constants.DataColumn] = strconv.Itoa(rand.Int())
		fileBlockOne.ColumnValues[event.HeaderFK] = headerOne.Id
		fileBlockOne.ColumnValues[event.LogFK] = logOneID
		fileBlockOneErr := event.PersistModels([]event.InsertionModel{fileBlockOne}, db)
		Expect(fileBlockOneErr).NotTo(HaveOccurred())

		headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)

		logBlockTwo := test_data.CreateTestLog(headerTwo.Id, db)
		fileBlockTwo := test_data.VatFileIlkDustModel()
		fileBlockTwo.ColumnValues[constants.IlkColumn] = ilkID
		fileBlockTwo.ColumnValues[constants.DataColumn] = strconv.Itoa(rand.Int())
		fileBlockTwo.ColumnValues[event.HeaderFK] = headerTwo.Id
		fileBlockTwo.ColumnValues[event.LogFK] = logBlockTwo.ID
		fileBlockTwoErr := event.PersistModels([]event.InsertionModel{fileBlockTwo}, db)
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
		var fileBlockOne, fileBlockTwo event.InsertionModel

		BeforeEach(func() {
			fileBlockOne = test_data.VatFileIlkDustModel()
			ilkID, createIlkError := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(createIlkError).NotTo(HaveOccurred())

			fileBlockOne.ColumnValues[constants.IlkColumn] = ilkID
			fileBlockOne.ColumnValues[constants.DataColumn] = strconv.Itoa(rand.Int())
			fileBlockOne.ColumnValues[event.HeaderFK] = headerOne.Id
			fileBlockOne.ColumnValues[event.LogFK] = logOneID
			fileBlockOneErr := event.PersistModels([]event.InsertionModel{fileBlockOne}, db)
			Expect(fileBlockOneErr).NotTo(HaveOccurred())

			headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
			logTwoID := test_data.CreateTestLog(headerOne.Id, db).ID

			fileBlockTwo = test_data.VatFileIlkDustModel()
			fileBlockTwo.ColumnValues[constants.IlkColumn] = ilkID
			fileBlockTwo.ColumnValues[constants.DataColumn] = strconv.Itoa(rand.Int())
			fileBlockTwo.ColumnValues[event.HeaderFK] = headerTwo.Id
			fileBlockTwo.ColumnValues[event.LogFK] = logTwoID
			fileBlockTwoErr := event.PersistModels([]event.InsertionModel{fileBlockTwo}, db)
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
		ilkID, createIlkError := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
		Expect(createIlkError).NotTo(HaveOccurred())

		relevantFile.ColumnValues[constants.IlkColumn] = ilkID
		relevantFile.ColumnValues[constants.DataColumn] = strconv.Itoa(rand.Int())
		relevantFile.ColumnValues[event.HeaderFK] = headerOne.Id
		relevantFile.ColumnValues[event.LogFK] = logOneID

		irrelevantLog := test_data.CreateTestLog(headerOne.Id, db)
		irrelevantFile := test_data.VatFileIlkDustModel()
		anotherIlkID, createIlkError := shared.GetOrCreateIlk(test_helpers.AnotherFakeIlk.Hex, db)
		Expect(createIlkError).NotTo(HaveOccurred())

		irrelevantFile.ColumnValues[constants.IlkColumn] = anotherIlkID
		irrelevantFile.ColumnValues[constants.DataColumn] = strconv.Itoa(rand.Int())
		irrelevantFile.ColumnValues[event.HeaderFK] = headerOne.Id
		irrelevantFile.ColumnValues[event.LogFK] = irrelevantLog.ID

		models := []event.InsertionModel{relevantFile, irrelevantFile}
		vatBlockOneErr := event.PersistModels(models, db)
		Expect(vatBlockOneErr).NotTo(HaveOccurred())

		var actualFiles []test_helpers.IlkFileEvent
		filesErr := db.Select(&actualFiles, `SELECT ilk_identifier, what, data FROM api.all_ilk_file_events($1)`, test_helpers.FakeIlk.Identifier)
		Expect(filesErr).NotTo(HaveOccurred())

		Expect(actualFiles).To(ConsistOf(
			test_helpers.IlkFileEvent{
				IlkIdentifier: relevantIlkIdentifier,
				What:          relevantFile.ColumnValues[constants.WhatColumn].(string),
				Data:          relevantFile.ColumnValues[constants.DataColumn].(string),
			},
		))
	})

	It("fails if no argument is supplied (STRICT)", func() {
		_, err := db.Exec(`SELECT * FROM api.all_ilk_file_events()`)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("function api.all_ilk_file_events() does not exist"))
	})
})
