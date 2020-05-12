package frob_test

import (
	"errors"
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/makerdao/vdb-mcd-transformers/backfill"
	"github.com/makerdao/vdb-mcd-transformers/backfill/frob"
	"github.com/makerdao/vdb-mcd-transformers/backfill/mocks"
	"github.com/makerdao/vdb-mcd-transformers/backfill/repository"
	"github.com/makerdao/vdb-mcd-transformers/backfill/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Frob BackFiller", func() {
	var (
		mockBlockChain        *fakes.MockBlockChain
		mockEventsRepository  *mocks.EventsRepository
		mockStorageRepository *mocks.StorageRepository
		backFiller            backfill.BackFiller
	)

	BeforeEach(func() {
		mockBlockChain = fakes.NewMockBlockChain()
		mockEventsRepository = &mocks.EventsRepository{}
		mockStorageRepository = &mocks.StorageRepository{}
		backFiller = frob.NewFrobBackFiller(
			mockBlockChain,
			mockEventsRepository,
			mockStorageRepository,
		)
	})

	Describe("BackFill", func() {
		It("gets vat frob events from starting block number onward", func() {
			startingBlock := rand.Int()

			err := backFiller.BackFill(startingBlock)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockEventsRepository.GetFrobsPassedStartingBlock).To(Equal(startingBlock))
		})

		It("returns error if getting frob events fails", func() {
			mockEventsRepository.GetFrobsError = fakes.FakeError

			err := backFiller.BackFill(0)

			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, fakes.FakeError)).To(BeTrue())
		})

		It("ignores frob if dink and dart are zero", func() {
			fakeFrob := repository.Frob{
				Dink: "0",
				Dart: "0",
			}
			mockEventsRepository.GetFrobsFrobsToReturn = []repository.Frob{fakeFrob}

			err := backFiller.BackFill(0)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockBlockChain.BatchGetStorageAtCalls).To(BeEmpty())
		})
		It("passes ilk ID and header ID to detect if ilk art exists", func() {
			fakeFrob := repository.Frob{
				HeaderID: rand.Int(),
				Dink:     "0",
				Dart:     "1",
			}
			mockEventsRepository.GetFrobsFrobsToReturn = []repository.Frob{fakeFrob}
			fakeUrn := repository.Urn{IlkID: rand.Int()}
			mockStorageRepository.GetUrnByIDUrnToReturn = fakeUrn
			mockStorageRepository.VatIlkArtExistsBoolToReturn = true
			mockStorageRepository.VatUrnArtExistsBoolToReturn = true
			mockStorageRepository.VatUrnInkExistsBoolToReturn = true

			err := backFiller.BackFill(0)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockStorageRepository.VatIlkArtExistsPassedIlkID).To(Equal(fakeUrn.IlkID))
			Expect(mockStorageRepository.VatIlkArtExistsPassedHeaderID).To(Equal(fakeFrob.HeaderID))
		})

		It("passes urn ID and header ID to detect if urn art exists", func() {
			fakeFrob := repository.Frob{
				HeaderID: rand.Int(),
				UrnID:    rand.Int(),
				Dink:     "0",
				Dart:     "1",
			}
			mockEventsRepository.GetFrobsFrobsToReturn = []repository.Frob{fakeFrob}
			mockStorageRepository.VatIlkArtExistsBoolToReturn = true
			mockStorageRepository.VatUrnArtExistsBoolToReturn = true
			mockStorageRepository.VatUrnInkExistsBoolToReturn = true

			err := backFiller.BackFill(0)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockStorageRepository.VatUrnArtExistsPassedHeaderID).To(Equal(fakeFrob.HeaderID))
			Expect(mockStorageRepository.VatUrnArtExistsPassedUrnID).To(Equal(fakeFrob.UrnID))
		})

		It("passes urn ID and header ID to detect if urn ink exists", func() {
			fakeFrob := repository.Frob{
				HeaderID: rand.Int(),
				UrnID:    rand.Int(),
				Dink:     "1",
				Dart:     "0",
			}
			mockEventsRepository.GetFrobsFrobsToReturn = []repository.Frob{fakeFrob}
			mockStorageRepository.VatIlkArtExistsBoolToReturn = true
			mockStorageRepository.VatUrnArtExistsBoolToReturn = true
			mockStorageRepository.VatUrnInkExistsBoolToReturn = true

			err := backFiller.BackFill(0)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockStorageRepository.VatUrnInkExistsPassedHeaderID).To(Equal(fakeFrob.HeaderID))
			Expect(mockStorageRepository.VatUrnInkExistsPassedUrnID).To(Equal(fakeFrob.UrnID))
		})

		It("ignores frob if transformed diffs already exists at header", func() {
			fakeFrob := repository.Frob{
				Dink: "1",
				Dart: "1",
			}
			mockEventsRepository.GetFrobsFrobsToReturn = []repository.Frob{fakeFrob}
			mockStorageRepository.VatIlkArtExistsBoolToReturn = true
			mockStorageRepository.VatUrnArtExistsBoolToReturn = true
			mockStorageRepository.VatUrnInkExistsBoolToReturn = true

			err := backFiller.BackFill(0)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockBlockChain.BatchGetStorageAtCalls).To(BeEmpty())
		})

		It("returns error if getting for urn for frob requiring back-fill fails", func() {
			fakeFrob := repository.Frob{
				UrnID: rand.Int(),
				Dink:  "0",
				Dart:  "1",
			}
			mockEventsRepository.GetFrobsFrobsToReturn = []repository.Frob{fakeFrob}
			mockStorageRepository.VatIlkArtExistsBoolToReturn = false
			mockStorageRepository.GetUrnByIDError = fakes.FakeError

			err := backFiller.BackFill(0)

			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, fakes.FakeError)).To(BeTrue())
		})

		It("returns error if getting header for frob requiring back-fill fails", func() {
			fakeFrob := repository.Frob{
				UrnID: rand.Int(),
				Dink:  "0",
				Dart:  "1",
			}
			mockEventsRepository.GetFrobsFrobsToReturn = []repository.Frob{fakeFrob}
			mockStorageRepository.VatIlkArtExistsBoolToReturn = false
			mockEventsRepository.GetHeaderByIDError = fakes.FakeError

			err := backFiller.BackFill(0)

			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, fakes.FakeError)).To(BeTrue())
		})

		It("looks up storage for frob when some values are non-zero and don't already exist", func() {
			fakeFrob := repository.Frob{
				Dink: "0",
				Dart: "1",
			}
			mockEventsRepository.GetFrobsFrobsToReturn = []repository.Frob{fakeFrob}
			mockStorageRepository.VatIlkArtExistsBoolToReturn = false
			mockStorageRepository.VatUrnArtExistsBoolToReturn = true
			mockStorageRepository.VatUrnInkExistsBoolToReturn = true
			fakeUrn := repository.Urn{Ilk: test_data.RandomString(64)}
			mockStorageRepository.GetUrnByIDUrnToReturn = fakeUrn
			fakeHeader := core.Header{BlockNumber: rand.Int63()}
			mockEventsRepository.GetHeaderByIDHeaderToReturn = fakeHeader

			err := backFiller.BackFill(0)

			Expect(err).NotTo(HaveOccurred())
			expectedIlkArtKey := storage.GetKeyForMapping(storage.IndexTwo, fakeUrn.Ilk)
			Expect(mockBlockChain.BatchGetStorageAtCalls).To(ConsistOf(fakes.BatchGetStorageAtCall{
				Account:     shared.VatAddress,
				Keys:        []common.Hash{expectedIlkArtKey},
				BlockNumber: big.NewInt(fakeHeader.BlockNumber),
			}))
		})

		It("inserts returned value", func() {
			fakeFrob := repository.Frob{
				Dink: "1",
				Dart: "0",
			}
			mockEventsRepository.GetFrobsFrobsToReturn = []repository.Frob{fakeFrob}
			fakeUrn := repository.Urn{
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
			mockBlockChain.SetStorageValuesToReturn(fakeHeader.BlockNumber, shared.VatAddress, fakeValue)

			err := backFiller.BackFill(0)

			Expect(err).NotTo(HaveOccurred())
			paddedUrn, padErr := utilities.PadAddress(fakeUrn.Urn)
			Expect(padErr).NotTo(HaveOccurred())
			expectedUrnInkKey := storage.GetKeyForNestedMapping(storage.IndexThree, fakeUrn.Ilk, paddedUrn)
			expectedDiff := types.RawDiff{
				HashedAddress: shared.HashedVatAddress,
				BlockHash:     common.HexToHash(fakeHeader.Hash),
				BlockHeight:   int(fakeHeader.BlockNumber),
				StorageKey:    crypto.Keccak256Hash(expectedUrnInkKey.Bytes()),
				StorageValue:  common.BytesToHash(fakeValue),
			}
			Expect(mockStorageRepository.InsertDiffPassedDiff).To(Equal(expectedDiff))
		})
	})
})
