package vow_heal_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vow_heal"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VowHeal transformer", func() {
	var (
		db          = test_config.NewTestDB(test_config.NewTestNode())
		transformer = vow_heal.Transformer{}
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("converts log to a model", func() {
		expectedModel := test_data.VowHealModel()
		msgSenderAddressID, msgSenderErr := shared.GetOrCreateAddress(test_data.VowHealEventLog.Log.Topics[1].Hex(), db)
		Expect(msgSenderErr).NotTo(HaveOccurred())
		expectedModel.ColumnValues[constants.MsgSenderColumn] = msgSenderAddressID

		models, err := transformer.ToModels("", []core.EventLog{test_data.VowHealEventLog}, db)

		Expect(err).NotTo(HaveOccurred())
		Expect(len(models)).To(Equal(1))
		Expect(models[0]).To(Equal(expectedModel))
	})

	It("Returns an error there are missing topics", func() {
		badLog := core.EventLog{
			Log: types.Log{
				Topics: []common.Hash{
					common.HexToHash("0x"),
				}},
		}

		_, err := transformer.ToModels("", []core.EventLog{badLog}, db)

		Expect(err).To(HaveOccurred())
	})
})
