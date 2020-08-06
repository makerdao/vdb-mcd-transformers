package pot_cage_test

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/pot_cage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pot cage transformer", func() {
	var (
		transformer = pot_cage.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("converts a pot cage log to a model", func() {
		models, err := transformer.ToModels(constants.PotABI(), []core.EventLog{test_data.PotCageEventLog}, db)

		Expect(err).NotTo(HaveOccurred())
		expectedModel := test_data.PotCageModel()
		test_data.AssignMessageSenderID(test_data.PotCageEventLog, expectedModel, db)
		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})

	It("returns err if log is missing topics", func() {
		badLog := core.EventLog{
			Log: types.Log{
				Data: []byte{1, 1, 1, 1, 1},
			},
		}

		_, err := transformer.ToModels(constants.PotABI(), []core.EventLog{badLog}, db)
		Expect(err).To(HaveOccurred())
	})
})
