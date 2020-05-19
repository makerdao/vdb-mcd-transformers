package repository_test

import (
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/backfill/repository"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Events repository", func() {
	var (
		db   = test_config.NewTestDB(test_config.NewTestNode())
		repo repository.EventsRepository
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		repo = repository.NewEventsRepository(db)
	})

	Describe("GetFrobs", func() {
		Describe("when there are frobs", func() {
			var (
				urnID, startingBlock, earlierBlock      int
				frobAtStartingBlock, frobAtEarlierBlock repository.Frob
			)

			BeforeEach(func() {
				fakeIlk := test_data.RandomString(64)
				fakeIlkIdentifier := "ETH-A"
				var ilkID int
				ilkErr := db.Get(&ilkID, `INSERT INTO maker.ilks (ilk, identifier) VALUES ($1, $2) RETURNING id`,
					fakeIlk, fakeIlkIdentifier)
				Expect(ilkErr).NotTo(HaveOccurred())

				fakeUrn := test_data.RandomString(40)
				urnErr := db.Get(&urnID, `INSERT INTO maker.urns (ilk_id, identifier) VALUES ($1, $2) RETURNING id`,
					ilkID, fakeUrn)
				Expect(urnErr).NotTo(HaveOccurred())

				startingBlock = rand.Int()
				var startingBlockID int
				startingBlockErr := db.Get(&startingBlockID, `
				INSERT INTO public.headers (block_number, hash, eth_node_id) VALUES ($1, $2, $3) RETURNING id`,
					startingBlock, test_data.RandomString(64), db.NodeID)
				Expect(startingBlockErr).NotTo(HaveOccurred())

				earlierBlock = startingBlock - 1
				var earlierBlockID int
				earlierBlockErr := db.Get(&earlierBlockID, `
				INSERT INTO public.headers (block_number, hash, eth_node_id) VALUES ($1, $2, $3) RETURNING id`,
					earlierBlock, test_data.RandomString(64), db.NodeID)
				Expect(earlierBlockErr).NotTo(HaveOccurred())

				frobAtStartingBlock = repository.Frob{
					HeaderID: startingBlockID,
					UrnID:    urnID,
					Dink:     strconv.Itoa(rand.Int()),
					Dart:     strconv.Itoa(rand.Int()),
				}
				frobAtStartingBlockLog := test_data.CreateTestLog(int64(frobAtStartingBlock.HeaderID), db)
				_, frobOneErr := db.Exec(`INSERT INTO maker.vat_frob (header_id, log_id, urn_id, dink, dart)
				VALUES ($1, $2, $3, $4, $5)`, frobAtStartingBlock.HeaderID, frobAtStartingBlockLog.ID, urnID,
					frobAtStartingBlock.Dink, frobAtStartingBlock.Dart)
				Expect(frobOneErr).NotTo(HaveOccurred())

				frobAtEarlierBlock = repository.Frob{
					HeaderID: earlierBlockID,
					UrnID:    urnID,
					Dink:     strconv.Itoa(rand.Int()),
					Dart:     strconv.Itoa(rand.Int()),
				}
				frobAtEarlierBlockLog := test_data.CreateTestLog(int64(frobAtStartingBlock.HeaderID), db)
				_, frobTwoErr := db.Exec(`INSERT INTO maker.vat_frob (header_id, log_id, urn_id, dink, dart)
				VALUES ($1, $2, $3, $4, $5)`, frobAtEarlierBlock.HeaderID, frobAtEarlierBlockLog.ID, urnID,
					frobAtEarlierBlock.Dink, frobAtEarlierBlock.Dart)
				Expect(frobTwoErr).NotTo(HaveOccurred())
			})

			It("returns frobs with block >= starting block", func() {
				frobs, err := repo.GetFrobs(startingBlock)

				Expect(err).NotTo(HaveOccurred())
				Expect(frobs).To(ConsistOf(frobAtStartingBlock))
			})

			It("orders results ascending by block_number", func() {
				frobs, err := repo.GetFrobs(earlierBlock)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(frobs)).To(Equal(2))
				Expect(frobs[0]).To(Equal(frobAtEarlierBlock))
			})
		})

		Describe("when there are no frobs", func() {
			It("returns empty list without error", func() {
				frobs, err := repo.GetFrobs(0)

				Expect(len(frobs)).To(BeZero())
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("GetGrabs", func() {
		Describe("when there are grabs", func() {
			var (
				urnID, startingBlock, earlierBlock      int
				grabAtStartingBlock, grabAtEarlierBlock repository.Grab
			)

			BeforeEach(func() {
				fakeIlk := test_data.RandomString(64)
				fakeIlkIdentifier := "ETH-A"
				var ilkID int
				ilkErr := db.Get(&ilkID, `INSERT INTO maker.ilks (ilk, identifier) VALUES ($1, $2) RETURNING id`,
					fakeIlk, fakeIlkIdentifier)
				Expect(ilkErr).NotTo(HaveOccurred())

				fakeUrn := test_data.RandomString(40)
				urnErr := db.Get(&urnID, `INSERT INTO maker.urns (ilk_id, identifier) VALUES ($1, $2) RETURNING id`,
					ilkID, fakeUrn)
				Expect(urnErr).NotTo(HaveOccurred())

				startingBlock = rand.Int()
				var startingBlockID int
				startingBlockErr := db.Get(&startingBlockID, `
				INSERT INTO public.headers (block_number, hash, eth_node_id) VALUES ($1, $2, $3) RETURNING id`,
					startingBlock, test_data.RandomString(64), db.NodeID)
				Expect(startingBlockErr).NotTo(HaveOccurred())

				earlierBlock = startingBlock - 1
				var earlierBlockID int
				earlierBlockErr := db.Get(&earlierBlockID, `
				INSERT INTO public.headers (block_number, hash, eth_node_id) VALUES ($1, $2, $3) RETURNING id`,
					earlierBlock, test_data.RandomString(64), db.NodeID)
				Expect(earlierBlockErr).NotTo(HaveOccurred())

				grabAtStartingBlock = repository.Grab{
					HeaderID: startingBlockID,
					UrnID:    urnID,
					Dink:     strconv.Itoa(rand.Int()),
					Dart:     strconv.Itoa(rand.Int()),
				}
				grabAtStartingBlockLog := test_data.CreateTestLog(int64(grabAtStartingBlock.HeaderID), db)
				_, grabOneErr := db.Exec(`INSERT INTO maker.vat_grab (header_id, log_id, urn_id, dink, dart)
				VALUES ($1, $2, $3, $4, $5)`, grabAtStartingBlock.HeaderID, grabAtStartingBlockLog.ID, urnID,
					grabAtStartingBlock.Dink, grabAtStartingBlock.Dart)
				Expect(grabOneErr).NotTo(HaveOccurred())

				grabAtEarlierBlock = repository.Grab{
					HeaderID: earlierBlockID,
					UrnID:    urnID,
					Dink:     strconv.Itoa(rand.Int()),
					Dart:     strconv.Itoa(rand.Int()),
				}
				grabAtEarlierBlockLog := test_data.CreateTestLog(int64(grabAtStartingBlock.HeaderID), db)
				_, grabTwoErr := db.Exec(`INSERT INTO maker.vat_grab (header_id, log_id, urn_id, dink, dart)
				VALUES ($1, $2, $3, $4, $5)`, grabAtEarlierBlock.HeaderID, grabAtEarlierBlockLog.ID, urnID,
					grabAtEarlierBlock.Dink, grabAtEarlierBlock.Dart)
				Expect(grabTwoErr).NotTo(HaveOccurred())
			})

			It("returns grabs with block >= starting block", func() {
				grabs, err := repo.GetGrabs(startingBlock)

				Expect(err).NotTo(HaveOccurred())
				Expect(grabs).To(ConsistOf(grabAtStartingBlock))
			})

			It("orders results ascending by block_number", func() {
				grabs, err := repo.GetGrabs(earlierBlock)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(grabs)).To(Equal(2))
				Expect(grabs[0]).To(Equal(grabAtEarlierBlock))
			})
		})

		Describe("when there are no grabs", func() {
			It("returns empty list without error", func() {
				grabs, err := repo.GetGrabs(0)

				Expect(len(grabs)).To(BeZero())
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("GetHeaderByID", func() {
		It("returns header with associated ID", func() {
			wantedHeader := core.Header{
				BlockNumber: rand.Int63(),
				Hash:        test_data.RandomString(64),
				Raw:         nil,
				Timestamp:   strconv.Itoa(rand.Int()),
			}
			var wantedHeaderID int
			wantedHeaderErr := db.Get(&wantedHeaderID, `
				INSERT INTO public.headers (block_number, hash, block_timestamp, eth_node_id) VALUES ($1, $2, $3, $4)
				RETURNING id`, wantedHeader.BlockNumber, wantedHeader.Hash, wantedHeader.Timestamp, db.NodeID)
			Expect(wantedHeaderErr).NotTo(HaveOccurred())
			wantedHeader.Id = int64(wantedHeaderID)

			_, anotherHeaderErr := db.Exec(`INSERT INTO public.headers (block_number, hash, block_timestamp,
                            eth_node_id) VALUES ($1, $2, $3, $4) RETURNING id`, rand.Int()-1, test_data.RandomString(64),
				strconv.Itoa(rand.Int()), db.NodeID)
			Expect(anotherHeaderErr).NotTo(HaveOccurred())

			header, err := repo.GetHeaderByID(wantedHeaderID)

			Expect(err).NotTo(HaveOccurred())
			Expect(header).To(Equal(wantedHeader))
		})
	})
})
