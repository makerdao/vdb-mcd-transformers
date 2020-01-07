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
		models, err := transformer.ToModels(constants.PotABI(), []core.HeaderSyncLog{test_data.PotExitHeaderSyncLog}, db)
		Expect(err).NotTo(HaveOccurred())

		var addressID int64
		addressErr := db.Get(&addressID, `SELECT id FROM addresses WHERE address = $1`,
			common.HexToAddress(test_data.PotExitHeaderSyncLog.Log.Topics[1].Hex()).Hex())
		Expect(addressErr).NotTo(HaveOccurred())
		expectedModel := test_data.PotExitModel()
		expectedModel.ColumnValues[constants.MsgSenderColumn] = addressID
		Expect(models).To(ConsistOf(expectedModel))
	})

	It("returns an error if there are missing topics", func() {
		invalidLog := test_data.PotExitHeaderSyncLog
		invalidLog.Log.Topics = []common.Hash{}

		_, err := transformer.ToModels(constants.PotABI(), []core.HeaderSyncLog{invalidLog}, db)

		Expect(err).To(MatchError(shared.ErrLogMissingTopics(3, 0)))
	})
})
