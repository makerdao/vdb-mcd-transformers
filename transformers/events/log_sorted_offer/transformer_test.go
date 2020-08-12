package log_sorted_offer_test

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_sorted_offer"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogSortedOffer Transformer", func() {
	var (
		transformer = log_sorted_offer.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("converts a raw log to a LogSortedOffer model", func() {
		models, err := transformer.ToModels(constants.OasisABI(), []core.EventLog{test_data.LogSortedOfferEventLog}, db)
		Expect(err).NotTo(HaveOccurred())

		expectedModel := test_data.LogSortedOfferModel()
		test_data.AssignAddressID(test_data.LogSortedOfferEventLog, expectedModel, db)

		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})

	It("returns an error if converting log to entity fails", func() {
		_, err := transformer.ToModels("error abi", []core.EventLog{test_data.LogSortedOfferEventLog}, db)

		Expect(err).To(HaveOccurred())
	})
})
