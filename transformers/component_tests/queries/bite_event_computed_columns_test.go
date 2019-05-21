package queries

import (
	"database/sql"
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/bite"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("Bite event computed columns", func() {
	var (
		db         *postgres.DB
		fakeBlock  int
		fakeGuy    string
		fakeHeader core.Header
		biteEvent  bite.BiteModel
		biteRepo   bite.BiteRepository
		headerId   int64
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

		biteRepo = bite.BiteRepository{}
		biteRepo.SetDB(db)
		biteEvent = test_data.BiteModel
		biteEvent.Ilk = test_helpers.FakeIlk.Hex
		biteEvent.Urn = fakeGuy
		insertBiteErr := biteRepo.Create(headerId, []interface{}{biteEvent})
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
				test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas)

			expectedIlk := test_helpers.IlkStateFromValues(test_helpers.FakeIlk.Hex, fakeHeader.Timestamp, fakeHeader.Timestamp, ilkValues)

			var result test_helpers.IlkState
			err := db.Get(&result, `
				SELECT ilk_name, rate, art, spot, line, dust, chop, lump, flip, rho, duty, created, updated
				FROM api.bite_event_ilk(
					(SELECT (ilk_name, urn_guy, ink, art, tab, block_height, tx_idx)::api.bite_event FROM api.all_bites($1))
				)`, test_helpers.FakeIlk.Name)

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
				SELECT urn_guy, ilk_name FROM api.bite_event_urn(
					(SELECT (ilk_name, urn_guy, ink, art, tab, block_height, tx_idx)::api.bite_event FROM api.all_bites($1)))`,
				test_helpers.FakeIlk.Name)
			Expect(err).NotTo(HaveOccurred())

			expectedUrn := test_helpers.UrnState{
				UrnGuy:  fakeGuy,
				IlkName: test_helpers.FakeIlk.Name,
			}

			test_helpers.AssertUrn(actualUrn, expectedUrn)
		})
	})

	Describe("bite_event_tx", func() {
		It("returns transaction for a bite_event", func() {
			expectedTx := Tx{
				TransactionHash:  test_helpers.GetValidNullString("txHash"),
				TransactionIndex: sql.NullInt64{Int64: int64(biteEvent.TransactionIndex), Valid: true},
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
					(SELECT (ilk_name, urn_guy, ink, art, tab, block_height, tx_idx)::api.bite_event FROM api.all_bites($1)))`,
				test_helpers.FakeIlk.Name)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualTx).To(Equal(expectedTx))
		})

		It("does not return transaction from same block with different index", func() {
			wrongTx := Tx{
				TransactionHash: test_helpers.GetValidNullString("wrongTxHash"),
				TransactionIndex: sql.NullInt64{
					Int64: int64(biteEvent.TransactionIndex) + 1,
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
					(SELECT (ilk_name, urn_guy, ink, art, tab, block_height, tx_idx)::api.bite_event FROM api.all_bites($1)))`,
				test_helpers.FakeIlk.Name)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualTx).To(BeZero())
		})
	})
})
