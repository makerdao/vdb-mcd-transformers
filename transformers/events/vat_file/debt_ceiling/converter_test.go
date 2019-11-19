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

package debt_ceiling_test

import (
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_file/debt_ceiling"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
)

var _ = Describe("Vat file debt ceiling converter", func() {
	var converter = debt_ceiling.VatFileDebtCeilingConverter{}
	It("returns err if log is missing topics", func() {
		badLog := core.HeaderSyncLog{
			Log: types.Log{
				Data: []byte{1, 1, 1, 1, 1},
			}}

		_, err := converter.ToModels(constants.VatABI(), []core.HeaderSyncLog{badLog})
		Expect(err).To(HaveOccurred())
	})

	It("converts a log to an model", func() {
		models, err := converter.ToModels(constants.VatABI(), []core.HeaderSyncLog{test_data.VatFileDebtCeilingHeaderSyncLog})

		Expect(err).NotTo(HaveOccurred())
		Expect(len(models)).To(Equal(1))
		Expect(models).To(Equal([]shared.InsertionModel{test_data.VatFileDebtCeilingModel()}))
	})
})
