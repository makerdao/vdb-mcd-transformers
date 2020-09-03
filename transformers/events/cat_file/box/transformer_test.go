package box_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/cat_file/box"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cat file box transformer", func() {
	var (
		transformer = box.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("converts a box log to a model", func() {
		models, toModelsErr := transformer.ToModels(constants.Cat110ABI(), []core.EventLog{test_data.CatFileBoxEventLog}, db)
		Expect(toModelsErr).NotTo(HaveOccurred())

		contractAddressID, msgSenderErr := shared.GetOrCreateAddress(test_data.Cat110Address(), db)
		Expect(msgSenderErr).NotTo(HaveOccurred())

		msgSender := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderID, msgSenderErr := shared.GetOrCreateAddress(msgSender, db)
		Expect(msgSenderErr).NotTo(HaveOccurred())

		expectedModel := test_data.CatFileBoxModel()
		expectedModel.ColumnValues[event.AddressFK] = contractAddressID
		expectedModel.ColumnValues[constants.MsgSenderColumn] = msgSenderID

		Expect(models).To(ConsistOf(expectedModel))
	})

	It("returns an err if the log is missing topics", func() {
		badLog := core.EventLog{
			Log: types.Log{
				Topics: []common.Hash{
					common.HexToHash("0xtest"),
				},
			},
		}
		_, err := transformer.ToModels(constants.Cat110ABI(), []core.EventLog{badLog}, nil)
		Expect(err).To(HaveOccurred())
	})
})
