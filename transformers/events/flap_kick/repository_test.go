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

			test_data.AssertDBRecordCount(db, "maker.flap_kick", 1)
			test_data.AssertDBRecordCount(db, "addresses", 1)

			var addressId string
			addressErr := db.Get(&addressId, `SELECT id FROM addresses`)
			Expect(addressErr).NotTo(HaveOccurred())

			var dbResult flap_kick.FlapKickModel
			getErr := db.Get(&dbResult, `SELECT bid, bid_id, lot, address_id, log_idx, tx_idx, raw_log FROM maker.flap_kick WHERE header_id = $1`, headerId)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(dbResult.Bid).To(Equal(test_data.FlapKickModel.Bid))
			Expect(dbResult.BidId).To(Equal(test_data.FlapKickModel.BidId))
			Expect(dbResult.Lot).To(Equal(test_data.FlapKickModel.Lot))
			Expect(dbResult.ContractAddress).To(Equal(addressId))
			Expect(dbResult.LogIndex).To(Equal(test_data.FlapKickModel.LogIndex))
			Expect(dbResult.TransactionIndex).To(Equal(test_data.FlapKickModel.TransactionIndex))
			Expect(dbResult.Raw).To(MatchJSON(test_data.FlapKickModel.Raw))
		})

		It("doesn't insert a new address if the flap kick insertion fails", func() {
			headerId, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())

			badFlapKick := test_data.FlapKickModel
			badFlapKick.Bid = ""
			err = flapKickRepository.Create(headerId, []interface{}{test_data.FlapKickModel, badFlapKick})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp("invalid input syntax for type numeric"))
			test_data.AssertDBRecordCount(db, "maker.flap_kick", 0)
			test_data.AssertDBRecordCount(db, "addresses", 0)
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
