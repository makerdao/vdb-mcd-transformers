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

package flip_kick_test

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/makerdao/vdb-mcd-transformers/transformers/events/flip_kick"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
)

var _ = Describe("FlipKick Transformer", func() {
	var (
		transformer flip_kick.Transformer
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		transformer = flip_kick.Transformer{}
		test_config.CleanTestDB(db)
	})

	It("converts a log to a model", func() {
		models, err := transformer.ToModels(constants.FlipV100ABI(), []core.EventLog{test_data.FlipKickEventLog}, db)
		Expect(err).NotTo(HaveOccurred())

		expectedKick := test_data.FlipKickModel()
		addressId, addressErr := shared.GetOrCreateAddress(test_data.FlipKickEventLog.Log.Address.Hex(), db)
		Expect(addressErr).NotTo(HaveOccurred())
		expectedKick.ColumnValues[event.AddressFK] = addressId

		Expect(models).To(ConsistOf(expectedKick))
	})

	It("returns an error if converting log to entity fails", func() {
		_, err := transformer.ToModels("error abi", []core.EventLog{test_data.FlipKickEventLog}, db)

		Expect(err).To(HaveOccurred())
	})
})
