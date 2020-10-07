package log_median_price_test

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_median_price"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogMedianPrice Transformer", func() {
	var (
		transformer = log_median_price.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("converts a raw to a LogMedianPrice model", func() {
		models, err := transformer.ToModels(constants.MedianABI(), []core.EventLog{test_data.EthLogMedianPriceEventLog}, db)
		Expect(err).NotTo(HaveOccurred())

		expectedModel := test_data.EthLogMedianPriceModel()
		test_data.AssignAddressID(test_data.EthLogMedianPriceEventLog, expectedModel, db)

		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})

	It("returns an error converting a log to an entity fails", func() {
		_, err := transformer.ToModels("error abi", []core.EventLog{test_data.EthLogMedianPriceEventLog}, db)
		Expect(err).To(HaveOccurred())
	})
})
