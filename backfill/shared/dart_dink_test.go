package shared_test

import (
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
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

var _ = Describe("Dart Dink helpers", func() {
	Describe("RetrieveDartDinkDiffs", func() {
		var (
			mockBlockChain        *fakes.MockBlockChain
			mockEventsRepository  *mocks.EventsRepository
			mockHeaderRepository  *fakes.MockHeaderRepository
			mockStorageRepository *mocks.StorageRepository
			dartDinkRetriever     shared.DartDinkRetriever
		)

		BeforeEach(func() {
			mockBlockChain = fakes.NewMockBlockChain()
			mockEventsRepository = &mocks.EventsRepository{}
			mockHeaderRepository = &fakes.MockHeaderRepository{}
			mockStorageRepository = &mocks.StorageRepository{}
			dartDinkRetriever = shared.NewDartDinkRetriever(mockBlockChain, mockEventsRepository, mockHeaderRepository, mockStorageRepository)
		})

		It("ignores row if dink and dart are zero", func() {
			fakeData := shared.DartDink{
				Dink: "0",
				Dart: "0",
			}

			err := dartDinkRetriever.RetrieveDartDinkDiffs(fakeData)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockBlockChain.BatchGetStorageAtCalls).To(BeEmpty())
		})

		It("passes ilk ID and header ID to detect if ilk art exists", func() {
			fakeData := shared.DartDink{
				HeaderID: rand.Int63(),
				Dink:     "0",
				Dart:     "1",
			}
			fakeUrn := repository.Urn{IlkID: rand.Int63()}
			mockStorageRepository.GetUrnByIDUrnToReturn = fakeUrn
			mockStorageRepository.VatIlkArtExistsBoolToReturn = true
			mockStorageRepository.VatUrnArtExistsBoolToReturn = true
			mockStorageRepository.VatUrnInkExistsBoolToReturn = true

			err := dartDinkRetriever.RetrieveDartDinkDiffs(fakeData)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockStorageRepository.VatIlkArtExistsPassedIlkID).To(Equal(fakeUrn.IlkID))
			Expect(mockStorageRepository.VatIlkArtExistsPassedHeaderID).To(Equal(fakeData.HeaderID))
		})

		It("passes urn ID and header ID to detect if urn art exists", func() {
			fakeData := shared.DartDink{
				HeaderID: rand.Int63(),
				UrnID:    rand.Int63(),
				Dink:     "0",
				Dart:     "1",
			}
			mockStorageRepository.VatIlkArtExistsBoolToReturn = true
			mockStorageRepository.VatUrnArtExistsBoolToReturn = true
			mockStorageRepository.VatUrnInkExistsBoolToReturn = true

			err := dartDinkRetriever.RetrieveDartDinkDiffs(fakeData)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockStorageRepository.VatUrnArtExistsPassedHeaderID).To(Equal(fakeData.HeaderID))
			Expect(mockStorageRepository.VatUrnArtExistsPassedUrnID).To(Equal(fakeData.UrnID))
		})

		It("passes urn ID and header ID to detect if urn ink exists", func() {
			fakeData := shared.DartDink{
				HeaderID: rand.Int63(),
				UrnID:    rand.Int63(),
				Dink:     "1",
				Dart:     "0",
			}
			mockStorageRepository.VatIlkArtExistsBoolToReturn = true
			mockStorageRepository.VatUrnArtExistsBoolToReturn = true
			mockStorageRepository.VatUrnInkExistsBoolToReturn = true

			err := dartDinkRetriever.RetrieveDartDinkDiffs(fakeData)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockStorageRepository.VatUrnInkExistsPassedHeaderID).To(Equal(fakeData.HeaderID))
			Expect(mockStorageRepository.VatUrnInkExistsPassedUrnID).To(Equal(fakeData.UrnID))
		})

		It("ignores data if transformed diffs already exists at header", func() {
			fakeData := shared.DartDink{
				Dink: "1",
				Dart: "1",
			}
			mockStorageRepository.VatIlkArtExistsBoolToReturn = true
			mockStorageRepository.VatUrnArtExistsBoolToReturn = true
			mockStorageRepository.VatUrnInkExistsBoolToReturn = true

			err := dartDinkRetriever.RetrieveDartDinkDiffs(fakeData)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockBlockChain.BatchGetStorageAtCalls).To(BeEmpty())
		})

		It("returns error if getting for urn for data requiring back-fill fails", func() {
			fakeData := shared.DartDink{
				UrnID: rand.Int63(),
				Dink:  "0",
				Dart:  "1",
			}
			mockStorageRepository.VatIlkArtExistsBoolToReturn = false
			mockStorageRepository.GetUrnByIDError = fakes.FakeError

			err := dartDinkRetriever.RetrieveDartDinkDiffs(fakeData)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("returns error if getting header for data requiring back-fill fails", func() {
			fakeData := shared.DartDink{
				UrnID: rand.Int63(),
				Dink:  "0",
				Dart:  "1",
			}
			mockStorageRepository.VatIlkArtExistsBoolToReturn = false
			mockHeaderRepository.GetHeaderByIDError = fakes.FakeError

			err := dartDinkRetriever.RetrieveDartDinkDiffs(fakeData)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("looks up storage for data when some values are non-zero and don't already exist", func() {
			fakeData := shared.DartDink{
				Dink: "0",
				Dart: "1",
			}
			mockStorageRepository.VatIlkArtExistsBoolToReturn = false
			mockStorageRepository.VatUrnArtExistsBoolToReturn = true
			mockStorageRepository.VatUrnInkExistsBoolToReturn = true
			fakeUrn := repository.Urn{Ilk: test_data.RandomString(64)}
			mockStorageRepository.GetUrnByIDUrnToReturn = fakeUrn
			fakeHeader := core.Header{BlockNumber: rand.Int63()}
			mockHeaderRepository.GetHeaderByIDHeaderToReturn = fakeHeader

			err := dartDinkRetriever.RetrieveDartDinkDiffs(fakeData)

			Expect(err).NotTo(HaveOccurred())
			expectedIlkArtKey := storage.GetKeyForMapping(storage.IndexTwo, fakeUrn.Ilk)
			Expect(mockBlockChain.BatchGetStorageAtCalls).To(ConsistOf(fakes.BatchGetStorageAtCall{
				Account:     shared.VatAddress,
				Keys:        []common.Hash{expectedIlkArtKey},
				BlockNumber: big.NewInt(fakeHeader.BlockNumber),
			}))
		})

		It("inserts returned value", func() {
			fakeData := shared.DartDink{
				Dink: "1",
				Dart: "0",
			}
			fakeUrn := repository.Urn{
				Ilk: test_data.RandomString(64),
				Urn: "0x" + test_data.RandomString(40),
			}
			mockStorageRepository.GetUrnByIDUrnToReturn = fakeUrn
			fakeHeader := core.Header{
				BlockNumber: rand.Int63(),
				Hash:        test_data.RandomString(64),
			}
			mockHeaderRepository.GetHeaderByIDHeaderToReturn = fakeHeader
			fakeValue := []byte{1, 2, 3, 4, 5}
			mockBlockChain.SetStorageValuesToReturn(fakeHeader.BlockNumber, shared.VatAddress, fakeValue)

			err := dartDinkRetriever.RetrieveDartDinkDiffs(fakeData)

			Expect(err).NotTo(HaveOccurred())
			paddedUrn, padErr := utilities.PadAddress(fakeUrn.Urn)
			Expect(padErr).NotTo(HaveOccurred())
			expectedUrnInkKey := storage.GetKeyForNestedMapping(storage.IndexThree, fakeUrn.Ilk, paddedUrn)
			expectedDiff := types.RawDiff{
				Address:      shared.VatAddress,
				BlockHash:    common.HexToHash(fakeHeader.Hash),
				BlockHeight:  int(fakeHeader.BlockNumber),
				StorageKey:   expectedUrnInkKey,
				StorageValue: common.BytesToHash(fakeValue),
			}
			Expect(mockStorageRepository.InsertDiffPassedDiff).To(Equal(expectedDiff))
		})
	})
})
