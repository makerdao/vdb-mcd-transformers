package grab_test

import (
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/backfill"
	"github.com/makerdao/vdb-mcd-transformers/backfill/grab"
	"github.com/makerdao/vdb-mcd-transformers/backfill/mocks"
	"github.com/makerdao/vdb-mcd-transformers/backfill/repository"
	"github.com/makerdao/vdb-mcd-transformers/backfill/shared"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Grab BackFiller", func() {
	Describe("BackFill", func() {
		var (
			backFiller            backfill.BackFiller
			mockBlockChain        *fakes.MockBlockChain
			mockDartDinkRetriever *mocks.MockDartDinkRetriever
			mockEventsRepository  *mocks.EventsRepository
			mockStorageRepository *mocks.StorageRepository
		)

		BeforeEach(func() {
			mockBlockChain = fakes.NewMockBlockChain()
			mockDartDinkRetriever = &mocks.MockDartDinkRetriever{}
			mockEventsRepository = &mocks.EventsRepository{}
			mockStorageRepository = &mocks.StorageRepository{}
			backFiller = grab.NewGrabBackFiller(mockBlockChain, mockEventsRepository, mockStorageRepository, mockDartDinkRetriever)
		})

		It("gets vat grab events from starting block number onward", func() {
			startingBlock := rand.Int()

			err := backFiller.BackFill(startingBlock)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockEventsRepository.GetGrabsPassedStartingBlock).To(Equal(startingBlock))
		})

		It("returns error if getting grab events fails", func() {
			mockEventsRepository.GetGrabsError = fakes.FakeError

			err := backFiller.BackFill(0)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("passes DartDink derived from grab to retriever", func() {
			fakeGrab := repository.Grab{
				HeaderID: rand.Int63(),
				UrnID:    rand.Int63(),
				Dink:     strconv.Itoa(rand.Int()),
				Dart:     strconv.Itoa(rand.Int()),
			}
			mockEventsRepository.GetGrabsGrabsToReturn = []repository.Grab{fakeGrab}

			err := backFiller.BackFill(0)

			Expect(err).NotTo(HaveOccurred())
			expectedDartDink := shared.DartDink{
				Dart:     fakeGrab.Dart,
				Dink:     fakeGrab.Dink,
				HeaderID: fakeGrab.HeaderID,
				UrnID:    fakeGrab.UrnID,
			}
			Expect(mockDartDinkRetriever.PassedDartDinks).To(ContainElement(expectedDartDink))
		})

		It("returns error if retrieving DartDink fails", func() {
			mockEventsRepository.GetGrabsGrabsToReturn = []repository.Grab{{}}
			mockDartDinkRetriever.RetrieveDartDinkDiffsError = fakes.FakeError

			err := backFiller.BackFill(0)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})
})
