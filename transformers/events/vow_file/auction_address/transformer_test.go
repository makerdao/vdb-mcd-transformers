package auction_address_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vow_file/auction_address"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vow file auction contract transformer", func() {
	var (
		transformer = auction_address.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("converts a flopper log to a model", func() {
		models, toModelsErr := transformer.ToModels(constants.VowABI(), []core.EventLog{test_data.VowFileAuctionAddressEventLog}, db)
		Expect(toModelsErr).NotTo(HaveOccurred())

		var dataAddressID int64
		dataAddressErr := db.Get(&dataAddressID, `SELECT id FROM addresses WHERE address = $1`,
			common.HexToAddress(test_data.VowFileAuctionAddressEventLog.Log.Topics[3].Hex()).Hex())
		Expect(dataAddressErr).NotTo(HaveOccurred())

		expectedModel := test_data.VowFileAuctionAddressModel()
		expectedModel.ColumnValues[constants.DataColumn] = dataAddressID

		Expect(models).To(ConsistOf(expectedModel))
	})

	It("returns err if the log is missing topics", func() {
		badLog := core.EventLog{
			Log: types.Log{
				Topics: []common.Hash{{}},
				Data:   []byte{1, 1, 1, 1, 1},
			}}

		_, err := transformer.ToModels(constants.VowABI(), []core.EventLog{badLog}, db)
		Expect(err).To(HaveOccurred())
	})
})
