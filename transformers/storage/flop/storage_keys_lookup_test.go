package flop_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage/flop"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("Flop storage mappings", func() {
	var (
		storageRepository *test_helpers.MockMakerStorageRepository
		mappings          flop.StorageKeysLookup
	)

	BeforeEach(func() {
		storageRepository = &test_helpers.MockMakerStorageRepository{}
		mappings = flop.StorageKeysLookup{StorageRepository: storageRepository, ContractAddress: "0x668001c75a9c02d6b10c7a17dbd8aa4afff95037"}
	})

	Describe("looking up static keys", func() {
		It("returns value metadata if key exists", func() {
			Expect(mappings.Lookup(flop.VatKey)).To(Equal(flop.VatMetadata))
			Expect(mappings.Lookup(flop.GemKey)).To(Equal(flop.GemMetadata))
			Expect(mappings.Lookup(flop.BegKey)).To(Equal(flop.BegMetadata))
			Expect(mappings.Lookup(flop.TtlAndTauKey)).To(Equal(flop.TtlAndTauMetadata))
			Expect(mappings.Lookup(flop.KicksKey)).To(Equal(flop.KicksMetadata))
			Expect(mappings.Lookup(flop.LiveKey)).To(Equal(flop.LiveMetadata))
		})

		It("returns error if key does not exist", func() {
			_, err := mappings.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrStorageKeyNotFound{Key: fakes.FakeHash.Hex()}))
		})
	})

	Describe("looking up dynamic keys", func() {
		It("refreshes mappings from repository if key not found", func() {
			_, _ = mappings.Lookup(fakes.FakeHash)

			Expect(storageRepository.GetFlopBidIdsCalledWith).To(Equal(mappings.ContractAddress))
		})

		It("returns error if bid ID lookup fails", func() {
			storageRepository.GetFlopBidIdsError = fakes.FakeError

			_, err := mappings.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})

	Describe("bid", func() {
		var fakeBidId string
		var bidBidKey common.Hash

		BeforeEach(func() {
			fakeBidId = "42"
			fakeHexBidId, conversionErr := shared.ConvertIntStringToHex(fakeBidId)

			Expect(conversionErr).NotTo(HaveOccurred())

			bidBidKey = common.BytesToHash(crypto.Keccak256(common.FromHex(fakeHexBidId + flop.BidsIndex)))
			storageRepository.FlopBidIds = []string{fakeBidId}
		})

		It("returns value metadata for bid bid", func() {
			expectedMetadata := utils.StorageValueMetadata{
				Name: flop.BidBid,
				Keys: map[utils.Key]string{constants.BidId: fakeBidId},
				Type: utils.Uint256,
			}
			Expect(mappings.Lookup(bidBidKey)).To(Equal(expectedMetadata))
		})

		It("returns value metadata for bid lot", func() {
			bidLotKey := storage.GetIncrementedKey(bidBidKey, 1)
			expectedMetadata := utils.StorageValueMetadata{
				Name: flop.BidLot,
				Keys: map[utils.Key]string{constants.BidId: fakeBidId},
				Type: utils.Uint256,
			}
			Expect(mappings.Lookup(bidLotKey)).To(Equal(expectedMetadata))
		})

		It("returns value metadata for bid guy", func() {
			bidGuyKey := storage.GetIncrementedKey(bidBidKey, 2)
			expectedMetadata := utils.StorageValueMetadata{
				Name: flop.BidGuy,
				Keys: map[utils.Key]string{constants.BidId: fakeBidId},
				Type: utils.Address,
			}
			Expect(mappings.Lookup(bidGuyKey)).To(Equal(expectedMetadata))
		})

		It("returns value metadata for bid tic", func() {
			bidTicKey := storage.GetIncrementedKey(bidBidKey, 3)
			expectedMetadata := utils.StorageValueMetadata{
				Name: flop.BidTic,
				Keys: map[utils.Key]string{constants.BidId: fakeBidId},
				Type: utils.Uint48,
			}
			Expect(mappings.Lookup(bidTicKey)).To(Equal(expectedMetadata))
		})

		It("returns value metadata for bid end", func() {
			bidEndKey := storage.GetIncrementedKey(bidBidKey, 4)
			expectedMetadata := utils.StorageValueMetadata{
				Name: flop.BidEnd,
				Keys: map[utils.Key]string{constants.BidId: fakeBidId},
				Type: utils.Uint48,
			}
			Expect(mappings.Lookup(bidEndKey)).To(Equal(expectedMetadata))
		})
	})
})
