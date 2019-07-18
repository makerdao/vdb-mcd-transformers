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

package flap_kick_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/flap_kick"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/mcd_transformers/transformers/test_data/shared_behaviors"
)

var _ = Describe("Flap Kick Repository", func() {
	var (
		db                 *postgres.DB
		flapKickRepository flap_kick.FlapKickRepository
		headerRepository   repositories.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		flapKickRepository = flap_kick.FlapKickRepository{}
		flapKickRepository.SetDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.FlapKickModel
		modelWithDifferentLogIdx.LogIndex = modelWithDifferentLogIdx.LogIndex + 1
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.FlapKickLabel,
			LogEventTableName:        "maker.flap_kick",
			TestModel:                test_data.FlapKickModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &flapKickRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("persists a flap kick record", func() {
			headerId, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())

			err = flapKickRepository.Create(headerId, []interface{}{test_data.FlapKickModel})
			Expect(err).NotTo(HaveOccurred())

			var count int
			countErr := db.Get(&count, `SELECT count(*) FROM maker.flap_kick`)
			Expect(countErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))

			var dbResult flap_kick.FlapKickModel
			getErr := db.Get(&dbResult, `SELECT bid, bid_id, gal, lot, contract_address, log_idx, tx_idx, raw_log FROM maker.flap_kick WHERE header_id = $1`, headerId)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(dbResult.Bid).To(Equal(test_data.FlapKickModel.Bid))
			Expect(dbResult.BidId).To(Equal(test_data.FlapKickModel.BidId))
			Expect(dbResult.Gal).To(Equal(test_data.FlapKickModel.Gal))
			Expect(dbResult.Lot).To(Equal(test_data.FlapKickModel.Lot))
			Expect(dbResult.ContractAddress).To(Equal(test_data.FlapKickModel.ContractAddress))
			Expect(dbResult.LogIndex).To(Equal(test_data.FlapKickModel.LogIndex))
			Expect(dbResult.TransactionIndex).To(Equal(test_data.FlapKickModel.TransactionIndex))
			Expect(dbResult.Raw).To(MatchJSON(test_data.FlapKickModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.FlapKickLabel,
			Repository:              &flapKickRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
