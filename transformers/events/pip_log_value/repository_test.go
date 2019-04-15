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

package pip_log_value_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/pip_log_value"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/mcd_transformers/transformers/test_data/shared_behaviors"
)

var _ = Describe("Pip LogValue repository", func() {
	var (
		db                    *postgres.DB
		pipLogValueRepository pip_log_value.PipLogValueRepository
		headerRepository      repositories.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		pipLogValueRepository = pip_log_value.PipLogValueRepository{}
		pipLogValueRepository.SetDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.PipLogValueModel
		modelWithDifferentLogIdx.LogIndex = modelWithDifferentLogIdx.LogIndex + 1
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.PipLogValueChecked,
			LogEventTableName:        "maker.pip_log_value",
			TestModel:                test_data.PipLogValueModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &pipLogValueRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("persists a pip log value", func() {
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err = pipLogValueRepository.Create(headerID, []interface{}{test_data.PipLogValueModel})

			Expect(err).NotTo(HaveOccurred())
			var dbPipLogValue pip_log_value.PipLogValueModel
			err = db.Get(&dbPipLogValue, `SELECT block_number, contract_address, val, log_idx, tx_idx, raw_log FROM maker.pip_log_value WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbPipLogValue.BlockNumber).To(Equal(test_data.PipLogValueModel.BlockNumber))
			Expect(dbPipLogValue.ContractAddress).To(Equal(test_data.PipLogValueModel.ContractAddress))
			Expect(dbPipLogValue.Value).To(Equal(test_data.PipLogValueModel.Value))
			Expect(dbPipLogValue.LogIndex).To(Equal(test_data.PipLogValueModel.LogIndex))
			Expect(dbPipLogValue.TransactionIndex).To(Equal(test_data.PipLogValueModel.TransactionIndex))
			Expect(dbPipLogValue.Raw).To(MatchJSON(test_data.PipLogValueModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.PipLogValueChecked,
			Repository:              &pipLogValueRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
