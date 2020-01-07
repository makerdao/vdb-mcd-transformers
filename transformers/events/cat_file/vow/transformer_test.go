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

package vow_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/cat_file/vow"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cat file vow transformer", func() {
	var transformer vow.Transformer

	BeforeEach(func() {
		transformer = vow.Transformer{}
	})

	It("returns err if log is missing topics", func() {
		badLog := core.HeaderSyncLog{
			Log: types.Log{
				Data: []byte{1, 1, 1, 1, 1},
			},
		}

		_, err := transformer.ToModels(constants.CatABI(), []core.HeaderSyncLog{badLog}, nil)
		Expect(err).To(HaveOccurred())
	})

	It("returns err if log is missing data", func() {
		badLog := core.HeaderSyncLog{
			Log: types.Log{
				Topics: []common.Hash{{}, {}, {}, {}},
			},
		}

		_, err := transformer.ToModels(constants.CatABI(), []core.HeaderSyncLog{badLog}, nil)
		Expect(err).To(HaveOccurred())
	})

	It("converts a log to an model", func() {
		models, err := transformer.ToModels(constants.CatABI(), []core.HeaderSyncLog{test_data.CatFileVowHeaderSyncLog}, nil)

		Expect(err).NotTo(HaveOccurred())
		Expect(models).To(Equal([]event.InsertionModel{test_data.CatFileVowModel()}))
	})
})
