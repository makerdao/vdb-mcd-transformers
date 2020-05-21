package fork

import (
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/backfill"
	"github.com/makerdao/vdb-mcd-transformers/backfill/mocks"
	"github.com/makerdao/vdb-mcd-transformers/backfill/repository"
	"github.com/makerdao/vdb-mcd-transformers/backfill/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Fork BackFiller", func() {
	var (
		mockBlockChain        *fakes.MockBlockChain
		mockDartDinkRetriever *mocks.MockDartDinkRetriever
		mockEventsRepository  *mocks.EventsRepository
		mockStorageRepository *mocks.StorageRepository
		backFiller            backfill.BackFiller
	)

	BeforeEach(func() {
		mockBlockChain = fakes.NewMockBlockChain()
		mockDartDinkRetriever = &mocks.MockDartDinkRetriever{}
		mockEventsRepository = &mocks.EventsRepository{}
		mockStorageRepository = &mocks.StorageRepository{}
		backFiller = NewForkBackFiller(mockBlockChain, mockEventsRepository, mockStorageRepository, mockDartDinkRetriever)
	})

	Describe("BackFill", func() {
		It("gets vat fork events from starting block number onward", func() {
			startingBlock := rand.Int()

			err := backFiller.BackFill(startingBlock)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockEventsRepository.GetForksPassedStartingBlock).To(Equal(startingBlock))
		})

		It("returns error if getting fork events fails", func() {
			mockEventsRepository.GetForksError = fakes.FakeError

			err := backFiller.BackFill(0)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("passes DartDink derived from fork src to retriever", func() {
			fakeFork := repository.Fork{
				HeaderID: rand.Int(),
				Ilk:      test_data.RandomString(64),
				Src:      test_data.RandomString(40),
				Dst:      "",
				Dink:     strconv.Itoa(rand.Int()),
				Dart:     strconv.Itoa(rand.Int()),
			}
			mockEventsRepository.GetForksForksToReturn = []repository.Fork{fakeFork}
			fakeUrnID := rand.Int63()
			mockStorageRepository.GetOrCreateUrnIDsToReturn = map[string]int64{fakeFork.Src: fakeUrnID}

			err := backFiller.BackFill(0)

			Expect(err).NotTo(HaveOccurred())
			expectedDartDink := shared.DartDink{
				Dart:     fakeFork.Dart,
				Dink:     fakeFork.Dink,
				HeaderID: fakeFork.HeaderID,
				UrnID:    int(fakeUrnID),
			}
			Expect(mockDartDinkRetriever.PassedDartDinks).To(ContainElement(expectedDartDink))
		})

		It("passes DartDink derived from fork dst to retriever", func() {
			fakeFork := repository.Fork{
				HeaderID: rand.Int(),
				Ilk:      test_data.RandomString(64),
				Src:      "",
				Dst:      test_data.RandomString(40),
				Dink:     strconv.Itoa(rand.Int()),
				Dart:     strconv.Itoa(rand.Int()),
			}
			mockEventsRepository.GetForksForksToReturn = []repository.Fork{fakeFork}
			fakeUrnID := rand.Int63()
			mockStorageRepository.GetOrCreateUrnIDsToReturn = map[string]int64{fakeFork.Dst: fakeUrnID}

			err := backFiller.BackFill(0)

			Expect(err).NotTo(HaveOccurred())
			expectedDartDink := shared.DartDink{
				Dart:     fakeFork.Dart,
				Dink:     fakeFork.Dink,
				HeaderID: fakeFork.HeaderID,
				UrnID:    int(fakeUrnID),
			}
			Expect(mockDartDinkRetriever.PassedDartDinks).To(ContainElement(expectedDartDink))
		})

		It("returns error if getting urn ID for src/dst fails", func() {
			mockEventsRepository.GetForksForksToReturn = []repository.Fork{{}}
			mockStorageRepository.GetOrCreateUrnError = fakes.FakeError

			err := backFiller.BackFill(0)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("returns error if retrieving DartDink fails", func() {
			mockEventsRepository.GetForksForksToReturn = []repository.Fork{{}}
			mockDartDinkRetriever.RetrieveDartDinkDiffsError = fakes.FakeError

			err := backFiller.BackFill(0)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})
})
