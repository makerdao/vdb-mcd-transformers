package pot_exit_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/pot_exit"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PotExit converter", func() {
	var converter = pot_exit.Converter{}

	It("converts log to a model", func() {
		models, err := converter.ToModels(constants.PotABI(), []core.HeaderSyncLog{test_data.PotExitHeaderSyncLog}, nil)

		Expect(err).NotTo(HaveOccurred())
		Expect(models).To(ConsistOf(test_data.PotExitModel()))
	})

	It("returns an error if there are missing topics", func() {
		invalidLog := test_data.PotExitHeaderSyncLog
		invalidLog.Log.Topics = []common.Hash{}

		_, err := converter.ToModels(constants.PotABI(), []core.HeaderSyncLog{invalidLog}, nil)

		Expect(err).To(MatchError(shared.ErrLogMissingTopics(3, 0)))
	})
})
