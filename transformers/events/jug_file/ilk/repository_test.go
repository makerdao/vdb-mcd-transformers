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

package ilk_test

import (
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/jug_file/ilk"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/mcd_transformers/transformers/test_data/shared_behaviors"
)

var _ = Describe("Jug file ilk repository", func() {
	var (
		db                   *postgres.DB
		jugFileIlkRepository ilk.JugFileIlkRepository
		headerRepository     datastore.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		jugFileIlkRepository = ilk.JugFileIlkRepository{}
		jugFileIlkRepository.SetDB(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.JugFileIlkModel
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.JugFileIlkChecked,
			LogEventTableName:        "maker.jug_file_ilk",
			TestModel:                test_data.JugFileIlkModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &jugFileIlkRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("adds a jug file ilk event", func() {
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err = jugFileIlkRepository.Create(headerID, []interface{}{test_data.JugFileIlkModel})

			Expect(err).NotTo(HaveOccurred())
			var dbJugFileIlk ilk.JugFileIlkModel
			err = db.Get(&dbJugFileIlk, `SELECT ilk_id, what, data, log_idx, tx_idx, raw_log FROM maker.jug_file_ilk WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(test_data.JugFileIlkModel.Ilk, db)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbJugFileIlk.Ilk).To(Equal(strconv.Itoa(ilkID)))
			Expect(dbJugFileIlk.What).To(Equal(test_data.JugFileIlkModel.What))
			Expect(dbJugFileIlk.Data).To(Equal(test_data.JugFileIlkModel.Data))
			Expect(dbJugFileIlk.LogIndex).To(Equal(test_data.JugFileIlkModel.LogIndex))
			Expect(dbJugFileIlk.TransactionIndex).To(Equal(test_data.JugFileIlkModel.TransactionIndex))
			Expect(dbJugFileIlk.Raw).To(MatchJSON(test_data.JugFileIlkModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.JugFileIlkChecked,
			Repository:              &jugFileIlkRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
