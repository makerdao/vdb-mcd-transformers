package hole_test

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/dog_file/hole"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DogFileHole transformer", func() {
	var (
		db          = test_config.NewTestDB(test_config.NewTestNode())
		transformer = hole.Transformer{}
	)

	It("returns an error if converting the log to an entity fails", func() {
		_, err := transformer.ToModels("wrong abi", []core.EventLog{test_data.DealEventLog}, db)

		Expect(err).To(HaveOccurred())
	})

	It("converts a log to a model", func() {
		models, err := transformer.ToModels(constants.DogABI(), []core.EventLog{test_data.DogFileHoleEventLog}, db)
		Expect(err).NotTo(HaveOccurred())
		expectedModel := test_data.DogFileHoleModel()
		test_data.AssignAddressID(test_data.DogFileHoleEventLog, expectedModel, db)
		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})
})
