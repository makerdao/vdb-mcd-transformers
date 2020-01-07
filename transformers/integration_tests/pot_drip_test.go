package integration_tests

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/pot_drip"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PotDrip Transformer", func() {
	var potDripConfig transformer.EventTransformerConfig

	BeforeEach(func() {
		test_config.CleanTestDB(db)

		potDripConfig = transformer.EventTransformerConfig{
			ContractAddresses: []string{test_data.PotAddress()},
			ContractAbi:       constants.PotABI(),
			Topic:             constants.PotDripSignature(),
		}
	})

	It("transforms PotDrip log events", func() {
		blockNumber := int64(15407546)
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
			transformer.HexStringsToAddresses(potDripConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(potDripConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		err = tr.Execute(headerSyncLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult potDripModel
		err = db.Get(&dbResult, `SELECT msg_sender from maker.pot_drip`)
		Expect(err).NotTo(HaveOccurred())

		addrID, addrErr := shared.GetOrCreateAddress("0x87e76b0a50efc20259cafE0530f75aE0e816aaF2", db)
		Expect(addrErr).NotTo(HaveOccurred())
		Expect(dbResult.MsgSender).To(Equal(strconv.FormatInt(addrID, 10)))
	})
})

type potDripModel struct {
	MsgSender string `db:"msg_sender"`
}
