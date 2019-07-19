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

package constants

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("event signature generator", func() {
	Describe("findSignatureInAbi", func() {
		It("returns the signature if it exists in the ABI", func() {
			signature, _ := findSignatureInAbi(catABI(), "file", []string{"bytes32", "bytes32", "uint256"})

			Expect(signature).To(Equal("file(bytes32,bytes32,uint256)"))
		})

		It("returns error if signature not found in ABI", func() {
			expectedError := errors.New("method file(bytes32,bytes32) does not exist in ABI")

			_, err := findSignatureInAbi(catABI(), "file", []string{"bytes32", "bytes32"})

			Expect(err).To(MatchError(expectedError))
		})
	})

	Describe("getOverloadedFunctionSignature", func() {
		It("panics if it encounters an error", func() {
			Expect(func() { getOverloadedFunctionSignature(catABI(), "file", []string{"bytes32", "bytes32"}) }).To(Panic())
		})
	})
})
