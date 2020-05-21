package frob_test

import (
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/backfill"
	"github.com/makerdao/vdb-mcd-transformers/backfill/frob"
	"github.com/makerdao/vdb-mcd-transformers/backfill/mocks"
	"github.com/makerdao/vdb-mcd-transformers/backfill/repository"
	"github.com/makerdao/vdb-mcd-transformers/backfill/shared"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Frob BackFiller", func() {
	var (
		mockBlockChain        *fakes.MockBlockChain
		mockEventsRepository  *mocks.EventsRepository
		mockStorageRepository *mocks.StorageRepository
		mockDartDinkRetriever *mocks.MockDartDinkRetriever
		backFiller            backfill.BackFiller
	)

	BeforeEach(func() {
		mockBlockChain = fakes.NewMockBlockChain()
		mockEventsRepository = &mocks.EventsRepository{}
		mockStorageRepository = &mocks.StorageRepository{}
		mockDartDinkRetriever = &mocks.MockDartDinkRetriever{}
		backFiller = frob.NewFrobBackFiller(
			mockBlockChain,
			mockEventsRepository,
			mockStorageRepository,
			mockDartDinkRetriever,
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
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("passes DartDink derived from frob to retriever", func() {
			fakeFrob := repository.Frob{
				HeaderID: rand.Int(),
				UrnID:    rand.Int(),
				Dink:     strconv.Itoa(rand.Int()),
				Dart:     strconv.Itoa(rand.Int()),
			}
			mockEventsRepository.GetFrobsFrobsToReturn = []repository.Frob{fakeFrob}

			err := backFiller.BackFill(0)

			Expect(err).NotTo(HaveOccurred())
			expectedDartDink := shared.DartDink{
				Dart:     fakeFrob.Dart,
				Dink:     fakeFrob.Dink,
				HeaderID: fakeFrob.HeaderID,
				UrnID:    fakeFrob.UrnID,
			}
			Expect(mockDartDinkRetriever.PassedDartDinks).To(ContainElement(expectedDartDink))
		})

		It("returns error if retrieving DartDink fails", func() {
			mockEventsRepository.GetFrobsFrobsToReturn = []repository.Frob{{}}
			mockDartDinkRetriever.RetrieveDartDinkDiffsError = fakes.FakeError

			err := backFiller.BackFill(0)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})
})
