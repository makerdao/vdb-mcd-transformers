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

package deal_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/mcd_transformers/transformers/events/deal"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("Flip Deal Converter", func() {
	var converter = deal.DealConverter{}
	It("converts logs to models", func() {
		models, err := converter.ToModels([]types.Log{test_data.DealLogNote})
		Expect(err).NotTo(HaveOccurred())
		Expect(models).To(Equal([]shared.InsertionModel{test_data.DealModel}))
	})

	It("returns an error if the expected amount of topics aren't in the log", func() {
		invalidLog := test_data.DealLogNote
		invalidLog.Topics = []common.Hash{}

		_, err := converter.ToModels([]types.Log{invalidLog})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError("deal log does not contain expected topics"))
	})
})
