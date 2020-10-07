package integration_tests

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/pot_drip"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PotDrip Transformer", func() {
	var potDripConfig event.TransformerConfig

	BeforeEach(func() {
		test_config.CleanTestDB(db)

		potDripConfig = event.TransformerConfig{
			ContractAddresses: []string{test_data.PotAddress()},
			ContractAbi:       constants.PotABI(),
			Topic:             constants.PotDripSignature(),
		}
	})

	It("transforms PotDrip log events", func() {
		blockNumber := int64(9127348)
		potDripConfig.StartingBlockNumber = blockNumber
		potDripConfig.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      potDripConfig,
			Transformer: pot_drip.Transformer{},
		}
		tr := initializer.NewTransformer(db)

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			event.HexStringsToAddresses(potDripConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(potDripConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		err = tr.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var msgSender int64
		err = db.Get(&msgSender, `SELECT msg_sender from maker.pot_drip`)
		Expect(err).NotTo(HaveOccurred())

		expectedAddrID, addrErr := shared.GetOrCreateAddress("0x825100c63933cABA16C8CE40814DAc88305D8810", db)
		Expect(addrErr).NotTo(HaveOccurred())
		Expect(msgSender).To(Equal(expectedAddrID))
	})
})
