package log_bump_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_bump"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogBump Transformer", func() {
	var (
		transformer = log_bump.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("converts a raw log to a LogBump model", func() {
		models, err := transformer.ToModels(constants.OasisABI(), []core.EventLog{test_data.LogBumpEventLog}, db)
		Expect(err).NotTo(HaveOccurred())

		expectedModel := test_data.LogBumpModel()
		addressID, addressErr := shared.GetOrCreateAddress(test_data.LogBumpEventLog.Log.Address.Hex(), db)
		Expect(addressErr).NotTo(HaveOccurred())
		expectedModel.ColumnValues[event.AddressFK] = addressID
		makerID, makerErr := shared.GetOrCreateAddress(common.HexToAddress(test_data.LogBumpEventLog.Log.Topics[3].Hex()).Hex(), db)
		Expect(makerErr).NotTo(HaveOccurred())
		expectedModel.ColumnValues[constants.MakerColumn] = makerID
		payGemID, payGemErr := shared.GetOrCreateAddress(test_data.LogBumpPayGemAddress.Hex(), db)
		Expect(payGemErr).NotTo(HaveOccurred())
		expectedModel.ColumnValues[constants.PayGemColumn] = payGemID
		buyGemID, buyGemErr := shared.GetOrCreateAddress(test_data.LogBumpBuyGemAddress.Hex(), db)
		Expect(buyGemErr).NotTo(HaveOccurred())
		expectedModel.ColumnValues[constants.BuyGemColumn] = buyGemID

		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})

	It("returns an error if converting log to entity fails", func() {
		_, err := transformer.ToModels("error abi", []core.EventLog{test_data.LogBumpEventLog}, db)

		Expect(err).To(HaveOccurred())
	})
})
