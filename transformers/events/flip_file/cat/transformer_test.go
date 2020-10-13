package cat_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/flip_file/cat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Flip file cat transformer", func() {
	var (
		transformer = cat.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("converts a flip file cat to a model", func() {
		models, err := transformer.ToModels(constants.FlipV110ABI(), []core.EventLog{test_data.FlipFileCatEventLog}, db)
		Expect(err).NotTo(HaveOccurred())

		dataColumnID, dataColumnErr := repository.GetOrCreateAddress(db, "0xA950524441892A31ebddF91d3cEEFa04Bf454466")
		Expect(dataColumnErr).NotTo(HaveOccurred())

		expectedFlipFileCat := test_data.FlipFileCatModel()
		test_data.AssignAddressID(test_data.FlipFileCatEventLog, expectedFlipFileCat, db)
		test_data.AssignMessageSenderID(test_data.FlipFileCatEventLog, expectedFlipFileCat, db)
		expectedFlipFileCat.ColumnValues[constants.DataColumn] = dataColumnID

		Expect(models[0]).To(Equal(expectedFlipFileCat))
	})

	It("returns an err if the log is missing topics", func() {
		badLog := core.EventLog{
			Log: types.Log{
				Topics: []common.Hash{
					common.HexToHash("0xtest"),
				},
			},
		}
		_, err := transformer.ToModels(constants.FlipV100ABI(), []core.EventLog{badLog}, nil)
		Expect(err).To(HaveOccurred())
	})

})
