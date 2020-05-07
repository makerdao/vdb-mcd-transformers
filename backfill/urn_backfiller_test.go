package backfill_test

import (
	"errors"
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/makerdao/vdb-mcd-transformers/backfill"
	"github.com/makerdao/vdb-mcd-transformers/backfill/mocks"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Urn BackFiller", func() {
	var (
		mockBlockChain       *fakes.MockBlockChain
		mockEventsRepository *mocks.EventsRepository
		mockUrnsRepository   *mocks.UrnsRepository
		backFiller           backfill.UrnBackFiller
	)

	BeforeEach(func() {
		mockBlockChain = fakes.NewMockBlockChain()
		mockEventsRepository = &mocks.EventsRepository{}
		mockUrnsRepository = &mocks.UrnsRepository{}
		backFiller = backfill.NewUrnBackFiller(
			mockBlockChain,
			mockEventsRepository,
			mockUrnsRepository,
		)
	})

	Describe("BackFillUrns", func() {
		It("gets urns", func() {
			err := backFiller.BackfillUrns(0)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockUrnsRepository.GetUrnsCalled).To(BeTrue())
		})

		It("returns error if getting urns fails", func() {
			mockUrnsRepository.GetUrnsErr = fakes.FakeError

			err := backFiller.BackfillUrns(0)

			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, fakes.FakeError)).To(BeTrue())
		})

		It("gets frobs for each urn", func() {
			fakeUrnOne := backfill.Urn{
				ID:  rand.Int(),
				Ilk: test_data.RandomString(64),
				Urn: test_data.RandomString(40),
			}
			fakeUrnTwo := backfill.Urn{
				ID:  rand.Int(),
				Ilk: test_data.RandomString(64),
				Urn: test_data.RandomString(40),
			}
			mockUrnsRepository.GetUrnsUrnsToReturn = []backfill.Urn{fakeUrnOne, fakeUrnTwo}

			err := backFiller.BackfillUrns(0)

			Expect(err).NotTo(HaveOccurred())

			Expect(err).NotTo(HaveOccurred())
			Expect(mockEventsRepository.GetFrobsPassedUrnIDs).To(ConsistOf(fakeUrnOne.ID, fakeUrnTwo.ID))
		})

		It("passes starting block when getting frobs to enable filtering results", func() {
			mockUrnsRepository.GetUrnsUrnsToReturn = []backfill.Urn{{}}
			startingBlock := rand.Int()

			err := backFiller.BackfillUrns(startingBlock)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockEventsRepository.GetFrobsPassedStartingBlock).To(Equal(startingBlock))
		})

		It("returns error if getting frobs fails", func() {
			mockUrnsRepository.GetUrnsUrnsToReturn = []backfill.Urn{{}}
			mockEventsRepository.GetFrobsError = fakes.FakeError

			err := backFiller.BackfillUrns(0)

			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, fakes.FakeError)).To(BeTrue())
		})

		It("gets header for each frob if dink or dart is not zero", func() {
			mockUrnsRepository.GetUrnsUrnsToReturn = []backfill.Urn{{Urn: "0x" + test_data.RandomString(40)}}
			fakeFrobOne := backfill.Frob{
				HeaderID: rand.Int(),
				Dink:     "1",
				Dart:     "0",
			}
			fakeFrobTwo := backfill.Frob{
				HeaderID: rand.Int(),
				Dink:     "0",
				Dart:     "1",
			}
			mockEventsRepository.GetFrobsFrobsToReturn = []backfill.Frob{fakeFrobOne, fakeFrobTwo}

			err := backFiller.BackfillUrns(0)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockEventsRepository.GetHeaderByIDPassedIDs).To(ConsistOf(fakeFrobOne.HeaderID, fakeFrobTwo.HeaderID))
		})

		It("does not get header if both dink and dart are zero", func() {
			mockUrnsRepository.GetUrnsUrnsToReturn = []backfill.Urn{{Urn: "0x" + test_data.RandomString(40)}}
			fakeFrob := backfill.Frob{
				HeaderID: rand.Int(),
				Dink:     "0",
				Dart:     "0",
			}
			mockEventsRepository.GetFrobsFrobsToReturn = []backfill.Frob{fakeFrob}

			err := backFiller.BackfillUrns(0)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockEventsRepository.GetHeaderByIDPassedIDs).To(BeEmpty())
		})

		It("returns error if getting header fails", func() {
			mockUrnsRepository.GetUrnsUrnsToReturn = []backfill.Urn{{Urn: "0x" + test_data.RandomString(40)}}
			fakeFrob := backfill.Frob{
				Dink: "1",
				Dart: "0",
			}
			mockEventsRepository.GetFrobsFrobsToReturn = []backfill.Frob{fakeFrob}
			mockEventsRepository.GetHeaderByIDError = fakes.FakeError

			err := backFiller.BackfillUrns(0)

			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, fakes.FakeError)).To(BeTrue())
		})

		It("gets art storage value for urn at header if corresponding vat_urn_art row doesn't exist", func() {
			fakeUrn := backfill.Urn{
				Ilk: "0x" + test_data.RandomString(64),
				Urn: "0x" + test_data.RandomString(40),
			}
			mockUrnsRepository.GetUrnsUrnsToReturn = []backfill.Urn{fakeUrn}
			fakeFrob := backfill.Frob{
				Dink: "0",
				Dart: "1",
			}
			mockEventsRepository.GetFrobsFrobsToReturn = []backfill.Frob{fakeFrob}
			fakeHeader := core.Header{BlockNumber: rand.Int63()}
			mockEventsRepository.GetHeaderByIDHeaderToReturn = fakeHeader

			err := backFiller.BackfillUrns(0)

			Expect(err).NotTo(HaveOccurred())
			paddedUrn, padErr := utilities.PadAddress(fakeUrn.Urn)
			Expect(padErr).NotTo(HaveOccurred())
			urnInkKey := storage.GetKeyForNestedMapping(storage.IndexThree, fakeUrn.Ilk, paddedUrn)
			expectedUrnArtKey := storage.GetIncrementedKey(urnInkKey, 1)
			Expect(mockBlockChain.BatchGetStorageAtCalls).To(ContainElement(
				fakes.BatchGetStorageAtCall{
					Account:     backfill.VatAddress,
					Keys:        []common.Hash{expectedUrnArtKey},
					BlockNumber: big.NewInt(fakeHeader.BlockNumber),
				},
			))
		})

		It("does not get art storage value if corresponding vat_urn_art row already exists", func() {
			mockUrnsRepository.GetUrnsUrnsToReturn = []backfill.Urn{{}}
			fakeFrob := backfill.Frob{
				Dink: "0",
				Dart: "1",
			}
			mockEventsRepository.GetFrobsFrobsToReturn = []backfill.Frob{fakeFrob}
			mockEventsRepository.GetHeaderByIDHeaderToReturn = core.Header{}
			mockUrnsRepository.VatUrnArtExistsBoolToReturn = true

			err := backFiller.BackfillUrns(0)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockBlockChain.BatchGetStorageAtCalls).To(BeEmpty())
		})

		It("gets ink storage value for urn at header if corresponding vat_urn_ink row doesn't exist", func() {
			fakeUrn := backfill.Urn{
				Ilk: "0x" + test_data.RandomString(64),
				Urn: "0x" + test_data.RandomString(40),
			}
			mockUrnsRepository.GetUrnsUrnsToReturn = []backfill.Urn{fakeUrn}
			fakeFrob := backfill.Frob{
				Dink: "1",
				Dart: "0",
			}
			mockEventsRepository.GetFrobsFrobsToReturn = []backfill.Frob{fakeFrob}
			fakeHeader := core.Header{BlockNumber: rand.Int63()}
			mockEventsRepository.GetHeaderByIDHeaderToReturn = fakeHeader

			err := backFiller.BackfillUrns(0)

			Expect(err).NotTo(HaveOccurred())
			paddedUrn, padErr := utilities.PadAddress(fakeUrn.Urn)
			Expect(padErr).NotTo(HaveOccurred())
			expectedUrnInkKey := storage.GetKeyForNestedMapping(storage.IndexThree, fakeUrn.Ilk, paddedUrn)
			Expect(mockBlockChain.BatchGetStorageAtCalls).To(ContainElement(
				fakes.BatchGetStorageAtCall{
					Account:     backfill.VatAddress,
					Keys:        []common.Hash{expectedUrnInkKey},
					BlockNumber: big.NewInt(fakeHeader.BlockNumber),
				},
			))
		})

		It("does not get ink storage value if corresponding vat_urn_ink row already exists", func() {
			mockUrnsRepository.GetUrnsUrnsToReturn = []backfill.Urn{{}}
			fakeFrob := backfill.Frob{
				Dink: "1",
				Dart: "0",
			}
			mockEventsRepository.GetFrobsFrobsToReturn = []backfill.Frob{fakeFrob}
			mockEventsRepository.GetHeaderByIDHeaderToReturn = core.Header{}
			mockUrnsRepository.VatUrnInkExistsBoolToReturn = true

			err := backFiller.BackfillUrns(0)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockBlockChain.BatchGetStorageAtCalls).To(BeEmpty())
		})

		It("returns error if getting storage value fails", func() {
			fakeUrn := backfill.Urn{
				Urn: "0x" + test_data.RandomString(40),
			}
			mockUrnsRepository.GetUrnsUrnsToReturn = []backfill.Urn{fakeUrn}
			fakeFrob := backfill.Frob{
				Dink: "1",
				Dart: "0",
			}
			mockEventsRepository.GetFrobsFrobsToReturn = []backfill.Frob{fakeFrob}
			mockEventsRepository.GetHeaderByIDHeaderToReturn = core.Header{BlockNumber: rand.Int63()}
			mockBlockChain.BatchGetStorageAtError = fakes.FakeError

			err := backFiller.BackfillUrns(0)

			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, fakes.FakeError)).To(BeTrue())
		})

		It("persists storage value returned by chain", func() {
			fakeUrn := backfill.Urn{
				Ilk: "0x" + test_data.RandomString(64),
				Urn: "0x" + test_data.RandomString(40),
			}
			mockUrnsRepository.GetUrnsUrnsToReturn = []backfill.Urn{fakeUrn}
			fakeFrob := backfill.Frob{
				Dink: "1",
				Dart: "0",
			}
			mockEventsRepository.GetFrobsFrobsToReturn = []backfill.Frob{fakeFrob}
			fakeHeader := core.Header{
				BlockNumber: rand.Int63(),
				Hash:        test_data.RandomString(64),
			}
			mockEventsRepository.GetHeaderByIDHeaderToReturn = fakeHeader
			fakeValue := []byte{0, 1, 2, 3, 4, 5}
			mockBlockChain.SetStorageValuesToReturn(fakeHeader.BlockNumber, backfill.VatAddress, fakeValue)

			err := backFiller.BackfillUrns(0)

			Expect(err).NotTo(HaveOccurred())
			paddedUrn, padErr := utilities.PadAddress(fakeUrn.Urn)
			Expect(padErr).NotTo(HaveOccurred())
			expectedUrnInkKey := storage.GetKeyForNestedMapping(storage.IndexThree, fakeUrn.Ilk, paddedUrn)
			expectedDiff := types.RawDiff{
				HashedAddress: backfill.HashedVatAddress,
				BlockHash:     common.HexToHash(fakeHeader.Hash),
				BlockHeight:   int(fakeHeader.BlockNumber),
				StorageKey:    crypto.Keccak256Hash(expectedUrnInkKey.Bytes()),
				StorageValue:  common.BytesToHash(fakeValue),
			}
			Expect(mockUrnsRepository.InsertUrnDiffPassedDiff).To(Equal(expectedDiff))
		})

		It("returns error if persisting diff fails", func() {
			fakeUrn := backfill.Urn{
				Urn: "0x" + test_data.RandomString(40),
			}
			mockUrnsRepository.GetUrnsUrnsToReturn = []backfill.Urn{fakeUrn}
			fakeFrob := backfill.Frob{
				Dink: "1",
				Dart: "0",
			}
			mockEventsRepository.GetFrobsFrobsToReturn = []backfill.Frob{fakeFrob}
			fakeHeader := core.Header{BlockNumber: rand.Int63()}
			mockEventsRepository.GetHeaderByIDHeaderToReturn = fakeHeader
			mockBlockChain.SetStorageValuesToReturn(fakeHeader.BlockNumber, backfill.VatAddress, []byte{0, 1, 2, 3, 4, 5})
			mockUrnsRepository.InsertUrnDiffErr = fakes.FakeError

			err := backFiller.BackfillUrns(0)

			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, fakes.FakeError)).To(BeTrue())
		})
	})
})
