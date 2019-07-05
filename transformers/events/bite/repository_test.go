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

package bite_test

import (
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/bite"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/mcd_transformers/transformers/test_data/shared_behaviors"
)

var _ = Describe("Bite repository", func() {
	var (
		biteRepository bite.BiteRepository
		db             *postgres.DB
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		biteRepository = bite.BiteRepository{}
		biteRepository.SetDB(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.BiteModel
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.BiteChecked,
			LogEventTableName:        "maker.bite",
			TestModel:                test_data.BiteModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &biteRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("persists a bite record", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())

			err = biteRepository.Create(headerID, []interface{}{test_data.BiteModel})
			Expect(err).NotTo(HaveOccurred())

			var dbBite bite.BiteModel
			err = db.Get(&dbBite, `SELECT urn_id, ink, art, tab, flip, bite_identifier, log_idx, tx_idx, raw_log FROM maker.bite WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			urnID, err := shared.GetOrCreateUrn(test_data.BiteModel.Urn, test_data.BiteModel.Ilk, db)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbBite.Urn).To(Equal(strconv.Itoa(urnID)))
			Expect(dbBite.Ink).To(Equal(test_data.BiteModel.Ink))
			Expect(dbBite.Art).To(Equal(test_data.BiteModel.Art))
			Expect(dbBite.Tab).To(Equal(test_data.BiteModel.Tab))
			Expect(dbBite.Flip).To(Equal(test_data.BiteModel.Flip))
			Expect(dbBite.Id).To(Equal(test_data.BiteModel.Id))
			Expect(dbBite.LogIndex).To(Equal(test_data.BiteModel.LogIndex))
			Expect(dbBite.TransactionIndex).To(Equal(test_data.BiteModel.TransactionIndex))
			Expect(dbBite.Raw).To(MatchJSON(test_data.BiteModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.BiteChecked,
			Repository:              &biteRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
