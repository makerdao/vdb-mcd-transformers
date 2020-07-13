package auction_file_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/auction_file"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Flip file transformer", func() {
	var (
		db          = test_config.NewTestDB(test_config.NewTestNode())
		transformer = auction_file.Transformer{}
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("returns err if log missing topics", func() {
		badLog := core.EventLog{
			Log: types.Log{
				Topics: []common.Hash{{}},
				Data:   []byte{1, 1, 1, 1, 1},
			}}

		_, err := transformer.ToModels(constants.FlipABI(), []core.EventLog{badLog}, nil)
		Expect(err).To(HaveOccurred())
	})

	It("converts a log to a model", func() {
		models, err := transformer.ToModels(constants.FlipABI(), []core.EventLog{test_data.AuctionFileEventLog}, db)

		expectedModel := test_data.AuctionFileModel()
		addressID, addressErr := shared.GetOrCreateAddress(test_data.AuctionFileEventLog.Log.Address.Hex(), db)
		Expect(addressErr).NotTo(HaveOccurred())
		msgSenderAddressID, msgSenderAddressErr := shared.GetOrCreateAddress(test_data.AuctionFileEventLog.Log.Topics[1].Hex(), db)
		Expect(msgSenderAddressErr).NotTo(HaveOccurred())
		expectedModel.ColumnValues[event.AddressFK] = addressID
		expectedModel.ColumnValues[constants.MsgSenderColumn] = msgSenderAddressID

		Expect(err).NotTo(HaveOccurred())
		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})
})