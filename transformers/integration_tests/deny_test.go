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

var _ = Describe("Deny transformer", func() {
	Context("Cat deny events", func() {
		usrAddress := "0xa9ee75d81d78c36c4163004e6cc7a988eec9433e"
		denyIntegrationTest(int64(8928165), test_data.CatAddress(), usrAddress, usrAddress)
	})

	Context("Flap deny events", func() {
		usrAddress := "0xd27a5f3416d8791fc238c148c93630d9e3c882e5"
		denyIntegrationTest(int64(8928163), test_data.FlapAddress(), usrAddress, usrAddress)
	})

	Context("Flip deny events", func() {
		usrAddress := "0xbab4fbea257abbfe84f4588d4eedc43656e46fc5"
		denyIntegrationTest(int64(8928180), test_data.EthFlipAddress(), usrAddress, usrAddress)
	})

	Context("Flop deny events", func() {
		usrAddress := "0xddb108893104de4e1c6d0e47c42237db4e617acc"
		denyIntegrationTest(int64(9008144), test_data.FlopAddress(), usrAddress, usrAddress)
	})

	Context("Jug deny events", func() {
		usrAddress := "0x45f0a929889ec8cc2d5b8cd79ab55e3279945cde"
		denyIntegrationTest(int64(8928160), test_data.JugAddress(), usrAddress, usrAddress)
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
