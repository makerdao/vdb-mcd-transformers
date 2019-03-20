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

package ilk_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_file/ilk"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("Vat file ilk converter", func() {
	It("returns err if log is missing topics", func() {
		converter := ilk.VatFileIlkConverter{}
		badLog := types.Log{
			Data: []byte{1, 1, 1, 1, 1},
		}

		_, err := converter.ToModels([]types.Log{badLog})

		Expect(err).To(HaveOccurred())
	})

	It("returns err if log is missing data", func() {
		converter := ilk.VatFileIlkConverter{}
		badLog := types.Log{
			Topics: []common.Hash{{}, {}, {}, {}},
		}

		_, err := converter.ToModels([]types.Log{badLog})

		Expect(err).To(HaveOccurred())
	})

	It("returns error if 'what' field is unknown", func() {
		invalidWhat := hexutil.Encode([]byte("invalid"))
		log := types.Log{
			Address: test_data.EthVatFileIlkLineLog.Address,
			Topics: []common.Hash{
				test_data.EthVatFileIlkLineLog.Topics[0],
				test_data.EthVatFileIlkLineLog.Topics[1],
				common.HexToHash(invalidWhat),
				test_data.EthVatFileIlkLineLog.Topics[3],
			},
			Data:        test_data.EthVatFileIlkLineLog.Data,
			BlockNumber: test_data.EthVatFileIlkLineLog.BlockNumber,
			TxHash:      test_data.EthVatFileIlkLineLog.TxHash,
			TxIndex:     test_data.EthVatFileIlkLineLog.TxIndex,
			BlockHash:   test_data.EthVatFileIlkLineLog.BlockHash,
			Index:       test_data.EthVatFileIlkLineLog.Index,
		}
		converter := ilk.VatFileIlkConverter{}

		_, err := converter.ToModels([]types.Log{log})

		Expect(err).To(HaveOccurred())
	})

	Describe("when log is valid", func() {
		It("converts to model with data converted to ray when what is 'spot'", func() {
			converter := ilk.VatFileIlkConverter{}

			models, err := converter.ToModels([]types.Log{test_data.EthVatFileIlkSpotLog})

			Expect(err).NotTo(HaveOccurred())
			Expect(len(models)).To(Equal(1))
			Expect(models[0].(ilk.VatFileIlkModel)).To(Equal(test_data.VatFileIlkSpotModel))
		})

		It("converts to model with data converted to wad when what is 'line'", func() {
			converter := ilk.VatFileIlkConverter{}

			models, err := converter.ToModels([]types.Log{test_data.EthVatFileIlkLineLog})

			Expect(err).NotTo(HaveOccurred())
			Expect(len(models)).To(Equal(1))
			Expect(models[0].(ilk.VatFileIlkModel)).To(Equal(test_data.VatFileIlkLineModel))
		})

		It("converts to model with data converted to rad when what is 'dust'", func() {
			converter := ilk.VatFileIlkConverter{}

			models, err := converter.ToModels([]types.Log{test_data.EthVatFileIlkDustLog})

			Expect(err).NotTo(HaveOccurred())
			Expect(len(models)).To(Equal(1))
			Expect(models[0].(ilk.VatFileIlkModel)).To(Equal(test_data.VatFileIlkDustModel))
		})
	})
})
