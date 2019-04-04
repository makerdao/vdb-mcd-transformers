// VulcanizeDB
// Copyright © 2018 Vulcanize

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

package vow_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/cat_file/vow"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/mcd_transformers/transformers/test_data/shared_behaviors"
)

var _ = Describe("Cat file vow repository", func() {
	var (
		catFileVowRepository vow.CatFileVowRepository
		db                   *postgres.DB
		headerRepository     datastore.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		catFileVowRepository = vow.CatFileVowRepository{}
		catFileVowRepository.SetDB(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.CatFileVowModel
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.CatFileVowChecked,
			LogEventTableName:        "maker.cat_file_vow",
			TestModel:                test_data.CatFileVowModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &catFileVowRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("adds a cat file vow event", func() {
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err = catFileVowRepository.Create(headerID, []interface{}{test_data.CatFileVowModel})

			Expect(err).NotTo(HaveOccurred())
			var dbResult vow.CatFileVowModel
			err = db.Get(&dbResult, `SELECT what, data, tx_idx, log_idx, raw_log FROM maker.cat_file_vow WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbResult.What).To(Equal(test_data.CatFileVowModel.What))
			Expect(dbResult.Data).To(Equal(test_data.CatFileVowModel.Data))
			Expect(dbResult.TransactionIndex).To(Equal(test_data.CatFileVowModel.TransactionIndex))
			Expect(dbResult.LogIndex).To(Equal(test_data.CatFileVowModel.LogIndex))
			Expect(dbResult.Raw).To(MatchJSON(test_data.CatFileVowModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.CatFileVowChecked,
			Repository:              &catFileVowRepository,
		}
		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
