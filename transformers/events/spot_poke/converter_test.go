//  VulcanizeDB
//  Copyright © 2019 Vulcanize
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU Affero General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  This program is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU Affero General Public License for more details.
//
//  You should have received a copy of the GNU Affero General Public License
//  along with this program.  If not, see <http://www.gnu.org/licenses/>.

package spot_poke_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/transformers/events/spot_poke"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

var _ = Describe("SpotPoke Converter", func() {
	var converter = spot_poke.SpotPokeConverter{}

	Describe("ToEntities", func() {
		It("converts eth logs into spot poke entities", func() {
			entities, err := converter.ToEntities(constants.SpotABI(), []core.HeaderSyncLog{test_data.SpotPokeHeaderSyncLog})
			Expect(err).NotTo(HaveOccurred())

			Expect(len(entities)).To(Equal(1))
			Expect(entities[0]).To(Equal(test_data.SpotPokeEntity))
		})

		It("returns an error converting a log to an entity fails", func() {
			_, err := converter.ToEntities("error abi", []core.HeaderSyncLog{test_data.SpotPokeHeaderSyncLog})

			Expect(err).To(HaveOccurred())
		})
	})

	Describe("ToModels", func() {
		It("converts spot poke entities to models", func() {
			models, err := converter.ToModels([]interface{}{test_data.SpotPokeEntity})
			Expect(err).NotTo(HaveOccurred())

			Expect(len(models)).To(Equal(1))
			Expect(models[0]).To(Equal(test_data.SpotPokeModel))
		})

		It("returns an error if the entity isn't the correct type", func() {
			_, err := converter.ToModels([]interface{}{test_data.WrongEntity{}})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("entity of type test_data.WrongEntity, not spot_poke.SpotPokeEntity"))
		})

		It("handles nil values", func() {
			expectedModel := spot_poke.SpotPokeModel{
				Ilk:   "0x0000000000000000000000000000000000000000000000000000000000000000",
				Value: "0.000000",
				Spot:  "",
			}
			models, err := converter.ToModels([]interface{}{spot_poke.SpotPokeEntity{}})

			Expect(err).NotTo(HaveOccurred())
			Expect(len(models)).To(Equal(1))
			model := models[0]
			Expect(model).To(Equal(expectedModel))
		})
	})
})
