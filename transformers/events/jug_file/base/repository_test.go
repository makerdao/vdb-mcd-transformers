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

package base_test

import (
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"math/big"
	"math/rand"

	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/jug_file/base"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/mcd_transformers/transformers/test_data/shared_behaviors"
)

var _ = Describe("Jug file base repository", func() {
	var (
		db                    *postgres.DB
		jugFileBaseRepository base.JugFileBaseRepository
		headerRepository      datastore.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		jugFileBaseRepository = base.JugFileBaseRepository{}
		jugFileBaseRepository.SetDB(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.JugFileBaseModel
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.JugFileBaseChecked,
			LogEventTableName:        "maker.jug_file_base",
			TestModel:                test_data.JugFileBaseModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &jugFileBaseRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("adds a jug file base event", func() {
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err = jugFileBaseRepository.Create(headerID, []interface{}{test_data.JugFileBaseModel})

			Expect(err).NotTo(HaveOccurred())
			var jugFileBase base.JugFileBaseModel
			err = db.Get(&jugFileBase, `SELECT what, data, log_idx, tx_idx, raw_log FROM maker.jug_file_base WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			assertJugFileBase(jugFileBase, test_data.JugFileBaseModel)
		})

		It("updates the what, data and raw_log when the header_id, tx_idx, log_idx unique constraint is violated", func() {
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err = jugFileBaseRepository.Create(headerID, []interface{}{test_data.JugFileBaseModel})
			Expect(err).NotTo(HaveOccurred())

			var jugFileBase base.JugFileBaseModel
			err = db.Get(&jugFileBase, `SELECT what, data, log_idx, tx_idx, raw_log FROM maker.jug_file_base WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			assertJugFileBase(jugFileBase, test_data.JugFileBaseModel)

			updatedJugFileBase := updatedJugFileBase()
			err = jugFileBaseRepository.Create(headerID, []interface{}{updatedJugFileBase})
			Expect(err).NotTo(HaveOccurred())

			var updatedJugFileBaseRecord base.JugFileBaseModel
			err = db.Get(&updatedJugFileBaseRecord, `SELECT what, data, log_idx, tx_idx, raw_log FROM maker.jug_file_base WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			assertJugFileBase(updatedJugFileBaseRecord, updatedJugFileBase)
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.JugFileBaseChecked,
			Repository:              &jugFileBaseRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})

func assertJugFileBase(actual, expected base.JugFileBaseModel) {
	Expect(actual.What).To(Equal(expected.What))
	Expect(actual.Data).To(Equal(expected.Data))
	Expect(actual.LogIndex).To(Equal(expected.LogIndex))
	Expect(actual.TransactionIndex).To(Equal(expected.TransactionIndex))
	Expect(actual.Raw).To(MatchJSON(expected.Raw))
}

func updatedJugFileBase() base.JugFileBaseModel {
	updatedJugFileBase := test_data.JugFileBaseModel
	updatedJugFileBase.What = "base"
	updatedJugFileBase.Data = big.NewInt(rand.Int63()).String()

	rawLog := test_data.EthJugFileBaseLog
	rawLog.Index = 1
	rawLogJson, _ := json.Marshal(rawLog)
	updatedJugFileBase.Raw = rawLogJson

	return updatedJugFileBase
}
