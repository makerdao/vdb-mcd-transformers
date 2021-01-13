package queries

import (
	"database/sql"
	"math/rand"
	"strconv"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/vat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bite event computed columns", func() {
	var (
		blockOne, timestampOne   int
		fakeGuy, fakeFlipAddress string
		headerOne                core.Header
		biteGethLog              types.Log
		biteEvent                event.InsertionModel
		vatRepository            vat.StorageRepository
		headerRepository         datastore.HeaderRepository
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)

		fakeGuy = fakes.RandomString(42)
		headerRepository = repositories.NewHeaderRepository(db)
		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepository)

		biteEventLog := test_data.CreateTestLog(headerOne.Id, db)
		biteGethLog = biteEventLog.Log

		fakeFlipAddress = fakes.FakeAddress.Hex()
		biteEvent = generateBite(test_helpers.FakeIlk.Hex, fakeGuy, fakeFlipAddress, headerOne.Id, biteEventLog.ID, db)
		insertBiteErr := event.PersistModels([]event.InsertionModel{biteEvent}, db)
		Expect(insertBiteErr).NotTo(HaveOccurred())
	})

	Describe("bite_event_ilk", func() {
		It("returns ilk_snapshot for a bite_event", func() {
			ilkValues := test_helpers.GetIlkValues(0)
			test_helpers.CreateIlk(db, headerOne, ilkValues, test_helpers.FakeIlkVatMetadatas, test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas, test_helpers.FakeIlkSpotMetadatas)

			expectedIlk := test_helpers.IlkSnapshotFromValues(test_helpers.FakeIlk.Hex, headerOne.Timestamp, headerOne.Timestamp, ilkValues)

			var result test_helpers.IlkSnapshot
			err := db.Get(&result, `
				SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, dunk, created, updated
				FROM api.bite_event_ilk(
					(SELECT (ilk_identifier, urn_identifier, bid_id, ink, art, tab, block_height, log_id, flip_address)::api.bite_event FROM api.all_bites($1))
				)`, test_helpers.FakeIlk.Identifier)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(expectedIlk))
		})
	})

	Describe("bite_event_urn", func() {
		It("returns urn_snapshot for a bite_event", func() {
			vatRepository.SetDB(db)
			urnSetupData := test_helpers.GetUrnSetupData()
			urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy)
			test_helpers.CreateUrn(db, urnSetupData, headerOne, urnMetadata, vatRepository)

			var actualUrn test_helpers.UrnState
			err := db.Get(&actualUrn, `
				SELECT urn_identifier, ilk_identifier FROM api.bite_event_urn(
					(SELECT (ilk_identifier, urn_identifier, bid_id, ink, art, tab, block_height, log_id, flip_address)::api.bite_event FROM api.all_bites($1)))`,
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
			dealt := false
			ilkId, urnId, ctxErr := test_helpers.SetUpFlipBidContext(
				test_helpers.FlipBidContextInput{
					DealCreationInput: test_helpers.DealCreationInput{
						DB:              db,
						BidId:           bidId,
						ContractAddress: fakeFlipAddress,
					},
					Dealt:            dealt,
					IlkHex:           test_helpers.FakeIlk.Hex,
					UrnGuy:           test_data.FakeUrn,
					FlipKickHeaderId: headerOne.Id,
				})
			Expect(ctxErr).NotTo(HaveOccurred())
			flipValues := test_helpers.GetFlipStorageValues(0, test_helpers.FakeIlk.Hex, bidId)
			flipMetadatas := test_helpers.GetFlipMetadatas(strconv.Itoa(bidId))
			test_helpers.CreateFlip(db, headerOne, flipValues, flipMetadatas, fakeFlipAddress)

			var actualBid test_helpers.FlipBid
			err := db.Get(&actualBid, `
				SELECT bid_id, ilk_id, urn_id, bid, lot, guy, tic, "end", gal, tab, dealt, flip_address, created, updated FROM api.bite_event_bid(
					(SELECT (ilk_identifier, urn_identifier, bid_id, ink, art, tab, block_height, log_id, flip_address)::api.bite_event FROM api.all_bites($1)))`,
				test_helpers.FakeIlk.Identifier)
			Expect(err).NotTo(HaveOccurred())

			expectedBid := test_helpers.FlipBidFromValues(strconv.Itoa(bidId), strconv.FormatInt(ilkId, 10), strconv.FormatInt(urnId, 10),
				strconv.FormatBool(dealt), headerOne.Timestamp, headerOne.Timestamp, flipValues)
			Expect(actualBid).To(Equal(expectedBid))
		})
	})

	Describe("bite_event_tx", func() {
		It("returns transaction for a bite_event", func() {
			expectedTx := Tx{
				TransactionHash:  test_helpers.GetValidNullString("txHash"),
				TransactionIndex: sql.NullInt64{Int64: int64(biteGethLog.TxIndex), Valid: true},
				BlockHeight:      sql.NullInt64{Int64: int64(blockOne), Valid: true},
				BlockHash:        test_helpers.GetValidNullString(headerOne.Hash),
				TxFrom:           test_helpers.GetValidNullString("fromAddress"),
				TxTo:             test_helpers.GetValidNullString("toAddress"),
			}

			_, err := db.Exec(`INSERT INTO public.transactions (header_id, hash, tx_from, tx_index, tx_to)
		        VALUES ($1, $2, $3, $4, $5)`, headerOne.Id, expectedTx.TransactionHash, expectedTx.TxFrom,
				expectedTx.TransactionIndex, expectedTx.TxTo)
			Expect(err).NotTo(HaveOccurred())

			var actualTx Tx
			err = db.Get(&actualTx, `
				SELECT * FROM api.bite_event_tx(
					(SELECT (ilk_identifier, urn_identifier, bid_id, ink, art, tab, block_height, log_id, flip_address)::api.bite_event FROM api.all_bites($1)))`,
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
				BlockHeight: sql.NullInt64{Int64: int64(blockOne), Valid: true},
				BlockHash:   test_helpers.GetValidNullString(headerOne.Hash),
				TxFrom:      test_helpers.GetValidNullString("wrongFromAddress"),
				TxTo:        test_helpers.GetValidNullString("wrongToAddress"),
			}

			_, insertErr := db.Exec(`INSERT INTO public.transactions (header_id, hash, tx_from, tx_index, tx_to)
				VALUES ($1, $2, $3, $4, $5)`, headerOne.Id, wrongTx.TransactionHash, wrongTx.TxFrom,
				wrongTx.TransactionIndex, wrongTx.TxTo)
			Expect(insertErr).NotTo(HaveOccurred())

			var actualTx Tx
			err := db.Get(&actualTx, `
				SELECT * FROM api.bite_event_tx(
					(SELECT (ilk_identifier, urn_identifier, bid_id, ink, art, tab, block_height, log_id, flip_address)::api.bite_event FROM api.all_bites($1)))`,
				test_helpers.FakeIlk.Identifier)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualTx).To(BeZero())
		})

		It("does not return transaction from different block with same index", func() {
			headerZero := createHeader(blockOne-1, timestampOne-1, headerRepository)
			wrongTx := Tx{
				TransactionHash: test_helpers.GetValidNullString("wrongTxHash"),
				TransactionIndex: sql.NullInt64{
					Int64: int64(biteGethLog.TxIndex),
					Valid: true,
				},
				BlockHeight: sql.NullInt64{Int64: headerZero.BlockNumber, Valid: true},
				BlockHash:   test_helpers.GetValidNullString(headerOne.Hash),
				TxFrom:      test_helpers.GetValidNullString("wrongFromAddress"),
				TxTo:        test_helpers.GetValidNullString("wrongToAddress"),
			}

			_, insertErr := db.Exec(`INSERT INTO public.transactions (header_id, hash, tx_from, tx_index, tx_to)
				VALUES ($1, $2, $3, $4, $5)`, headerZero.Id, wrongTx.TransactionHash, wrongTx.TxFrom,
				wrongTx.TransactionIndex, wrongTx.TxTo)
			Expect(insertErr).NotTo(HaveOccurred())

			var actualTx Tx
			err := db.Get(&actualTx, `
				SELECT * FROM api.bite_event_tx(
					(SELECT (ilk_identifier, urn_identifier, bid_id, ink, art, tab, block_height, log_id, flip_address)::api.bite_event FROM api.all_bites($1)))`,
				test_helpers.FakeIlk.Identifier)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualTx).To(BeZero())
		})
	})
})
