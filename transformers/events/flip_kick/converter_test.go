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

package flip_kick_test

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/mcd_transformers/transformers/events/flip_kick"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("FlipKick Converter", func() {
	var converter = flip_kick.FlipKickConverter{}
	It("converts a log to a model", func() {
		models, err := converter.ToModels(constants.FlipABI(), []core.HeaderSyncLog{test_data.FlipKickHeaderSyncLog})

		Expect(err).NotTo(HaveOccurred())

		models[0].ForeignKeyValues[constants.AddressFK] = strings.ToLower(models[0].ForeignKeyValues[constants.AddressFK])
		Expect(models).To(Equal([]shared.InsertionModel{test_data.FlipKickModel()}))
	})

	It("returns an error if converting log to entity fails", func() {
		_, err := converter.ToModels("error abi", []core.HeaderSyncLog{test_data.FlipKickHeaderSyncLog})

		Expect(err).To(HaveOccurred())
	})
})
