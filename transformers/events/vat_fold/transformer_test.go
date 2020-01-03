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

package vat_fold_test

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_fold"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vat fold transformer", func() {
	var transformer = vat_fold.Transformer{}
	var db = test_config.NewTestDB(test_config.NewTestNode())

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("returns err if log missing topics", func() {
		badLog := core.HeaderSyncLog{}

		_, err := transformer.ToModels(constants.VatABI(), []core.HeaderSyncLog{badLog}, db)
		Expect(err).To(HaveOccurred())
	})

	It("converts a log with positive rate to an model", func() {
		models, err := transformer.ToModels(constants.VatABI(), []core.HeaderSyncLog{test_data.VatFoldHeaderSyncLogWithPositiveRate}, db)
		Expect(err).NotTo(HaveOccurred())

		var ilkID int64
		ilkErr := db.Get(&ilkID, `SELECT id FROM maker.ilks`)
		Expect(ilkErr).NotTo(HaveOccurred())
		Expect(len(models)).To(Equal(1))
		expectedModel := test_data.VatFoldModelWithPositiveRate()
		expectedModel.ColumnValues[constants.IlkColumn] = ilkID
		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})

	It("converts a log with negative rate to an model", func() {
		models, err := transformer.ToModels(constants.VatABI(), []core.HeaderSyncLog{test_data.VatFoldHeaderSyncLogWithNegativeRate}, db)
		Expect(err).NotTo(HaveOccurred())

		var ilkID int64
		ilkErr := db.Get(&ilkID, `SELECT id FROM maker.ilks`)
		Expect(ilkErr).NotTo(HaveOccurred())
		Expect(len(models)).To(Equal(1))

		expectedModel := test_data.VatFoldModelWithNegativeRate()
		expectedModel.ColumnValues[constants.IlkColumn] = ilkID
		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})
})
