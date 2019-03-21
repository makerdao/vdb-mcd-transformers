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

package vow_test

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vow"
)

var _ = Describe("Vow storage mappings", func() {
	Describe("looking up static keys", func() {
		It("returns value metadata if key exists", func() {
			storageRepository := &test_helpers.MockMakerStorageRepository{}

			mappings := vow.VowMappings{StorageRepository: storageRepository}

			Expect(mappings.Lookup(vow.VatKey)).To(Equal(vow.VatMetadata))
			Expect(mappings.Lookup(vow.CowKey)).To(Equal(vow.CowMetadata))
			Expect(mappings.Lookup(vow.RowKey)).To(Equal(vow.RowMetadata))
			Expect(mappings.Lookup(vow.SinKey)).To(Equal(vow.SinMetadata))
			Expect(mappings.Lookup(vow.AshKey)).To(Equal(vow.AshMetadata))
			Expect(mappings.Lookup(vow.WaitKey)).To(Equal(vow.WaitMetadata))
			Expect(mappings.Lookup(vow.SumpKey)).To(Equal(vow.SumpMetadata))
			Expect(mappings.Lookup(vow.BumpKey)).To(Equal(vow.BumpMetadata))
			Expect(mappings.Lookup(vow.HumpKey)).To(Equal(vow.HumpMetadata))
		})

		It("returns error if key does not exist", func() {
			storageRepository := &test_helpers.MockMakerStorageRepository{}

			mappings := vow.VowMappings{StorageRepository: storageRepository}
			_, err := mappings.Lookup(common.HexToHash(fakes.FakeHash.Hex()))

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrStorageKeyNotFound{Key: fakes.FakeHash.Hex()}))
		})
	})
})
