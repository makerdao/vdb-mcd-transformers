package set_min_sell_test

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/set_min_sell"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SetMinSell Note Transformer", func() {
	var (
		transformer = set_min_sell.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("converts a raw log to a SetMinSell model", func() {
		models, err := transformer.ToModels("irrelevant", []core.EventLog{test_data.SetMinSellEventLog}, db)
		Expect(err).NotTo(HaveOccurred())

		expectedModel := test_data.SetMinSellModel()
		addressID, addressErr := shared.GetOrCreateAddress(test_data.SetMinSellEventLog.Log.Address.Hex(), db)
		Expect(addressErr).NotTo(HaveOccurred())
		expectedModel.ColumnValues[event.AddressFK] = addressID

		payGemID, payGemErr := shared.GetOrCreateAddress(test_data.SetMinSellPayGemAddress.Hex(), db)
		Expect(payGemErr).NotTo(HaveOccurred())
		expectedModel.ColumnValues[constants.PayGemColumn] = payGemID

		senderAddressID, senderAddressErr := shared.GetOrCreateAddress(test_data.SetMinSellMsgSenderAddress.Hex(), db)
		Expect(senderAddressErr).NotTo(HaveOccurred())
		expectedModel.ColumnValues[constants.MsgSenderColumn] = senderAddressID

		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})

	It("returns an error if missing log topics", func() {
		badLog := core.EventLog{}
		_, err := transformer.ToModels("irrelevant", []core.EventLog{badLog}, db)

		Expect(err).To(MatchError(ContainSubstring("log missing topics: has 0, want 4")))
	})
})
