package osm_change_test

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/osm_change"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("OsmChange transformer", func() {
	var (
		transformer = osm_change.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("returns err if log is missing topics", func() {
		incompleteLog := core.EventLog{}
		_, err := transformer.ToModels(constants.OsmABI(), []core.EventLog{incompleteLog}, db)
		Expect(err).To(HaveOccurred())
	})

	It("convert a log to an insertion model", func() {
		models, err := transformer.ToModels(constants.OsmABI(), []core.EventLog{test_data.OsmChangeEventLog}, db)
		Expect(err).NotTo(HaveOccurred())

		expectedModel := test_data.OsmChangeModel()
		test_data.AssignAddressID(test_data.OsmChangeEventLog, expectedModel, db)
		test_data.AssignMessageSenderID(test_data.OsmChangeEventLog, expectedModel, db)
		srcAddressID, srcAddressErr := shared.GetOrCreateAddress(test_data.OsmChangeEventLog.Log.Topics[2].Hex(), db)
		Expect(srcAddressErr).NotTo(HaveOccurred())
		expectedModel.ColumnValues[constants.SrcColumn] = srcAddressID

		Expect(models).To(ConsistOf(expectedModel))
	})
})
