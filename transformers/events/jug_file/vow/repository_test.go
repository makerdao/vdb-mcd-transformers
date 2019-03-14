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

package vow_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/jug_file/vow"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/mcd_transformers/transformers/test_data/shared_behaviors"
)

var _ = Describe("Jug file vow repository", func() {
	var (
		db                   *postgres.DB
		jugFileVowRepository vow.JugFileVowRepository
		headerRepository     datastore.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		jugFileVowRepository = vow.JugFileVowRepository{}
		jugFileVowRepository.SetDB(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.JugFileVowModel
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.JugFileVowChecked,
			LogEventTableName:        "maker.jug_file_vow",
			TestModel:                test_data.JugFileVowModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &jugFileVowRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("adds a jug file vow event", func() {
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err = jugFileVowRepository.Create(headerID, []interface{}{test_data.JugFileVowModel})

			Expect(err).NotTo(HaveOccurred())
			var dbJugFileVow vow.JugFileVowModel
			err = db.Get(&dbJugFileVow, `SELECT what, data, log_idx, tx_idx, raw_log FROM maker.jug_file_vow WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbJugFileVow.What).To(Equal(test_data.JugFileVowModel.What))
			Expect(dbJugFileVow.Data).To(Equal(test_data.JugFileVowModel.Data))
			Expect(dbJugFileVow.LogIndex).To(Equal(test_data.JugFileVowModel.LogIndex))
			Expect(dbJugFileVow.TransactionIndex).To(Equal(test_data.JugFileVowModel.TransactionIndex))
			Expect(dbJugFileVow.Raw).To(MatchJSON(test_data.JugFileVowModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.JugFileVowChecked,
			Repository:              &jugFileVowRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
