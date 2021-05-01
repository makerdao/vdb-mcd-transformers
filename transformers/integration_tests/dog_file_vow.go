package integration_tests

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/dog_file/vow"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func fetchDogFileVowLogsFromChain(config event.TransformerConfig, header core.Header) []core.EventLog {
	logFetcher := fetcher.NewLogFetcher(blockChain)
	logs, err := logFetcher.FetchLogs(
		[]common.Address{common.HexToAddress(config.ContractAddresses[0])},
		[]common.Hash{common.HexToHash(config.Topic)},
		header)
	Expect(err).NotTo(HaveOccurred())

	return test_data.CreateLogs(header.Id, logs, db)
}

var _ = Describe("Dog File Vow Transformer", func() {
	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	XIt("fetches and transforms a Dog File Vow event", func() {
		blockNumber := int64(1)

		dogFileVowConfig := event.TransformerConfig{
			TransformerName:     constants.DogFileVowTable,
			ContractAddresses:   []string{test_data.Dog130Address()},
			ContractAbi:         constants.DogABI(),
			Topic:               constants.DogFileVowSignature(),
			StartingBlockNumber: blockNumber,
			EndingBlockNumber:   blockNumber,
		}

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      dogFileVowConfig,
			Transformer: vow.Transformer{},
		}
		transformer := initializer.NewTransformer(db)

		eventLogs := fetchDogFileVowLogsFromChain(dogFileVowConfig, header)

		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())
	})
})
