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

package vat_suck_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"

	"github.com/vulcanize/mcd_transformers/transformers/events/vat_suck"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("VatSuck converter", func() {
	It("Converts log to a model", func() {
		converter := vat_suck.VatSuckConverter{}
		models, err := converter.ToModels([]core.HeaderSyncLog{test_data.VatSuckHeaderSyncLog})

		Expect(err).NotTo(HaveOccurred())
		Expect(len(models)).To(Equal(1))
		Expect(models[0]).To(Equal(test_data.VatSuckModel))
	})

	It("Returns an error if there are missing topics", func() {
		converter := vat_suck.VatSuckConverter{}
		badLog := core.HeaderSyncLog{
			Log: types.Log{
				Topics: []common.Hash{
					common.HexToHash("0x"),
					common.HexToHash("0x"),
					common.HexToHash("0x"),
				}},
		}
		_, err := converter.ToModels([]core.HeaderSyncLog{badLog})

		Expect(err).To(HaveOccurred())
	})
})
