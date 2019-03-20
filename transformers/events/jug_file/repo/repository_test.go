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

package repo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/jug_file/repo"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/mcd_transformers/transformers/test_data/shared_behaviors"
)

var _ = Describe("Jug file repo repository", func() {
	var (
		db                    *postgres.DB
		jugFileRepoRepository repo.JugFileRepoRepository
		headerRepository      datastore.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		jugFileRepoRepository = repo.JugFileRepoRepository{}
		jugFileRepoRepository.SetDB(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.JugFileRepoModel
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.JugFileRepoChecked,
			LogEventTableName:        "maker.jug_file_repo",
			TestModel:                test_data.JugFileRepoModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &jugFileRepoRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("adds a jug file repo event", func() {
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err = jugFileRepoRepository.Create(headerID, []interface{}{test_data.JugFileRepoModel})

			Expect(err).NotTo(HaveOccurred())
			var jugFileRepo repo.JugFileRepoModel
			err = db.Get(&jugFileRepo, `SELECT what, data, log_idx, tx_idx, raw_log FROM maker.jug_file_repo WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(jugFileRepo.What).To(Equal(test_data.JugFileRepoModel.What))
			Expect(jugFileRepo.Data).To(Equal(test_data.JugFileRepoModel.Data))
			Expect(jugFileRepo.LogIndex).To(Equal(test_data.JugFileRepoModel.LogIndex))
			Expect(jugFileRepo.TransactionIndex).To(Equal(test_data.JugFileRepoModel.TransactionIndex))
			Expect(jugFileRepo.Raw).To(MatchJSON(test_data.JugFileRepoModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.JugFileRepoChecked,
			Repository:              &jugFileRepoRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
