package integration_tests

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vow_heal"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VowHeal Transformer", func() {
	vowHealConfig := event.TransformerConfig{
		TransformerName:   constants.VowHealTable,
		ContractAddresses: []string{test_data.VowAddress()},
		ContractAbi:       constants.VowABI(),
		Topic:             constants.VowHealSignature(),
	}

	It("transforms VowHeal log events", func() {
		blockNumber := int64(9724446)
		vowHealConfig.StartingBlockNumber = blockNumber
		vowHealConfig.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)
		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			event.HexStringsToAddresses(vowHealConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(vowHealConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		tr := event.ConfiguredTransformer{
			Config:      vowHealConfig,
			Transformer: vow_heal.Transformer{},
		}.NewTransformer(db)

		err = tr.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult vowHealModel
		err = db.Get(&dbResult, `SELECT msg_sender, rad from maker.vow_heal`)
		Expect(err).NotTo(HaveOccurred())

		msgSender, msgSenderErr := shared.GetOrCreateAddress("0x233a1b1A5381D7EAa3e9b373C392aB48A47bA596", db)
		Expect(msgSenderErr).NotTo(HaveOccurred())
		expectedModel := vowHealModel{
			MsgSender: msgSender,
			Rad:       "50000000000000000000000000000000000000000000000000",
		}

		Expect(dbResult).To(Equal(expectedModel))
	})
})

type vowHealModel struct {
	MsgSender int64 `db:"msg_sender"`
	Rad       string
}
