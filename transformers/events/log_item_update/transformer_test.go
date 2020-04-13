package log_item_update_test

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_item_update"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogItemUpdate Transformer", func() {
	var (
		transformer = log_item_update.Transformer{}
		db = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("converts a raw log to a LogItemUpdate model", func() {
		models, err := transformer.ToModels(constants.OasisABI(), []core.EventLog{test_data.LogItemUpdateEventLog}, db)
		Expect(err).NotTo(HaveOccurred())

		expectedModel := test_data.LogItemUpdateModel()
		addressId, addressErr := shared.GetOrCreateAddress(test_data.LogItemUpdateEventLog.Log.Address.Hex(), db)
		Expect(addressErr).NotTo(HaveOccurred())
		expectedModel.ColumnValues[event.AddressFK] = addressId

		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})

	It("returns an error if converting log to entity fails", func() {
		_, err := transformer.ToModels("error abi", []core.EventLog{test_data.LogItemUpdateEventLog}, db)

		Expect(err).To(HaveOccurred())
	})
})
