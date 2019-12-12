package pot_drip_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/pot_drip"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pot drip converter", func() {
	var (
		converter = pot_drip.Converter{}
		db        = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("returns err if log missing topics", func() {
		badLog := core.HeaderSyncLog{}

		_, err := converter.ToModels(constants.PotABI(), []core.HeaderSyncLog{badLog}, db)
		Expect(err).To(HaveOccurred())
	})

	It("converts a log to model", func() {
		model, err := converter.ToModels(constants.PotABI(), []core.HeaderSyncLog{test_data.PotDripHeaderSyncLog}, db)

		Expect(err).NotTo(HaveOccurred())
		var addrID int64
		addrErr := db.Get(&addrID, `SELECT id FROM addresses WHERE address = $1`, common.HexToAddress(test_data.PotDripHeaderSyncLog.Log.Topics[1].Hex()).Hex())
		Expect(addrErr).NotTo(HaveOccurred())
		expectedModel := test_data.PotDripModel()
		expectedModel.ColumnValues[constants.MsgSenderColumn] = addrID
		Expect(model).To(Equal([]event.InsertionModel{expectedModel}))
	})
})
