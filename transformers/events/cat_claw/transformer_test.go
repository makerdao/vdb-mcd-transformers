package cat_claw_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/cat_claw"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cat Claw Transformer", func() {
	var (
		transformer = cat_claw.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("converts claw logs to models", func() {
		models, err := transformer.ToModels(constants.Cat110ABI(), []core.EventLog{test_data.CatClawEventLog}, db)
		Expect(err).NotTo(HaveOccurred())

		addressID, addressErr := repository.GetOrCreateAddress(db, test_data.Cat110Address())
		Expect(addressErr).NotTo(HaveOccurred())

		msgSender := "0xF32836B9E1f47a0515c6Ec431592D5EbC276407f"
		msgSenderID, msgSenderErr := repository.GetOrCreateAddress(db, msgSender)
		Expect(msgSenderErr).NotTo(HaveOccurred())

		expectedModel := test_data.CatClawModel()
		expectedModel.ColumnValues[event.AddressFK] = addressID
		expectedModel.ColumnValues[constants.MsgSenderColumn] = msgSenderID

		Expect(models).To(ContainElement(expectedModel))
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
