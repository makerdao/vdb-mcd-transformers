package integration_tests

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/auction_file"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auction file transformer", func() {
	Context("Flap file events", func() {
		msgSender := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		what := "beg"
		data := "1020000000000000000"
		auctionFileIntegrationTest(int64(9529100), test_data.FlapV100Address(), msgSender, what, data)
	})

	Context("Flip BAT file events", func() {
		msgSender := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		what := "tau"
		data := "259200"
		auctionFileIntegrationTest(int64(8928412), test_data.FlipBatV100Address(), msgSender, what, data)
	})

	Context("Flip ETH file events", func() {
		msgSender := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		what := "ttl"
		data := "600"
		auctionFileIntegrationTest(int64(8928402), test_data.FlipEthAV100Address(), msgSender, what, data)
	})

	Context("Flip GUSD file events", func() {
		msgSender := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		what := "beg"
		data := "1030000000000000000"
		auctionFileIntegrationTest(int64(11314893), test_data.FlipGusdAV115Address(), msgSender, what, data)
	})

	Context("Flip RENBTC file events", func() {
		msgSender := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		what := "tau"
		data := "21600"
		auctionFileIntegrationTest(int64(11451553), test_data.FlipRenbtcA121Address(), msgSender, what, data)
	})

	Context("Flip TUSD file events", func() {
		msgSender := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		what := "ttl"
		data := "21600"
		auctionFileIntegrationTest(int64(10201136), test_data.FlipTusdAV107Address(), msgSender, what, data)
	})

	Context("Flip UNI-A file events", func() {
		msgSender := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		what := "tau"
		data := "21600"
		auctionFileIntegrationTest(int64(11451553), test_data.FlipUniAV121Address(), msgSender, what, data)
	})

	Context("Flip USDC-A file events", func() {
		msgSender := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		what := "beg"
		data := "1030000000000000000"
		auctionFileIntegrationTest(int64(9686502), test_data.FlipUsdcAV104Address(), msgSender, what, data)
	})

	Context("Flip USDC-B file events", func() {
		msgSender := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		what := "beg"
		data := "1030000000000000000"
		auctionFileIntegrationTest(int64(10201136), test_data.FlipUsdcBV107Address(), msgSender, what, data)
	})

	Context("Flip WBTC file events", func() {
		msgSender := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		what := "beg"
		data := "1030000000000000000"
		auctionFileIntegrationTest(int64(9990976), test_data.FlipWbtcAV106Address(), msgSender, what, data)
	})

	Context("Flop file events", func() {
		msgSender := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		what := "pad"
		data := "1200000000000000000"
		auctionFileIntegrationTest(int64(9017707), test_data.FlopV101Address(), msgSender, what, data)
	})
})

func auctionFileIntegrationTest(blockNumber int64, contractAddressHex, msgSenderAddressHex, what, data string) {
	It("persists event", func() {
		test_config.CleanTestDB(db)
		auctionFileConfig := event.TransformerConfig{
			ContractAddresses:   []string{contractAddressHex},
			EndingBlockNumber:   blockNumber,
			StartingBlockNumber: blockNumber,
			TransformerName:     constants.AuctionFileTable,
			Topic:               constants.AuctionFileSignature(),
		}

		header, headerErr := persistHeader(db, blockNumber, blockChain)
		Expect(headerErr).NotTo(HaveOccurred())

		address := common.HexToAddress(contractAddressHex)
		topics := []common.Hash{common.HexToHash(auctionFileConfig.Topic)}

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, logsErr := logFetcher.FetchLogs([]common.Address{address}, topics, header)
		Expect(logsErr).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		initializer := event.ConfiguredTransformer{
			Config:      auctionFileConfig,
			Transformer: auction_file.Transformer{},
		}
		denyTransformer := initializer.NewTransformer(db)
		transformErr := denyTransformer.Execute(eventLogs)
		Expect(transformErr).NotTo(HaveOccurred())

		var dbResults []auctionFileModel
		err := db.Select(&dbResults, `SELECT address_id, msg_sender, what, data FROM maker.auction_file`)
		Expect(err).NotTo(HaveOccurred())

		contractAddressID, contractAddressErr := shared.GetOrCreateAddress(contractAddressHex, db)
		Expect(contractAddressErr).NotTo(HaveOccurred())
		msgSenderAddressID, msgSenderAddressErr := shared.GetOrCreateAddress(msgSenderAddressHex, db)
		Expect(msgSenderAddressErr).NotTo(HaveOccurred())

		var matchFound bool
		for _, result := range dbResults {
			if result.AddressID == contractAddressID &&
				result.MsgSender == msgSenderAddressID &&
				result.What == what &&
				result.Data == data {
				matchFound = true
			}
		}
		Expect(matchFound).To(BeTrue(), getAuctionFileFailureMessage(contractAddressHex, blockNumber))
	})
}

type auctionFileModel struct {
	AddressID int64 `db:"address_id"`
	MsgSender int64 `db:"msg_sender"`
	What      string
	Data      string
}

func getAuctionFileFailureMessage(contractAddress string, blockNumber int64) string {
	failureMsgToFmt := "no matching auction file event found for contract %s at block %d"
	return fmt.Sprintf(failureMsgToFmt, contractAddress, blockNumber)
}
