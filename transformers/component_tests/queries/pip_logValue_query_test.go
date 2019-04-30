package queries

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
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
		expectedValues := []string{test_data.PipLogValueModel.Value, test_data.PipLogValueModel.Value}
		headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
		Expect(err).NotTo(HaveOccurred())

		err = pipLogValueRepository.Create(headerID, []interface{}{test_data.PipLogValueModel})
		Expect(err).NotTo(HaveOccurred())

		anotherFakeHeader := fakes.FakeHeader
		anotherFakeHeader.BlockNumber = anotherFakeHeader.BlockNumber + 1
		anotherFakeHeader.Timestamp = "111111112"
		anotherHeaderID, err := headerRepository.CreateOrUpdateHeader(anotherFakeHeader)
		Expect(err).NotTo(HaveOccurred())

		err = pipLogValueRepository.Create(anotherHeaderID, []interface{}{test_data.PipLogValueModel})
		Expect(err).NotTo(HaveOccurred())

		var dbPipLogValue []string
		err = db.Select(&dbPipLogValue, `SELECT val FROM maker.pip_log_value 
											JOIN public.headers ON public.headers.id = maker.pip_log_value.header_id
											WHERE public.headers.block_timestamp BETWEEN $1 AND $2`, beginningTimeRange, endingTimeRange)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbPipLogValue).To(Equal(expectedValues))
	})
})
