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

package mat_test

import (
	"encoding/json"
	"math/big"
	"math/rand"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/spot_file/mat"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/mcd_transformers/transformers/test_data/shared_behaviors"
)

var _ = Describe("Spot file mat repository", func() {
	var (
		db                    *postgres.DB
		spotFileMatRepository mat.SpotFileMatRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		spotFileMatRepository = mat.SpotFileMatRepository{}
		spotFileMatRepository.SetDB(db)
	})

	Describe("Create", func() {
		var (
			headerID         int64
			headerRepository datastore.HeaderRepository
		)

		BeforeEach(func() {
			headerRepository = repositories.NewHeaderRepository(db)
			var err error
			headerID, err = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
		})

		modelWithDifferentLogIdx := test_data.SpotFileMatModel
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.SpotFileMatChecked,
			LogEventTableName:        "maker.spot_file_mat",
			TestModel:                test_data.SpotFileMatModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &spotFileMatRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("adds a spot file mat", func() {
			createErr := spotFileMatRepository.Create(headerID, []interface{}{test_data.SpotFileMatModel})

			Expect(createErr).NotTo(HaveOccurred())
			var dbSpotFileMat mat.SpotFileMatModel
			getErr := db.Get(&dbSpotFileMat, `SELECT ilk_id, what, data, log_idx, tx_idx, raw_log FROM maker.spot_file_mat WHERE header_id = $1`, headerID)

			Expect(getErr).NotTo(HaveOccurred())
			assertSpotFileMat(dbSpotFileMat, test_data.SpotFileMatModel, db)
		})

		It("updates the ilk, what, data and raw_log when the header_id, tx_idx, log_idx unique constraint is violated", func() {
			createErr := spotFileMatRepository.Create(headerID, []interface{}{test_data.SpotFileMatModel})

			Expect(createErr).NotTo(HaveOccurred())
			var dbSpotFileMat mat.SpotFileMatModel
			getErr := db.Get(&dbSpotFileMat, `SELECT ilk_id, what, data, log_idx, tx_idx, raw_log FROM maker.spot_file_mat WHERE header_id = $1`, headerID)

			Expect(getErr).NotTo(HaveOccurred())
			assertSpotFileMat(dbSpotFileMat, test_data.SpotFileMatModel, db)

			updatedSpotFileMat := updatedSpotFileMat()
			updateErr := spotFileMatRepository.Create(headerID, []interface{}{updatedSpotFileMat})
			Expect(updateErr).NotTo(HaveOccurred())

			var dbUpdatedSpotFileMat mat.SpotFileMatModel
			getUpdatedErr := db.Get(&dbUpdatedSpotFileMat, `SELECT ilk_id, what, data, log_idx, tx_idx, raw_log FROM maker.spot_file_mat WHERE header_id = $1`, headerID)
			Expect(getUpdatedErr).NotTo(HaveOccurred())
			assertSpotFileMat(dbUpdatedSpotFileMat, updatedSpotFileMat, db)
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.SpotFileMatChecked,
			Repository:              &spotFileMatRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})

func assertSpotFileMat(actual, expected mat.SpotFileMatModel, db *postgres.DB) {
	ilkID, err := shared.GetOrCreateIlk(expected.Ilk, db)
	Expect(err).NotTo(HaveOccurred())
	Expect(actual.Ilk).To(Equal(strconv.Itoa(ilkID)))
	Expect(actual.What).To(Equal(expected.What))
	Expect(actual.Data).To(Equal(expected.Data))
	Expect(actual.LogIndex).To(Equal(expected.LogIndex))
	Expect(actual.TransactionIndex).To(Equal(expected.TransactionIndex))
	Expect(actual.Raw).To(MatchJSON(expected.Raw))
}

func updatedSpotFileMat() mat.SpotFileMatModel {
	updatedSpotFileMat := test_data.SpotFileMatModel
	updatedSpotFileMat.What = "mat"
	updatedSpotFileMat.Data = big.NewInt(rand.Int63()).String()

	rawLog := test_data.EthSpotFileMatLog
	rawLog.Index = rawLog.Index + 1
	rawLogJson, _ := json.Marshal(rawLog)
	updatedSpotFileMat.Raw = rawLogJson

	return updatedSpotFileMat
}
