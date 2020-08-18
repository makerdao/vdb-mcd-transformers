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

package flip_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/cat_file/flip"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cat file flip transformer", func() {
	var (
		transformer = flip.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("returns err if log is missing topics", func() {
		badLog := core.EventLog{
			Log: types.Log{
				Data: []byte{1, 1, 1, 1, 2},
			},
		}

		_, err := transformer.ToModels(constants.CatABI(), []core.EventLog{badLog}, db)
		Expect(err).To(HaveOccurred())
	})

	It("returns err if log is missing data", func() {
		badLog := core.EventLog{
			Log: types.Log{
				Topics: []common.Hash{{}, {}, {}, {}},
			},
		}

		_, err := transformer.ToModels(constants.CatABI(), []core.EventLog{badLog}, db)
		Expect(err).To(HaveOccurred())
	})

	It("converts a log to an model", func() {
		models, err := transformer.ToModels(constants.CatABI(), []core.EventLog{test_data.CatFileFlipEventLog}, db)
		Expect(err).NotTo(HaveOccurred())

		ilkID, ilkIDErr := shared.GetOrCreateIlk(test_data.CatFileFlipEventLog.Log.Topics[2].Hex(), db)
		Expect(ilkIDErr).NotTo(HaveOccurred())

		expectedModel := test_data.CatFileFlipModel()
		test_data.AssignAddressID(test_data.CatFileFlipEventLog, expectedModel, db)
		test_data.AssignMessageSenderID(test_data.CatFileFlipEventLog, expectedModel, db)
		expectedModel.ColumnValues[constants.IlkColumn] = ilkID

		Expect(len(models)).To(Equal(1))
		Expect(models[0]).To(Equal(expectedModel))
	})
})
