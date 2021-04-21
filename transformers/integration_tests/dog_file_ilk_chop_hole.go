package integration_tests

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/dog_file/ilk_chop_hole"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func fetchDogFileIlkChopHoleLogsFromChain(config event.TransformerConfig, header core.Header) []core.EventLog {
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

	XIt("fetches and transforms a Dog File Ilk Chop Hole event", func() {
		blockNumber := int64(1) //TODO: update this when there are Dog File Ilk Chop Hole events on the chain

		dogFileIlkChopHoleConfig := event.TransformerConfig{
			TransformerName:     constants.DogFileIlkChopHoleTable,
			ContractAddresses:   []string{test_data.Dog130Address()},
			ContractAbi:         constants.DogABI(),
			Topic:               constants.DogFileIlkChopHoleSignature(),
			StartingBlockNumber: blockNumber,
			EndingBlockNumber:   blockNumber,
		}

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      dogFileIlkChopHoleConfig,
			Transformer: ilk_chop_hole.Transformer{},
		}
		transformer := initializer.NewTransformer(db)

		eventLogs := fetchDogFileIlkChopHoleLogsFromChain(dogFileIlkChopHoleConfig, header)

		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())
	})
})
