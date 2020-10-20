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

package shared_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	shared2 "github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Shared utilities", func() {
	Describe("getting log note data bytes at index", func() {
		Describe("extracting Vat Note data", func() {
			It("returns error if index less than two (arguments 0 and 1 are always in topics)", func() {
				_, err := shared.GetLogNoteArgumentAtIndex(1, []byte{})

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(shared.ErrInvalidIndex(1)))
			})

			It("returns error if index greater than five (no functions with > 6 arguments)", func() {
				_, err := shared.GetLogNoteArgumentAtIndex(6, []byte{})

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(shared.ErrInvalidIndex(6)))
			})

			It("extracts fourth argument of four arguments", func() {
				wadBytes, err := shared.GetLogNoteArgumentAtIndex(3, test_data.VatFluxEventLog.Log.Data)

				Expect(err).NotTo(HaveOccurred())
				wadInt := shared2.ConvertUint256HexToBigInt(hexutil.Encode(wadBytes))
				Expect(wadInt.String()).To(Equal(test_data.VatFluxModel().ColumnValues["wad"].(string)))
			})

			It("extracts fourth of five arguments", func() {
				dinkBytes, err := shared.GetLogNoteArgumentAtIndex(3, test_data.VatForkEventLogWithNegativeDinkDart.Log.Data)

				Expect(err).NotTo(HaveOccurred())
				dinkInt := shared2.ConvertInt256HexToBigInt(hexutil.Encode(dinkBytes))
				Expect(dinkInt.String()).To(Equal(test_data.VatForkModelWithNegativeDinkDart().ColumnValues["dink"].(string)))
			})

			It("extracts fifth of five arguments", func() {
				dartBytes, err := shared.GetLogNoteArgumentAtIndex(4, test_data.VatForkEventLogWithNegativeDinkDart.Log.Data)

				Expect(err).NotTo(HaveOccurred())
				dartInt := shared2.ConvertInt256HexToBigInt(hexutil.Encode(dartBytes))
				Expect(dartInt.String()).To(Equal(test_data.VatForkModelWithNegativeDinkDart().ColumnValues["dart"].(string)))
			})

			It("extracts the fourth of six arguments", func() {
				wBytes, err := shared.GetLogNoteArgumentAtIndex(3, test_data.VatGrabEventLogWithPositiveDink.Log.Data)

				Expect(err).NotTo(HaveOccurred())
				wAddress := common.BytesToAddress(wBytes)
				Expect(wAddress.String()).To(Equal(test_data.VatGrabModelWithPositiveDink().ColumnValues["w"]))
			})

			It("extracts the fifth of six arguments", func() {
				dinkBytes, err := shared.GetLogNoteArgumentAtIndex(4, test_data.VatGrabEventLogWithPositiveDink.Log.Data)

				Expect(err).NotTo(HaveOccurred())
				dinkInt := shared2.ConvertInt256HexToBigInt(hexutil.Encode(dinkBytes))
				Expect(dinkInt.String()).To(Equal(test_data.VatGrabModelWithPositiveDink().ColumnValues["dink"]))
			})

			It("extracts the sixth of six arguments", func() {
				dartBytes, err := shared.GetLogNoteArgumentAtIndex(5, test_data.VatGrabEventLogWithPositiveDink.Log.Data)

				Expect(err).NotTo(HaveOccurred())
				dartInt := shared2.ConvertInt256HexToBigInt(hexutil.Encode(dartBytes))
				Expect(dartInt.String()).To(Equal(test_data.VatGrabModelWithPositiveDink().ColumnValues["dart"]))
			})
		})
	})
})
