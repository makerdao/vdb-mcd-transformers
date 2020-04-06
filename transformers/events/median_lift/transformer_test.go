package median_lift_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lib/pq"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/median_lift"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Median lift transformer", func() {
	var (
		transformer = median_lift.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("returns err if log is missing topics", func() {
		incompleteLog := core.EventLog{}
		_, err := transformer.ToModels(constants.MedianABI(), []core.EventLog{incompleteLog}, db)
		Expect(err).To(HaveOccurred())
	})

	It("returns err if log is missing data", func() {
		badLog := core.EventLog{
			Log: types.Log{
				Topics: []common.Hash{{}, {}, {}, {}},
			}}

		_, err := transformer.ToModels(constants.MedianABI(), []core.EventLog{badLog}, db)
		Expect(err).To(HaveOccurred())
	})

	It("convert a log with 4 accounts to an insertion model", func() {
		models, err := transformer.ToModels(constants.MedianABI(), []core.EventLog{test_data.MedianLiftLogWithFourAccounts}, db)
		Expect(err).NotTo(HaveOccurred())

		a0Bytes, aErr := shared.GetLogNoteArgumentAtIndex(2, test_data.MedianLiftLogWithFourAccounts.Log.Data)
		Expect(aErr).NotTo(HaveOccurred())
		address0 := common.BytesToAddress(a0Bytes).String()

		a1Bytes, a2Err := shared.GetLogNoteArgumentAtIndex(3, test_data.MedianLiftLogWithFourAccounts.Log.Data)
		Expect(a2Err).NotTo(HaveOccurred())
		address1 := common.BytesToAddress(a1Bytes).String()

		a2Bytes, a3Err := shared.GetLogNoteArgumentAtIndex(4, test_data.MedianLiftLogWithFourAccounts.Log.Data)
		Expect(a3Err).NotTo(HaveOccurred())
		address2 := common.BytesToAddress(a2Bytes).String()

		a3Bytes, a4Err := shared.GetLogNoteArgumentAtIndex(5, test_data.MedianLiftLogWithFourAccounts.Log.Data)
		Expect(a4Err).NotTo(HaveOccurred())
		address3 := common.BytesToAddress(a3Bytes).String()

		expectedModel := test_data.MedianLiftModelWithFourAccounts()
		contractAddressID, contractAddressErr := shared.GetOrCreateAddress(test_data.MedianLiftLogWithFourAccounts.Log.Address.String(), db)
		Expect(contractAddressErr).NotTo(HaveOccurred())
		msgSenderAddressID, msgSenderAddressErr := shared.GetOrCreateAddress(test_data.MedianLiftLogWithFourAccounts.Log.Topics[1].Hex(), db)
		Expect(msgSenderAddressErr).NotTo(HaveOccurred())
		expectedModel.ColumnValues[event.AddressFK] = contractAddressID
		expectedModel.ColumnValues[constants.MsgSenderColumn] = msgSenderAddressID
		expectedModel.ColumnValues[constants.AColumn] = pq.Array([]string{address0, address1, address2, address3})

		Expect(models[0]).To(Equal(expectedModel))
	})

	It("convert a log with 1 account to an insertion model", func() {
		models, err := transformer.ToModels(constants.MedianABI(), []core.EventLog{test_data.MedianLiftLogWithOneAccount}, db)
		Expect(err).NotTo(HaveOccurred())

		a0Bytes, aErr := shared.GetLogNoteArgumentAtIndex(2, test_data.MedianLiftLogWithOneAccount.Log.Data)
		Expect(aErr).NotTo(HaveOccurred())
		address0 := common.BytesToAddress(a0Bytes).Hex()

		expectedModel := test_data.MedianLiftModelWithOneAccount()
		contractAddressID, contractAddressErr := shared.GetOrCreateAddress(test_data.MedianLiftLogWithOneAccount.Log.Address.String(), db)
		Expect(contractAddressErr).NotTo(HaveOccurred())
		msgSenderAddressID, msgSenderAddressErr := shared.GetOrCreateAddress(test_data.MedianLiftLogWithOneAccount.Log.Topics[1].Hex(), db)
		Expect(msgSenderAddressErr).NotTo(HaveOccurred())
		expectedModel.ColumnValues[event.AddressFK] = contractAddressID
		expectedModel.ColumnValues[constants.MsgSenderColumn] = msgSenderAddressID
		expectedModel.ColumnValues[constants.AColumn] = pq.Array([]string{address0})

		Expect(models[0]).To(Equal(expectedModel))
	})
})
