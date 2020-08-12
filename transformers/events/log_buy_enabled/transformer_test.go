package log_buy_enabled_test

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_buy_enabled"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogBuyEnabled Transformer", func() {
	var (
		transformer = log_buy_enabled.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("converts a raw log to a LogBuyEnabled model", func() {
		models, err := transformer.ToModels(constants.OasisABI(), []core.EventLog{test_data.LogBuyEnabledEventLog}, db)
		Expect(err).NotTo(HaveOccurred())

		expectedModel := test_data.LogBuyEnabledModel()
		test_data.AssignAddressID(test_data.LogBuyEnabledEventLog, expectedModel, db)

		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})

	It("returns an error if converting log to entity fails", func() {
		_, err := transformer.ToModels("error abi", []core.EventLog{test_data.LogBuyEnabledEventLog}, db)

		Expect(err).To(HaveOccurred())
	})
})
