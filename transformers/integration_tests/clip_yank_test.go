package integration_tests

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/dog_bark"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func fetchClipYankLogsFromChain(config event.TransformerConfig, header core.Header) []core.EventLog {
	//TODO: when there are real Clip Yank events on chain, use the following code to fetch real events from the chain

	//logFetcher := fetcher.NewLogFetcher(blockChain)
	//logs, err := logFetcher.FetchLogs(
	//	[]common.Address{common.HexToAddress(config.ContractAddresses[0])},
	//	[]common.Hash{common.HexToHash(config.Topic)},
	//	header)
	//Expect(err).NotTo(HaveOccurred())

	logs := []types.Log{test_data.RawClipYankLog}
	return test_data.CreateLogs(header.Id, logs, db)
}

var _ = Describe("Clip Yank Transformer", func() {
	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	XIt("fetches and transforms a Clip Yank event", func() {
		blockNumber := int64(1) //TODO: update this when there are Clip Yank events on the chain

		clipYankConfig := event.TransformerConfig{
			TransformerName:     constants.ClipYankTable,
			ContractAddresses:   []string{test_data.Clip130Address()},
			ContractAbi:         constants.ClipABI(),
			Topic:               constants.ClipYankSignature(),
			StartingBlockNumber: blockNumber,
			EndingBlockNumber:   blockNumber,
		}

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      clipYankConfig,
			Transformer: dog_bark.Transformer{},
		}
		transformer := initializer.NewTransformer(db)

		eventLogs := fetchClipYankLogsFromChain(clipYankConfig, header)

		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())
	})
})
