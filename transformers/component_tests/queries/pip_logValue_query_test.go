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

	It("returns 2 pip log values between a time range", func() {
		beginningTimeRange := 111111111
		endingTimeRange := 111111112

		fakeHeaderOne := fakes.GetFakeHeaderWithTimestamp(111111111, 10606964)
		headerID, err := headerRepository.CreateOrUpdateHeader(fakeHeaderOne)
		Expect(err).NotTo(HaveOccurred())

		err = pipLogValueRepository.Create(headerID, []interface{}{test_data.PipLogValueModel})
		Expect(err).NotTo(HaveOccurred())

		fakeHeaderTwo := fakes.GetFakeHeaderWithTimestamp(111111112, 10606965)
		anotherHeaderID, err := headerRepository.CreateOrUpdateHeader(fakeHeaderTwo)
		Expect(err).NotTo(HaveOccurred())

		anotherPipLogValue := test_data.GetFakePipLogValue(10606965, 3, "123456789")
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
		Expect(dbPipLogValue).To(ConsistOf(expectedValues))
	})
})
