package integration_tests

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/flip_file/cat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FlipFileCat Transformer", func() {
	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	flipFileCatConfig := event.TransformerConfig{
		TransformerName:   constants.FlipFileCatTable,
		ContractAddresses: test_data.FlipV110Addresses(),
		ContractAbi:       constants.FlipV110ABI(),
		Topic:             constants.FlipFileCatSignature(),
	}

	It("fetches and transforms a Flip File Cat event", func() {
		//TODO: Needs real integration test data when it shows up on chain

		blockNumber := int64(10769102)
		flipFileCatConfig.StartingBlockNumber = blockNumber
		flipFileCatConfig.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		tr := event.ConfiguredTransformer{
			Config:      flipFileCatConfig,
			Transformer: cat.Transformer{},
		}.NewTransformer(db)

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			event.HexStringsToAddresses(flipFileCatConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(flipFileCatConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		err = tr.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())
	})
})
