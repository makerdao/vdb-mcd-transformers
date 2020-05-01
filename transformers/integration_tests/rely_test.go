package integration_tests

import (
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

var _ = Describe("Rely transformer", func() {
	Context("Cat rely events", func() {
		usrAddress := "0xbaa65281c2fa2baacb2cb550ba051525a480d3f4"
		msgSenderAddress := "0xa9ee75d81d78c36c4163004e6cc7a988eec9433e"
		relyIntegrationTest(int64(8928165), test_data.CatAddress(), msgSenderAddress, usrAddress)
	})

	Context("Flap rely events", func() {
		usrAddress := "0xbaa65281c2fa2baacb2cb550ba051525a480d3f4"
		msgSenderAddress := "0xd27a5f3416d8791fc238c148c93630d9e3c882e5"
		relyIntegrationTest(int64(8928163), test_data.FlapAddress(), msgSenderAddress, usrAddress)
	})

	Context("Flip ETH rely events", func() {
		usrAddress := "0xbaa65281c2FA2baAcb2cb550BA051525A480D3F4"
		msgSenderAddress := "0xBAB4FbeA257ABBfe84F4588d4Eedc43656E46Fc5"
		relyIntegrationTest(int64(8928180), test_data.FlipEthAddress(), msgSenderAddress, usrAddress)
	})

	Context("Flip WBTC rely events", func() {
		usrAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		msgSenderAddress := "0xBAB4FbeA257ABBfe84F4588d4Eedc43656E46Fc5"
		relyIntegrationTest(int64(9975625), test_data.FlipWbtcAddress(), msgSenderAddress, usrAddress)
	})

	Context("Flop rely events", func() {
		usrAddress := "0xbe8e3e3618f7474f8cb1d074a26affef007e98fb"
		msgSenderAddress := "0xddb108893104de4e1c6d0e47c42237db4e617acc"
		relyIntegrationTest(int64(9008136), test_data.FlopAddress(), msgSenderAddress, usrAddress)
	})

	Context("Jug rely events", func() {
		usrAddress := "0xbaa65281c2fa2baacb2cb550ba051525a480d3f4"
		msgSenderAddress := "0x45f0a929889ec8cc2d5b8cd79ab55e3279945cde"
		relyIntegrationTest(int64(8928160), test_data.JugAddress(), msgSenderAddress, usrAddress)
	})

	Context("Median BAT rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		relyIntegrationTest(int64(8956961), test_data.MedianBatAddress(), msgSenderAddress, usrAddress)
	})

	Context("Median ETH rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		relyIntegrationTest(int64(8956896), test_data.MedianEthAddress(), msgSenderAddress, usrAddress)
	})

	Context("Median WBTC rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		relyIntegrationTest(int64(8956963), test_data.MedianWbtcAddress(), msgSenderAddress, usrAddress)
	})

	Context("OSM BAT rely events", func() {
		usrAddress := "0x76416A4d5190d071bfed309861527431304aA14f"
		msgSenderAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		relyIntegrationTest(int64(9529100), test_data.OsmBatAddress(), msgSenderAddress, usrAddress)
	})

	Context("OSM ETH rely events", func() {
		usrAddress := "0x76416A4d5190d071bfed309861527431304aA14f"
		msgSenderAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		relyIntegrationTest(int64(9529100), test_data.OsmEthAddress(), msgSenderAddress, usrAddress)
	})

	// "OSM_USDC" is actually just a DS value contract right now and does not have "rely"
	XContext("OSM USDC rely events", func() {
		usrAddress := ""
		msgSenderAddress := ""
		relyIntegrationTest(int64(1), test_data.OsmUsdcAddress(), msgSenderAddress, usrAddress)
	})

	Context("OSM WBTC rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		relyIntegrationTest(int64(9975536), test_data.OsmWbtcAddress(), msgSenderAddress, usrAddress)
	})

	Context("Pot rely events", func() {
		usrAddress := "0xbaa65281c2fa2baacb2cb550ba051525a480d3f4"
		msgSenderAddress := "0x1a5ee7c64cf874c735968e3a42fa13f1c03427f9"
		relyIntegrationTest(int64(8928160), test_data.PotAddress(), msgSenderAddress, usrAddress)
	})

	Context("Spot rely events", func() {
		usrAddress := "0xbaa65281c2fa2baacb2cb550ba051525a480d3f4"
		msgSenderAddress := "0xdedd12bcb045c02b2fe11031c2b269bcde457410"
		relyIntegrationTest(int64(8928152), test_data.SpotAddress(), msgSenderAddress, usrAddress)
	})

	Context("Vow rely events", func() {
		usrAddress := "0xbaa65281c2fa2baacb2cb550ba051525a480d3f4"
		msgSenderAddress := "0x68322ca1a9aeb8c1d610b5fc8a8920aa0fba423b"
		relyIntegrationTest(int64(8928163), test_data.VowAddress(), msgSenderAddress, usrAddress)
	})
})

func relyIntegrationTest(blockNumber int64, contractAddressHex, msgSenderAddressHex, usrAddressHex string) {
	It("persists event", func() {
		test_config.CleanTestDB(db)
		logFetcher := fetcher.NewLogFetcher(blockChain)
		relyConfig := event.TransformerConfig{
			TransformerName: constants.RelyTable,
			Topic:           constants.RelySignature(),
		}
		initializer := event.ConfiguredTransformer{
			Config:      relyConfig,
			Transformer: auth.Transformer{TableName: constants.RelyTable},
		}

		header, headerErr := persistHeader(db, blockNumber, blockChain)
		Expect(headerErr).NotTo(HaveOccurred())

		initializer.Config.ContractAddresses = []string{contractAddressHex}
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		address := common.HexToAddress(contractAddressHex)
		topics := []common.Hash{common.HexToHash(relyConfig.Topic)}

		logs, logsErr := logFetcher.FetchLogs([]common.Address{address}, topics, header)
		Expect(logsErr).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		relyTransformer := initializer.NewTransformer(db)
		transformErr := relyTransformer.Execute(eventLogs)
		Expect(transformErr).NotTo(HaveOccurred())

		var dbResult []relyModel
		err := db.Select(&dbResult, `SELECT address_id, msg_sender, usr FROM maker.rely ORDER BY log_id`)
		Expect(err).NotTo(HaveOccurred())

		contractAddressID, contractAddressErr := shared.GetOrCreateAddress(contractAddressHex, db)
		Expect(contractAddressErr).NotTo(HaveOccurred())
		msgSenderAddressID, msgSenderAddressErr := shared.GetOrCreateAddress(msgSenderAddressHex, db)
		Expect(msgSenderAddressErr).NotTo(HaveOccurred())
		usrAddressID, usrAddressErr := shared.GetOrCreateAddress(usrAddressHex, db)
		Expect(usrAddressErr).NotTo(HaveOccurred())

		Expect(dbResult[0].AddressID).To(Equal(contractAddressID))
		Expect(dbResult[0].MsgSender).To(Equal(msgSenderAddressID))
		Expect(dbResult[0].Usr).To(Equal(usrAddressID))
	})
}

type relyModel struct {
	Usr       int64 `db:"usr"`
	MsgSender int64 `db:"msg_sender"`
	AddressID int64 `db:"address_id"`
}
