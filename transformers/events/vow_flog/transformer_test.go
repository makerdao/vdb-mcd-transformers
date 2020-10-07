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

package vow_flog_test

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vow_flog"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vow flog transformer", func() {
	var (
		transformer = vow_flog.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	It("returns err if log is missing topics", func() {
		badLog := core.EventLog{
			Log: types.Log{
				Data: []byte{1, 1, 1, 1, 1},
			}}

		_, err := transformer.ToModels(constants.VowABI(), []core.EventLog{badLog}, db)
		Expect(err).To(HaveOccurred())
	})

	It("converts a log to a model", func() {
		models, err := transformer.ToModels(constants.VowABI(), []core.EventLog{test_data.VowFlogEventLog}, db)

		expectedModel := test_data.VowFlogModel
		test_data.AssignMessageSenderID(test_data.VowFlogEventLog, expectedModel, db)

		Expect(err).NotTo(HaveOccurred())
		Expect(len(models)).To(Equal(1))
		Expect(models[0]).To(Equal(expectedModel))
	})
})
