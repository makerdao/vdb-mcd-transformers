package queries

import (
	"database/sql"
	"math/rand"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/pip_log_value"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
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
		rand.Seed(GinkgoRandomSeed())
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	It("returns 2 pip log values in different blocks between a time range", func() {
		var (
			anotherBlockNumber int64 = 10606965
			transactionIdx           = 3
		)

		logValue := strconv.Itoa(test_helpers.GetRandomInt(100, 10000))
		beginningTimeRange := int64(test_helpers.GetRandomInt(1558710000, 1558720000))
		endingTimeRange := int64(test_helpers.GetRandomInt(1558720001, 1558730000))

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
				Val:             test_data.PipLogValueModel.Value,
				BlockNumber:     test_data.PipLogValueModel.BlockNumber,
				TxIdx:           test_data.PipLogValueModel.TransactionIndex,
				ContractAddress: test_data.PipLogValueModel.ContractAddress,
			},
			{
				Val:             anotherPipLogValue.Value,
				BlockNumber:     anotherPipLogValue.BlockNumber,
				TxIdx:           anotherPipLogValue.TransactionIndex,
				ContractAddress: anotherPipLogValue.ContractAddress,
			},
		}

		var dbPipLogValue []test_helpers.LogValue
		err = db.Select(&dbPipLogValue, `SELECT * FROM api.log_values($1, $2)`, beginningTimeRange, endingTimeRange)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbPipLogValue).To(Equal(expectedValues))
	})

	It("returns a transaction from a logValue", func() {
		var (
			anotherBlockNumber int64 = 10606965
			transactionIdx           = 3
		)

		logValue := strconv.Itoa(test_helpers.GetRandomInt(100, 10000))
		beginningTimeRange := int64(test_helpers.GetRandomInt(1558710000, 1558720000))
		endingTimeRange := int64(test_helpers.GetRandomInt(1558720001, 1558730000))

		fakeHeaderTwo := fakes.GetFakeHeaderWithTimestamp(endingTimeRange, anotherBlockNumber)
		anotherHeaderID, err := headerRepository.CreateOrUpdateHeader(fakeHeaderTwo)
		Expect(err).NotTo(HaveOccurred())

		anotherPipLogValue := test_data.GetFakePipLogValue(anotherBlockNumber, transactionIdx, logValue)
		err = pipLogValueRepository.Create(anotherHeaderID, []interface{}{anotherPipLogValue})
		Expect(err).NotTo(HaveOccurred())

		expectedTx := Tx{
			TransactionHash:  test_helpers.GetValidNullString(fakeHeaderTwo.Hash),
			TransactionIndex: sql.NullInt64{Int64: int64(transactionIdx), Valid: true},
			BlockHeight:      sql.NullInt64{Int64: anotherBlockNumber, Valid: true},
			BlockHash:        test_helpers.GetValidNullString(fakeHeaderTwo.Hash),
			TxFrom:           test_helpers.GetValidNullString("fromAddress"),
			TxTo:             test_helpers.GetValidNullString("toAddress"),
		}
		_, err = db.Exec(`INSERT INTO header_sync_transactions (header_id, hash, tx_from, tx_index, tx_to)
                VALUES ($1, $2, $3, $4, $5)`, anotherHeaderID, expectedTx.TransactionHash, expectedTx.TxFrom,
			expectedTx.TransactionIndex, expectedTx.TxTo)
		Expect(err).NotTo(HaveOccurred())

		var actualTx []Tx
		err = db.Select(&actualTx, `SELECT * FROM api.log_value_tx(
			(SELECT (val, block_number, tx_idx, contract_address)::api.log_value 
				FROM api.log_values($1, $2)))`, beginningTimeRange, endingTimeRange)

		Expect(err).NotTo(HaveOccurred())
		Expect(actualTx[0]).To(Equal(expectedTx))
	})

	It("returns 2 pip log values with transactions in the same block", func() {
		var (
			anotherBlockNumber int64 = 10606964
			transactionIdx           = 3
		)

		logValue := strconv.Itoa(test_helpers.GetRandomInt(100, 10000))
		beginningTimeRange := int64(test_helpers.GetRandomInt(1558710000, 1558720000))
		endingTimeRange := int64(test_helpers.GetRandomInt(1558720001, 1558730000))

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
				Val:             test_data.PipLogValueModel.Value,
				BlockNumber:     test_data.PipLogValueModel.BlockNumber,
				TxIdx:           test_data.PipLogValueModel.TransactionIndex,
				ContractAddress: test_data.PipLogValueModel.ContractAddress,
			},
			{
				Val:             anotherPipLogValue.Value,
				BlockNumber:     anotherPipLogValue.BlockNumber,
				TxIdx:           anotherPipLogValue.TransactionIndex,
				ContractAddress: anotherPipLogValue.ContractAddress,
			},
		}

		var dbPipLogValue []test_helpers.LogValue
		err = db.Select(&dbPipLogValue, `SELECT * FROM api.log_values($1, $2)`, beginningTimeRange, endingTimeRange)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbPipLogValue).To(Equal(expectedValues))
	})

	It("returns 1 pip log value between a time range", func() {
		var (
			anotherBlockNumber int64 = 10606965
			transactionIdx           = 3
		)

		logValue := strconv.Itoa(test_helpers.GetRandomInt(100, 10000))
		beginningTimeRange := int64(test_helpers.GetRandomInt(1558710000, 1558720000))
		endingTimeRange := int64(test_helpers.GetRandomInt(1558720001, 1558730000))
		outsideTimeRange := int64(test_helpers.GetRandomInt(1558730001, 1558740000))

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
				Val:             test_data.PipLogValueModel.Value,
				BlockNumber:     test_data.PipLogValueModel.BlockNumber,
				TxIdx:           test_data.PipLogValueModel.TransactionIndex,
				ContractAddress: test_data.PipLogValueModel.ContractAddress,
			},
		}

		var dbPipLogValue []test_helpers.LogValue
		err = db.Select(&dbPipLogValue, `SELECT * FROM api.log_values($1, $2)`, beginningTimeRange, endingTimeRange)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbPipLogValue).To(Equal(expectedValues))
	})

	It("uses default arguments when none are passed in", func() {
		_, err := db.Exec(`SELECT * FROM api.log_values()`)
		Expect(err).NotTo(HaveOccurred())
	})
})
