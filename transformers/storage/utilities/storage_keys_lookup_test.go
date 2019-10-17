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

package utilities_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/storage/utilities"
)

var _ = Describe("Storage keys lookup utilities", func() {
	Describe("padding an address", func() {
		It("returns error if input not 20 bytes", func() {
			shortAddress := "0x23"

			_, err := utilities.PadAddress(shortAddress)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utilities.ErrInvalidAddress(shortAddress)))
		})

		It("returns address with padding to 32 bytes", func() {
			res, err := utilities.PadAddress(test_helpers.FakeAddress)

			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(Equal("0x000000000000000000000000" + test_helpers.FakeAddress[2:]))
		})
	})
})
