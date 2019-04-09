package queries

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_frob"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
	"strconv"
)

var _ = Describe("Extension function", func() {
	var (
		db         *postgres.DB
		fakeBlock  int
		fakeHeader core.Header
		headerId   int64
		err        error
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		fakeBlock = rand.Int()
		fakeHeader = fakes.GetFakeHeader(int64(fakeBlock))
		headerId, err = headerRepository.CreateOrUpdateHeader(fakeHeader)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("using ilk and frob", func() {
		var (
			frobRepo  vat_frob.VatFrobRepository
			frobEvent vat_frob.VatFrobModel
			ilkValues map[string]string
		)

		BeforeEach(func() {
			frobRepo = vat_frob.VatFrobRepository{}
			frobRepo.SetDB(db)
			vatRepository.SetDB(db)
			catRepository.SetDB(db)
			jugRepository.SetDB(db)

			// Create an ilk
			ilkValues = test_helpers.GetIlkValues(0)
			createIlkAtBlock(fakeHeader, ilkValues, test_helpers.FakeIlkVatMetadatas,
				test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas)

			// Create a frob event
			frobEvent = test_data.VatFrobModel
			frobEvent.Ilk = test_helpers.FakeIlk
			err = frobRepo.Create(headerId, []interface{}{frobEvent})
			Expect(err).NotTo(HaveOccurred())
		})

		Describe("ilk_state_frobs", func() {
			It("returns relevant frobs for an ilk_state", func() {
				irrelevantFrob := test_data.VatFrobModel
				irrelevantFrob.Ilk = "irrelevantIlk"
				irrelevantFrob.Urn = "anotherGuy"
				irrelevantFrob.TransactionIndex = frobEvent.TransactionIndex + 1
				err = frobRepo.Create(headerId, []interface{}{irrelevantFrob})
				Expect(err).NotTo(HaveOccurred())

				var ilkId int
				err = db.Get(&ilkId, `SELECT id FROM maker.ilks WHERE ilk = $1`, test_helpers.FakeIlk)
				Expect(err).NotTo(HaveOccurred())

				var actualFrobs []test_helpers.FrobEvent
				err = db.Select(&actualFrobs,
					`SELECT ilkid, urnid, dink, dart FROM maker.ilk_state_frobs(
                        (SELECT (ilk_id, ilk, block_number, rate, art, spot, line, dust, chop, lump, flip, rho, tax, created, updated)::maker.ilk_state
                         FROM maker.get_ilk_at_block_number($1, $2))
                    )`, fakeBlock, ilkId)
				Expect(err).NotTo(HaveOccurred())

				expectedFrobs := []test_helpers.FrobEvent{{
					IlkId: test_helpers.FakeIlk,
					UrnId: frobEvent.Urn,
					Dink:  frobEvent.Dink,
					Dart:  frobEvent.Dart,
				}}

				Expect(actualFrobs).To(Equal(expectedFrobs))
			})
		})

		Describe("frob_event_ilk", func() {
			It("returns ilk_state for a frob_event", func() {
				expectedIlk := test_helpers.IlkStateFromValues(test_helpers.FakeIlk, fakeHeader.Timestamp, fakeHeader.Timestamp, ilkValues)

				var result test_helpers.IlkState
				err = db.Get(&result,
					`SELECT ilk, rate, art, spot, line, dust, chop, lump, flip, rho, tax, created, updated
                    FROM maker.frob_event_ilk(
                        (SELECT (ilkid, urnid, dink, dart, block_number)::maker.frob_event FROM maker.all_frobs($1))
                    )`, test_helpers.FakeIlk)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal(expectedIlk))
			})
		})
	})

	Describe("using urn and frob", func() {
		var (
			frobRepo  vat_frob.VatFrobRepository
			frobEvent vat_frob.VatFrobModel
			fakeGuy   = "fakeAddress"
		)

		BeforeEach(func() {
			frobRepo = vat_frob.VatFrobRepository{}
			frobRepo.SetDB(db)
			vatRepository.SetDB(db)

			urnSetupData := test_helpers.GetUrnSetupData(fakeBlock, 1)
			urnSetupData.Header.Hash = fakeHeader.Hash
			urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk, fakeGuy)
			test_helpers.CreateUrn(urnSetupData, urnMetadata, vatRepository, headerRepository)

			frobEvent = test_data.VatFrobModel
			frobEvent.Urn = fakeGuy
			frobEvent.Ilk = test_helpers.FakeIlk
			err = frobRepo.Create(headerId, []interface{}{frobEvent})
			Expect(err).NotTo(HaveOccurred())
		})

		Describe("frob_event_urn", func() {
			It("returns urn_state for a frob_event", func() {
				var actualUrn test_helpers.UrnState
				err = db.Get(&actualUrn,
					`SELECT urnId, ilkId FROM maker.frob_event_urn(
                        (SELECT (ilkid, urnid, dink, dart, block_number)::maker.frob_event FROM maker.all_frobs($1)))`,
					test_helpers.FakeIlk)
				Expect(err).NotTo(HaveOccurred())

				expectedUrn := test_helpers.UrnState{
					UrnId: fakeGuy,
					IlkId: test_helpers.FakeIlk,
				}

				test_helpers.AssertUrn(actualUrn, expectedUrn)
			})
		})

		Describe("urn_state_frobs", func() {
			It("returns relevant frobs for an urn_state", func() {
				var actualFrobs test_helpers.FrobEvent
				err = db.Get(&actualFrobs,
					`SELECT ilkid, urnid, dink, dart FROM maker.urn_state_frobs(
                        (SELECT (urnid, ilkid, blockheight, ink, art, ratio, safe, created, updated)::maker.urn_state
                         FROM maker.get_all_urn_states_at_block($1))
                    )`, fakeBlock)
				Expect(err).NotTo(HaveOccurred())

				expectedFrobs := test_helpers.FrobEvent{
					IlkId: test_helpers.FakeIlk,
					UrnId: fakeGuy,
					Dink:  frobEvent.Dink,
					Dart:  frobEvent.Dart,
				}

				Expect(actualFrobs).To(Equal(expectedFrobs))
			})
		})
	})

	Describe("tx_era", func() {
		It("returns an era object for a transaction", func() {
			txFrom := "fromAddress"
			txTo := "toAddress"
			txIndex := rand.Intn(10)
			_, err = db.Exec(`INSERT INTO light_sync_transactions (header_id, hash, tx_from, tx_index, tx_to)
                VALUES ($1, $2, $3, $4, $5)`, headerId, fakeHeader.Hash, txFrom, txIndex, txTo)
			Expect(err).NotTo(HaveOccurred())

			var actualEra Era
			err = db.Get(&actualEra, `SELECT * FROM maker.tx_era(
                    (SELECT (txs.hash, txs.tx_index, h.block_number, h.hash, txs.tx_from, txs.tx_to)::maker.tx
			        FROM light_sync_transactions txs
			        LEFT JOIN headers h ON h.id = txs.header_id)
			    )`)
			Expect(err).NotTo(HaveOccurred())

			expectedEra := Era{
				Epoch: fakeHeader.Timestamp,
				Iso:   "1973-07-10T00:11:51Z", // Z for Zulu, meaning UTC
			}
			Expect(actualEra).To(Equal(expectedEra))
		})
	})

	Describe("frob_event_tx", func() {
		It("returns transaction for a frob_event", func() {
			frobRepo := vat_frob.VatFrobRepository{}
			frobRepo.SetDB(db)
			frobEvent := test_data.VatFrobModel
			frobEvent.Ilk = test_helpers.FakeIlk
			err = frobRepo.Create(headerId, []interface{}{frobEvent})
			Expect(err).NotTo(HaveOccurred())

			expectedTx := Tx{
				TransactionHash:  "txHash",
				TransactionIndex: strconv.Itoa(rand.Intn(10)),
				BlockNumber:      strconv.Itoa(fakeBlock),
				BlockHash:        fakeHeader.Hash,
				TxFrom:           "fromAddress",
				TxTo:             "toAddress",
			}

			_, err = db.Exec(`INSERT INTO light_sync_transactions (header_id, hash, tx_from, tx_index, tx_to)
                VALUES ($1, $2, $3, $4, $5)`, headerId, expectedTx.TransactionHash, expectedTx.TxFrom,
				expectedTx.TransactionIndex, expectedTx.TxTo)
			Expect(err).NotTo(HaveOccurred())

			var actualTx Tx
			err = db.Get(&actualTx, `SELECT * FROM maker.frob_event_tx(
			    (SELECT (ilkid, urnid, dink, dart, block_number)::maker.frob_event FROM maker.all_frobs($1)))`,
				test_helpers.FakeIlk)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualTx).To(Equal(expectedTx))
		})
	})
})

type Era struct {
	Epoch string
	Iso   string
}

type Tx struct {
	TransactionHash  string `db:"transaction_hash"`
	TransactionIndex string `db:"transaction_index"`
	BlockNumber      string `db:"block_number"`
	BlockHash        string `db:"block_hash"`
	TxFrom           string `db:"tx_from"`
	TxTo             string `db:"tx_to"`
}
