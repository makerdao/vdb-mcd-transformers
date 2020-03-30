package median_drop_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/median_drop"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Median drop transformer", func() {
	var (
		transformer = median_drop.Transformer{}
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
		models, err := transformer.ToModels(constants.MedianABI(), []core.EventLog{test_data.MedianDropLogWithFourAccounts}, db)
		Expect(err).NotTo(HaveOccurred())

		var addressID1, addressID2, addressID3, addressID4 int64
		a1Bytes, aErr := shared.GetLogNoteArgumentAtIndex(2, test_data.MedianDropLogWithFourAccounts.Log.Data)
		Expect(aErr).NotTo(HaveOccurred())
		address1 := common.BytesToAddress(a1Bytes).String()
		aErr = db.Get(&addressID1, `SELECT id FROM public.addresses where address = $1`, address1)
		Expect(aErr).NotTo(HaveOccurred())

		a2Bytes, a2Err := shared.GetLogNoteArgumentAtIndex(3, test_data.MedianDropLogWithFourAccounts.Log.Data)
		Expect(a2Err).NotTo(HaveOccurred())
		address2 := common.BytesToAddress(a2Bytes).String()
		a2Err = db.Get(&addressID2, `SELECT id FROM public.addresses where address = $1`, address2)
		Expect(a2Err).NotTo(HaveOccurred())

		a3Bytes, a3Err := shared.GetLogNoteArgumentAtIndex(4, test_data.MedianDropLogWithFourAccounts.Log.Data)
		Expect(a3Err).NotTo(HaveOccurred())
		address3 := common.BytesToAddress(a3Bytes).String()
		a3Err = db.Get(&addressID3, `SELECT id FROM public.addresses where address = $1`, address3)
		Expect(a3Err).NotTo(HaveOccurred())

		a4Bytes, a4Err := shared.GetLogNoteArgumentAtIndex(5, test_data.MedianDropLogWithFourAccounts.Log.Data)
		Expect(a4Err).NotTo(HaveOccurred())
		address4 := common.BytesToAddress(a4Bytes).String()
		a4Err = db.Get(&addressID4, `SELECT id FROM public.addresses where address = $1`, address4)
		Expect(a4Err).NotTo(HaveOccurred())

		expectedModel := test_data.MedianDropModelWithFourAccounts()
		contractAddressID, contractAddressErr := shared.GetOrCreateAddress(test_data.MedianDropLogWithFourAccounts.Log.Address.String(), db)
		Expect(contractAddressErr).NotTo(HaveOccurred())
		msgSenderAddressID, msgSenderAddressErr := shared.GetOrCreateAddress(test_data.MedianDropLogWithFourAccounts.Log.Topics[1].Hex(), db)
		Expect(msgSenderAddressErr).NotTo(HaveOccurred())
		expectedModel.ColumnValues[event.AddressFK] = contractAddressID
		expectedModel.ColumnValues[constants.MsgSenderColumn] = msgSenderAddressID
		expectedModel.ColumnValues[constants.A0Column] = addressID1
		expectedModel.ColumnValues[constants.A1Column] = addressID2
		expectedModel.ColumnValues[constants.A2Column] = addressID3
		expectedModel.ColumnValues[constants.A3Column] = addressID4

		Expect(models[0]).To(Equal(expectedModel))
	})

	It("convert a log with 1 account to an insertion model", func() {
		models, err := transformer.ToModels(constants.MedianABI(), []core.EventLog{test_data.MedianDropLogWithOneAccount}, db)
		Expect(err).NotTo(HaveOccurred())

		var addressID1, emptySlot int64
		aBytes, aErr := shared.GetLogNoteArgumentAtIndex(2, test_data.MedianDropLogWithOneAccount.Log.Data)
		Expect(aErr).NotTo(HaveOccurred())
		address1 := common.BytesToAddress(aBytes).String()
		aErr = db.Get(&addressID1, `SELECT id FROM public.addresses where address = $1`, address1)
		Expect(aErr).NotTo(HaveOccurred())

		expectedModel := test_data.MedianDropModelWithOneAccount()
		contractAddressID, contractAddressErr := shared.GetOrCreateAddress(test_data.MedianDropLogWithOneAccount.Log.Address.String(), db)
		Expect(contractAddressErr).NotTo(HaveOccurred())
		msgSenderAddressID, msgSenderAddressErr := shared.GetOrCreateAddress(test_data.MedianDropLogWithOneAccount.Log.Topics[1].Hex(), db)
		Expect(msgSenderAddressErr).NotTo(HaveOccurred())
		expectedModel.ColumnValues[event.AddressFK] = contractAddressID
		expectedModel.ColumnValues[constants.MsgSenderColumn] = msgSenderAddressID
		expectedModel.ColumnValues[constants.A0Column] = addressID1
		expectedModel.ColumnValues[constants.A1Column] = emptySlot
		expectedModel.ColumnValues[constants.A2Column] = emptySlot
		expectedModel.ColumnValues[constants.A3Column] = emptySlot

		Expect(models[0]).To(Equal(expectedModel))
	})
})
