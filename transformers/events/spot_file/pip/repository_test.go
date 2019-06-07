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

package pip_test

import (
	"encoding/json"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/spot_file/pip"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/mcd_transformers/transformers/test_data/shared_behaviors"
)

var _ = Describe("Spot file pip repository", func() {
	var (
		db                    *postgres.DB
		headerID              int64
		spotFilePipRepository pip.SpotFilePipRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		spotFilePipRepository = pip.SpotFilePipRepository{}
		spotFilePipRepository.SetDB(db)
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		headerID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.SpotFilePipModel
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.SpotFilePipChecked,
			LogEventTableName:        "maker.spot_file_pip",
			TestModel:                test_data.SpotFilePipModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &spotFilePipRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("adds a spot file pip", func() {
			createErr := spotFilePipRepository.Create(headerID, []interface{}{test_data.SpotFilePipModel})
			Expect(createErr).NotTo(HaveOccurred())
			var dbSpotFilePip pip.SpotFilePipModel
			getErr := db.Get(&dbSpotFilePip, `SELECT ilk_id, pip, log_idx, tx_idx, raw_log FROM maker.spot_file_pip WHERE header_id = $1`, headerID)

			Expect(getErr).NotTo(HaveOccurred())
			assertSpotFilePip(dbSpotFilePip, test_data.SpotFilePipModel, db)
		})

		It("updates the ilk, pip, and raw_log when the header_id, tx_idx, log_idx unique constraint is violated", func() {
			createErr := spotFilePipRepository.Create(headerID, []interface{}{test_data.SpotFilePipModel})

			Expect(createErr).NotTo(HaveOccurred())
			var dbSpotFilePip pip.SpotFilePipModel
			getErr := db.Get(&dbSpotFilePip, `SELECT ilk_id, pip, log_idx, tx_idx, raw_log FROM maker.spot_file_pip WHERE header_id = $1`, headerID)

			Expect(getErr).NotTo(HaveOccurred())
			assertSpotFilePip(dbSpotFilePip, test_data.SpotFilePipModel, db)

			updatedSpotFilePip := updatedSpotFilePip()
			updateErr := spotFilePipRepository.Create(headerID, []interface{}{updatedSpotFilePip})
			Expect(updateErr).NotTo(HaveOccurred())

			var dbUpdatedSpotFilePip pip.SpotFilePipModel
			getUpdatedErr := db.Get(&dbUpdatedSpotFilePip, `SELECT ilk_id, pip, log_idx, tx_idx, raw_log FROM maker.spot_file_pip WHERE header_id = $1`, headerID)
			Expect(getUpdatedErr).NotTo(HaveOccurred())
			assertSpotFilePip(dbUpdatedSpotFilePip, updatedSpotFilePip, db)
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.SpotFilePipChecked,
			Repository:              &spotFilePipRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})

func assertSpotFilePip(actual, expected pip.SpotFilePipModel, db *postgres.DB) {
	ilkID, err := shared.GetOrCreateIlk(expected.Ilk, db)
	Expect(err).NotTo(HaveOccurred())
	Expect(actual.Ilk).To(Equal(strconv.Itoa(ilkID)))
	Expect(actual.Pip).To(Equal(expected.Pip))
	Expect(actual.LogIndex).To(Equal(expected.LogIndex))
	Expect(actual.TransactionIndex).To(Equal(expected.TransactionIndex))
	Expect(actual.Raw).To(MatchJSON(expected.Raw))
}

func updatedSpotFilePip() pip.SpotFilePipModel {
	updatedSpotFilePip := test_data.SpotFilePipModel
	updatedSpotFilePip.Pip = test_data.RandomString(20)

	rawLog := test_data.EthSpotFilePipLog
	rawLog.Index = rawLog.Index + 1
	rawLogJson, _ := json.Marshal(rawLog)
	updatedSpotFilePip.Raw = rawLogJson

	return updatedSpotFilePip
}
