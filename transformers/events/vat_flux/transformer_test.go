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

package vat_flux_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_flux"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VatFlux transformer", func() {
	var (
		transformer = vat_flux.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)
	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("Converts logs to models", func() {
		log := []core.HeaderSyncLog{test_data.VatFluxHeaderSyncLog}
		models, err := transformer.ToModels(constants.VatABI(), log, db)
		Expect(err).NotTo(HaveOccurred())

		ilk := log[0].Log.Topics[1].Hex()
		ilkID, ilkErr := shared.GetOrCreateIlk(ilk, db)
		Expect(ilkErr).NotTo(HaveOccurred())

		expectedModel := test_data.VatFluxModel()
		expectedModel.ColumnValues[constants.IlkColumn] = ilkID
		Expect(models[0]).To(Equal(expectedModel))
	})

	It("Returns an error there are missing topics", func() {
		badLog := core.HeaderSyncLog{
			Log: types.Log{
				Topics: []common.Hash{
					common.HexToHash("0x"),
					common.HexToHash("0x"),
					common.HexToHash("0x"),
				}},
		}

		_, err := transformer.ToModels(constants.VatABI(), []core.HeaderSyncLog{badLog}, db)
		Expect(err).To(HaveOccurred())
	})
})
