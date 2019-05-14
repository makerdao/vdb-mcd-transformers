package queries

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/pip_log_value"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"strconv"
)

var _ = Describe("Pip logValue query", func() {
	var (
		db                    *postgres.DB
		pipLogValueRepository pip_log_value.PipLogValueRepository
		headerRepository      repositories.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		pipLogValueRepository = pip_log_value.PipLogValueRepository{}
		pipLogValueRepository.SetDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
	})

	It("returns 2 pip log values in different blocks between a time range", func() {
		var (
			anotherBlockNumber int64 = 10606965
			beginningTimeRange int64 = 111111111
			endingTimeRange    int64 = 111111112
			logValue                 = "123456789"
			transactionIdx           = 3
		)

		fakeHeaderOne := fakes.GetFakeHeaderWithTimestamp(beginningTimeRange, int64(test_data.PipLogValueModel.BlockNumber))
		headerID, err := headerRepository.CreateOrUpdateHeader(fakeHeaderOne)
		Expect(err).NotTo(HaveOccurred())

		err = pipLogValueRepository.Create(headerID, []interface{}{test_data.PipLogValueModel})
		Expect(err).NotTo(HaveOccurred())

		fakeHeaderTwo := fakes.GetFakeHeaderWithTimestamp(endingTimeRange, anotherBlockNumber)
		anotherHeaderID, err := headerRepository.CreateOrUpdateHeader(fakeHeaderTwo)
		Expect(err).NotTo(HaveOccurred())

		anotherPipLogValue := test_data.GetFakePipLogValue(anotherBlockNumber, transactionIdx, logValue)
		err = pipLogValueRepository.Create(anotherHeaderID, []interface{}{anotherPipLogValue})
		Expect(err).NotTo(HaveOccurred())

		expectedValues := []test_helpers.LogValue{
			{
				Val:         test_data.PipLogValueModel.Value,
				BlockNumber: test_data.PipLogValueModel.BlockNumber,
				TxIdx:       test_data.PipLogValueModel.TransactionIndex,
			},
			{
				Val:         anotherPipLogValue.Value,
				BlockNumber: anotherPipLogValue.BlockNumber,
				TxIdx:       anotherPipLogValue.TransactionIndex,
			},
		}

		var dbPipLogValue []test_helpers.LogValue
		err = db.Select(&dbPipLogValue, `SELECT val, maker.pip_log_value.block_number, tx_idx FROM maker.pip_log_value 
                                            JOIN public.headers ON public.headers.id = maker.pip_log_value.header_id
                                            WHERE public.headers.block_timestamp BETWEEN $1 AND $2`, beginningTimeRange, endingTimeRange)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbPipLogValue).To(Equal(expectedValues))
	})

	It("returns a transaction from a logValue", func() {
		var (
			anotherBlockNumber int64 = 10606965
			endingTimeRange    int64 = 111111112
			logValue                 = "123456789"
			transactionIdx           = 3
		)

		fakeHeaderTwo := fakes.GetFakeHeaderWithTimestamp(endingTimeRange, anotherBlockNumber)
		anotherHeaderID, err := headerRepository.CreateOrUpdateHeader(fakeHeaderTwo)
		Expect(err).NotTo(HaveOccurred())

		anotherPipLogValue := test_data.GetFakePipLogValue(anotherBlockNumber, transactionIdx, logValue)
		err = pipLogValueRepository.Create(anotherHeaderID, []interface{}{anotherPipLogValue})
		Expect(err).NotTo(HaveOccurred())

		expectedTx := LogValueTx{
			TransactionHash:  fakeHeaderTwo.Hash,
			TransactionIndex: strconv.Itoa(transactionIdx),
			BlockHeight:      "10606965",
			BlockHash:        "",
			TxFrom:           "fromAddress",
			TxTo:             "toAddress",
		}

		_, err = db.Exec(`INSERT INTO header_sync_transactions (header_id, hash, tx_from, tx_index, tx_to)
                VALUES ($1, $2, $3, $4, $5)`, anotherHeaderID, expectedTx.TransactionHash, expectedTx.TxFrom,
			expectedTx.TransactionIndex, expectedTx.TxTo)
		Expect(err).NotTo(HaveOccurred())

		var actualTx []LogValueTx
		err = db.Select(&actualTx, `SELECT txs.hash, txs.tx_index, headers.block_number, headers.hash, txs.tx_from, txs.tx_to
                                        FROM maker.pip_log_value plv 
                                        LEFT JOIN public.header_sync_transactions txs ON plv.header_id = txs.header_id
                                        LEFT JOIN headers ON plv.header_id = headers.id
                                        WHERE headers.block_number = $1 AND txs.tx_index = $2
                                        ORDER BY headers.block_number DESC`, expectedTx.BlockHeight, expectedTx.TransactionIndex)
		Expect(err).NotTo(HaveOccurred())
		Expect(actualTx[0]).To(Equal(expectedTx))
	})

	It("returns 2 pip log values with transactions in the same block", func() {
		var (
			beginningTimeRange int64 = 111111111
			endingTimeRange    int64 = 111111112
			anotherBlockNumber int64 = 10606964
			logValue                 = "123456789"
			transactionIdx           = 3
		)

		fakeHeaderOne := fakes.GetFakeHeaderWithTimestamp(beginningTimeRange, anotherBlockNumber)
		headerID, err := headerRepository.CreateOrUpdateHeader(fakeHeaderOne)
		Expect(err).NotTo(HaveOccurred())

		err = pipLogValueRepository.Create(headerID, []interface{}{test_data.PipLogValueModel})
		Expect(err).NotTo(HaveOccurred())

		anotherPipLogValue := test_data.GetFakePipLogValue(anotherBlockNumber, transactionIdx, logValue)
		err = pipLogValueRepository.Create(headerID, []interface{}{anotherPipLogValue})
		Expect(err).NotTo(HaveOccurred())

		expectedValues := []test_helpers.LogValue{
			{
				Val:         test_data.PipLogValueModel.Value,
				BlockNumber: test_data.PipLogValueModel.BlockNumber,
				TxIdx:       test_data.PipLogValueModel.TransactionIndex,
			},
			{
				Val:         anotherPipLogValue.Value,
				BlockNumber: anotherPipLogValue.BlockNumber,
				TxIdx:       anotherPipLogValue.TransactionIndex,
			},
		}

		var dbPipLogValue []test_helpers.LogValue
		err = db.Select(&dbPipLogValue, `SELECT val, maker.pip_log_value.block_number, tx_idx FROM maker.pip_log_value 
                                            JOIN public.headers ON public.headers.id = maker.pip_log_value.header_id
                                            WHERE public.headers.block_timestamp BETWEEN $1 AND $2`, beginningTimeRange, endingTimeRange)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbPipLogValue).To(Equal(expectedValues))
	})

	It("returns 1 pip log value between a time range", func() {
		var (
			anotherBlockNumber int64 = 10606965
			beginningTimeRange int64 = 111111111
			endingTimeRange    int64 = 111111113
			outsideTimeRange   int64 = 111111200
			logValue                 = "123456789"
			transactionIdx           = 3
		)

		fakeHeaderOne := fakes.GetFakeHeaderWithTimestamp(beginningTimeRange, int64(test_data.PipLogValueModel.BlockNumber))
		headerID, err := headerRepository.CreateOrUpdateHeader(fakeHeaderOne)
		Expect(err).NotTo(HaveOccurred())

		err = pipLogValueRepository.Create(headerID, []interface{}{test_data.PipLogValueModel})
		Expect(err).NotTo(HaveOccurred())

		fakeHeaderTwo := fakes.GetFakeHeaderWithTimestamp(outsideTimeRange, anotherBlockNumber)
		anotherHeaderID, err := headerRepository.CreateOrUpdateHeader(fakeHeaderTwo)
		Expect(err).NotTo(HaveOccurred())

		anotherPipLogValue := test_data.GetFakePipLogValue(anotherBlockNumber, transactionIdx, logValue)
		err = pipLogValueRepository.Create(anotherHeaderID, []interface{}{anotherPipLogValue})
		Expect(err).NotTo(HaveOccurred())

		expectedValues := []test_helpers.LogValue{
			{
				Val:         test_data.PipLogValueModel.Value,
				BlockNumber: test_data.PipLogValueModel.BlockNumber,
				TxIdx:       test_data.PipLogValueModel.TransactionIndex,
			},
		}

		var dbPipLogValue []test_helpers.LogValue
		err = db.Select(&dbPipLogValue, `SELECT val, maker.pip_log_value.block_number, tx_idx FROM maker.pip_log_value 
                                            JOIN public.headers ON public.headers.id = maker.pip_log_value.header_id
                                            WHERE public.headers.block_timestamp BETWEEN $1 AND $2`, beginningTimeRange, endingTimeRange)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbPipLogValue).To(Equal(expectedValues))
	})

})

type LogValueTx struct {
	TransactionHash  string `db:"hash"`
	TransactionIndex string `db:"tx_index"`
	BlockHeight      string `db:"block_number"`
	BlockHash        string `db:"block_hash"`
	TxFrom           string `db:"tx_from"`
	TxTo             string `db:"tx_to"`
}
