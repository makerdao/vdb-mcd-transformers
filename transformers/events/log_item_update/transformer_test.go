package log_item_update_test

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_item_update"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogItemUpdate Transformer", func() {
	var transformer = log_item_update.Transformer{}
	It("converts a raw log to a LogItemUpdate model", func() {
		models, err := transformer.ToModels(constants.OasisABI(), []core.EventLog{test_data.LogItemUpdateEventLog}, nil)
		Expect(err).NotTo(HaveOccurred())

		Expect(models).To(Equal([]event.InsertionModel{test_data.LogItemUpdateModel()}))
	})

	It("returns an error if converting log to entity fails", func() {
		_, err := transformer.ToModels("error abi", []core.EventLog{test_data.LogItemUpdateEventLog}, nil)

		Expect(err).To(HaveOccurred())
	})
})
