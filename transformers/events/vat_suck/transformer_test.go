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
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_suck"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VatSuck transformer", func() {
	var (
		transformer = vat_suck.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	It("Converts log to a model", func() {
		models, err := transformer.ToModels(constants.VatABI(), []core.EventLog{test_data.VatSuckEventLog}, db)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(models)).To(Equal(1))

		u := common.BytesToAddress(test_data.VatSuckEventLog.Log.Topics[1].Bytes()).String()
		uID, uErr := shared.GetOrCreateAddress(u, db)
		Expect(uErr).NotTo(HaveOccurred())
		v := common.BytesToAddress(test_data.VatSuckEventLog.Log.Topics[2].Bytes()).String()
		vID, vErr := shared.GetOrCreateAddress(v, db)
		Expect(vErr).NotTo(HaveOccurred())

		expectedModel := test_data.VatSuckModel()
		expectedModel.ColumnValues[constants.UColumn] = uID
		expectedModel.ColumnValues[constants.VColumn] = vID
		Expect(models[0]).To(Equal(expectedModel))
	})

	It("Returns an error if there are missing topics", func() {
		badLog := core.EventLog{
			Log: types.Log{
				Topics: []common.Hash{
					common.HexToHash("0x"),
					common.HexToHash("0x"),
					common.HexToHash("0x"),
				}},
		}

		_, err := transformer.ToModels(constants.VatABI(), []core.EventLog{badLog}, nil)

		Expect(err).To(HaveOccurred())
	})
})
