package queries

import (
	"database/sql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/bite"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
	"strconv"
)

var _ = Describe("Bites query", func() {
	var (
		db         *postgres.DB
		biteRepo   bite.BiteRepository
		headerRepo repositories.HeaderRepository
		fakeUrn    = test_data.RandomString(5)
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		biteRepo = bite.BiteRepository{}
		biteRepo.SetDB(db)
	})

	Describe("all_bites", func() {
		It("returns bites for an ilk", func() {
			headerOne := fakes.GetFakeHeader(1)
			headerOneId, err := headerRepo.CreateOrUpdateHeader(headerOne)
			Expect(err).NotTo(HaveOccurred())

			biteOne := test_data.BiteModel
			biteOne.Ilk = test_helpers.FakeIlk.Hex
			biteOne.Urn = fakeUrn
			biteOne.Ink = strconv.Itoa(rand.Int())
			biteOne.Art = strconv.Itoa(rand.Int())
			biteOne.Tab = strconv.Itoa(rand.Int())

			err = biteRepo.Create(headerOneId, []interface{}{biteOne})
			Expect(err).NotTo(HaveOccurred())

			var actualBites []test_helpers.BiteEvent
			err = db.Select(&actualBites, `SELECT ilk_name, urn_id, ink, art, tab FROM api.all_bites($1)`, test_helpers.FakeIlk.Name)
			Expect(err).NotTo(HaveOccurred())

			Expect(actualBites).To(ConsistOf(
				test_helpers.BiteEvent{IlkName: test_helpers.FakeIlk.Name, UrnId: fakeUrn, Ink: biteOne.Ink, Art: biteOne.Art, Tab: biteOne.Tab},
			))
		})

		It("returns bites from multiple blocks", func() {
			headerOne := fakes.GetFakeHeader(1)
			headerOneId, err := headerRepo.CreateOrUpdateHeader(headerOne)
			Expect(err).NotTo(HaveOccurred())

			biteBlockOne := test_data.BiteModel
			biteBlockOne.Ilk = test_helpers.FakeIlk.Hex
			biteBlockOne.Urn = fakeUrn
			biteBlockOne.Ink = strconv.Itoa(rand.Int())
			biteBlockOne.Art = strconv.Itoa(rand.Int())
			biteBlockOne.Tab = strconv.Itoa(rand.Int())

			err = biteRepo.Create(headerOneId, []interface{}{biteBlockOne})
			Expect(err).NotTo(HaveOccurred())

			// New block
			headerTwo := fakes.GetFakeHeader(2)
			headerTwo.Hash = "anotherHash"
			headerTwoId, err := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(err).NotTo(HaveOccurred())

			biteBlockTwo := test_data.BiteModel
			biteBlockTwo.Ilk = test_helpers.FakeIlk.Hex
			biteBlockTwo.Urn = fakeUrn
			biteBlockTwo.Ink = strconv.Itoa(rand.Int())
			biteBlockTwo.Art = strconv.Itoa(rand.Int())
			biteBlockTwo.Tab = strconv.Itoa(rand.Int())

			err = biteRepo.Create(headerTwoId, []interface{}{biteBlockTwo})
			Expect(err).NotTo(HaveOccurred())

			var actualBites []test_helpers.BiteEvent
			err = db.Select(&actualBites, `SELECT ilk_name, urn_id, ink, art, tab FROM api.all_bites($1)`, test_helpers.FakeIlk.Name)
			Expect(err).NotTo(HaveOccurred())

			Expect(actualBites).To(ConsistOf(
				test_helpers.BiteEvent{IlkName: test_helpers.FakeIlk.Name, UrnId: fakeUrn, Ink: biteBlockTwo.Ink, Art: biteBlockTwo.Art, Tab: biteBlockTwo.Tab},
				test_helpers.BiteEvent{IlkName: test_helpers.FakeIlk.Name, UrnId: fakeUrn, Ink: biteBlockOne.Ink, Art: biteBlockOne.Art, Tab: biteBlockOne.Tab},
			))
		})

		It("ignores bites from irrelevant ilks", func() {
			headerOne := fakes.GetFakeHeader(1)
			headerOneId, err := headerRepo.CreateOrUpdateHeader(headerOne)
			Expect(err).NotTo(HaveOccurred())

			bite := test_data.BiteModel
			bite.Ilk = test_helpers.FakeIlk.Hex
			bite.Urn = fakeUrn
			bite.Ink = strconv.Itoa(rand.Int())
			bite.Art = strconv.Itoa(rand.Int())
			bite.Tab = strconv.Itoa(rand.Int())

			irrelevantBite := test_data.BiteModel
			irrelevantBite.Ilk = test_helpers.AnotherFakeIlk.Hex
			irrelevantBite.Urn = fakeUrn
			irrelevantBite.Ink = strconv.Itoa(rand.Int())
			irrelevantBite.Art = strconv.Itoa(rand.Int())
			irrelevantBite.Tab = strconv.Itoa(rand.Int())
			irrelevantBite.TransactionIndex = bite.TransactionIndex + 1

			err = biteRepo.Create(headerOneId, []interface{}{bite, irrelevantBite})
			Expect(err).NotTo(HaveOccurred())

			var actualBites []test_helpers.BiteEvent
			err = db.Select(&actualBites, `SELECT ilk_name, urn_id, ink, art, tab FROM api.all_bites($1)`, test_helpers.FakeIlk.Name)
			Expect(err).NotTo(HaveOccurred())

			Expect(actualBites).To(ConsistOf(
				test_helpers.BiteEvent{IlkName: test_helpers.FakeIlk.Name, UrnId: fakeUrn, Ink: bite.Ink, Art: bite.Art, Tab: bite.Tab},
			))
		})
	})

	Describe("Extending bite_event", func() {
		var (
			fakeBlock  int
			fakeHeader core.Header
			headerId   int64
			err        error
		)

		BeforeEach(func() {
			fakeBlock = rand.Int()
			fakeHeader = fakes.GetFakeHeader(int64(fakeBlock))
			headerId, err = headerRepo.CreateOrUpdateHeader(fakeHeader)
			Expect(err).NotTo(HaveOccurred())
		})

		Describe("bite_event_ilk", func() {
			It("returns ilk_state for a bite_event", func() {
				vatRepository.SetDB(db)
				catRepository.SetDB(db)
				jugRepository.SetDB(db)
				ilkValues := test_helpers.GetIlkValues(0)
				createIlkAtBlock(fakeHeader, ilkValues, test_helpers.FakeIlkVatMetadatas,
					test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas)

				biteEvent := test_data.BiteModel
				biteEvent.Ilk = test_helpers.FakeIlk.Hex
				err = biteRepo.Create(headerId, []interface{}{biteEvent})
				Expect(err).NotTo(HaveOccurred())

				expectedIlk := test_helpers.IlkStateFromValues(test_helpers.FakeIlk.Hex, fakeHeader.Timestamp, fakeHeader.Timestamp, ilkValues)

				var result test_helpers.IlkState
				err = db.Get(&result,
					`SELECT ilk_name, rate, art, spot, line, dust, chop, lump, flip, rho, duty, created, updated
						FROM api.bite_event_ilk(
							(SELECT (ilk_name, urn_id, ink, art, tab, block_height, tx_idx)::api.bite_event FROM api.all_bites($1))
						)`, test_helpers.FakeIlk.Name)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal(expectedIlk))
			})
		})

		Describe("bite_event_urn", func() {
			It("returns urn_state for a bite_event", func() {
				vatRepository.SetDB(db)
				fakeGuy := "fakeGuy"
				urnSetupData := test_helpers.GetUrnSetupData(fakeBlock, 1)
				urnSetupData.Header.Hash = fakeHeader.Hash
				urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy)
				test_helpers.CreateUrn(urnSetupData, urnMetadata, vatRepository, headerRepo)

				biteEvent := test_data.BiteModel
				biteEvent.Urn = fakeGuy
				biteEvent.Ilk = test_helpers.FakeIlk.Hex
				err = biteRepo.Create(headerId, []interface{}{biteEvent})
				Expect(err).NotTo(HaveOccurred())

				var actualUrn test_helpers.UrnState
				err = db.Get(&actualUrn,
					`SELECT urn_id, ilk_name FROM api.bite_event_urn(
							(SELECT (ilk_name, urn_id, ink, art, tab, block_height, tx_idx)::api.bite_event FROM api.all_bites($1)))`,
					test_helpers.FakeIlk.Name)
				Expect(err).NotTo(HaveOccurred())

				expectedUrn := test_helpers.UrnState{
					UrnId:   fakeGuy,
					IlkName: test_helpers.FakeIlk.Name,
				}

				test_helpers.AssertUrn(actualUrn, expectedUrn)
			})
		})

		Describe("bite_event_tx", func() {
			It("returns transaction for a bite_event", func() {
				biteEvent := test_data.BiteModel
				biteEvent.Ilk = test_helpers.FakeIlk.Hex
				err = biteRepo.Create(headerId, []interface{}{biteEvent})
				Expect(err).NotTo(HaveOccurred())

				expectedTx := Tx{
					TransactionHash:  sql.NullString{String: "txHash", Valid: true},
					TransactionIndex: sql.NullInt64{Int64: int64(biteEvent.TransactionIndex), Valid: true},
					BlockHeight:      sql.NullInt64{Int64: int64(fakeBlock), Valid: true},
					BlockHash:        sql.NullString{String: fakeHeader.Hash, Valid: true},
					TxFrom:           sql.NullString{String: "fromAddress", Valid: true},
					TxTo:             sql.NullString{String: "toAddress", Valid: true},
				}

				_, err = db.Exec(`INSERT INTO header_sync_transactions (header_id, hash, tx_from, tx_index, tx_to)
		        VALUES ($1, $2, $3, $4, $5)`, headerId, expectedTx.TransactionHash, expectedTx.TxFrom,
					expectedTx.TransactionIndex, expectedTx.TxTo)
				Expect(err).NotTo(HaveOccurred())

				var actualTx Tx
				err = db.Get(&actualTx, `SELECT * FROM api.bite_event_tx(
			    (SELECT (ilk_name, urn_id, ink, art, tab, block_height, tx_idx)::api.bite_event FROM api.all_bites($1)))`,
					test_helpers.FakeIlk.Name)

				Expect(err).NotTo(HaveOccurred())
				Expect(actualTx).To(Equal(expectedTx))
			})

			It("does not return transaction from same block with different index", func() {
				biteEvent := test_data.BiteModel
				biteEvent.Ilk = test_helpers.FakeIlk.Hex
				err = biteRepo.Create(headerId, []interface{}{biteEvent})
				Expect(err).NotTo(HaveOccurred())

				wrongTx := Tx{
					TransactionHash: sql.NullString{String: "wrongTxHash", Valid: true},
					TransactionIndex: sql.NullInt64{
						Int64: int64(biteEvent.TransactionIndex) + 1,
						Valid: true,
					},
					BlockHeight: sql.NullInt64{Int64: int64(fakeBlock), Valid: true},
					BlockHash:   sql.NullString{String: fakeHeader.Hash, Valid: true},
					TxFrom:      sql.NullString{String: "wrongFromAddress", Valid: true},
					TxTo:        sql.NullString{String: "wrongToAddress", Valid: true},
				}

				_, insertErr := db.Exec(`INSERT INTO header_sync_transactions (header_id, hash, tx_from, tx_index, tx_to)
				VALUES ($1, $2, $3, $4, $5)`, headerId, wrongTx.TransactionHash, wrongTx.TxFrom,
					wrongTx.TransactionIndex, wrongTx.TxTo)
				Expect(insertErr).NotTo(HaveOccurred())

				var actualTx Tx
				err := db.Get(&actualTx, `SELECT * FROM api.bite_event_tx(
			    (SELECT (ilk_name, urn_id, ink, art, tab, block_height, tx_idx)::api.bite_event FROM api.all_bites($1)))`,
					test_helpers.FakeIlk.Name)

				Expect(err).NotTo(HaveOccurred())
				Expect(actualTx).To(BeZero())
			})
		})
	})
})
