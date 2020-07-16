package integration_tests

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_buy_enabled"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogBuyEnabled Transformer", func() {
	logBuyEnabledConfig := event.TransformerConfig{
		TransformerName:   constants.LogBuyEnabledTable,
		ContractAddresses: test_data.OasisAddresses(),
		ContractAbi:       constants.OasisABI(),
		Topic:             constants.LogBuyEnabledSignature(),
	}
	//TODO: Update test with blocknumber and event when available on mainnet
	XIt("fetches and transforms a LogBuyEnabled event for OASIS_MATCHING_MARKET_ONE contract", func() {
		blockNumber := int64(9613377)
		logBuyEnabledConfig.StartingBlockNumber = blockNumber
		logBuyEnabledConfig.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      logBuyEnabledConfig,
			Transformer: log_buy_enabled.Transformer{},
		}
		transformer := initializer.NewTransformer(db)

		oasis_one_address := constants.GetContractAddress("OASIS_MATCHING_MARKET_ONE")
		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(oasis_one_address)},
			[]common.Hash{common.HexToHash(logBuyEnabledConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []logBuyEnabledModel
		err = db.Select(&dbResult, `SELECT is_enabled from maker.log_buy_enabled`)
		Expect(err).NotTo(HaveOccurred())

		Expect(dbResult).To(Equal(true)) //TODO: update with test
	})
})

type logBuyEnabledModel struct {
	isEnabled bool `db:"is_Enabled"`
	HeaderID  int64
	LogID     int64 `db:"log_id"`
}
