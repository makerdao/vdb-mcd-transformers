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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/mcd_transformers/transformers/events/tick"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

var _ = Describe("TickConverter", func() {
	converter := tick.TickConverter{}

	Describe("ToModels", func() {
		It("converts an eth log to a db model", func() {
			models, err := converter.ToModels(constants.FlipABI(), []core.HeaderSyncLog{test_data.FlipTickHeaderSyncLog})

			Expect(err).NotTo(HaveOccurred())
			Expect(models).To(Equal([]shared.InsertionModel{test_data.TickModel}))
		})

		It("returns an error if the expected amount of topics aren't in the log", func() {
			invalidLog := test_data.FlipTickHeaderSyncLog
			invalidLog.Log.Topics = []common.Hash{}
			_, err := converter.ToModels(constants.FlipABI(), []core.HeaderSyncLog{invalidLog})

			Expect(err).To(MatchError(shared.ErrLogMissingTopics(3, 0)))
		})
	})
})
