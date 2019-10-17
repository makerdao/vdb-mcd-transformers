package flip_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	vdbStorage "github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/flip"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
)

var _ = Describe("Flip storage mappings", func() {
	var (
		storageRepository *test_helpers.MockMakerStorageRepository
		storageKeysLookup flip.StorageKeysLookup
	)

	BeforeEach(func() {
		storageRepository = &test_helpers.MockMakerStorageRepository{}
		storageKeysLookup = flip.StorageKeysLookup{StorageRepository: storageRepository, ContractAddress: constants.GetContractAddress("MCD_FLIP_ETH_A")}
	})

	Describe("looking up static keys", func() {
		It("returns value metadata if key exists", func() {
			Expect(storageKeysLookup.Lookup(flip.VatKey)).To(Equal(flip.VatMetadata))
			Expect(storageKeysLookup.Lookup(flip.IlkKey)).To(Equal(flip.IlkMetadata))
			Expect(storageKeysLookup.Lookup(flip.BegKey)).To(Equal(flip.BegMetadata))
			Expect(storageKeysLookup.Lookup(flip.TtlAndTauStorageKey)).To(Equal(flip.TtlAndTauMetadata))
			Expect(storageKeysLookup.Lookup(flip.KicksKey)).To(Equal(flip.KicksMetadata))
		})

		It("returns value metadata if keccak hashed key exists", func() {
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(flip.VatKey[:]))).To(Equal(flip.VatMetadata))
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(flip.IlkKey[:]))).To(Equal(flip.IlkMetadata))
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(flip.BegKey[:]))).To(Equal(flip.BegMetadata))
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(flip.TtlAndTauStorageKey[:]))).To(Equal(flip.TtlAndTauMetadata))
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(flip.KicksKey[:]))).To(Equal(flip.KicksMetadata))
		})

		It("returns error if key does not exist", func() {
			_, err := storageKeysLookup.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrStorageKeyNotFound{Key: fakes.FakeHash.Hex()}))
		})
	})

	Describe("looking up dynamic keys", func() {
		It("refreshes mappings from repository if key not found", func() {
			_, _ = storageKeysLookup.Lookup(fakes.FakeHash)

			Expect(storageRepository.GetFlipBidIdsCalledWith).To(Equal(storageKeysLookup.ContractAddress))
		})

		It("returns error if bid ID lookup fails", func() {
			storageRepository.GetFlipBidIdsError = fakes.FakeError

			_, err := storageKeysLookup.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		Describe("bid", func() {
			fakeBidId := "1"
			fakeHexBidId, _ := shared.ConvertIntStringToHex(fakeBidId)
			var bidBidKey = common.BytesToHash(crypto.Keccak256(common.FromHex(fakeHexBidId + flip.BidsMappingIndex)))

			BeforeEach(func() {
				storageRepository.FlipBidIds = []string{fakeBidId}
			})

			It("returns value metadata for bid bid", func() {
				expectedMetadata := utils.StorageValueMetadata{
					Name: storage.BidBid,
					Keys: map[utils.Key]string{constants.BidId: fakeBidId},
					Type: utils.Uint256,
				}
				Expect(storageKeysLookup.Lookup(bidBidKey)).To(Equal(expectedMetadata))
			})

			It("returns value metadata for bid lot", func() {
				bidLotKey := vdbStorage.GetIncrementedKey(bidBidKey, 1)
				expectedMetadata := utils.StorageValueMetadata{
					Name: storage.BidLot,
					Keys: map[utils.Key]string{constants.BidId: fakeBidId},
					Type: utils.Uint256,
				}
				Expect(storageKeysLookup.Lookup(bidLotKey)).To(Equal(expectedMetadata))
			})

			It("returns value metadata for bid guy + tic + end packed slot", func() {
				bidGuyKey := vdbStorage.GetIncrementedKey(bidBidKey, 2)
				expectedMetadata := utils.StorageValueMetadata{
					Name:        storage.Packed,
					Keys:        map[utils.Key]string{constants.BidId: fakeBidId},
					Type:        utils.PackedSlot,
					PackedTypes: map[int]utils.ValueType{0: utils.Address, 1: utils.Uint48, 2: utils.Uint48},
					PackedNames: map[int]string{0: storage.BidGuy, 1: storage.BidTic, 2: storage.BidEnd},
				}
				Expect(storageKeysLookup.Lookup(bidGuyKey)).To(Equal(expectedMetadata))
			})

			It("returns value metadata for bid usr", func() {
				bidUsrKey := vdbStorage.GetIncrementedKey(bidBidKey, 3)
				expectedMetadata := utils.StorageValueMetadata{
					Name: storage.BidUsr,
					Keys: map[utils.Key]string{constants.BidId: fakeBidId},
					Type: utils.Address,
				}
				Expect(storageKeysLookup.Lookup(bidUsrKey)).To(Equal(expectedMetadata))
			})

			It("returns value metadata for bid gal", func() {
				bidGalKey := vdbStorage.GetIncrementedKey(bidBidKey, 4)
				expectedMetadata := utils.StorageValueMetadata{
					Name: storage.BidGal,
					Keys: map[utils.Key]string{constants.BidId: fakeBidId},
					Type: utils.Address,
				}
				Expect(storageKeysLookup.Lookup(bidGalKey)).To(Equal(expectedMetadata))
			})

			It("returns value metadata for bid tab", func() {
				bidTabKey := vdbStorage.GetIncrementedKey(bidBidKey, 5)
				expectedMetadata := utils.StorageValueMetadata{
					Name: storage.BidTab,
					Keys: map[utils.Key]string{constants.BidId: fakeBidId},
					Type: utils.Uint256,
				}
				Expect(storageKeysLookup.Lookup(bidTabKey)).To(Equal(expectedMetadata))
			})
		})
	})
})
