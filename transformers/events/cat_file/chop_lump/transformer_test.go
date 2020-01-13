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

package chop_lump_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/cat_file/chop_lump"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cat file chop lump transformer", func() {
	var (
		transformer = chop_lump.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	Context("chop events", func() {
		It("converts a chop log to a model", func() {
			models, err := transformer.ToModels(constants.CatABI(), []core.EventLog{test_data.CatFileChopEventLog}, db)
			Expect(err).NotTo(HaveOccurred())

			var ilkID int64
			ilkErr := db.Get(&ilkID, `SELECT id FROM maker.ilks where ilk = $1`, test_data.CatFileChopEventLog.Log.Topics[2].Hex())
			Expect(ilkErr).NotTo(HaveOccurred())
			expectedModel := test_data.CatFileChopModel()
			expectedModel.ColumnValues[constants.IlkColumn] = ilkID

			Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
		})
	})

	Context("lump events", func() {
		It("converts a lump log to a model", func() {
			models, err := transformer.ToModels(constants.CatABI(), []core.EventLog{test_data.CatFileLumpEventLog}, db)
			Expect(err).NotTo(HaveOccurred())

			var ilkID int64
			ilkErr := db.Get(&ilkID, `SELECT id FROM maker.ilks where ilk = $1`, test_data.CatFileLumpEventLog.Log.Topics[2].Hex())
			Expect(ilkErr).NotTo(HaveOccurred())
			expectedModel := test_data.CatFileLumpModel()
			expectedModel.ColumnValues[constants.IlkColumn] = ilkID

			Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
		})
	})

	It("returns err if log is missing topics", func() {
		badLog := core.EventLog{
			Log: types.Log{
				Data: []byte{1, 1, 1, 1, 1},
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
})
