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

	It("convert a log to an insertion model", func() {
		models, err := transformer.ToModels(constants.MedianABI(), []core.EventLog{test_data.MedianDropLog}, db)
		Expect(err).NotTo(HaveOccurred())

		var aAddressID, a2AddressID, a3AddressID, a4AddressID int64
		aBytes, aErr := shared.GetLogNoteArgumentAtIndex(2, test_data.MedianDropLog.Log.Data)
		Expect(aErr).NotTo(HaveOccurred())
		aAddress := common.BytesToAddress(aBytes).String()
		aErr = db.Get(&aAddressID, `SELECT id FROM public.addresses where address = $1`, aAddress)
		Expect(aErr).NotTo(HaveOccurred())

		a2Bytes, a2Err := shared.GetLogNoteArgumentAtIndex(3, test_data.MedianDropLog.Log.Data)
		Expect(a2Err).NotTo(HaveOccurred())
		a2Address := common.BytesToAddress(a2Bytes).String()
		a2Err = db.Get(&a2AddressID, `SELECT id FROM public.addresses where address = $1`, a2Address)
		Expect(a2Err).NotTo(HaveOccurred())

		a3Bytes, a3Err := shared.GetLogNoteArgumentAtIndex(4, test_data.MedianDropLog.Log.Data)
		Expect(a3Err).NotTo(HaveOccurred())
		a3Address := common.BytesToAddress(a3Bytes).String()
		a3Err = db.Get(&a3AddressID, `SELECT id FROM public.addresses where address = $1`, a3Address)
		Expect(a3Err).NotTo(HaveOccurred())

		a4Bytes, a4Err := shared.GetLogNoteArgumentAtIndex(5, test_data.MedianDropLog.Log.Data)
		Expect(a4Err).NotTo(HaveOccurred())
		a4Address := common.BytesToAddress(a4Bytes).String()
		a4Err = db.Get(&a4AddressID, `SELECT id FROM public.addresses where address = $1`, a4Address)
		Expect(a4Err).NotTo(HaveOccurred())

		expectedModel := test_data.MedianDropModel()
		contractAddressID, contractAddressErr := shared.GetOrCreateAddress(test_data.MedianDropLog.Log.Address.String(), db)
		Expect(contractAddressErr).NotTo(HaveOccurred())
		msgSenderAddressID, msgSenderAddressErr := shared.GetOrCreateAddress(test_data.MedianDropLog.Log.Topics[1].Hex(), db)
		Expect(msgSenderAddressErr).NotTo(HaveOccurred())
		expectedModel.ColumnValues[event.AddressFK] = contractAddressID
		expectedModel.ColumnValues[constants.MsgSenderColumn] = msgSenderAddressID
		expectedModel.ColumnValues[constants.AColumn] = aAddressID
		expectedModel.ColumnValues[constants.A2Column] = a2AddressID
		expectedModel.ColumnValues[constants.A3Column] = a3AddressID
		expectedModel.ColumnValues[constants.A4Column] = a4AddressID

		Expect(models).To(ConsistOf(expectedModel))
	})
})
