package pot_exit_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/pot_exit"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PotExit transformer", func() {
	var (
		transformer = pot_exit.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("converts log to a model", func() {
		models, err := transformer.ToModels(constants.PotABI(), []core.EventLog{test_data.PotExitEventLog}, db)
		Expect(err).NotTo(HaveOccurred())

		expectedModel := test_data.PotExitModel()
		test_data.AssignMessageSenderID(test_data.PotExitEventLog, expectedModel, db)

		Expect(models).To(ConsistOf(expectedModel))
	})

	It("returns an error if there are missing topics", func() {
		invalidLog := test_data.PotExitEventLog
		invalidLog.Log.Topics = []common.Hash{}

		_, err := transformer.ToModels(constants.PotABI(), []core.EventLog{invalidLog}, db)

		Expect(err).To(MatchError(shared.ErrLogMissingTopics(3, 0)))
	})
})
