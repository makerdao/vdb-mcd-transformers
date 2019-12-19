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
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PotExit Transformer", func() {
	potExitConfig := transformer.EventTransformerConfig{
		TransformerName:   constants.PotExitTable,
		ContractAddresses: []string{test_data.PotAddress()},
		ContractAbi:       constants.PotABI(),
		Topic:             constants.PotExitSignature(),
	}

	It("transforms PotExit log events", func() {
		blockNumber := int64(15506076)
		potExitConfig.StartingBlockNumber = blockNumber
		potExitConfig.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)
		header, headerErr := persistHeader(db, blockNumber, blockChain)
		Expect(headerErr).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, fetchErr := logFetcher.FetchLogs(
			transformer.HexStringsToAddresses(potExitConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(potExitConfig.Topic)},
			header)
		Expect(fetchErr).NotTo(HaveOccurred())

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		tr := event.Transformer{
			Config:    potExitConfig,
			Converter: pot_exit.Converter{},
		}.NewTransformer(db)

		transformErr := tr.Execute(headerSyncLogs)
		Expect(transformErr).NotTo(HaveOccurred())

		var dbResult potExitModel
		queryErr := db.Get(&dbResult, `SELECT msg_sender, wad from maker.pot_exit`)
		Expect(queryErr).NotTo(HaveOccurred())

		addressID, addressErr := shared.GetOrCreateAddress("0xe7bc397dbd069fc7d0109c0636d06888bb50668c", db)
		Expect(addressErr).NotTo(HaveOccurred())
		Expect(dbResult.MsgSender).To(Equal(strconv.FormatInt(addressID, 10)))
		Expect(dbResult.Wad).To(Equal("22957121481043076331"))
	})
})

type potExitModel struct {
	MsgSender string `db:"msg_sender"`
	Wad       string
}
