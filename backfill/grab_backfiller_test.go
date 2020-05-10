package backfill_test

import (
	"errors"
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/makerdao/vdb-mcd-transformers/backfill"
	"github.com/makerdao/vdb-mcd-transformers/backfill/mocks"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Grab BackFiller", func() {
	Describe("BackFillGrabStorage", func() {
		var (
			backFiller            backfill.GrabBackFiller
			db                    = test_config.NewTestDB(test_config.NewTestNode())
			mockBlockChain        *fakes.MockBlockChain
			mockEventsRepository  *mocks.EventsRepository
			mockStorageRepository *mocks.StorageRepository
		)

		BeforeEach(func() {
			test_config.CleanTestDB(db)
			mockBlockChain = fakes.NewMockBlockChain()
			mockEventsRepository = &mocks.EventsRepository{}
			mockStorageRepository = &mocks.StorageRepository{}
			backFiller = backfill.NewGrabBackFiller(mockBlockChain, mockEventsRepository, mockStorageRepository)
		})

		It("gets vat grab events from starting block number onward", func() {
			startingBlock := rand.Int()

			err := backFiller.BackFillGrabStorage(startingBlock)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockEventsRepository.GetGrabsPassedStartingBlock).To(Equal(startingBlock))
		})

		It("returns error if getting grab events fails", func() {
			mockEventsRepository.GetGrabsError = fakes.FakeError

			err := backFiller.BackFillGrabStorage(0)

			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, fakes.FakeError)).To(BeTrue())
		})

		It("ignores grab if dink and dart are zero", func() {
			fakeGrab := backfill.Grab{
				Dink: "0",
				Dart: "0",
			}
			mockEventsRepository.GetGrabsGrabsToReturn = []backfill.Grab{fakeGrab}

			err := backFiller.BackFillGrabStorage(0)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockBlockChain.BatchGetStorageAtCalls).To(BeEmpty())
		})

		It("passes ilk ID and header ID to detect if ilk art exists", func() {
			fakeGrab := backfill.Grab{
				HeaderID: rand.Int(),
				Dink:     "0",
				Dart:     "1",
			}
			mockEventsRepository.GetGrabsGrabsToReturn = []backfill.Grab{fakeGrab}
			fakeUrn := backfill.Urn{IlkID: rand.Int()}
			mockStorageRepository.GetUrnByIDUrnToReturn = fakeUrn
			mockStorageRepository.VatIlkArtExistsBoolToReturn = true
			mockStorageRepository.VatUrnArtExistsBoolToReturn = true
			mockStorageRepository.VatUrnInkExistsBoolToReturn = true

			err := backFiller.BackFillGrabStorage(0)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockStorageRepository.VatIlkArtExistsPassedIlkID).To(Equal(fakeUrn.IlkID))
			Expect(mockStorageRepository.VatIlkArtExistsPassedHeaderID).To(Equal(fakeGrab.HeaderID))
		})

		It("passes urn ID and header ID to detect if urn art exists", func() {
			fakeGrab := backfill.Grab{
				HeaderID: rand.Int(),
				UrnID:    rand.Int(),
				Dink:     "0",
				Dart:     "1",
			}
			mockEventsRepository.GetGrabsGrabsToReturn = []backfill.Grab{fakeGrab}
			mockStorageRepository.VatIlkArtExistsBoolToReturn = true
			mockStorageRepository.VatUrnArtExistsBoolToReturn = true
			mockStorageRepository.VatUrnInkExistsBoolToReturn = true

			err := backFiller.BackFillGrabStorage(0)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockStorageRepository.VatUrnArtExistsPassedHeaderID).To(Equal(fakeGrab.HeaderID))
			Expect(mockStorageRepository.VatUrnArtExistsPassedUrnID).To(Equal(fakeGrab.UrnID))
		})

		It("passes urn ID and header ID to detect if urn ink exists", func() {
			fakeGrab := backfill.Grab{
				HeaderID: rand.Int(),
				UrnID:    rand.Int(),
				Dink:     "1",
				Dart:     "0",
			}
			mockEventsRepository.GetGrabsGrabsToReturn = []backfill.Grab{fakeGrab}
			mockStorageRepository.VatIlkArtExistsBoolToReturn = true
			mockStorageRepository.VatUrnArtExistsBoolToReturn = true
			mockStorageRepository.VatUrnInkExistsBoolToReturn = true

			err := backFiller.BackFillGrabStorage(0)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockStorageRepository.VatUrnInkExistsPassedHeaderID).To(Equal(fakeGrab.HeaderID))
			Expect(mockStorageRepository.VatUrnInkExistsPassedUrnID).To(Equal(fakeGrab.UrnID))
		})

		It("ignores grab if transformed diffs already exists at header", func() {
			fakeGrab := backfill.Grab{
				Dink: "1",
				Dart: "1",
			}
			mockEventsRepository.GetGrabsGrabsToReturn = []backfill.Grab{fakeGrab}
			mockStorageRepository.VatIlkArtExistsBoolToReturn = true
			mockStorageRepository.VatUrnArtExistsBoolToReturn = true
			mockStorageRepository.VatUrnInkExistsBoolToReturn = true

			err := backFiller.BackFillGrabStorage(0)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockBlockChain.BatchGetStorageAtCalls).To(BeEmpty())
		})

		It("returns error if getting for urn for grab requiring back-fill fails", func() {
			fakeGrab := backfill.Grab{
				UrnID: rand.Int(),
				Dink:  "0",
				Dart:  "1",
			}
			mockEventsRepository.GetGrabsGrabsToReturn = []backfill.Grab{fakeGrab}
			mockStorageRepository.VatIlkArtExistsBoolToReturn = false
			mockStorageRepository.GetUrnByIDError = fakes.FakeError

			err := backFiller.BackFillGrabStorage(0)

			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, fakes.FakeError)).To(BeTrue())
		})

		It("returns error if getting header for grab requiring back-fill fails", func() {
			fakeGrab := backfill.Grab{
				UrnID: rand.Int(),
				Dink:  "0",
				Dart:  "1",
			}
			mockEventsRepository.GetGrabsGrabsToReturn = []backfill.Grab{fakeGrab}
			mockStorageRepository.VatIlkArtExistsBoolToReturn = false
			mockEventsRepository.GetHeaderByIDError = fakes.FakeError

			err := backFiller.BackFillGrabStorage(0)

			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, fakes.FakeError)).To(BeTrue())
		})

		It("looks up storage for grab when some values are non-zero and don't already exist", func() {
			fakeGrab := backfill.Grab{
				Dink: "0",
				Dart: "1",
			}
			mockEventsRepository.GetGrabsGrabsToReturn = []backfill.Grab{fakeGrab}
			mockStorageRepository.VatIlkArtExistsBoolToReturn = false
			mockStorageRepository.VatUrnArtExistsBoolToReturn = true
			mockStorageRepository.VatUrnInkExistsBoolToReturn = true
			fakeUrn := backfill.Urn{Ilk: test_data.RandomString(64)}
			mockStorageRepository.GetUrnByIDUrnToReturn = fakeUrn
			fakeHeader := core.Header{BlockNumber: rand.Int63()}
			mockEventsRepository.GetHeaderByIDHeaderToReturn = fakeHeader

			err := backFiller.BackFillGrabStorage(0)

			Expect(err).NotTo(HaveOccurred())
			expectedIlkArtKey := storage.GetKeyForMapping(storage.IndexTwo, fakeUrn.Ilk)
			Expect(mockBlockChain.BatchGetStorageAtCalls).To(ConsistOf(fakes.BatchGetStorageAtCall{
				Account:     backfill.VatAddress,
				Keys:        []common.Hash{expectedIlkArtKey},
				BlockNumber: big.NewInt(fakeHeader.BlockNumber),
			}))
		})

		It("inserts returned value", func() {
			fakeGrab := backfill.Grab{
				Dink: "1",
				Dart: "0",
			}
			mockEventsRepository.GetGrabsGrabsToReturn = []backfill.Grab{fakeGrab}
			fakeUrn := backfill.Urn{
				Ilk: test_data.RandomString(64),
				Urn: "0x" + test_data.RandomString(40),
			}
			mockStorageRepository.GetUrnByIDUrnToReturn = fakeUrn
			fakeHeader := core.Header{
				BlockNumber: rand.Int63(),
				Hash:        test_data.RandomString(64),
			}
			mockEventsRepository.GetHeaderByIDHeaderToReturn = fakeHeader
			fakeValue := []byte{1, 2, 3, 4, 5}
			mockBlockChain.SetStorageValuesToReturn(fakeHeader.BlockNumber, backfill.VatAddress, fakeValue)

			err := backFiller.BackFillGrabStorage(0)

			Expect(err).NotTo(HaveOccurred())
			paddedUrn, padErr := utilities.PadAddress(fakeUrn.Urn)
			Expect(padErr).NotTo(HaveOccurred())
			expectedUrnInkKey := storage.GetKeyForNestedMapping(storage.IndexThree, fakeUrn.Ilk, paddedUrn)
			expectedDiff := types.RawDiff{
				HashedAddress: crypto.Keccak256Hash(backfill.VatAddress.Bytes()),
				BlockHash:     common.HexToHash(fakeHeader.Hash),
				BlockHeight:   int(fakeHeader.BlockNumber),
				StorageKey:    crypto.Keccak256Hash(expectedUrnInkKey.Bytes()),
				StorageValue:  common.BytesToHash(fakeValue),
			}
			Expect(mockStorageRepository.InsertDiffPassedDiff).To(Equal(expectedDiff))
		})
	})
})
