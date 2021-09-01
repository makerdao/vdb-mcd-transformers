package dog_bark_test

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/dog_bark"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DogBark transformer", func() {
	var (
		db               = test_config.NewTestDB(test_config.NewTestNode())
		ilkEthIdentifier = "ETH-A"
		transformer      = dog_bark.Transformer{}
	)

	It("converts a log to a Model", func() {
		models, err := transformer.ToModels(constants.DogABI(), []core.EventLog{test_data.DogBarkEventLog}, db)
		Expect(err).NotTo(HaveOccurred())
		expectedModel := test_data.DogBarkModel()
		test_data.AssignIlkID(expectedModel, ilkEthIdentifier, db)
		test_data.AssignUrnID(expectedModel, db)
		test_data.AssignAddressID(test_data.DogBarkEventLog, expectedModel, db)
		test_data.AssignClip(test_data.ClipAddress, expectedModel, db)
		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})

	It("returns an error if converting the log to an entity fails", func() {
		_, err := transformer.ToModels("error abi", []core.EventLog{test_data.BiteEventLog}, db)

		Expect(err).To(HaveOccurred())
	})
})
