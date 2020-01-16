package integration_tests

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/pot_cage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//TODO: Need tu update when there are events on mainnet
var _ = XDescribe("PotCage EventTransformer", func() {
	var cageDeploymentPotAddress = "0x52ca216f93836eea1ee605cf6aa41127134b9754"

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	Describe("Pot cage", func() {
		var (
			addresses   []common.Address
			blockNumber int64
			header      core.Header
			logs        []types.Log
			topics      []common.Hash
		)

		BeforeEach(func() {
			blockNumber = int64(14681713)
			var insertHeaderErr error
			header, insertHeaderErr = persistHeader(db, blockNumber, blockChain)
			Expect(insertHeaderErr).NotTo(HaveOccurred())

			potCageConfig := transformer.EventTransformerConfig{
				TransformerName:     constants.PotCageTable,
				ContractAddresses:   []string{cageDeploymentPotAddress},
				ContractAbi:         constants.PotABI(),
				Topic:               constants.PotCageSignature(),
				StartingBlockNumber: blockNumber,
				EndingBlockNumber:   blockNumber,
			}

			addresses = transformer.HexStringsToAddresses(potCageConfig.ContractAddresses)
			topics = []common.Hash{common.HexToHash(potCageConfig.Topic)}

			initializer := event.ConfiguredTransformer{
				Config:      potCageConfig,
				Transformer: pot_cage.Transformer{},
			}

			logFetcher := fetcher.NewLogFetcher(blockChain)
			var fetcherErr error
			logs, fetcherErr = logFetcher.FetchLogs(addresses, topics, header)
			Expect(fetcherErr).NotTo(HaveOccurred())

			eventLogs := test_data.CreateLogs(header.Id, logs, db)

			tr := initializer.NewTransformer(db)
			executeErr := tr.Execute(eventLogs)
			Expect(executeErr).NotTo(HaveOccurred())
		})

		It("fetches and transforms a Pot.cage event from Kovan", func() {
			var dbResult []potCageModel
			getFileErr := db.Select(&dbResult, `SELECT id FROM maker.pot_cage`)

			Expect(getFileErr).NotTo(HaveOccurred())
			Expect(len(dbResult)).To(Equal(1))
		})
	})
})

type potCageModel struct {
	Id string
}
