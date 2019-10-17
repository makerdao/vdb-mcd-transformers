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
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vow"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("Vow storage mappings", func() {
	var (
		storageRepository *test_helpers.MockMakerStorageRepository
		storageKeysLookup vow.StorageKeysLookup
	)

	BeforeEach(func() {
		storageRepository = &test_helpers.MockMakerStorageRepository{}
		storageKeysLookup = vow.StorageKeysLookup{StorageRepository: storageRepository}
	})

	Describe("looking up static keys", func() {
		It("returns value metadata if key exists", func() {
			Expect(storageKeysLookup.Lookup(vow.VatKey)).To(Equal(vow.VatMetadata))
			Expect(storageKeysLookup.Lookup(vow.FlapperKey)).To(Equal(vow.FlapperMetadata))
			Expect(storageKeysLookup.Lookup(vow.FlopperKey)).To(Equal(vow.FlopperMetadata))
			Expect(storageKeysLookup.Lookup(vow.SinIntegerKey)).To(Equal(vow.SinIntegerMetadata))
			Expect(storageKeysLookup.Lookup(vow.AshKey)).To(Equal(vow.AshMetadata))
			Expect(storageKeysLookup.Lookup(vow.WaitKey)).To(Equal(vow.WaitMetadata))
			Expect(storageKeysLookup.Lookup(vow.DumpKey)).To(Equal(vow.DumpMetadata))
			Expect(storageKeysLookup.Lookup(vow.SumpKey)).To(Equal(vow.SumpMetadata))
			Expect(storageKeysLookup.Lookup(vow.BumpKey)).To(Equal(vow.BumpMetadata))
			Expect(storageKeysLookup.Lookup(vow.HumpKey)).To(Equal(vow.HumpMetadata))
		})

		It("returns value metadata if keccak of key exists", func() {
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(vow.VatKey[:]))).To(Equal(vow.VatMetadata))
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(vow.FlapperKey[:]))).To(Equal(vow.FlapperMetadata))
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(vow.FlopperKey[:]))).To(Equal(vow.FlopperMetadata))
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(vow.SinIntegerKey[:]))).To(Equal(vow.SinIntegerMetadata))
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(vow.AshKey[:]))).To(Equal(vow.AshMetadata))
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(vow.WaitKey[:]))).To(Equal(vow.WaitMetadata))
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(vow.SumpKey[:]))).To(Equal(vow.SumpMetadata))
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(vow.BumpKey[:]))).To(Equal(vow.BumpMetadata))
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(vow.HumpKey[:]))).To(Equal(vow.HumpMetadata))
		})

		It("returns error if key does not exist", func() {
			_, err := storageKeysLookup.Lookup(common.HexToHash(fakes.FakeHash.Hex()))

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrStorageKeyNotFound{Key: fakes.FakeHash.Hex()}))
		})
	})

	Describe("looking up dynamic keys", func() {
		It("refreshes mappings from repository if key not found", func() {
			storageKeysLookup.Lookup(fakes.FakeHash)

			Expect(storageRepository.GetVowSinKeysCalled).To(BeTrue())
		})

		It("returns error if sin keys lookup fails", func() {
			storageRepository.GetVowSinKeysError = fakes.FakeError

			_, err := storageKeysLookup.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("returns value metadata for sin with vow flog event", func() {
			fakeTimestamp := "1538558052"
			storageRepository.SinKeys = []string{fakeTimestamp}
			sinKey := common.HexToHash("0x409bb97b2bc2657d61f96ef15378c58e2a7d5a67559d3718cbad711b817d9000")
			// key found at https://github.com/8thlight/maker-vulcanizedb/pull/132/files#diff-fe4d48373094a6c01df6ca0e35c677c3R1360
			expectedKeys := map[utils.Key]string{constants.Timestamp: fakeTimestamp}
			expectedMetadata := utils.GetStorageValueMetadata(vow.SinMapping, expectedKeys, utils.Uint256)

			Expect(storageKeysLookup.Lookup(sinKey)).To(Equal(expectedMetadata))
		})

		It("returns value metadata for sin with vow fess event", func() {
			fakeTimestamp := "1540893520"
			storageRepository.SinKeys = []string{fakeTimestamp}
			sinKey := common.HexToHash("0x37f4e61f380b4127c877057bc12214bd6b243aa33839584689548356b019d8b8")
			// key found at https://github.com/8thlight/maker-vulcanizedb/pull/132/files#diff-fe4d48373094a6c01df6ca0e35c677c3R2058
			expectedKeys := map[utils.Key]string{constants.Timestamp: fakeTimestamp}
			expectedMetadata := utils.GetStorageValueMetadata(vow.SinMapping, expectedKeys, utils.Uint256)

			Expect(storageKeysLookup.Lookup(sinKey)).To(Equal(expectedMetadata))
		})
	})
})
