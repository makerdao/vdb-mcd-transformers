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

package pot_join_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/pot_join"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PotJoin transformer", func() {
	var (
		transformer = pot_join.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("converts log to a model", func() {
		models, err := transformer.ToModels(constants.PotABI(), []core.EventLog{test_data.PotJoinEventLog}, db)
		Expect(err).NotTo(HaveOccurred())

		expectedModel := test_data.PotJoinModel()
		test_data.AssignMessageSenderID(test_data.PotJoinEventLog, expectedModel, db)
		Expect(models).To(ConsistOf(expectedModel))
	})

	It("returns an error if there are missing topics", func() {
		invalidLog := test_data.PotJoinEventLog
		invalidLog.Log.Topics = []common.Hash{}

		_, err := transformer.ToModels(constants.PotABI(), []core.EventLog{invalidLog}, db)

		Expect(err).To(MatchError(shared.ErrLogMissingTopics(3, 0)))
	})
})
