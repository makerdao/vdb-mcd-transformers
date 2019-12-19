package cat_rely_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/cat_rely"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cat Rely converter", func() {
	var (
		converter event.Converter
		db        *postgres.DB
	)

	BeforeEach(func() {
		converter = cat_rely.Converter{}
		db = test_config.NewTestDB(test_config.NewTestNode())
	})

	It("converts a cat rely log to a model", func() {
		models, err := converter.ToModels(constants.CatABI(), []core.HeaderSyncLog{test_data.CatRelyHeaderSyncLog}, db)
		Expect(err).NotTo(HaveOccurred())
		expectedModel := test_data.CatRelyModel()
		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})

	It("returns an error if there are missing topics", func() {
		invalidLog := test_data.CatRelyHeaderSyncLog
		invalidLog.Log.Topics = []common.Hash{}

		_, err := converter.ToModels(constants.CatABI(), []core.HeaderSyncLog{invalidLog}, nil)

		Expect(err).To(MatchError(shared.ErrLogMissingTopics(2, 0)))
	})
})
