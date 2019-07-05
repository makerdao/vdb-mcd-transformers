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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/cat_file/chop_lump"
	"github.com/vulcanize/mcd_transformers/transformers/events/cat_file/flip"
	ilk2 "github.com/vulcanize/mcd_transformers/transformers/events/jug_file/ilk"
	"github.com/vulcanize/mcd_transformers/transformers/events/spot_file/mat"
	"github.com/vulcanize/mcd_transformers/transformers/events/spot_file/pip"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_file/ilk"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("Ilk File Events Query", func() {
	var (
		catFileChopLumpRepo   chop_lump.CatFileChopLumpRepository
		catFileFlipRepo       flip.CatFileFlipRepository
		db                    *postgres.DB
		err                   error
		headerOneId           int64
		headerRepo            datastore.HeaderRepository
		jugFileRepo           ilk2.JugFileIlkRepository
		relevantIlkIdentifier = test_helpers.GetValidNullString(test_helpers.FakeIlk.Identifier)
		spotFileMatRepo       mat.SpotFileMatRepository
		spotFilePipRepo       pip.SpotFilePipRepository
		vatFileRepo           ilk.VatFileIlkRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		catFileChopLumpRepo = chop_lump.CatFileChopLumpRepository{}
		catFileChopLumpRepo.SetDB(db)
		catFileFlipRepo = flip.CatFileFlipRepository{}
		catFileFlipRepo.SetDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		headerOne := fakes.GetFakeHeader(1)
		headerOneId, err = headerRepo.CreateOrUpdateHeader(headerOne)
		Expect(err).NotTo(HaveOccurred())
		jugFileRepo = ilk2.JugFileIlkRepository{}
		jugFileRepo.SetDB(db)
		spotFileMatRepo = mat.SpotFileMatRepository{}
		spotFileMatRepo.SetDB(db)
		spotFilePipRepo = pip.SpotFilePipRepository{}
		spotFilePipRepo.SetDB(db)
		vatFileRepo = ilk.VatFileIlkRepository{}
		vatFileRepo.SetDB(db)
	})

	It("returns all ilk file events for ilk", func() {
		catFileChopLump := test_data.CopyModel(test_data.CatFileChopModel)
		catFileChopLump.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
		chopLumpErr := catFileChopLumpRepo.Create(headerOneId, []shared.InsertionModel{catFileChopLump})
		Expect(chopLumpErr).NotTo(HaveOccurred())

		catFileFlip := test_data.CatFileFlipModel
		catFileFlip.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
		flipErr := catFileFlipRepo.Create(headerOneId, []shared.InsertionModel{catFileFlip})
		Expect(flipErr).NotTo(HaveOccurred())

		jugFile := test_data.CopyModel(test_data.JugFileIlkModel)
		jugFile.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
		jugErr := jugFileRepo.Create(headerOneId, []shared.InsertionModel{jugFile})
		Expect(jugErr).NotTo(HaveOccurred())

		spotFileMat := test_data.CopyModel(test_data.SpotFileMatModel)
		spotFileMat.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
		spotFileMatErr := spotFileMatRepo.Create(headerOneId, []shared.InsertionModel{spotFileMat})
		Expect(spotFileMatErr).NotTo(HaveOccurred())

		spotFilePip := test_data.CopyModel(test_data.SpotFilePipModel)
		spotFilePip.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
		spotFilePipErr := spotFilePipRepo.Create(headerOneId, []shared.InsertionModel{spotFilePip})
		Expect(spotFilePipErr).NotTo(HaveOccurred())

		vatFile := test_data.CopyModel(test_data.VatFileIlkDustModel)
		vatFile.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
		vatErr := vatFileRepo.Create(headerOneId, []shared.InsertionModel{vatFile})
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
				What:          spotFileMat.ColumnValues["what"].(string),
				Data:          spotFileMat.ColumnValues["data"].(string),
			},
			test_helpers.IlkFileEvent{
				IlkIdentifier: relevantIlkIdentifier,
				What:          "pip",
				Data:          spotFilePip.ColumnValues["pip"].(string),
			},
			test_helpers.IlkFileEvent{
				IlkIdentifier: relevantIlkIdentifier,
				What:          vatFile.ColumnValues["what"].(string),
				Data:          vatFile.ColumnValues["data"].(string),
			},
		))
	})

	It("includes results across blocks", func() {
		fileBlockOne := test_data.CopyModel(test_data.VatFileIlkDustModel)
		fileBlockOne.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
		fileBlockOne.ColumnValues["data"] = strconv.Itoa(rand.Int())
		fileBlockOneErr := vatFileRepo.Create(headerOneId, []shared.InsertionModel{fileBlockOne})
		Expect(fileBlockOneErr).NotTo(HaveOccurred())

		headerTwo := fakes.GetFakeHeader(2)
		headerTwo.Hash = "anotherHash"
		headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
		Expect(headerTwoErr).NotTo(HaveOccurred())

		fileBlockTwo := test_data.CopyModel(test_data.VatFileIlkDustModel)
		fileBlockTwo.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex

		fileBlockTwo.ColumnValues["data"] = strconv.Itoa(rand.Int())
		fileBlockTwoErr := vatFileRepo.Create(headerTwoId, []shared.InsertionModel{fileBlockTwo})
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

	It("does not include ilk file events for a different ilk", func() {
		relevantFile := test_data.CopyModel(test_data.VatFileIlkDustModel)
		relevantFile.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
		relevantFile.ColumnValues["data"] = strconv.Itoa(rand.Int())

		irrelevantFile := test_data.CopyModel(test_data.VatFileIlkDustModel)
		irrelevantFile.ForeignKeyValues[constants.IlkFK] = test_helpers.AnotherFakeIlk.Hex
		irrelevantFile.ColumnValues["data"] = strconv.Itoa(rand.Int())
		irrelevantFile.ColumnValues["tx_idx"] = test_data.VatFileIlkDustModel.ColumnValues["tx_idx"].(uint) + 1

		models := []shared.InsertionModel{relevantFile, irrelevantFile}
		vatBlockOneErr := vatFileRepo.Create(headerOneId, models)
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
