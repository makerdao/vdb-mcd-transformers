package log_matching_enabled_test

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_matching_enabled"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogMatchingEnabled Transformer", func() {
	var (
		transformer = log_matching_enabled.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("converts a raw log to a LogMatchingEnabled model", func() {
		models, err := transformer.ToModels(constants.OasisABI(), []core.EventLog{test_data.LogMatchingEnabledEventLog}, db)
		Expect(err).NotTo(HaveOccurred())

		expectedModel := test_data.LogMatchingEnabledModel()
		addressId, addressErr := shared.GetOrCreateAddress(test_data.LogMatchingEnabledEventLog.Log.Address.Hex(), db)
		Expect(addressErr).NotTo(HaveOccurred())
		expectedModel.ColumnValues[event.AddressFK] = addressId

		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})

	It("returns an error if converting log to entity fails", func() {
		_, err := transformer.ToModels("error abi", []core.EventLog{test_data.LogMatchingEnabledEventLog}, db)

		Expect(err).To(HaveOccurred())
	})

})
