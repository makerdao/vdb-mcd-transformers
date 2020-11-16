package batch_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lib/pq"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/median_diss/batch"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Median diss (batch) transformer", func() {
	var (
		transformer = batch.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("returns err if log is missing topics", func() {
		incompleteLog := core.EventLog{}
		_, err := transformer.ToModels(constants.MedianV100ABI(), []core.EventLog{incompleteLog}, db)
		Expect(err).To(HaveOccurred())
	})

	It("returns err if log is missing data", func() {
		badLog := core.EventLog{
			Log: types.Log{
				Topics: []common.Hash{{}, {}, {}, {}},
			}}

		_, err := transformer.ToModels(constants.MedianV100ABI(), []core.EventLog{badLog}, db)
		Expect(err).To(HaveOccurred())
	})

	It("converts a log with a 1-address array to an insertion model", func() {
		models, err := transformer.ToModels(constants.MedianV100ABI(), []core.EventLog{test_data.MedianDissBatchLogOneAddress}, db)
		Expect(err).NotTo(HaveOccurred())

		expectedModel := test_data.MedianDissBatchModelOneAddress()
		contractAddressID, contractAddressErr := shared.GetOrCreateAddress(test_data.MedianDissBatchLogOneAddress.Log.Address.String(), db)
		Expect(contractAddressErr).NotTo(HaveOccurred())
		msgSenderAddressID, msgSenderAddressErr := shared.GetOrCreateAddress(test_data.MedianDissBatchLogOneAddress.Log.Topics[1].Hex(), db)
		Expect(msgSenderAddressErr).NotTo(HaveOccurred())

		a0Bytes, a0Err := shared.GetLogNoteArgumentAtIndex(2, test_data.MedianDissBatchLogOneAddress.Log.Data)
		Expect(a0Err).NotTo(HaveOccurred())
		address0 := common.BytesToAddress(a0Bytes).Hex()

		expectedModel.ColumnValues[event.AddressFK] = contractAddressID
		expectedModel.ColumnValues[constants.MsgSenderColumn] = msgSenderAddressID
		expectedModel.ColumnValues[constants.AColumn] = pq.Array([]string{address0})

		Expect(models).To(ConsistOf(expectedModel))
	})

	It("converts a log with a 5-address array to an insertion model, truncating the array to 4 addresses", func() {
		models, err := transformer.ToModels(constants.MedianV100ABI(), []core.EventLog{test_data.MedianDissBatchLogFiveAddresses}, db)
		Expect(err).NotTo(HaveOccurred())

		expectedModel := test_data.MedianDissBatchModelFiveAddresses()
		contractAddressID, contractAddressErr := shared.GetOrCreateAddress(test_data.MedianDissBatchLogFiveAddresses.Log.Address.String(), db)
		Expect(contractAddressErr).NotTo(HaveOccurred())
		msgSenderAddressID, msgSenderAddressErr := shared.GetOrCreateAddress(test_data.MedianDissBatchLogFiveAddresses.Log.Topics[1].Hex(), db)
		Expect(msgSenderAddressErr).NotTo(HaveOccurred())

		a0Bytes, a0Err := shared.GetLogNoteArgumentAtIndex(2, test_data.MedianDissBatchLogFiveAddresses.Log.Data)
		Expect(a0Err).NotTo(HaveOccurred())
		address0 := common.BytesToAddress(a0Bytes).Hex()
		a1Bytes, a1Err := shared.GetLogNoteArgumentAtIndex(3, test_data.MedianDissBatchLogFiveAddresses.Log.Data)
		Expect(a1Err).NotTo(HaveOccurred())
		address1 := common.BytesToAddress(a1Bytes).Hex()
		a2Bytes, a2Err := shared.GetLogNoteArgumentAtIndex(4, test_data.MedianDissBatchLogFiveAddresses.Log.Data)
		Expect(a2Err).NotTo(HaveOccurred())
		address2 := common.BytesToAddress(a2Bytes).Hex()
		a3Bytes, a3Err := shared.GetLogNoteArgumentAtIndex(5, test_data.MedianDissBatchLogFiveAddresses.Log.Data)
		Expect(a3Err).NotTo(HaveOccurred())
		address3 := common.BytesToAddress(a3Bytes).Hex()

		expectedModel.ColumnValues[event.AddressFK] = contractAddressID
		expectedModel.ColumnValues[constants.MsgSenderColumn] = msgSenderAddressID
		expectedModel.ColumnValues[constants.AColumn] = pq.Array([]string{address0, address1, address2, address3})

		Expect(models).To(ConsistOf(expectedModel))
	})
})
