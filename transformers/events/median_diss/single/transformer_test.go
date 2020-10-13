package single_test

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/median_diss/single"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Median diss (single) transformer", func() {
	var (
		transformer = single.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("returns err if log is missing topics", func() {
		incompleteLog := core.EventLog{}
		_, err := transformer.ToModels(constants.MedianABI(), []core.EventLog{incompleteLog}, db)
		Expect(err).To(HaveOccurred())
	})

	It("convert a log to an insertion model", func() {
		models, err := transformer.ToModels(constants.MedianABI(), []core.EventLog{test_data.MedianDissSingleLog}, db)
		Expect(err).NotTo(HaveOccurred())

		expectedModel := test_data.MedianDissSingleModel()
		test_data.AssignAddressID(test_data.MedianDissSingleLog, expectedModel, db)
		test_data.AssignMessageSenderID(test_data.MedianDissSingleLog, expectedModel, db)
		aAddressID, aAddressErr := repository.GetOrCreateAddress(db, test_data.MedianDissSingleLog.Log.Topics[2].Hex())
		Expect(aAddressErr).NotTo(HaveOccurred())
		expectedModel.ColumnValues[constants.AColumn] = aAddressID

		Expect(models).To(ConsistOf(expectedModel))
	})
})
