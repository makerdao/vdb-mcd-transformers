package log_delete_test

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_delete"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogDelete Transformer", func() {
	var (
		transformer = log_delete.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("converts a raw log to a LogDelete model", func() {
		models, err := transformer.ToModels(constants.OasisABI(), []core.EventLog{test_data.LogDeleteEventLog}, db)
		Expect(err).NotTo(HaveOccurred())

		expectedModel := test_data.LogDeleteModel()
		test_data.AssignAddressID(test_data.LogDeleteEventLog, expectedModel, db)

		keeperID, keeperErr := shared.GetOrCreateAddress(test_data.LogDeleteKeeperAddress.Hex(), db)
		Expect(keeperErr).NotTo(HaveOccurred())
		expectedModel.ColumnValues[constants.KeeperColumn] = keeperID

		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})

	It("returns an error if converting log to entity fails", func() {
		_, err := transformer.ToModels("error abi", []core.EventLog{test_data.LogDeleteEventLog}, db)

		Expect(err).To(HaveOccurred())
	})
})
