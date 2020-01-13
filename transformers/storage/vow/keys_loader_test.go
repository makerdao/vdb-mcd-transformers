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
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/vow"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vow storage keys loader", func() {
	var (
		storageRepository *test_helpers.MockMakerStorageRepository
		storageKeysLoader storage.KeysLoader
	)

	BeforeEach(func() {
		storageRepository = &test_helpers.MockMakerStorageRepository{}
		storageKeysLoader = vow.NewKeysLoader(storageRepository, test_data.VowAddress())
	})

	It("loads value metadata for static keys", func() {
		mappings, err := storageKeysLoader.LoadMappings()

		Expect(err).NotTo(HaveOccurred())
		Expect(mappings[vow.VatKey]).To(Equal(vow.VatMetadata))
		Expect(mappings[vow.FlapperKey]).To(Equal(vow.FlapperMetadata))
		Expect(mappings[vow.FlopperKey]).To(Equal(vow.FlopperMetadata))
		Expect(mappings[vow.SinIntegerKey]).To(Equal(vow.SinIntegerMetadata))
		Expect(mappings[vow.AshKey]).To(Equal(vow.AshMetadata))
		Expect(mappings[vow.WaitKey]).To(Equal(vow.WaitMetadata))
		Expect(mappings[vow.DumpKey]).To(Equal(vow.DumpMetadata))
		Expect(mappings[vow.SumpKey]).To(Equal(vow.SumpMetadata))
		Expect(mappings[vow.BumpKey]).To(Equal(vow.BumpMetadata))
		Expect(mappings[vow.HumpKey]).To(Equal(vow.HumpMetadata))
	})

	It("returns value metadata for sin with vow flog event", func() {
		fakeTimestamp := "1538558052"
		storageRepository.SinKeys = []string{fakeTimestamp}
		sinKey := common.HexToHash("0x409bb97b2bc2657d61f96ef15378c58e2a7d5a67559d3718cbad711b817d9000")
		// key found at https://github.com/8thlight/maker-vulcanizedb/pull/132/files#diff-fe4d48373094a6c01df6ca0e35c677c3R1360
		expectedKeys := map[vdbStorage.Key]string{constants.Timestamp: fakeTimestamp}
		expectedMetadata := vdbStorage.GetValueMetadata(vow.SinMapping, expectedKeys, vdbStorage.Uint256)

		mappings, err := storageKeysLoader.LoadMappings()

		Expect(err).NotTo(HaveOccurred())
		Expect(mappings[sinKey]).To(Equal(expectedMetadata))
	})

	It("returns value metadata for sin with vow fess event", func() {
		fakeTimestamp := "1540893520"
		storageRepository.SinKeys = []string{fakeTimestamp}
		sinKey := common.HexToHash("0x37f4e61f380b4127c877057bc12214bd6b243aa33839584689548356b019d8b8")
		// key found at https://github.com/8thlight/maker-vulcanizedb/pull/132/files#diff-fe4d48373094a6c01df6ca0e35c677c3R2058
		expectedKeys := map[vdbStorage.Key]string{constants.Timestamp: fakeTimestamp}
		expectedMetadata := vdbStorage.GetValueMetadata(vow.SinMapping, expectedKeys, vdbStorage.Uint256)

		mappings, err := storageKeysLoader.LoadMappings()

		Expect(err).NotTo(HaveOccurred())
		Expect(mappings[sinKey]).To(Equal(expectedMetadata))
	})

	Describe("wards", func() {
		It("returns value metadata for wards", func() {
			wardsUser := fakes.FakeAddress.Hex()
			storageRepository.WardsKeys = []string{wardsUser}
			paddedWardsUser := "0x000000000000000000000000" + wardsUser[2:]
			wardsKey := common.BytesToHash(crypto.Keccak256(common.FromHex(paddedWardsUser + wards.WardsMappingIndex)))
			expectedMetadata := vdbStorage.ValueMetadata{
				Name: wards.Wards,
				Keys: map[vdbStorage.Key]string{constants.User: fakes.FakeAddress.Hex()},
				Type: vdbStorage.Uint256,
			}

			mappings, err := storageKeysLoader.LoadMappings()
			Expect(err).NotTo(HaveOccurred())

			Expect(storageRepository.GetWardsKeysCalledWith).To(Equal(test_data.VowAddress()))
			Expect(mappings[wardsKey]).To(Equal(expectedMetadata))
		})

		Describe("when getting users fails", func() {
			It("returns error", func() {
				storageRepository.GetWardsKeysError = fakes.FakeError

				_, err := storageKeysLoader.LoadMappings()

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(fakes.FakeError))
			})
		})
	})
})
