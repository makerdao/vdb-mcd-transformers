package grab_test

import (
	"errors"
	"math/rand"

	"github.com/makerdao/vdb-mcd-transformers/backfill"
	"github.com/makerdao/vdb-mcd-transformers/backfill/grab"
	"github.com/makerdao/vdb-mcd-transformers/backfill/mocks"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Grab BackFiller", func() {
	Describe("BackFill", func() {
		var (
			backFiller            backfill.BackFiller
			mockBlockChain        *fakes.MockBlockChain
			mockEventsRepository  *mocks.EventsRepository
			mockStorageRepository *mocks.StorageRepository
		)

		BeforeEach(func() {
			mockBlockChain = fakes.NewMockBlockChain()
			mockEventsRepository = &mocks.EventsRepository{}
			mockStorageRepository = &mocks.StorageRepository{}
			backFiller = grab.NewGrabBackFiller(mockBlockChain, mockEventsRepository, mockStorageRepository)
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
			Expect(errors.Is(err, fakes.FakeError)).To(BeTrue())
		})
	})
})
