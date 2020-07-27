package par_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/spot_file/par"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Spot file par transformer", func() {
	var transformer = par.Transformer{}
	db := test_config.NewTestDB(test_config.NewTestNode())

	It("returns err if log missing topics", func() {
		badLog := core.EventLog{
			Log: types.Log{
				Topics: []common.Hash{{}},
				Data:   []byte{1, 1, 1, 1, 1},
			}}

		_, err := transformer.ToModels(constants.SpotABI(), []core.EventLog{badLog}, db)
		Expect(err).To(HaveOccurred())
	})

	It("converts a log to a model", func() {
		models, err := transformer.ToModels(constants.SpotABI(), []core.EventLog{test_data.SpotFileParEventLog}, db)

		expectedModel := test_data.SpotFileParModel()
		test_data.AssignMessageSenderID(test_data.SpotFileParEventLog, expectedModel, db)
		Expect(err).NotTo(HaveOccurred())
		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})
})
