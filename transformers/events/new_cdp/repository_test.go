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

package new_cdp_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/new_cdp"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/mcd_transformers/transformers/test_data/shared_behaviors"
)

var _ = Describe("NewCdp Repository", func() {
	var (
		db               *postgres.DB
		newCdpRepository new_cdp.NewCdpRepository
		headerRepository repositories.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		newCdpRepository = new_cdp.NewCdpRepository{}
		newCdpRepository.SetDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.NewCdpModel
		modelWithDifferentLogIdx.LogIndex = modelWithDifferentLogIdx.LogIndex + 1
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.NewCdpLabel,
			LogEventTableName:        "maker.new_cdp",
			TestModel:                test_data.NewCdpModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &newCdpRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("persists a new_cdp record", func() {
			headerId, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())

			err = newCdpRepository.Create(headerId, []interface{}{test_data.NewCdpModel})
			Expect(err).NotTo(HaveOccurred())

			var count int
			countErr := db.Get(&count, `SELECT count(*) FROM maker.new_cdp`)
			Expect(countErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))

			var dbNewCdp new_cdp.NewCdpModel
			newCdpQuery := `SELECT usr, own, cdp, tx_idx, log_idx, raw_log FROM maker.new_cdp WHERE header_id = $1`
			getErr := db.Get(&dbNewCdp, newCdpQuery, headerId)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(dbNewCdp.Usr).To(Equal(test_data.NewCdpModel.Usr))
			Expect(dbNewCdp.Own).To(Equal(test_data.NewCdpModel.Own))
			Expect(dbNewCdp.Cdp).To(Equal(test_data.NewCdpModel.Cdp))
			Expect(dbNewCdp.LogIndex).To(Equal(test_data.NewCdpModel.LogIndex))
			Expect(dbNewCdp.TransactionIndex).To(Equal(test_data.NewCdpModel.TransactionIndex))
			Expect(dbNewCdp.Raw).To(MatchJSON(test_data.NewCdpModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.NewCdpLabel,
			Repository:              &newCdpRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
