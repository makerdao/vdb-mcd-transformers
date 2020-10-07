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

package tend_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/tend"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tend transformer", func() {
	var (
		transformer = tend.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	Describe("ToModels", func() {
		It("converts an eth log to a db model", func() {
			models, err := transformer.ToModels(constants.FlipV100ABI(), []core.EventLog{test_data.TendEventLog}, db)
			Expect(err).NotTo(HaveOccurred())

			expectedModel := test_data.TendModel()
			test_data.AssignAddressID(test_data.TendEventLog, expectedModel, db)
			test_data.AssignMessageSenderID(test_data.TendEventLog, expectedModel, db)

			Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
		})

		It("returns an error if the log data is empty", func() {
			emptyDataLog := test_data.TendEventLog
			emptyDataLog.Log.Data = []byte{}
			_, err := transformer.ToModels(constants.FlipV100ABI(), []core.EventLog{emptyDataLog}, db)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(shared.ErrLogMissingData))
		})

		It("returns an error if the expected amount of topics aren't in the log", func() {
			invalidLog := test_data.TendEventLog
			invalidLog.Log.Topics = []common.Hash{}
			_, err := transformer.ToModels(constants.FlipV100ABI(), []core.EventLog{invalidLog}, db)

			Expect(err).To(MatchError(shared.ErrLogMissingTopics(4, 0)))
		})
	})
})
