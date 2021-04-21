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

func fetchDogDenyLogsFromChain(config event.TransformerConfig, header core.Header) []core.EventLog {
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

		dogDenyConfig := event.TransformerConfig{
			TransformerName:     constants.DogDenyTable,
			ContractAddresses:   []string{test_data.Dog130Address()},
			ContractAbi:         constants.DogABI(),
			Topic:               constants.DogDenySignature(),
			StartingBlockNumber: blockNumber,
			EndingBlockNumber:   blockNumber,
		}

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      dogDenyConfig,
			Transformer: dog_rely.Transformer{},
		}
		transformer := initializer.NewTransformer(db)

		eventLogs := fetchDogDenyLogsFromChain(dogDenyConfig, header)

		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())
	})
})
