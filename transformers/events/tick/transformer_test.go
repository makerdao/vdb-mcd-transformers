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

package tick_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/tick"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tick transformer", func() {
	transformer := tick.Transformer{}
	db := test_config.NewTestDB(test_config.NewTestNode())

	Describe("ToModels", func() {
		It("converts an eth log to a db model", func() {
			models, err := transformer.ToModels(constants.FlipABI(), []core.HeaderSyncLog{test_data.FlipTickHeaderSyncLog}, db)
			Expect(err).NotTo(HaveOccurred())

			var addressID int64
			addrErr := db.Get(&addressID, `SELECT id FROM public.addresses`)
			Expect(addrErr).NotTo(HaveOccurred())

			expectedModel := test_data.TickModel()
			expectedModel.ColumnValues[event.AddressFK] = addressID

			Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
		})

		It("returns an error if the expected amount of topics aren't in the log", func() {
			invalidLog := test_data.FlipTickHeaderSyncLog
			invalidLog.Log.Topics = []common.Hash{}
			_, err := transformer.ToModels(constants.FlipABI(), []core.HeaderSyncLog{invalidLog}, db)

			Expect(err).To(MatchError(shared.ErrLogMissingTopics(3, 0)))
		})
	})
})
