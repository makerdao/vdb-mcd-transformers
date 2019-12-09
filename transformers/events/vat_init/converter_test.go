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

package vat_init_test

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_init"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
)

var _ = Describe("Vat init converter", func() {
	var (
		converter vat_init.Converter
		db        *postgres.DB
	)

	BeforeEach(func() {
		converter = vat_init.Converter{}
		db = test_config.NewTestDB(test_config.NewTestNode())
	})

	It("returns err if log missing topics", func() {
		badLog := core.HeaderSyncLog{}
		_, err := converter.ToModels(constants.VatABI(), []core.HeaderSyncLog{badLog}, db)

		Expect(err).To(HaveOccurred())
	})

	It("converts a log to a model", func() {
		log := []core.HeaderSyncLog{test_data.VatInitHeaderSyncLog}
		models, err := converter.ToModels(constants.VatABI(), log, db)
		Expect(err).NotTo(HaveOccurred())

		ilk := log[0].Log.Topics[1].Hex()
		ilkID, ilkErr := shared.GetOrCreateIlk(ilk, db)
		Expect(ilkErr).NotTo(HaveOccurred())

		expectedModel := test_data.VatInitModel
		expectedModel.ColumnValues[constants.IlkColumn] = ilkID

		Expect(len(models)).To(Equal(1))
		Expect(models[0]).To(Equal(expectedModel))
	})
})
