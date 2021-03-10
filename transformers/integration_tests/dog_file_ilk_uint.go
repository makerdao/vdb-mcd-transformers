package integration_tests

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/dog_file/ilk_uint"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func fetchDogFileIlkUintLogsFromChain(config event.TransformerConfig, header core.Header) []core.EventLog {
	//TODO: when there are real Dog File Uint events on chain, use the following code to fetch real events from the chain

	//logFetcher := fetcher.NewLogFetcher(blockChain)
	//logs, err := logFetcher.FetchLogs(
	//	[]common.Address{common.HexToAddress(config.ContractAddresses[0])},
	//	[]common.Hash{common.HexToHash(config.Topic)},
	//	header)
	//Expect(err).NotTo(HaveOccurred())

	logs := []types.Log{test_data.RawDogBarkLog}
	return test_data.CreateLogs(header.Id, logs, db)
}

var _ = Describe("Dog File Ilk Uint Transformer", func() {
	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	XIt("fetches and transforms a Dog File Ilk Uint event", func() {
		blockNumber := int64(1) //TODO: update this when there are Dog File Ilk Uint events on the chain

		dogFileIlkUintConfig := event.TransformerConfig{
			TransformerName:     constants.DogFileIlkUintTable,
			ContractAddresses:   []string{test_data.Dog1xxAddress()},
			ContractAbi:         constants.DogABI(),
			Topic:               constants.DogFileIlkUintSignature(),
			StartingBlockNumber: blockNumber,
			EndingBlockNumber:   blockNumber,
		}

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      dogFileIlkUintConfig,
			Transformer: ilk_uint.Transformer{},
		}
		transformer := initializer.NewTransformer(db)

		eventLogs := fetchDogFileIlkUintLogsFromChain(dogFileIlkUintConfig, header)

		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())
	})
})
