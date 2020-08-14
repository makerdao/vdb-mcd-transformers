package log_min_sell_test

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_min_sell"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogMinSell Transformer", func() {
	var (
		transformer = log_min_sell.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("converts a raw log to a LogMinSell model", func() {
		models, err := transformer.ToModels(constants.OasisABI(), []core.EventLog{test_data.LogMinSellEventLog}, db)
		Expect(err).NotTo(HaveOccurred())

		expectedModel := test_data.LogMinSellModel()
		test_data.AssignAddressID(test_data.LogMinSellEventLog, expectedModel, db)

		payGemID, payGemErr := shared.GetOrCreateAddress(test_data.LogMinSellPayGemAddress.Hex(), db)
		Expect(payGemErr).NotTo(HaveOccurred())
		expectedModel.ColumnValues[constants.PayGemColumn] = payGemID

		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})

	It("returns an error if converting log to entity fails", func() {
		_, err := transformer.ToModels("error abi", []core.EventLog{test_data.LogMinSellEventLog}, db)

		Expect(err).To(HaveOccurred())
	})
})
