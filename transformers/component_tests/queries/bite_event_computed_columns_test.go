package queries

import (
	"database/sql"
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/event"
	"math/rand"
	"strconv"

	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/bite"
	"github.com/vulcanize/mcd_transformers/transformers/events/flip_kick"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vat"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("Bite event computed columns", func() {
	var (
		db               *postgres.DB
		fakeBlock        int
		fakeGuy          string
		fakeHeader       core.Header
		biteGethLog      types.Log
		biteEvent        event.InsertionModel
		biteRepo         bite.Repository
		headerId         int64
		vatRepository    vat.VatStorageRepository
		headerRepository repositories.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)

		fakeGuy = "fakeGuy"
		headerRepository = repositories.NewHeaderRepository(db)
		fakeBlock = rand.Int()
		fakeHeader = fakes.GetFakeHeader(int64(fakeBlock))
		var insertHeaderErr error
		headerId, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())

		biteHeaderSyncLog := test_data.CreateTestLog(headerId, db)
		biteGethLog = biteHeaderSyncLog.Log

		biteRepo = bite.Repository{}
		biteRepo.SetDB(db)
		biteEvent = generateBite(test_helpers.FakeIlk.Hex, fakeGuy, headerId, biteHeaderSyncLog.ID, db)
		insertBiteErr := biteRepo.Create([]event.InsertionModel{biteEvent})
		Expect(insertBiteErr).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("bite_event_ilk", func() {
		It("returns ilk_state for a bite_event", func() {
			ilkValues := test_helpers.GetIlkValues(0)
			test_helpers.CreateIlk(db, fakeHeader, ilkValues, test_helpers.FakeIlkVatMetadatas,
				test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas, test_helpers.FakeIlkSpotMetadatas)

			expectedIlk := test_helpers.IlkStateFromValues(test_helpers.FakeIlk.Hex, fakeHeader.Timestamp, fakeHeader.Timestamp, ilkValues)

			var result test_helpers.IlkState
			err := db.Get(&result, `
				SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated
				FROM api.bite_event_ilk(
					(SELECT (ilk_identifier, urn_identifier, bid_id, ink, art, tab, block_height, log_id)::api.bite_event FROM api.all_bites($1))
				)`, test_helpers.FakeIlk.Identifier)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(expectedIlk))
		})
	})

	Describe("bite_event_urn", func() {
		It("returns urn_state for a bite_event", func() {
			vatRepository.SetDB(db)
			urnSetupData := test_helpers.GetUrnSetupData(fakeBlock, 1)
			urnSetupData.Header.Hash = fakeHeader.Hash
			urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy)
			test_helpers.CreateUrn(urnSetupData, urnMetadata, vatRepository, headerRepository)

			var actualUrn test_helpers.UrnState
			err := db.Get(&actualUrn, `
				SELECT urn_identifier, ilk_identifier FROM api.bite_event_urn(
					(SELECT (ilk_identifier, urn_identifier, bid_id, ink, art, tab, block_height, log_id)::api.bite_event FROM api.all_bites($1)))`,
				test_helpers.FakeIlk.Identifier)
			Expect(err).NotTo(HaveOccurred())

			expectedUrn := test_helpers.UrnState{
				UrnIdentifier: fakeGuy,
				IlkIdentifier: test_helpers.FakeIlk.Identifier,
			}

			test_helpers.AssertUrn(actualUrn, expectedUrn)
		})
	})

	Describe("bite_event_bid", func() {
		It("returns flip_state for a bite_event", func() {
			bidId, convErr := strconv.Atoi(biteEvent.ColumnValues["bid_id"].(string))
			Expect(convErr).NotTo(HaveOccurred())
			address := fakes.FakeAddress
			dealt := false
			flipKickRepo := flip_kick.FlipKickRepository{}
			flipKickRepo.SetDB(db)
			ilkId, urnId, ctxErr := test_helpers.SetUpFlipBidContext(
				test_helpers.FlipBidContextInput{
					DealCreationInput: test_helpers.DealCreationInput{
						Db:              db,
						BidId:           bidId,
						ContractAddress: address.Hex(),
					},
					Dealt:            dealt,
					IlkHex:           test_helpers.FakeIlk.Hex,
					UrnGuy:           test_data.FakeUrn,
					FlipKickRepo:     flipKickRepo,
					FlipKickHeaderId: headerId,
				})
			Expect(ctxErr).NotTo(HaveOccurred())
			flipValues := test_helpers.GetFlipStorageValues(0, test_helpers.FakeIlk.Hex, bidId)
			flipMetadatas := test_helpers.GetFlipMetadatas(strconv.Itoa(bidId))
			test_helpers.CreateFlip(db, fakeHeader, flipValues, flipMetadatas, address.Hex())

			var actualBid test_helpers.FlipBid
			err := db.Get(&actualBid, `
				SELECT bid_id, ilk_id, urn_id, bid, lot, guy, tic, "end", gal, tab, dealt, created, updated FROM api.bite_event_bid(
					(SELECT (ilk_identifier, urn_identifier, bid_id, ink, art, tab, block_height, log_id)::api.bite_event FROM api.all_bites($1)))`,
				test_helpers.FakeIlk.Identifier)
			Expect(err).NotTo(HaveOccurred())

			expectedBid := test_helpers.FlipBidFromValues(strconv.Itoa(bidId), strconv.FormatInt(ilkId, 10), strconv.FormatInt(urnId, 10),
				strconv.FormatBool(dealt), fakeHeader.Timestamp, fakeHeader.Timestamp, flipValues)
			Expect(actualBid).To(Equal(expectedBid))
		})
	})

	Describe("bite_event_tx", func() {
		It("returns transaction for a bite_event", func() {
			expectedTx := Tx{
				TransactionHash:  test_helpers.GetValidNullString("txHash"),
				TransactionIndex: sql.NullInt64{Int64: int64(biteGethLog.TxIndex), Valid: true},
				BlockHeight:      sql.NullInt64{Int64: int64(fakeBlock), Valid: true},
				BlockHash:        test_helpers.GetValidNullString(fakeHeader.Hash),
				TxFrom:           test_helpers.GetValidNullString("fromAddress"),
				TxTo:             test_helpers.GetValidNullString("toAddress"),
			}

			_, err := db.Exec(`INSERT INTO header_sync_transactions (header_id, hash, tx_from, tx_index, tx_to)
		        VALUES ($1, $2, $3, $4, $5)`, headerId, expectedTx.TransactionHash, expectedTx.TxFrom,
				expectedTx.TransactionIndex, expectedTx.TxTo)
			Expect(err).NotTo(HaveOccurred())

			var actualTx Tx
			err = db.Get(&actualTx, `
				SELECT * FROM api.bite_event_tx(
					(SELECT (ilk_identifier, urn_identifier, bid_id, ink, art, tab, block_height, log_id)::api.bite_event FROM api.all_bites($1)))`,
				test_helpers.FakeIlk.Identifier)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualTx).To(Equal(expectedTx))
		})

		It("does not return transaction from same block with different index", func() {
			wrongTx := Tx{
				TransactionHash: test_helpers.GetValidNullString("wrongTxHash"),
				TransactionIndex: sql.NullInt64{
					Int64: int64(biteGethLog.TxIndex) + 1,
					Valid: true,
				},
				BlockHeight: sql.NullInt64{Int64: int64(fakeBlock), Valid: true},
				BlockHash:   test_helpers.GetValidNullString(fakeHeader.Hash),
				TxFrom:      test_helpers.GetValidNullString("wrongFromAddress"),
				TxTo:        test_helpers.GetValidNullString("wrongToAddress"),
			}

			_, insertErr := db.Exec(`INSERT INTO header_sync_transactions (header_id, hash, tx_from, tx_index, tx_to)
				VALUES ($1, $2, $3, $4, $5)`, headerId, wrongTx.TransactionHash, wrongTx.TxFrom,
				wrongTx.TransactionIndex, wrongTx.TxTo)
			Expect(insertErr).NotTo(HaveOccurred())

			var actualTx Tx
			err := db.Get(&actualTx, `
				SELECT * FROM api.bite_event_tx(
					(SELECT (ilk_identifier, urn_identifier, bid_id, ink, art, tab, block_height, log_id)::api.bite_event FROM api.all_bites($1)))`,
				test_helpers.FakeIlk.Identifier)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualTx).To(BeZero())
		})
	})
})
