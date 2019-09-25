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

package flap_kick_test

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/events/flap_kick"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("Flap kick converter", func() {
	var converter = flap_kick.FlapKickConverter{}

	Describe("ToEntity", func() {
		It("converts an Eth Log to a FlapKickEntity", func() {
			entities, err := converter.ToEntities(constants.FlapABI(), []core.HeaderSyncLog{test_data.FlapKickHeaderSyncLog})

			Expect(err).NotTo(HaveOccurred())
			Expect(len(entities)).To(Equal(1))
			Expect(entities[0]).To(Equal(test_data.FlapKickEntity))
		})
	})

	Describe("ToModel", func() {
		It("returns an error if converting log to entity fails", func() {
			_, err := converter.ToEntities("error abi", []core.HeaderSyncLog{test_data.FlapKickHeaderSyncLog})

			Expect(err).To(HaveOccurred())
		})

		It("converts a log to a Model", func() {
			models, err := converter.ToModels(constants.FlapABI(), []core.HeaderSyncLog{test_data.FlapKickHeaderSyncLog})

			Expect(err).NotTo(HaveOccurred())
			Expect(len(models)).To(Equal(1))
			Expect(models[0]).To(Equal(test_data.FlapKickModel))
		})
	})
})
