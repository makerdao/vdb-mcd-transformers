package clip_redo_test

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/clip_redo"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ClipKick transformer", func() {
	var (
		db          = test_config.NewTestDB(test_config.NewTestNode())
		transformer = clip_redo.Transformer{}
	)

	It("returns an error if converting the log to an entity fails", func() {
		_, err := transformer.ToModels("wrong abi", []core.EventLog{test_data.DealEventLog}, db)

		Expect(err).To(HaveOccurred())
	})

	It("converts a log to a Model", func() {
		models, err := transformer.ToModels(constants.ClipV150ABI(), []core.EventLog{test_data.ClipRedoEventLog}, db)
		Expect(err).NotTo(HaveOccurred())
		expectedModel := test_data.ClipRedoModel()
		test_data.AssignAddressID(test_data.ClipRedoEventLog, expectedModel, db)
		test_data.AssignClipUsrID(test_data.ClipRedoEventLog, expectedModel, db)
		test_data.AssignKprAddressID(test_data.ClipRedoEventLog, expectedModel, db)
		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})
})
