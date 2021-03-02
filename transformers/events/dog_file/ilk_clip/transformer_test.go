package ilk_clip_test

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/dog_file/ilk_clip"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)
var _ = Describe("Dog file ilk clip transformer", func() {
	var (
		transformer = ilk_clip.Transformer{}
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
		models, err := transformer.ToModels(constants.DogABI(), []core.EventLog{test_data.DogFileIlkClipEventLog}, db)
		Expect(err).NotTo(HaveOccurred())
		expectedModel := test_data.DogFileIlkClipModel()
		test_data.AssignIlkID(expectedModel, db)
		test_data.AssignAddressID(test_data.DogFileClipEventLog, expectedModel, db)
		assignClipAddressID(test_data.ClipAddress, expectedModel, db)
		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})

})
