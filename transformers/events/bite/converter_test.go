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

package bite_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/event"

	"github.com/vulcanize/mcd_transformers/transformers/events/bite"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

var _ = Describe("Bite Converter", func() {
	db := test_config.NewTestDB(test_config.NewTestNode())
	var converter = bite.Converter{}
	converter.SetDB(db)

	It("converts a log to a Model", func() {
		models, err := converter.ToModels(constants.CatABI(), []core.HeaderSyncLog{test_data.BiteHeaderSyncLog})
		Expect(err).NotTo(HaveOccurred())

		var urnID int64
		urnErr := db.Get(&urnID, `SELECT id FROM maker.urns`)
		Expect(urnErr).NotTo(HaveOccurred())
		expectedModel := test_data.BiteModel()
		expectedModel.ColumnValues[constants.UrnColumn] = urnID

		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})

	It("returns an error if converting log to entity fails", func() {
		_, err := converter.ToModels("error abi", []core.HeaderSyncLog{test_data.BiteHeaderSyncLog})

		Expect(err).To(HaveOccurred())
	})
})
