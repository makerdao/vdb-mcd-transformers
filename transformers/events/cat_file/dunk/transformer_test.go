package dunk_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/cat_file/dunk"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cat file dunk transformer", func() {
	var (
		db          = test_config.NewTestDB(test_config.NewTestNode())
		transformer = dunk.Transformer{}
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("returns an error if log is missing topics", func() {
		badLog := core.EventLog{
			Log: types.Log{
				Data: []byte{1, 1, 1, 1, 1},
			},
		}

		_, err := transformer.ToModels(constants.Cat110ABI(), []core.EventLog{badLog}, nil)
		Expect(err).To(HaveOccurred())

		It("returns err if log is missing data", func() {
			badLog := core.EventLog{
				Log: types.Log{
					Topics: []common.Hash{{}, {}, {}, {}},
				},
			}

			_, err := transformer.ToModels(constants.Cat110ABI(), []core.EventLog{badLog}, nil)
			Expect(err).To(HaveOccurred())
		})
	})
})
