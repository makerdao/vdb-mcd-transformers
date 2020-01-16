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
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Rely transformer", func() {
	XContext("Cat rely events", func() {
		usrAddress := "0x0e4725db88bb038bba4c4723e91ba183be11edf3"
		msgSenderAddress := "0x13141b8a5e4a82ebc6b636849dd6a515185d6236"
		relyIntegrationTest(int64(14764552), test_data.CatAddress(), msgSenderAddress, usrAddress)
	})

	Context("Flap rely events", func() {
		usrAddress := "0xbaa65281c2fa2baacb2cb550ba051525a480d3f4"
		msgSenderAddress := "0xd27a5f3416d8791fc238c148c93630d9e3c882e5"
		relyIntegrationTest(int64(8928163), test_data.FlapAddress(), msgSenderAddress, usrAddress)
	})

	Context("Flip rely events", func() {
		usrAddress := "0x13141b8a5e4a82ebc6b636849dd6a515185d6236"
		msgSenderAddress := "0xffb0382ca7cfdc4fc4d5cc8913af1393d7ee1ef1"
		relyIntegrationTest(int64(8928180), test_data.EthFlipAddress(), msgSenderAddress, usrAddress)
	})

	Context("Flop rely events", func() {
		usrAddress := "0xbaa65281c2fa2baacb2cb550ba051525a480d3f4"
		msgSenderAddress := "0xc41c4759f67ff54c7a7314d155f40fc6504f5d28"
		relyIntegrationTest(int64(8928163), test_data.FlopAddress(), msgSenderAddress, usrAddress)
	})

	XContext("Jug rely events", func() {
		usrAddress := "0x0e4725db88bb038bba4c4723e91ba183be11edf3"
		msgSenderAddress := "0x13141b8a5e4a82ebc6b636849dd6a515185d6236"
		relyIntegrationTest(int64(14764552), test_data.JugAddress(), msgSenderAddress, usrAddress)
	})

	XContext("Pot rely events", func() {
		usrAddress := "0x0e4725db88bb038bba4c4723e91ba183be11edf3"
		msgSenderAddress := "0x13141b8a5e4a82ebc6b636849dd6a515185d6236"
		relyIntegrationTest(int64(14764552), test_data.PotAddress(), msgSenderAddress, usrAddress)
	})

	XContext("Spot rely events", func() {
		usrAddress := "0x0e4725db88bb038bba4c4723e91ba183be11edf3"
		msgSenderAddress := "0x13141b8a5e4a82ebc6b636849dd6a515185d6236"
		relyIntegrationTest(int64(14764552), test_data.SpotAddress(), msgSenderAddress, usrAddress)
	})

	XContext("Vow rely events", func() {
		usrAddress := "0xbaa65281c2fa2baacb2cb550ba051525a480d3f4"
		msgSenderAddress := "0xc41c4759f67ff54c7a7314d155f40fc6504f5d28"
		relyIntegrationTest(int64(14764552), test_data.VowAddress(), msgSenderAddress, usrAddress)
	})
})

func relyIntegrationTest(blockNumber int64, contractAddressHex, msgSenderAddressHex, usrAddressHex string) {
	It("persists event", func() {
		test_config.CleanTestDB(db)
		logFetcher := fetcher.NewLogFetcher(blockChain)
		relyConfig := transformer.EventTransformerConfig{
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
