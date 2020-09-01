package dunk_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/cat_file/dunk"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cat file dunk transformer", func() {
	var (
		db          = test_config.NewTestDB(test_config.NewTestNode())
		transformer = dunk.Transformer{}
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("returns an error if log is missing topics", func() {
		badLog := core.EventLog{
			Log: types.Log{
				Data: []byte{1, 1, 1, 1, 1},
			},
		}

		_, err := transformer.ToModels(constants.Cat110ABI(), []core.EventLog{badLog}, nil)
		Expect(err).To(HaveOccurred())

	})

	It("returns err if log is missing data", func() {
		badLog := core.EventLog{
			Log: types.Log{
				Topics: []common.Hash{{}, {}, {}, {}},
			},
		}

		_, err := transformer.ToModels("", []core.EventLog{badLog}, nil)
		Expect(err).To(HaveOccurred())
	})

	It("converts a log to a model", func() {
		models, err := transformer.ToModels(constants.Cat110ABI(), []core.EventLog{test_data.CatFileDunkEventLog}, db)
		Expect(err).NotTo(HaveOccurred())

		ilkID, ilkErr := shared.GetOrCreateIlk(test_data.CatFileDunkEventLog.Log.Topics[2].Hex(), db)
		Expect(ilkErr).NotTo(HaveOccurred())
		addressID, addressErr := shared.GetOrCreateAddress(test_data.CatFileDunkEventLog.Log.Address.Hex(), db)
		Expect(addressErr).NotTo(HaveOccurred())
		msgSenderID, msgSenderErr := shared.GetOrCreateAddress(common.HexToAddress(test_data.CatFileDunkEventLog.Log.Topics[1].Hex()).Hex(), db)
		Expect(msgSenderErr).NotTo(HaveOccurred())

		expectedModel := test_data.CatFileDunkModel()
		expectedModel.ColumnValues[constants.IlkColumn] = ilkID
		expectedModel.ColumnValues[event.AddressFK] = addressID
		expectedModel.ColumnValues[constants.MsgSenderColumn] = msgSenderID

		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})
})
