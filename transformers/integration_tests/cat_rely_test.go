package integration_tests

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/cat_rely"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CatRely Transformer", func() {
	catRelyConfig := transformer.EventTransformerConfig{
		TransformerName:   constants.CatRelyTable,
		ContractAddresses: []string{test_data.CatAddress()},
		ContractAbi:       constants.CatABI(),
		Topic:             constants.CatRelySignature(),
	}

	It("transforms CatRely log events", func() {
		var (
			expectedID int64
			dbResult   catRelyModel
		)
		blockNumber := int64(14764546)
		catRelyConfig.StartingBlockNumber = blockNumber
		catRelyConfig.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)
		header, headerErr := persistHeader(db, blockNumber, blockChain)
		Expect(headerErr).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, fetchErr := logFetcher.FetchLogs(
			transformer.HexStringsToAddresses(catRelyConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(catRelyConfig.Topic)},
			header)
		Expect(fetchErr).NotTo(HaveOccurred())

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		tr := event.Transformer{
			Config:    catRelyConfig,
			Converter: cat_rely.Converter{},
		}.NewTransformer(db)

		transformErr := tr.Execute(headerSyncLogs)
		Expect(transformErr).NotTo(HaveOccurred())
		expectedID, addrErr := shared.GetOrCreateAddress(test_data.CatRelyHeaderSyncLog.Log.Topics[1].Hex(), db)
		Expect(addrErr).NotTo(HaveOccurred())

		queryErr := db.Get(&dbResult, `SELECT address_id from maker.cat_rely`)
		Expect(queryErr).NotTo(HaveOccurred())
		Expect(dbResult.AddressId).To(Equal(expectedID))
	})
})

type catRelyModel struct {
	AddressId int64 `db:"address_id"`
}
