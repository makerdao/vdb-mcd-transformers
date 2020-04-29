package integration_tests

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_matching_enabled"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogMatchingEnabled Transformer", func() {
	LogMatchingEnabledConfig := event.TransformerConfig{
		TransformerName:   constants.LogMatchingEnabledTable,
		ContractAddresses: test_data.OasisAddresses(),
		ContractAbi:       constants.OasisABI(),
		Topic:             constants.LogMatchingEnabledSignature(),
	}
	//TODO: Update test with blocknumber and event when available on mainnet
	XIt("fetches and transforms a LogMatchingEnabled event for OASIS_MATCHING_MARKET_ONE contract", func() {
		blockNumber := int64(9613577)
		LogMatchingEnabledConfig.StartingBlockNumber = blockNumber
		LogMatchingEnabledConfig.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      LogMatchingEnabledConfig,
			Transformer: log_matching_enabled.Transformer{},
		}
		transformer := initializer.NewTransformer(db)

		oasis_one_address := constants.GetContractAddress("OASIS_MATCHING_MARKET_ONE")
		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(oasis_one_address)},
			[]common.Hash{common.HexToHash(LogMatchingEnabledConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []LogMatchingEnabledModel
		err = db.Select(&dbResult, `SELECT is_enabled from maker.log_matching_enabled`)
		Expect(err).NotTo(HaveOccurred())

		Expect(dbResult).To(Equal(true)) //TODO: update with test
	})
})

type LogMatchingEnabledModel struct {
	isEnabled bool `db:"is_Enabled"`
	HeaderID  int64
	LogID     int64 `db:"log_id"`
}
