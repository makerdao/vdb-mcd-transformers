package frob_test

import (
	"errors"
	"math/rand"

	"github.com/makerdao/vdb-mcd-transformers/backfill"
	"github.com/makerdao/vdb-mcd-transformers/backfill/frob"
	"github.com/makerdao/vdb-mcd-transformers/backfill/mocks"
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
	})
})
