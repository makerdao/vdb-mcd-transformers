package vat_lite_test

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/vat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/vat_lite"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vat lite storage keys loader", func() {
	var (
		storageRepository *test_helpers.MockMakerStorageRepository
		storageKeysLoader storage.KeysLoader
	)

	BeforeEach(func() {
		storageRepository = &test_helpers.MockMakerStorageRepository{}
		storageKeysLoader = vat_lite.NewKeysLoader(storageRepository)
	})

	Describe("looking up static keys", func() {
		It("returns value metadata if key exists", func() {
			mappings, err := storageKeysLoader.LoadMappings()

			Expect(err).NotTo(HaveOccurred())
			Expect(mappings[vat.DebtKey]).To(Equal(vat.DebtMetadata))
			Expect(mappings[vat.ViceKey]).To(Equal(vat.ViceMetadata))
			Expect(mappings[vat.LineKey]).To(Equal(vat.LineMetadata))
			Expect(mappings[vat.LiveKey]).To(Equal(vat.LiveMetadata))
		})

	})

	Describe("looking up dynamic keys", func() {
		It("returns error if wards keys lookup fails", func() {
			storageRepository.GetVatWardsKeysError = fakes.FakeError

			_, err := storageKeysLoader.LoadMappings()

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("returns error if gem keys lookup fails", func() {
			storageRepository.GetGemKeysError = fakes.FakeError

			_, err := storageKeysLoader.LoadMappings()

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("returns error if ilks lookup fails", func() {
			storageRepository.GetIlksError = fakes.FakeError

			_, err := storageKeysLoader.LoadMappings()

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("returns error if sin keys lookup fails", func() {
			storageRepository.GetVatSinKeysError = fakes.FakeError

			_, err := storageKeysLoader.LoadMappings()

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		Describe("wards", func() {
			It("returns value metadata for wards", func() {
				wardsUser := fakes.FakeAddress.Hex()
				storageRepository.VatWardsKeys = []string{wardsUser}
				paddedWardsUser := "0x000000000000000000000000" + wardsUser[2:]
				wardsKey := common.BytesToHash(crypto.Keccak256(common.FromHex(paddedWardsUser + wards.WardsMappingIndex)))
				expectedMetadata := types.ValueMetadata{
					Name: wards.Wards,
					Keys: map[types.Key]string{constants.User: fakes.FakeAddress.Hex()},
					Type: types.Uint256,
				}

				mappings, err := storageKeysLoader.LoadMappings()
				Expect(err).NotTo(HaveOccurred())

				Expect(storageRepository.GetVatWardsKeysCalled).To(Equal(true))
				Expect(mappings[wardsKey]).To(Equal(expectedMetadata))
			})
		})

		Describe("ilk", func() {
			var (
				ilkArtKey   common.Hash
				ilkArtAsInt *big.Int
				mappings    map[common.Hash]types.ValueMetadata
			)

			BeforeEach(func() {
				storageRepository.Ilks = []string{test_helpers.FakeIlk}
				ilkArtKey = common.BytesToHash(crypto.Keccak256(common.FromHex(test_helpers.FakeIlk + vat.IlksMappingIndex)))
				ilkArtAsInt = big.NewInt(0).SetBytes(ilkArtKey.Bytes())
				var err error
				mappings, err = storageKeysLoader.LoadMappings()
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns value metadata for ilk Art", func() {
				expectedMetadata := types.ValueMetadata{
					Name: vat.IlkArt,
					Keys: map[types.Key]string{constants.Ilk: test_helpers.FakeIlk},
					Type: types.Uint256,
				}

				Expect(mappings[ilkArtKey]).To(Equal(expectedMetadata))
			})

			It("returns value metadata for ilk rate", func() {
				incrementedIlkArt := big.NewInt(0).Add(ilkArtAsInt, big.NewInt(1))
				ilkRateKey := common.BytesToHash(incrementedIlkArt.Bytes())
				expectedMetadata := types.ValueMetadata{
					Name: vat.IlkRate,
					Keys: map[types.Key]string{constants.Ilk: test_helpers.FakeIlk},
					Type: types.Uint256,
				}

				Expect(mappings[ilkRateKey]).To(Equal(expectedMetadata))
			})

			It("returns value metadata for ilk spot", func() {
				incrementedIlkArt := big.NewInt(0).Add(ilkArtAsInt, big.NewInt(2))
				ilkSpotKey := common.BytesToHash(incrementedIlkArt.Bytes())
				expectedMetadata := types.ValueMetadata{
					Name: vat.IlkSpot,
					Keys: map[types.Key]string{constants.Ilk: test_helpers.FakeIlk},
					Type: types.Uint256,
				}

				Expect(mappings[ilkSpotKey]).To(Equal(expectedMetadata))
			})

			It("returns value metadata for ilk line", func() {
				incrementedIlkArt := big.NewInt(0).Add(ilkArtAsInt, big.NewInt(3))
				ilkLineKey := common.BytesToHash(incrementedIlkArt.Bytes())
				expectedMetadata := types.ValueMetadata{
					Name: vat.IlkLine,
					Keys: map[types.Key]string{constants.Ilk: test_helpers.FakeIlk},
					Type: types.Uint256,
				}

				Expect(mappings[ilkLineKey]).To(Equal(expectedMetadata))
			})

			It("returns value metadata for ilk dust", func() {
				incrementedIlkArt := big.NewInt(0).Add(ilkArtAsInt, big.NewInt(4))
				ilkDustKey := common.BytesToHash(incrementedIlkArt.Bytes())
				expectedMetadata := types.ValueMetadata{
					Name: vat.IlkDust,
					Keys: map[types.Key]string{constants.Ilk: test_helpers.FakeIlk},
					Type: types.Uint256,
				}

				Expect(mappings[ilkDustKey]).To(Equal(expectedMetadata))
			})
		})

		Describe("gem", func() {
			It("returns value metadata for gem", func() {
				storageRepository.GemKeys = []mcdStorage.Urn{{Ilk: test_helpers.FakeIlk, Identifier: test_helpers.FakeAddress}}
				encodedPrimaryMapIndex := crypto.Keccak256(common.FromHex(test_helpers.FakeIlk + vat.GemsMappingIndex))
				paddedGemAddress := common.FromHex("0x000000000000000000000000" + test_helpers.FakeAddress[2:])
				encodedSecondaryMapIndex := crypto.Keccak256(paddedGemAddress, encodedPrimaryMapIndex)
				gemKey := common.BytesToHash(encodedSecondaryMapIndex)
				expectedMetadata := types.ValueMetadata{
					Name: vat.Gem,
					Keys: map[types.Key]string{constants.Ilk: test_helpers.FakeIlk, constants.Guy: test_helpers.FakeAddress},
					Type: types.Uint256,
				}
				mappings, err := storageKeysLoader.LoadMappings()

				Expect(err).NotTo(HaveOccurred())
				Expect(mappings[gemKey]).To(Equal(expectedMetadata))
			})
		})

		Describe("sin", func() {
			It("returns value metadata for sin", func() {
				storageRepository.SinKeys = []string{test_helpers.FakeAddress}
				paddedSinAddress := "0x000000000000000000000000" + test_helpers.FakeAddress[2:]
				sinKey := common.BytesToHash(crypto.Keccak256(common.FromHex(paddedSinAddress + vat.SinMappingIndex)))
				expectedMetadata := types.ValueMetadata{
					Name: vat.Sin,
					Keys: map[types.Key]string{constants.Guy: test_helpers.FakeAddress},
					Type: types.Uint256,
				}
				mappings, err := storageKeysLoader.LoadMappings()

				Expect(err).NotTo(HaveOccurred())
				Expect(mappings[sinKey]).To(Equal(expectedMetadata))
			})
		})
	})
})
