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

package yank_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/yank"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Yank Transformer", func() {
	var (
		transformer = yank.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	It("converts logs to models", func() {
		models, err := transformer.ToModels(constants.FlipABI(), []core.EventLog{test_data.YankEventLog}, db)
		var addressID int64
		addrErr := db.Get(&addressID, `SELECT id FROM public.addresses`)
		Expect(addrErr).NotTo(HaveOccurred())
		expectedModel := test_data.YankModel()
		expectedModel.ColumnValues[event.AddressFK] = addressID
		Expect(err).NotTo(HaveOccurred())
		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})

	It("returns an error if the expected topics aren't in the log", func() {
		invalidLog := test_data.YankEventLog
		invalidLog.Log.Topics = []common.Hash{}

		_, err := transformer.ToModels(constants.FlipABI(), []core.EventLog{invalidLog}, db)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(shared.ErrLogMissingTopics(3, 0)))
	})
})
