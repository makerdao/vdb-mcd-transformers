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

package flop_kick_test

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"

	"github.com/vulcanize/mcd_transformers/transformers/events/flop_kick"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("FlopKick Converter", func() {
	var converter flop_kick.FlopKickConverter

	Describe("ToEntities", func() {
		It("converts a log to a FlopKick entity", func() {
			entities, err := converter.ToEntities(constants.FlopABI(), []types.Log{test_data.EthFlopKickLog})

			Expect(err).NotTo(HaveOccurred())
			Expect(len(entities)).To(Equal(1))
			Expect(entities[0]).To(Equal(test_data.FlopKickEntity))
		})

		It("returns an error if converting the log to an entity fails", func() {
			entities, err := converter.ToEntities("error abi", []types.Log{test_data.EthFlopKickLog})

			Expect(err).To(HaveOccurred())
			Expect(entities).To(BeNil())
		})
	})

	Describe("ToModels", func() {
		It("converts an Entity to a Model", func() {
			models, err := converter.ToModels([]interface{}{test_data.FlopKickEntity})

			Expect(err).NotTo(HaveOccurred())
			Expect(models[0]).To(Equal(test_data.FlopKickModel))
		})

		It("returns error if wrong entity", func() {
			_, err := converter.ToModels([]interface{}{test_data.WrongEntity{}})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("entity of type test_data.WrongEntity, not flop_kick.Entity"))
		})

		It("handles nil values", func() {
			emptyAddressHex := "0x0000000000000000000000000000000000000000"
			emptyString := ""
			emptyEntity := flop_kick.Entity{}
			emptyLog, err := json.Marshal(types.Log{})
			expectedModel := flop_kick.Model{
				BidId:            emptyString,
				Lot:              emptyString,
				Bid:              emptyString,
				Gal:              emptyAddressHex,
				ContractAddress:  emptyAddressHex,
				TransactionIndex: 0,
				LogIndex:         0,
				Raw:              emptyLog,
			}

			models, err := converter.ToModels([]interface{}{emptyEntity})

			Expect(err).NotTo(HaveOccurred())
			Expect(len(models)).To(Equal(1))
			Expect(models[0]).To(Equal(expectedModel))
		})
	})
})
