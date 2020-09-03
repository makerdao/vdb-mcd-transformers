package integration_tests

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/auth"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Deny transformer", func() {
	Context("Cat deny events", func() {
		usrAddress := "0xa9ee75d81d78c36c4163004e6cc7a988eec9433e"
		denyIntegrationTest(int64(8928165), test_data.Cat100Address(), usrAddress, usrAddress)
	})

	Context("Flap v1.0.0 deny events", func() {
		usrAddress := "0xd27a5f3416d8791fc238c148c93630d9e3c882e5"
		denyIntegrationTest(int64(8928163), test_data.FlapV100Address(), usrAddress, usrAddress)
	})

	Context("Flap v1.0.9 deny events", func() {
		usrAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		denyIntegrationTest(int64(10510886), test_data.FlapV109Address(), usrAddress, usrAddress)
	})

	Context("Flip BAT v1.0.0 deny events", func() {
		usrAddress := "0xbaa65281c2FA2baAcb2cb550BA051525A480D3F4"
		denyIntegrationTest(int64(8928267), test_data.FlipBatV100Address(), usrAddress, usrAddress)
	})

	Context("Flip BAT v1.0.9 deny events", func() {
		usrAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		denyIntegrationTest(int64(10510886), test_data.FlipBatV109Address(), usrAddress, usrAddress)
	})

	Context("Flip ETH v1.0.0 deny events", func() {
		usrAddress := "0xBAB4FbeA257ABBfe84F4588d4Eedc43656E46Fc5"
		denyIntegrationTest(int64(8928180), test_data.FlipEthV100Address(), usrAddress, usrAddress)
	})

	Context("Flip ETH v1.0.9 deny events", func() {
		usrAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		denyIntegrationTest(int64(10510886), test_data.FlipEthV109Address(), usrAddress, usrAddress)
	})

	Context("Flip KNC-A v1.0.8 deny events", func() {
		usrAddress := "0xBAB4FbeA257ABBfe84F4588d4Eedc43656E46Fc5"
		denyIntegrationTest(int64(10323245), test_data.FlipKncAV108Address(), usrAddress, usrAddress)
	})

	Context("Flip KNC-A v1.0.9 deny events", func() {
		usrAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		denyIntegrationTest(int64(10510886), test_data.FlipKncAV109Address(), usrAddress, usrAddress)
	})

	Context("Flip MANA-A v1.0.9 deny events", func() {
		usrAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		denyIntegrationTest(int64(10510886), test_data.FlipManaAV109Address(), usrAddress, usrAddress)
	})

	Context("Flip TUSD-A v1.0.7 deny events", func() {
		usrAddress := "0xBAB4FbeA257ABBfe84F4588d4Eedc43656E46Fc5"
		denyIntegrationTest(int64(10144451), test_data.FlipTusdAV107Address(), usrAddress, usrAddress)
	})

	Context("Flip TUSD-A v1.0.9 deny events", func() {
		usrAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		denyIntegrationTest(int64(10510886), test_data.FlipTusdAV109Address(), usrAddress, usrAddress)
	})

	Context("Flip USDC-A v1.0.4 deny events", func() {
		usrAddress := "0xBAB4FbeA257ABBfe84F4588d4Eedc43656E46Fc5"
		denyIntegrationTest(int64(9686502), test_data.FlipUsdcAV104Address(), usrAddress, usrAddress)
	})

	Context("Flip USDC-A v1.0.9 deny events", func() {
		usrAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		denyIntegrationTest(int64(10510886), test_data.FlipUsdcAV109Address(), usrAddress, usrAddress)
	})

	Context("Flip USDC-B v1.0.7 deny events", func() {
		usrAddress := "0xBAB4FbeA257ABBfe84F4588d4Eedc43656E46Fc5"
		denyIntegrationTest(int64(10144450), test_data.FlipUsdcBV107Address(), usrAddress, usrAddress)
	})

	Context("Flip USDC-B v1.0.9 deny events", func() {
		usrAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		denyIntegrationTest(int64(10510886), test_data.FlipUsdcBV109Address(), usrAddress, usrAddress)
	})

	Context("Flip WBTC v1.0.6 deny events", func() {
		usrAddress := "0xBAB4FbeA257ABBfe84F4588d4Eedc43656E46Fc5"
		denyIntegrationTest(int64(9975625), test_data.FlipWbtcAV106Address(), usrAddress, usrAddress)
	})

	Context("Flip WBTC v1.0.9 deny events", func() {
		usrAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		denyIntegrationTest(int64(10510886), test_data.FlipWbtcAV109Address(), usrAddress, usrAddress)
	})

	Context("Flip ZRX-A v1.0.8 deny events", func() {
		usrAddress := "0xBAB4FbeA257ABBfe84F4588d4Eedc43656E46Fc5"
		denyIntegrationTest(int64(10323245), test_data.FlipZrxAV108Address(), usrAddress, usrAddress)
	})

	Context("Flip ZRX-A v1.0.9 deny events", func() {
		usrAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		denyIntegrationTest(int64(10510886), test_data.FlipZrxAV109Address(), usrAddress, usrAddress)
	})

	Context("Flop v1.0.1 deny events", func() {
		usrAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		denyIntegrationTest(int64(9008144), test_data.FlopV101Address(), usrAddress, usrAddress)
	})

	Context("Flop v1.0.9 deny events", func() {
		usrAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		denyIntegrationTest(int64(10510886), test_data.FlopV109Address(), usrAddress, usrAddress)
	})

	Context("Jug deny events", func() {
		usrAddress := "0x45f0a929889ec8cc2d5b8cd79ab55e3279945cde"
		denyIntegrationTest(int64(8928160), test_data.JugAddress(), usrAddress, usrAddress)
	})

	Context("Median BAT deny events", func() {
		usrAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		denyIntegrationTest(int64(8957024), test_data.MedianBatAddress(), usrAddress, usrAddress)
	})

	Context("Median ETH deny events", func() {
		usrAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		denyIntegrationTest(int64(8957020), test_data.MedianEthAddress(), usrAddress, usrAddress)
	})

	Context("Median KNC deny events", func() {
		usrAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		denyIntegrationTest(int64(10350821), test_data.MedianKncAddress(), usrAddress, usrAddress)
	})

	Context("Median MANA deny events", func() {
		usrAddress := "0x0048d6225D1F3eA4385627eFDC5B4709Cab4A21c"
		denyIntegrationTest(int64(10542736), test_data.MedianManaAddress(), usrAddress, usrAddress)
	})

	Context("Median WBTC deny events", func() {
		usrAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		denyIntegrationTest(int64(8957027), test_data.MedianWbtcAddress(), usrAddress, usrAddress)
	})

	Context("Median ZRX deny events", func() {
		usrAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		denyIntegrationTest(int64(10350821), test_data.MedianZrxAddress(), usrAddress, usrAddress)
	})

	Context("OSM BAT deny events", func() {
		usrAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		denyIntegrationTest(int64(8957031), test_data.OsmBatAddress(), usrAddress, usrAddress)
	})

	Context("OSM ETH deny events", func() {
		usrAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		denyIntegrationTest(int64(8957029), test_data.OsmEthAddress(), usrAddress, usrAddress)
	})

	Context("OSM KNC deny events", func() {
		usrAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		denyIntegrationTest(int64(10350821), test_data.OsmKncAddress(), usrAddress, usrAddress)
	})

	Context("OSM MANA deny events", func() {
		usrAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		denyIntegrationTest(int64(10524041), test_data.OsmManaAddress(), usrAddress, usrAddress)
	})

	Context("OSM WBTC deny events", func() {
		usrAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		denyIntegrationTest(int64(9975543), test_data.OsmWbtcAddress(), usrAddress, usrAddress)
	})

	Context("OSM ZRX deny events", func() {
		usrAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		denyIntegrationTest(int64(10350821), test_data.OsmZrxAddress(), usrAddress, usrAddress)
	})

	Context("Pot deny events", func() {
		usrAddress := "0x1a5ee7c64cf874c735968e3a42fa13f1c03427f9"
		denyIntegrationTest(int64(8928160), test_data.PotAddress(), usrAddress, usrAddress)
	})

	Context("Spot deny events", func() {
		usrAddress := "0xdedd12bcb045c02b2fe11031c2b269bcde457410"
		denyIntegrationTest(int64(8928152), test_data.SpotAddress(), usrAddress, usrAddress)
	})

	Context("Vow deny events", func() {
		usrAddress := "0x68322ca1a9aeb8c1d610b5fc8a8920aa0fba423b"
		denyIntegrationTest(int64(8928163), test_data.VowAddress(), usrAddress, usrAddress)
	})
})

func denyIntegrationTest(blockNumber int64, contractAddressHex, msgSenderAddressHex, usrAddressHex string) {
	It("persists event", func() {
		test_config.CleanTestDB(db)
		logFetcher := fetcher.NewLogFetcher(blockChain)
		denyConfig := event.TransformerConfig{
			TransformerName: constants.DenyTable,
			Topic:           constants.DenySignature(),
		}
		initializer := event.ConfiguredTransformer{
			Config:      denyConfig,
			Transformer: auth.Transformer{TableName: constants.DenyTable},
		}

		header, headerErr := persistHeader(db, blockNumber, blockChain)
		Expect(headerErr).NotTo(HaveOccurred())

		initializer.Config.ContractAddresses = []string{contractAddressHex}
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		address := common.HexToAddress(contractAddressHex)
		topics := []common.Hash{common.HexToHash(denyConfig.Topic)}

		logs, logsErr := logFetcher.FetchLogs([]common.Address{address}, topics, header)
		Expect(logsErr).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		denyTransformer := initializer.NewTransformer(db)
		transformErr := denyTransformer.Execute(eventLogs)
		Expect(transformErr).NotTo(HaveOccurred())

		var dbResult []denyModel
		err := db.Select(&dbResult, `SELECT address_id, msg_sender, usr FROM maker.deny`)
		Expect(err).NotTo(HaveOccurred())

		contractAddressID, contractAddressErr := shared.GetOrCreateAddress(contractAddressHex, db)
		Expect(contractAddressErr).NotTo(HaveOccurred())
		msgSenderAddressID, msgSenderAddressErr := shared.GetOrCreateAddress(msgSenderAddressHex, db)
		Expect(msgSenderAddressErr).NotTo(HaveOccurred())
		usrAddressID, usrAddressErr := shared.GetOrCreateAddress(usrAddressHex, db)
		Expect(usrAddressErr).NotTo(HaveOccurred())

		var matchFound bool
		for _, result := range dbResult {
			if result.AddressID == contractAddressID &&
				result.MsgSender == msgSenderAddressID &&
				result.Usr == usrAddressID {
				matchFound = true
			}
		}

		Expect(matchFound).To(BeTrue(), getDenyFailureMessage(contractAddressHex, blockNumber))
	})
}

type denyModel struct {
	Usr       int64 `db:"usr"`
	MsgSender int64 `db:"msg_sender"`
	AddressID int64 `db:"address_id"`
}

func getDenyFailureMessage(contractAddress string, blockNumber int64) string {
	failureMsgToFmt := "no matching deny event found for contract %s at block %d"
	return fmt.Sprintf(failureMsgToFmt, contractAddress, blockNumber)
}
