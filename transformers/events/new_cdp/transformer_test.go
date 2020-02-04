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

package new_cdp_test

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/new_cdp"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewCdp Transformer", func() {
	var (
		transformer = new_cdp.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("converts a log to a Model", func() {
		models, err := transformer.ToModels(constants.CdpManagerABI(), []core.EventLog{test_data.NewCdpEventLog}, db)
		Expect(err).NotTo(HaveOccurred())

		var usrAddressID, ownAddressID int64
		usrAddressID, usrAddressErr := shared.GetOrCreateAddress(test_data.NewCdpEventLog.Log.Topics[1].Hex(), db)
		Expect(usrAddressErr).NotTo(HaveOccurred())
		ownAddressID, ownAddressErr := shared.GetOrCreateAddress(test_data.NewCdpEventLog.Log.Topics[2].Hex(), db)
		Expect(ownAddressErr).NotTo(HaveOccurred())
		expectedModel := test_data.NewCdpModel()
		expectedModel.ColumnValues[constants.UsrColumn] = usrAddressID
		expectedModel.ColumnValues[constants.OwnColumn] = ownAddressID

		Expect(models).To(ConsistOf(expectedModel))
	})

	It("returns an error if converting log to entity fails", func() {
		_, err := transformer.ToModels("error abi", []core.EventLog{test_data.NewCdpEventLog}, nil)
		Expect(err).To(HaveOccurred())
	})
})
