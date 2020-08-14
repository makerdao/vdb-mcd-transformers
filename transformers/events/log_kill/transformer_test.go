package log_kill_test

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_kill"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogKill Transformer", func() {
	var (
		transformer = log_kill.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("converts a raw log to a LogKill model", func() {
		models, err := transformer.ToModels(constants.OasisABI(), []core.EventLog{test_data.LogKillEventLog}, db)
		Expect(err).NotTo(HaveOccurred())

		expectedModel := test_data.LogKillModel()
		test_data.AssignAddressID(test_data.LogKillEventLog, expectedModel, db)
		makerAddressId, makerAddressErr := shared.GetOrCreateAddress(test_data.LogKillEventLog.Log.Topics[3].Hex(), db)
		Expect(makerAddressErr).NotTo(HaveOccurred())
		payGemAddressId, payGemAddressErr := shared.GetOrCreateAddress(test_data.PayGemAddress.Hex(), db)
		Expect(payGemAddressErr).NotTo(HaveOccurred())
		buyGemAddressId, buyGemAddressErr := shared.GetOrCreateAddress(test_data.BuyGemAddress.Hex(), db)
		Expect(buyGemAddressErr).NotTo(HaveOccurred())

		expectedModel.ColumnValues[constants.MakerColumn] = makerAddressId
		expectedModel.ColumnValues[constants.PayGemColumn] = payGemAddressId
		expectedModel.ColumnValues[constants.BuyGemColumn] = buyGemAddressId

		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})

	It("returns an error if converting log to entity fails", func() {
		_, err := transformer.ToModels("error abi", []core.EventLog{test_data.LogKillEventLog}, db)

		Expect(err).To(HaveOccurred())
	})
})
