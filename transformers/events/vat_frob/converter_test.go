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

package vat_frob_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_frob"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("Frob converter", func() {
	converter := vat_frob.VatFrobConverter{}

	It("returns err if log is missing topics", func() {
		badLog := types.Log{
			Data: []byte{1, 1, 1, 1, 1},
		}

		_, err := converter.ToModels([]types.Log{badLog})
		Expect(err).To(HaveOccurred())
	})

	It("returns err if log is missing data", func() {
		badLog := types.Log{
			Topics: []common.Hash{{}, {}, {}, {}},
		}

		_, err := converter.ToModels([]types.Log{badLog})
		Expect(err).To(HaveOccurred())
	})

	It("converts a log with positive dart to a model", func() {
		models, err := converter.ToModels([]types.Log{test_data.EthVatFrobLogWithPositiveDart})

		Expect(err).NotTo(HaveOccurred())
		Expect(len(models)).To(Equal(1))
		Expect(models).To(Equal([]shared.InsertionModel{test_data.VatFrobModelWithPositiveDart}))
	})

	It("converts a log with negative dink to a model", func() {
		models, err := converter.ToModels([]types.Log{test_data.EthVatFrobLogWithNegativeDink})

		Expect(err).NotTo(HaveOccurred())
		Expect(len(models)).To(Equal(1))
		Expect(models).To(Equal([]shared.InsertionModel{test_data.VatFrobModelWithNegativeDink}))
	})
})
