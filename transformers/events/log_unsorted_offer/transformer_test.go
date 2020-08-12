package log_unsorted_offer_test

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_unsorted_offer"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogUnsortedOffer Transformer", func() {
	var (
		transformer = log_unsorted_offer.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("converts a raw log to a LogUnsortedOffer model", func() {
		models, err := transformer.ToModels(constants.OasisABI(), []core.EventLog{test_data.LogUnsortedOfferEventLog}, db)
		Expect(err).NotTo(HaveOccurred())

		expectedModel := test_data.LogUnsortedOfferModel()
		test_data.AssignAddressID(test_data.LogUnsortedOfferEventLog, expectedModel, db)

		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})

	It("returns an error if converting log to entity fails", func() {
		_, err := transformer.ToModels("error abi", []core.EventLog{test_data.LogUnsortedOfferEventLog}, db)

		Expect(err).To(HaveOccurred())
	})
})
