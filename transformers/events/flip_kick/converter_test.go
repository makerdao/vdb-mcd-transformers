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
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/transformers/events/flip_kick"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"strings"
)

var _ = Describe("FlipKick Converter", func() {
	var converter = flip_kick.FlipKickConverter{}

	Describe("ToEntity", func() {
		It("converts an Eth Log to a FlipKickEntity", func() {
			entities, err := converter.ToEntities(constants.FlipABI(), []types.Log{test_data.EthFlipKickLog})

			Expect(err).NotTo(HaveOccurred())
			Expect(len(entities)).To(Equal(1))
			entity := entities[0]
			Expect(entity).To(Equal(test_data.FlipKickEntity))
		})
	})

	Describe("ToModel", func() {
		It("converts a log to a model", func() {
			models, err := converter.ToModels(constants.FlipABI(), []types.Log{test_data.EthFlipKickLog})

			Expect(err).NotTo(HaveOccurred())

			// TODO Why is the converter capitalising some chars in the address?!
			models[0].ForeignKeyValues[constants.AddressFK] = strings.ToLower(models[0].ForeignKeyValues[constants.AddressFK])
			Expect(models).To(Equal([]shared.InsertionModel{test_data.FlipKickModel}))
		})

		It("returns an error if converting log to entity fails", func() {
			_, err := converter.ToEntities("error abi", []types.Log{test_data.EthFlipKickLog})

			Expect(err).To(HaveOccurred())
		})
	})
})

// Old test compared lowercase. The testing.toml lists the ETH_FLIP_A lowercased, yet here it appears mixed case in "actual". Where does the capitalisation happen?
/*
func expectEqualModels(actual interface{}, expected flip_kick.FlipKickModel) {
	actualFlipKick := actual.(flip_kick.FlipKickModel)
	Expect(actualFlipKick.BidId).To(Equal(expected.BidId))
	Expect(actualFlipKick.Lot).To(Equal(expected.Lot))
	Expect(actualFlipKick.Bid).To(Equal(expected.Bid))
	Expect(actualFlipKick.Tab).To(Equal(expected.Tab))
	Expect(actualFlipKick.Usr).To(Equal(expected.Usr))
	Expect(actualFlipKick.Gal).To(Equal(expected.Gal))
	Expect(strings.ToLower(actualFlipKick.ContractAddress)).To(Equal(strings.ToLower(expected.ContractAddress)))
	Expect(actualFlipKick.TransactionIndex).To(Equal(expected.TransactionIndex))
	Expect(actualFlipKick.LogIndex).To(Equal(expected.LogIndex))
	Expect(actualFlipKick.Raw).To(Equal(expected.Raw))
}
*/