package vow_test

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/pot_file/vow"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pot file vow converter", func() {
	var (
		converter event.Converter
		db        *postgres.DB
	)

	BeforeEach(func() {
		converter = &vow.Converter{}
		db = test_config.NewTestDB(test_config.NewTestNode())
	})

	It("converts a pot file vow log to a model", func() {
		models, err := converter.ToModels(constants.PotABI(), []core.HeaderSyncLog{test_data.PotFileVowHeaderSyncLog}, db)

		Expect(err).NotTo(HaveOccurred())
		expectedModel := test_data.PotFileVowModel()
		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})

	It("returns err if log is missing topics", func() {
		badLog := core.HeaderSyncLog{
			Log: types.Log{
				Data: []byte{1, 1, 1, 1, 1},
			},
		}

		_, err := converter.ToModels(constants.PotABI(), []core.HeaderSyncLog{badLog}, db)
		Expect(err).To(HaveOccurred())
	})
})
