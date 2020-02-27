package val_void_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/val_void"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Val Void Transformer", func() {
	var (
		transformer = val_void.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("converts logs to models", func() {
		models, err := transformer.ToModels(constants.ValABI(), []core.EventLog{test_data.ValVoidEventLog}, db)
		Expect(err).NotTo(HaveOccurred())

		addressID, addressErr := shared.GetOrCreateAddress(test_data.ValVoidEventLog.Log.Address.Hex(), db)
		Expect(addressErr).NotTo(HaveOccurred())

		msgSenderAddress := common.HexToAddress(test_data.ValVoidEventLog.Log.Topics[1].Hex()).Hex()
		msgSenderID, msgSenderErr := shared.GetOrCreateAddress(msgSenderAddress, db)
		Expect(msgSenderErr).NotTo(HaveOccurred())

		expectedModel := test_data.ValVoidModel()
		expectedModel.ColumnValues[event.AddressFK] = addressID
		expectedModel.ColumnValues[constants.MsgSenderColumn] = msgSenderID
		Expect(models).To(ConsistOf(expectedModel))
	})

	It("returns an error if the expected amount of topics aren't in the log", func() {
		invalidLog := test_data.ValVoidEventLog
		invalidLog.Log.Topics = []common.Hash{}

		_, err := transformer.ToModels(constants.ValABI(), []core.EventLog{invalidLog}, db)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(shared.ErrLogMissingTopics(2, 0)))
	})
})
