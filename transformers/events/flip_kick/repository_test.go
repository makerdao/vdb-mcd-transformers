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

package flip_kick_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/flip_kick"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/mcd_transformers/transformers/test_data/shared_behaviors"
)

var _ = Describe("FlipKick Repository", func() {
	var db *postgres.DB
	var flipKickRepository flip_kick.FlipKickRepository

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		flipKickRepository = flip_kick.FlipKickRepository{}
		flipKickRepository.SetDB(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.FlipKickModel
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.FlipKickLabel,
			LogEventTableName:        "maker.flip_kick",
			TestModel:                test_data.FlipKickModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &flipKickRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("persists flip_kick records", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			headerId, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())

			err = flipKickRepository.Create(headerId, []interface{}{test_data.FlipKickModel})
			Expect(err).NotTo(HaveOccurred())

			test_data.AssertDBRecordCount(db, "maker.flip_kick", 1)

			var addressId string
			addressErr := db.Get(&addressId, `SELECT id FROM addresses`)
			Expect(addressErr).NotTo(HaveOccurred())

			dbResult := test_data.FlipKickDBRow{}
			err = db.QueryRowx(`SELECT * FROM maker.flip_kick`).StructScan(&dbResult)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbResult.HeaderId).To(Equal(headerId))
			Expect(dbResult.BidId).To(Equal(test_data.FlipKickModel.BidId))
			Expect(dbResult.Lot).To(Equal(test_data.FlipKickModel.Lot))
			Expect(dbResult.Bid).To(Equal(test_data.FlipKickModel.Bid))
			Expect(dbResult.Tab).To(Equal(test_data.FlipKickModel.Tab))
			Expect(dbResult.Usr).To(Equal(test_data.FlipKickModel.Usr))
			Expect(dbResult.Gal).To(Equal(test_data.FlipKickModel.Gal))
			Expect(dbResult.ContractAddress).To(Equal(addressId))
			Expect(dbResult.TransactionIndex).To(Equal(test_data.FlipKickModel.TransactionIndex))
			Expect(dbResult.LogIndex).To(Equal(test_data.FlipKickModel.LogIndex))
			Expect(dbResult.Raw).To(MatchJSON(test_data.FlipKickModel.Raw))
		})

		It("doesn't insert a new address if the flip kick insertion fails", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			headerId, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())

			badFlipKick := test_data.FlipKickModel
			badFlipKick.Bid = ""
			err = flipKickRepository.Create(headerId, []interface{}{test_data.FlipKickModel, badFlipKick})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp("invalid input syntax for type numeric"))

			test_data.AssertDBRecordCount(db, "maker.flip_kick", 0)
			test_data.AssertDBRecordCount(db, "addresses", 0)
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.FlipKickLabel,
			Repository:              &flipKickRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
