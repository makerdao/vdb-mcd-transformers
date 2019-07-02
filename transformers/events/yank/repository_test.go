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

package yank_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/yank"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/mcd_transformers/transformers/test_data/shared_behaviors"
)

var _ = Describe("Yank repository", func() {
	var (
		db             *postgres.DB
		yankRepository yank.YankRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		yankRepository = yank.YankRepository{}
		yankRepository.SetDB(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.YankModel
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.YankChecked,
			LogEventTableName:        "maker.yank",
			TestModel:                test_data.YankModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &yankRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("adds a yank", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())

			err = yankRepository.Create(headerID, []interface{}{test_data.YankModel})
			Expect(err).NotTo(HaveOccurred())
			var dbYank yank.YankModel
			err = db.Get(&dbYank, `SELECT bid_id, contract_address, log_idx, tx_idx, raw_log FROM maker.yank WHERE header_id = $1`, headerID)

			Expect(err).NotTo(HaveOccurred())
			Expect(dbYank.BidId).To(Equal(test_data.YankModel.BidId))
			Expect(dbYank.ContractAddress).To(Equal(test_data.YankModel.ContractAddress))
			Expect(dbYank.LogIndex).To(Equal(test_data.YankModel.LogIndex))
			Expect(dbYank.TransactionIndex).To(Equal(test_data.YankModel.TransactionIndex))
			Expect(dbYank.Raw).To(MatchJSON(test_data.YankModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.YankChecked,
			Repository:              &yankRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
