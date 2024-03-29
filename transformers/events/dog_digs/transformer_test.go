package dog_digs_test

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/dog_digs"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DogDigs transformer", func() {
	var (
		db               = test_config.NewTestDB(test_config.NewTestNode())
		ilkEthIdentifier = "ETH-A"
		transformer      = dog_digs.Transformer{}
	)
	It("returns an error if converting the log to an entity fails", func() {
		_, err := transformer.ToModels("wrong abi", []core.EventLog{test_data.DealEventLog}, db)

		Expect(err).To(HaveOccurred())
	})

	It("converts a log to a Model", func() {
		models, err := transformer.ToModels(constants.DogABI(), []core.EventLog{test_data.DogDigsEventLog}, db)
		Expect(err).NotTo(HaveOccurred())
		expectedModel := test_data.DogDigsModel()
		test_data.AssignIlkID(expectedModel, ilkEthIdentifier, db)
		test_data.AssignAddressID(test_data.DogDigsEventLog, expectedModel, db)
		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})
})
