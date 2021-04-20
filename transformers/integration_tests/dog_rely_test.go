package integration_tests

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/dog_rely"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func fetchDogRelyLogsFromChain(config event.TransformerConfig, header core.Header) []core.EventLog {
	logFetcher := fetcher.NewLogFetcher(blockChain)
	logs, err := logFetcher.FetchLogs(
		[]common.Address{common.HexToAddress(config.ContractAddresses[0])},
		[]common.Hash{common.HexToHash(config.Topic)},
		header)
	Expect(err).NotTo(HaveOccurred())

	return test_data.CreateLogs(header.Id, logs, db)
}

var _ = Describe("Dog Rely Transformer", func() {
	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("fetches and transforms a Dog Rely event", func() {
		blockNumber := int64(12246358)

		dogRelyConfig := event.TransformerConfig{
			TransformerName:     constants.DogRelyTable,
			ContractAddresses:   []string{test_data.Dog1xxAddress()},
			ContractAbi:         constants.DogABI(),
			Topic:               constants.DogRelySignature(),
			StartingBlockNumber: blockNumber,
			EndingBlockNumber:   blockNumber,
		}

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      dogRelyConfig,
			Transformer: dog_rely.Transformer{},
		}
		transformer := initializer.NewTransformer(db)

		eventLogs := fetchDogRelyLogsFromChain(dogRelyConfig, header)

		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())
	})
})
