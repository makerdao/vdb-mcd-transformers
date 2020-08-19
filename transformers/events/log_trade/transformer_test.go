package log_trade_test

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_trade"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogTrade Transformer", func() {
	var (
		transformer = log_trade.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("converts a raw log to a LogTrade model", func() {
		models, err := transformer.ToModels(constants.OasisABI(), []core.EventLog{test_data.LogTradeEventLog}, db)
		Expect(err).NotTo(HaveOccurred())

		expectedModel := test_data.LogTradeModel()
		test_data.AssignAddressID(test_data.LogTradeEventLog, expectedModel, db)
		payGemID, payGemErr := shared.GetOrCreateAddress(test_data.LogTradeEventLog.Log.Topics[1].Hex(), db)
		Expect(payGemErr).NotTo(HaveOccurred())
		buyGemID, buyGemErr := shared.GetOrCreateAddress(test_data.LogTradeEventLog.Log.Topics[2].Hex(), db)
		Expect(buyGemErr).NotTo(HaveOccurred())

		expectedModel.ColumnValues[constants.PayGemColumn] = payGemID
		expectedModel.ColumnValues[constants.BuyGemColumn] = buyGemID

		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})

	It("returns an error if converting log to entity fails", func() {
		_, err := transformer.ToModels("error abi", []core.EventLog{test_data.LogTradeEventLog}, db)

		Expect(err).To(HaveOccurred())
	})
})
