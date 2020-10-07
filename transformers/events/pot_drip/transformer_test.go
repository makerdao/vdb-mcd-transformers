package pot_drip_test

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/pot_drip"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pot drip transformer", func() {
	var (
		transformer = pot_drip.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("returns err if log missing topics", func() {
		badLog := core.EventLog{}

		_, err := transformer.ToModels(constants.PotABI(), []core.EventLog{badLog}, db)
		Expect(err).To(HaveOccurred())
	})

	It("converts a log to model", func() {
		model, err := transformer.ToModels(constants.PotABI(), []core.EventLog{test_data.PotDripEventLog}, db)

		Expect(err).NotTo(HaveOccurred())
		expectedModel := test_data.PotDripModel()
		test_data.AssignMessageSenderID(test_data.PotDripEventLog, expectedModel, db)
		Expect(model).To(Equal([]event.InsertionModel{expectedModel}))
	})
})
