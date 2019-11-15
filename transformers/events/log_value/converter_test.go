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

package log_value_test

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_value"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogValue Converter", func() {
	db := test_config.NewTestDB(test_config.NewTestNode())
	var converter = log_value.Converter{}
	converter.SetDB(db)

	It("converts a log to a Model", func() {
		models, err := converter.ToModels(constants.OsmABI(), []core.HeaderSyncLog{test_data.LogValueHeaderSyncLog})
		Expect(err).NotTo(HaveOccurred())

		Expect(models).To(Equal([]event.InsertionModel{test_data.LogValueModel()}))
	})

	It("returns an error if converting log to entity fails", func() {
		_, err := converter.ToModels("error abi", []core.HeaderSyncLog{test_data.LogValueHeaderSyncLog})

		Expect(err).To(HaveOccurred())
	})
})
