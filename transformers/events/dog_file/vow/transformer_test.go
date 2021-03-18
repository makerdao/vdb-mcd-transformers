package vow_test

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/dog_file/vow"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dog file vow transformer", func() {
	var (
		transformer = vow.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("returns err if log is missing topics", func() {
		badLog := core.EventLog{
			Log: types.Log{
				Data: []byte{1, 1, 1, 1, 1},
			}}

		_, err := transformer.ToModels(constants.DogABI(), []core.EventLog{badLog}, nil)
		Expect(err).To(HaveOccurred())
	})

	It("converts a log to a model", func() {
		models, err := transformer.ToModels(constants.DogABI(), []core.EventLog{test_data.DogFileVowEventLog}, db)
		Expect(err).NotTo(HaveOccurred())
		expectedModel := test_data.DogFileVowModel()
		test_data.AssignAddressID(test_data.DogFileVowEventLog, expectedModel, db)
		test_data.AssignDataAddressID(test_data.DataAddress, expectedModel, db)
		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})
})
