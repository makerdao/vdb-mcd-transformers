//  VulcanizeDB
//  Copyright © 2019 Vulcanize
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU Affero General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  This program is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU Affero General Public License for more details.
//
//  You should have received a copy of the GNU Affero General Public License
//  along with this program.  If not, see <http://www.gnu.org/licenses/>.

package spot_poke_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/spot_poke"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/mcd_transformers/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"strconv"
)

var _ = Describe("Spot Poke repository", func() {
	var (
		repository spot_poke.SpotPokeRepository
		db         *postgres.DB
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repository = spot_poke.SpotPokeRepository{}
		repository.SetDB(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.SpotPokeModel
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.SpotPokeLabel,
			LogEventTableName:        "maker.spot_poke",
			TestModel:                test_data.SpotPokeModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &repository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("persists a spot poke record", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())

			err = repository.Create(headerID, []interface{}{test_data.SpotPokeModel})
			Expect(err).NotTo(HaveOccurred())

			ilkID, err := shared.GetOrCreateIlk(test_data.SpotPokeModel.Ilk, db)
			Expect(err).NotTo(HaveOccurred())

			var dbSpotPoke spot_poke.SpotPokeModel
			err = db.Get(&dbSpotPoke, `SELECT ilk_id, value, spot, log_idx, tx_idx, raw_log FROM maker.spot_poke WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbSpotPoke.Ilk).To(Equal(strconv.Itoa(ilkID)))
			Expect(dbSpotPoke.Value).To(Equal(test_data.SpotPokeModel.Value))
			Expect(dbSpotPoke.Spot).To(Equal(test_data.SpotPokeModel.Spot))
			Expect(dbSpotPoke.LogIndex).To(Equal(test_data.SpotPokeModel.LogIndex))
			Expect(dbSpotPoke.TransactionIndex).To(Equal(test_data.SpotPokeModel.TransactionIndex))
			Expect(dbSpotPoke.Raw).To(MatchJSON(test_data.SpotPokeModel.Raw))
		})

		It("rolls back the transaction if insertion fails", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())

			badSpotPokeModel := spot_poke.SpotPokeModel{}

			err = repository.Create(headerID, []interface{}{test_data.SpotPokeModel, badSpotPokeModel})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid input syntax for type numeric"))

			var spotPokeCount int
			err = db.Get(&spotPokeCount, `SELECT count(*) FROM maker.spot_poke`)
			Expect(err).NotTo(HaveOccurred())
			Expect(spotPokeCount).To(Equal(0))

			var ilkCount int
			err = db.Get(&ilkCount, `SELECT count(*) FROM maker.ilks`)
			Expect(err).NotTo(HaveOccurred())
			Expect(ilkCount).To(Equal(0))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.BiteLabel,
			Repository:              &repository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
