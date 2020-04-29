package integration_tests

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_delete"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// TODO: test against real log from mainnet if detected
var _ = XDescribe("LogSortedOffer Transformer", func() {
	logDeleteConfig := event.TransformerConfig{
		TransformerName:   constants.LogDeleteTable,
		ContractAddresses: test_data.OasisAddresses(),
		ContractAbi:       constants.OasisABI(),
		Topic:             constants.LogDeleteSignature(),
	}

	It("transforms and persists a LogDelete event", func() {
		log := test_data.LogDeleteEventLog.Log
		blockNumber := int64(log.BlockNumber)
		logDeleteConfig.StartingBlockNumber = blockNumber
		logDeleteConfig.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      logDeleteConfig,
			Transformer: log_delete.Transformer{},
		}
		transformer := initializer.NewTransformer(db)

		eventLogs := test_data.CreateLogs(header.Id, []types.Log{log}, db)

		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult logDeleteModel
		err = db.Get(&dbResult, `SELECT keeper, offer_id FROM maker.log_delete`)
		Expect(err).NotTo(HaveOccurred())

		keeperID, keeperErr := shared.GetOrCreateAddress(test_data.LogDeleteKeeperAddress.Hex(), db)
		Expect(keeperErr).NotTo(HaveOccurred())

		Expect(dbResult.Keeper).To(Equal(keeperID))
		Expect(dbResult.OfferID).To(Equal(test_data.LogDeleteModel().ColumnValues[constants.OfferId]))
	})
})

type logDeleteModel struct {
	Keeper  int64
	OfferID string `db:"offer_id"`
}
