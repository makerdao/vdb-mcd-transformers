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

package vow_file_test

import (
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/vow_file"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/mcd_transformers/transformers/test_data/shared_behaviors"
)

var _ = Describe("Vow file repository", func() {
	var (
		db                *postgres.DB
		vowFileRepository vow_file.VowFileRepository
		headerRepository  datastore.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		vowFileRepository = vow_file.VowFileRepository{}
		vowFileRepository.SetDB(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.VowFileModel
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.VowFileChecked,
			LogEventTableName:        "maker.vow_file",
			TestModel:                test_data.VowFileModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &vowFileRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("adds a vow file event", func() {
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err = vowFileRepository.Create(headerID, []interface{}{test_data.VowFileModel})
			Expect(err).NotTo(HaveOccurred())

			var vowFile vow_file.VowFileModel
			err = db.Get(&vowFile, `SELECT what, data, log_idx, tx_idx, raw_log FROM maker.vow_file WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			assertVowFile(vowFile, test_data.VowFileModel)
		})

		It("updates the what, data and raw_log when the header_id, tx_idx, log_idx unique constraint is violated", func() {
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err = vowFileRepository.Create(headerID, []interface{}{test_data.VowFileModel})
			Expect(err).NotTo(HaveOccurred())

			var vowFile vow_file.VowFileModel
			err = db.Get(&vowFile, `SELECT what, data, log_idx, tx_idx, raw_log FROM maker.vow_file WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			assertVowFile(vowFile, test_data.VowFileModel)

			updatedVowFile := updatedVowFile()
			err = vowFileRepository.Create(headerID, []interface{}{updatedVowFile})
			Expect(err).NotTo(HaveOccurred())

			var newVowFile vow_file.VowFileModel
			err = db.Get(&newVowFile, `SELECT what, data, log_idx, tx_idx, raw_log FROM maker.vow_file WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			assertVowFile(newVowFile, updatedVowFile)
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.VowFileChecked,
			Repository:              &vowFileRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})

func assertVowFile(actual, expected vow_file.VowFileModel) {
	Expect(actual.What).To(Equal(expected.What))
	Expect(actual.Data).To(Equal(expected.Data))
	Expect(actual.LogIndex).To(Equal(expected.LogIndex))
	Expect(actual.TransactionIndex).To(Equal(expected.TransactionIndex))
	Expect(actual.Raw).To(MatchJSON(expected.Raw))
}

func updatedVowFile() vow_file.VowFileModel {
	updatedVowFileModel := test_data.VowFileModel
	updatedVowFileModel.What = "bump"
	updatedVowFileModel.Data = "1"

	rawLog := test_data.EthVowFileLog
	rawLog.Index = 1
	rawLogJson, _ := json.Marshal(rawLog)
	updatedVowFileModel.Raw = rawLogJson

	return updatedVowFileModel
}
