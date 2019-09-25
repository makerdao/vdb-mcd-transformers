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

package new_cdp_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"

	"github.com/vulcanize/mcd_transformers/transformers/events/new_cdp"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("NewCdp Converter", func() {
	var converter = new_cdp.NewCdpConverter{}

	Describe("ToEntity", func() {
		It("converts an Eth Log to a NewCdpEntity", func() {
			entities, err := converter.ToEntities(constants.CdpManagerABI(), []core.HeaderSyncLog{test_data.NewCdpHeaderSyncLog})

			Expect(err).NotTo(HaveOccurred())
			Expect(len(entities)).To(Equal(1))
			Expect(entities[0]).To(Equal(test_data.NewCdpEntity))
		})

		It("returns an error if converting log to entity fails", func() {
			_, err := converter.ToEntities("error abi", []core.HeaderSyncLog{test_data.NewCdpHeaderSyncLog})

			Expect(err).To(HaveOccurred())
		})
	})

	Describe("ToModel", func() {
		It("converts an Entity to a Model", func() {
			models, err := converter.ToModels([]interface{}{test_data.NewCdpEntity})

			Expect(err).NotTo(HaveOccurred())
			Expect(len(models)).To(Equal(1))
			Expect(models[0]).To(Equal(test_data.NewCdpModel))
		})

		It("handles nil values", func() {
			emptyAddressHex := "0x0000000000000000000000000000000000000000"
			emptyString := ""
			emptyEntity := new_cdp.NewCdpEntity{}

			models, err := converter.ToModels([]interface{}{emptyEntity})
			Expect(err).NotTo(HaveOccurred())

			Expect(len(models)).To(Equal(1))
			model := models[0].(new_cdp.NewCdpModel)
			Expect(model.Usr).To(Equal(emptyAddressHex))
			Expect(model.Own).To(Equal(emptyAddressHex))
			Expect(model.Cdp).To(Equal(emptyString))
		})

		It("returns an error if the wrong entity type is passed in", func() {
			_, err := converter.ToModels([]interface{}{test_data.WrongEntity{}})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("entity of type"))
		})
	})
})
