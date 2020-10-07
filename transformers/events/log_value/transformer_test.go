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

package log_value_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_value"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogValue Transformer", func() {
	db := test_config.NewTestDB(test_config.NewTestNode())
	var transformer = log_value.Transformer{}

	It("converts a log to a Model", func() {
		models, err := transformer.ToModels(constants.OsmABI(), []core.EventLog{test_data.LogValueEventLog}, db)
		Expect(err).NotTo(HaveOccurred())

		expectedModel := test_data.LogValueModel()
		test_data.AssignAddressID(test_data.LogValueEventLog, expectedModel, db)

		Expect(models[0]).To(Equal(expectedModel))
	})

	It("returns an error if the log is missing a topic", func() {
		incompleteLog := core.EventLog{}
		_, err := transformer.ToModels(constants.OsmABI(), []core.EventLog{incompleteLog}, db)
		Expect(err).To(HaveOccurred())
	})

	It("returns err if log is missing data", func() {
		badLog := core.EventLog{
			Log: types.Log{
				Topics: []common.Hash{{}, {}, {}, {}},
			}}

		_, err := transformer.ToModels(constants.OsmABI(), []core.EventLog{badLog}, db)
		Expect(err).To(HaveOccurred())
	})
})
