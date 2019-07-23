package flap_test

import (
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"math/rand"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage/flap"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
)

var _ = Describe("Flap storage mappings", func() {
	var (
		storageRepository *test_helpers.MockMakerStorageRepository
		mapping           flap.StorageKeysLookup
	)

	BeforeEach(func() {
		storageRepository = &test_helpers.MockMakerStorageRepository{}
		mapping = flap.StorageKeysLookup{
			StorageRepository: storageRepository,
			ContractAddress:   constants.FlapperContractAddress(),
		}
	})

	Describe("looks up static keys", func() {
		It("returns storage value mapping if storage key exists", func() {
			Expect(mapping.Lookup(flap.VatStorageKey)).To(Equal(utils.StorageValueMetadata{
				Name: flap.Vat,
				Keys: nil,
				Type: utils.Address,
			}))
			Expect(mapping.Lookup(flap.GemStorageKey)).To(Equal(utils.StorageValueMetadata{
				Name: flap.Gem,
				Keys: nil,
				Type: utils.Address,
			}))
			Expect(mapping.Lookup(flap.BegStorageKey)).To(Equal(utils.StorageValueMetadata{
				Name: flap.Beg,
				Keys: nil,
				Type: utils.Uint256,
			}))
			Expect(mapping.Lookup(flap.TtlAndTauStorageKey)).To(Equal(utils.StorageValueMetadata{
				Name:        flap.Packed,
				Type:        utils.PackedSlot,
				PackedTypes: map[int]utils.ValueType{0: utils.Uint48, 1: utils.Uint48},
				PackedNames: map[int]string{0: flap.Ttl, 1: flap.Tau},
			}))
			Expect(mapping.Lookup(flap.KicksStorageKey)).To(Equal(utils.StorageValueMetadata{
				Name: flap.Kicks,
				Keys: nil,
				Type: utils.Uint256,
			}))
			Expect(mapping.Lookup(flap.LiveStorageKey)).To(Equal(utils.StorageValueMetadata{
				Name: flap.Live,
				Keys: nil,
				Type: utils.Uint256,
			}))
		})

		It("returns an error if the key doesn't exist", func() {
			emptyMetadata, err := mapping.Lookup(fakes.FakeHash)

			Expect(emptyMetadata).To(Equal(utils.StorageValueMetadata{}))
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrStorageKeyNotFound{Key: fakes.FakeHash.Hex()}))
		})
	})

	Describe("looks up dynamic keys", func() {
		It("returns error if getting bid ids fails", func() {
			storageRepository.GetFlapBidIdsError = fakes.FakeError

			_, err := mapping.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("refreshes the bid keys if the given key is not found", func() {
			mapping.Lookup(fakes.FakeHash)

			Expect(storageRepository.GetFlapBidIdsCalled).To(BeTrue())
		})

		Describe("flap bids", func() {
			var (
				bidId         = strconv.FormatInt(rand.Int63(), 10)
				bidIdHex, _   = shared.ConvertIntStringToHex(bidId)
				flapBidBidKey = common.BytesToHash(
					crypto.Keccak256(
						common.FromHex(bidIdHex + flap.BidsIndex)))
			)

			BeforeEach(func() {
				storageRepository.FlapBidIds = []string{bidId}
			})

			It("gets bid metadata", func() {
				metadata, err := mapping.Lookup(flapBidBidKey)
				Expect(err).NotTo(HaveOccurred())

				Expect(metadata).To(Equal(utils.StorageValueMetadata{
					Name: flap.BidBid,
					Keys: map[utils.Key]string{constants.BidId: bidId},
					Type: utils.Uint256,
				}))
			})

			It("gets lot metadata", func() {
				flapBidLotKey := storage.GetIncrementedKey(flapBidBidKey, 1)
				metadata, err := mapping.Lookup(flapBidLotKey)
				Expect(err).NotTo(HaveOccurred())

				Expect(metadata).To(Equal(utils.StorageValueMetadata{
					Name: flap.BidLot,
					Keys: map[utils.Key]string{constants.BidId: bidId},
					Type: utils.Uint256,
				}))
			})

			It("gets guy metadata", func() {
				flapBidGuyKey := storage.GetIncrementedKey(flapBidBidKey, 2)
				metadata, err := mapping.Lookup(flapBidGuyKey)
				Expect(err).NotTo(HaveOccurred())

				Expect(metadata).To(Equal(utils.StorageValueMetadata{
					Name: flap.BidGuy,
					Keys: map[utils.Key]string{constants.BidId: bidId},
					Type: utils.Address,
				}))
			})

			It("gets tic metadata", func() {
				flapBidTicKey := storage.GetIncrementedKey(flapBidBidKey, 3)
				metadata, err := mapping.Lookup(flapBidTicKey)
				Expect(err).NotTo(HaveOccurred())

				Expect(metadata).To(Equal(utils.StorageValueMetadata{
					Name: flap.BidTic,
					Keys: map[utils.Key]string{constants.BidId: bidId},
					Type: utils.Uint48,
				}))
			})

			It("gets end metadata", func() {
				flapBidEndKey := storage.GetIncrementedKey(flapBidBidKey, 4)
				metadata, err := mapping.Lookup(flapBidEndKey)
				Expect(err).NotTo(HaveOccurred())

				Expect(metadata).To(Equal(utils.StorageValueMetadata{
					Name: flap.BidEnd,
					Keys: map[utils.Key]string{constants.BidId: bidId},
					Type: utils.Uint48,
				}))
			})

			It("gets gal metadata", func() {
				flapBidGalKey := storage.GetIncrementedKey(flapBidBidKey, 5)
				metadata, err := mapping.Lookup(flapBidGalKey)
				Expect(err).NotTo(HaveOccurred())

				Expect(metadata).To(Equal(utils.StorageValueMetadata{
					Name: flap.BidGal,
					Keys: map[utils.Key]string{constants.BidId: bidId},
					Type: utils.Address,
				}))
			})
		})
	})
})
