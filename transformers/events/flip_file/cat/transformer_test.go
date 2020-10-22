package cat_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/flip_file/cat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Flip file cat transformer", func() {
	var (
		transformer = cat.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("converts a flip file cat to a model", func() {
		models, err := transformer.ToModels(constants.FlipV110ABI(), []core.EventLog{test_data.FlipFileCatEventLog}, db)
		Expect(err).NotTo(HaveOccurred())

		addressID, addressErr := shared.GetOrCreateAddress(test_data.FlipEthAV110Address(), db)
		Expect(addressErr).NotTo(HaveOccurred())

		msgSenderID, msgSenderErr := shared.GetOrCreateAddress("0xbe8e3e3618f7474f8cb1d074a26affef007e98fb", db)
		Expect(msgSenderErr).NotTo(HaveOccurred())

		dataColumnID, dataColumnErr := shared.GetOrCreateAddress("0xA950524441892A31ebddF91d3cEEFa04Bf454466", db)
		Expect(dataColumnErr).NotTo(HaveOccurred())

		expectedFlipFileCat := test_data.FlipFileCatModel()
		expectedFlipFileCat.ColumnValues[event.AddressFK] = addressID
		expectedFlipFileCat.ColumnValues[constants.MsgSenderColumn] = msgSenderID
		expectedFlipFileCat.ColumnValues[constants.DataColumn] = dataColumnID

		Expect(models[0]).To(Equal(expectedFlipFileCat))
	})

	It("returns an err if the log is missing topics", func() {
		badLog := core.EventLog{
			Log: types.Log{
				Topics: []common.Hash{
					common.HexToHash("0xtest"),
				},
			},
		}
		_, err := transformer.ToModels(constants.FlipV100ABI(), []core.EventLog{badLog}, nil)
		Expect(err).To(HaveOccurred())
	})

})
