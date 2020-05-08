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

var _ = Describe("Frob BackFiller", func() {
	var (
		mockBlockChain        *fakes.MockBlockChain
		mockEventsRepository  *mocks.EventsRepository
		mockStorageRepository *mocks.StorageRepository
		backFiller            backfill.FrobBackFiller
	)

	BeforeEach(func() {
		mockBlockChain = fakes.NewMockBlockChain()
		mockEventsRepository = &mocks.EventsRepository{}
		mockStorageRepository = &mocks.StorageRepository{}
		backFiller = backfill.NewFrobBackFiller(
			mockBlockChain,
			mockEventsRepository,
			mockStorageRepository,
		)
	})

	Describe("BackFillFrobStorage", func() {
		It("gets urns", func() {
			err := backFiller.BackFillFrobStorage(0)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockStorageRepository.GetUrnsCalled).To(BeTrue())
		})

		It("returns error if getting urns fails", func() {
			mockStorageRepository.GetUrnsErr = fakes.FakeError

			err := backFiller.BackFillFrobStorage(0)

			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, fakes.FakeError)).To(BeTrue())
		})

		It("gets frobs for each urn", func() {
			fakeUrnOne := backfill.Urn{
				UrnID: rand.Int(),
				Ilk:   test_data.RandomString(64),
				Urn:   test_data.RandomString(40),
			}
			fakeUrnTwo := backfill.Urn{
				UrnID: rand.Int(),
				Ilk:   test_data.RandomString(64),
				Urn:   test_data.RandomString(40),
			}
			mockStorageRepository.GetUrnsUrnsToReturn = []backfill.Urn{fakeUrnOne, fakeUrnTwo}

			err := backFiller.BackFillFrobStorage(0)

			Expect(err).NotTo(HaveOccurred())

			Expect(err).NotTo(HaveOccurred())
			Expect(mockEventsRepository.GetFrobsPassedUrnIDs).To(ConsistOf(fakeUrnOne.UrnID, fakeUrnTwo.UrnID))
		})

		It("passes starting block when getting frobs to enable filtering results", func() {
			mockStorageRepository.GetUrnsUrnsToReturn = []backfill.Urn{{}}
			startingBlock := rand.Int()

			err := backFiller.BackFillFrobStorage(startingBlock)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockEventsRepository.GetFrobsPassedStartingBlock).To(Equal(startingBlock))
		})

		It("returns error if getting frobs fails", func() {
			mockStorageRepository.GetUrnsUrnsToReturn = []backfill.Urn{{}}
			mockEventsRepository.GetFrobsError = fakes.FakeError

			err := backFiller.BackFillFrobStorage(0)

			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, fakes.FakeError)).To(BeTrue())
		})

		It("gets header for each frob if dink or dart is not zero", func() {
			mockStorageRepository.GetUrnsUrnsToReturn = []backfill.Urn{{Urn: "0x" + test_data.RandomString(40)}}
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

			err := backFiller.BackFillFrobStorage(0)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockEventsRepository.GetHeaderByIDPassedIDs).To(ConsistOf(fakeFrobOne.HeaderID, fakeFrobTwo.HeaderID))
		})

		It("does not get header if both dink and dart are zero", func() {
			mockStorageRepository.GetUrnsUrnsToReturn = []backfill.Urn{{Urn: "0x" + test_data.RandomString(40)}}
			fakeFrob := backfill.Frob{
				HeaderID: rand.Int(),
				Dink:     "0",
				Dart:     "0",
			}
			mockEventsRepository.GetFrobsFrobsToReturn = []backfill.Frob{fakeFrob}

			err := backFiller.BackFillFrobStorage(0)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockEventsRepository.GetHeaderByIDPassedIDs).To(BeEmpty())
		})

		It("returns error if getting header fails", func() {
			mockStorageRepository.GetUrnsUrnsToReturn = []backfill.Urn{{Urn: "0x" + test_data.RandomString(40)}}
			fakeFrob := backfill.Frob{
				Dink: "1",
				Dart: "0",
			}
			mockEventsRepository.GetFrobsFrobsToReturn = []backfill.Frob{fakeFrob}
			mockEventsRepository.GetHeaderByIDError = fakes.FakeError

			err := backFiller.BackFillFrobStorage(0)

			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, fakes.FakeError)).To(BeTrue())
		})

		It("gets urn_art storage value at header if corresponding vat_urn_art row doesn't exist", func() {
			fakeUrn := backfill.Urn{
				Ilk: "0x" + test_data.RandomString(64),
				Urn: "0x" + test_data.RandomString(40),
			}
			mockStorageRepository.GetUrnsUrnsToReturn = []backfill.Urn{fakeUrn}
			fakeFrob := backfill.Frob{
				Dink: "0",
				Dart: "1",
			}
			mockEventsRepository.GetFrobsFrobsToReturn = []backfill.Frob{fakeFrob}
			fakeHeader := core.Header{BlockNumber: rand.Int63()}
			mockEventsRepository.GetHeaderByIDHeaderToReturn = fakeHeader
			mockStorageRepository.VatIlkArtExistsBoolToReturn = true

			err := backFiller.BackFillFrobStorage(0)

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

		It("does not get urn_art storage value if corresponding vat_urn_art row already exists", func() {
			mockStorageRepository.GetUrnsUrnsToReturn = []backfill.Urn{{}}
			fakeFrob := backfill.Frob{
				Dink: "0",
				Dart: "1",
			}
			mockEventsRepository.GetFrobsFrobsToReturn = []backfill.Frob{fakeFrob}
			mockEventsRepository.GetHeaderByIDHeaderToReturn = core.Header{}
			mockStorageRepository.VatUrnArtExistsBoolToReturn = true
			mockStorageRepository.VatIlkArtExistsBoolToReturn = true

			err := backFiller.BackFillFrobStorage(0)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockBlockChain.BatchGetStorageAtCalls).To(BeEmpty())
		})

		It("gets ilk_art storage value at header if corresponding vat_urn_art row doesn't exist", func() {
			fakeUrn := backfill.Urn{
				Ilk:   test_data.RandomString(64),
				IlkID: rand.Int(),
			}
			mockStorageRepository.GetUrnsUrnsToReturn = []backfill.Urn{fakeUrn}
			fakeFrob := backfill.Frob{
				Dink: "0",
				Dart: "1",
			}
			mockEventsRepository.GetFrobsFrobsToReturn = []backfill.Frob{fakeFrob}
			fakeHeader := core.Header{BlockNumber: rand.Int63()}
			mockEventsRepository.GetHeaderByIDHeaderToReturn = fakeHeader
			mockStorageRepository.VatUrnArtExistsBoolToReturn = true

			err := backFiller.BackFillFrobStorage(0)

			Expect(err).NotTo(HaveOccurred())

			expectedIlkArtKey := storage.GetKeyForMapping(storage.IndexTwo, fakeUrn.Ilk)
			Expect(mockBlockChain.BatchGetStorageAtCalls).To(ContainElement(
				fakes.BatchGetStorageAtCall{
					Account:     backfill.VatAddress,
					Keys:        []common.Hash{expectedIlkArtKey},
					BlockNumber: big.NewInt(fakeHeader.BlockNumber),
				},
			))
		})

		It("does not get ilk_art storage value if corresponding vat_ilk_art row already exists", func() {
			mockStorageRepository.GetUrnsUrnsToReturn = []backfill.Urn{{}}
			fakeFrob := backfill.Frob{
				Dink: "0",
				Dart: "1",
			}
			mockEventsRepository.GetFrobsFrobsToReturn = []backfill.Frob{fakeFrob}
			mockEventsRepository.GetHeaderByIDHeaderToReturn = core.Header{}
			mockStorageRepository.VatUrnArtExistsBoolToReturn = true
			mockStorageRepository.VatIlkArtExistsBoolToReturn = true

			err := backFiller.BackFillFrobStorage(0)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockBlockChain.BatchGetStorageAtCalls).To(BeEmpty())
		})

		It("gets urn_ink storage value for urn at header if corresponding vat_urn_ink row doesn't exist", func() {
			fakeUrn := backfill.Urn{
				Ilk: "0x" + test_data.RandomString(64),
				Urn: "0x" + test_data.RandomString(40),
			}
			mockStorageRepository.GetUrnsUrnsToReturn = []backfill.Urn{fakeUrn}
			fakeFrob := backfill.Frob{
				Dink: "1",
				Dart: "0",
			}
			mockEventsRepository.GetFrobsFrobsToReturn = []backfill.Frob{fakeFrob}
			fakeHeader := core.Header{BlockNumber: rand.Int63()}
			mockEventsRepository.GetHeaderByIDHeaderToReturn = fakeHeader

			err := backFiller.BackFillFrobStorage(0)

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

		It("does not get urn_ink storage value if corresponding vat_urn_ink row already exists", func() {
			mockStorageRepository.GetUrnsUrnsToReturn = []backfill.Urn{{}}
			fakeFrob := backfill.Frob{
				Dink: "1",
				Dart: "0",
			}
			mockEventsRepository.GetFrobsFrobsToReturn = []backfill.Frob{fakeFrob}
			mockEventsRepository.GetHeaderByIDHeaderToReturn = core.Header{}
			mockStorageRepository.VatUrnInkExistsBoolToReturn = true

			err := backFiller.BackFillFrobStorage(0)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockBlockChain.BatchGetStorageAtCalls).To(BeEmpty())
		})

		It("returns error if getting storage value fails", func() {
			fakeUrn := backfill.Urn{
				Urn: "0x" + test_data.RandomString(40),
			}
			mockStorageRepository.GetUrnsUrnsToReturn = []backfill.Urn{fakeUrn}
			fakeFrob := backfill.Frob{
				Dink: "1",
				Dart: "0",
			}
			mockEventsRepository.GetFrobsFrobsToReturn = []backfill.Frob{fakeFrob}
			mockEventsRepository.GetHeaderByIDHeaderToReturn = core.Header{BlockNumber: rand.Int63()}
			mockBlockChain.BatchGetStorageAtError = fakes.FakeError

			err := backFiller.BackFillFrobStorage(0)

			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, fakes.FakeError)).To(BeTrue())
		})

		It("persists storage value returned by chain", func() {
			fakeUrn := backfill.Urn{
				Ilk: "0x" + test_data.RandomString(64),
				Urn: "0x" + test_data.RandomString(40),
			}
			mockStorageRepository.GetUrnsUrnsToReturn = []backfill.Urn{fakeUrn}
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

			err := backFiller.BackFillFrobStorage(0)

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
			Expect(mockStorageRepository.InsertDiffPassedDiff).To(Equal(expectedDiff))
		})

		It("returns error if persisting diff fails", func() {
			fakeUrn := backfill.Urn{
				Urn: "0x" + test_data.RandomString(40),
			}
			mockStorageRepository.GetUrnsUrnsToReturn = []backfill.Urn{fakeUrn}
			fakeFrob := backfill.Frob{
				Dink: "1",
				Dart: "0",
			}
			mockEventsRepository.GetFrobsFrobsToReturn = []backfill.Frob{fakeFrob}
			fakeHeader := core.Header{BlockNumber: rand.Int63()}
			mockEventsRepository.GetHeaderByIDHeaderToReturn = fakeHeader
			mockBlockChain.SetStorageValuesToReturn(fakeHeader.BlockNumber, backfill.VatAddress, []byte{0, 1, 2, 3, 4, 5})
			mockStorageRepository.InsertDiffErr = fakes.FakeError

			err := backFiller.BackFillFrobStorage(0)

			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, fakes.FakeError)).To(BeTrue())
		})
	})
})
