package integration_tests

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/clip_redo"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func fetchClipRedoLogsFromChain(config event.TransformerConfig, header core.Header) []core.EventLog {
	//TODO: when there are real Clip Redo events on chain, use the following code to fetch real events from the chain

	//logFetcher := fetcher.NewLogFetcher(blockChain)
	//logs, err := logFetcher.FetchLogs(
	//	[]common.Address{common.HexToAddress(config.ContractAddresses[0])},
	//	[]common.Hash{common.HexToHash(config.Topic)},
	//	header)
	//Expect(err).NotTo(HaveOccurred())

	logs := []types.Log{test_data.RawClipRedoLog}
	return test_data.CreateLogs(header.Id, logs, db)
}

var _ = Describe("Clip Redo Transformer", func() {
	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	XIt("fetches and transforms a Clip Redo event", func() {
		blockNumber := int64(1)

		clipRedoConfig := event.TransformerConfig{
			TransformerName:     constants.ClipRedoTable,
			ContractAddresses:   []string{test_data.ClipLinkAV130Address()},
			ContractAbi:         constants.ClipV150ABI(),
			Topic:               constants.ClipRedoSignature(),
			StartingBlockNumber: blockNumber,
			EndingBlockNumber:   blockNumber,
		}

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      clipRedoConfig,
			Transformer: clip_redo.Transformer{},
		}
		transformer := initializer.NewTransformer(db)

		eventLogs := fetchClipRedoLogsFromChain(clipRedoConfig, header)

		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())
	})
})
