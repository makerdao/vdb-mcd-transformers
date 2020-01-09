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

var _ = Describe("Deny transformer", func() {
	Context("Cat deny events", func() {
		usrAddress := "0x13141b8a5e4a82ebc6b636849dd6a515185d6236"
		denyIntegrationTest(int64(14764643), test_data.CatAddress(), usrAddress, usrAddress)
	})

	Context("Flap deny events", func() {
		usrAddress := "0x13141b8a5e4a82ebc6b636849dd6a515185d6236"
		denyIntegrationTest(int64(14764643), test_data.FlapAddress(), usrAddress, usrAddress)
	})

	Context("Flip deny events", func() {
		usrAddress := "0x13141b8a5e4a82ebc6b636849dd6a515185d6236"
		denyIntegrationTest(int64(14764652), test_data.EthFlipAddress(), usrAddress, usrAddress)
	})

	Context("Flop deny events", func() {
		usrAddress := "0xdb33dfd3d61308c33c63209845dad3e6bfb2c674"
		denyIntegrationTest(int64(15196525), test_data.FlopAddress(), usrAddress, usrAddress)
	})

	Context("Jug deny events", func() {
		usrAddress := "0x13141b8a5e4a82ebc6b636849dd6a515185d6236"
		denyIntegrationTest(int64(14764643), test_data.JugAddress(), usrAddress, usrAddress)
	})

	Context("Pot deny events", func() {
		usrAddress := "0x13141b8a5e4a82ebc6b636849dd6a515185d6236"
		denyIntegrationTest(int64(14764643), test_data.PotAddress(), usrAddress, usrAddress)
	})

	Context("Spot deny events", func() {
		usrAddress := "0x13141b8a5e4a82ebc6b636849dd6a515185d6236"
		denyIntegrationTest(int64(14764643), test_data.SpotAddress(), usrAddress, usrAddress)
	})

	Context("Vow deny events", func() {
		usrAddress := "0x13141b8a5e4a82ebc6b636849dd6a515185d6236"
		denyIntegrationTest(int64(14764643), test_data.VowAddress(), usrAddress, usrAddress)
	})
})

func denyIntegrationTest(blockNumber int64, contractAddressHex, msgSenderAddressHex, usrAddressHex string) {
	It("persists event", func() {
		test_config.CleanTestDB(db)
		logFetcher := fetcher.NewLogFetcher(blockChain)
		denyConfig := transformer.EventTransformerConfig{
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

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		denyTransformer := initializer.NewTransformer(db)
		transformErr := denyTransformer.Execute(headerSyncLogs)
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

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].AddressID).To(Equal(contractAddressID))
		Expect(dbResult[0].MsgSender).To(Equal(msgSenderAddressID))
		Expect(dbResult[0].Usr).To(Equal(usrAddressID))
	})
}

type denyModel struct {
	Usr       int64 `db:"usr"`
	MsgSender int64 `db:"msg_sender"`
	AddressID int64 `db:"address_id"`
}
