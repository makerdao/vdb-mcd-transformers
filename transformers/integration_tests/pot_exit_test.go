package integration_tests

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/pot_exit"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PotExit Transformer", func() {
	potExitConfig := event.TransformerConfig{
		TransformerName:   constants.PotExitTable,
		ContractAddresses: []string{test_data.PotAddress()},
		ContractAbi:       constants.PotABI(),
		Topic:             constants.PotExitSignature(),
	}

	It("transforms PotExit log events", func() {
		blockNumber := int64(9132846)
		potExitConfig.StartingBlockNumber = blockNumber
		potExitConfig.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)
		header, headerErr := persistHeader(db, blockNumber, blockChain)
		Expect(headerErr).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, fetchErr := logFetcher.FetchLogs(
			event.HexStringsToAddresses(potExitConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(potExitConfig.Topic)},
			header)
		Expect(fetchErr).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		tr := event.ConfiguredTransformer{
			Config:      potExitConfig,
			Transformer: pot_exit.Transformer{},
		}.NewTransformer(db)

		transformErr := tr.Execute(eventLogs)
		Expect(transformErr).NotTo(HaveOccurred())

		var dbResult potExitModel
		queryErr := db.Get(&dbResult, `SELECT msg_sender, wad from maker.pot_exit`)
		Expect(queryErr).NotTo(HaveOccurred())

		addressID, addressErr := shared.GetOrCreateAddress("0x2D8990618391aD152c336649D27F164b2618bf60", db)
		Expect(addressErr).NotTo(HaveOccurred())
		Expect(dbResult.MsgSender).To(Equal(strconv.FormatInt(addressID, 10)))
		Expect(dbResult.Wad).To(Equal("200473120597989120478"))
	})
})

type potExitModel struct {
	MsgSender string `db:"msg_sender"`
	Wad       string
}
