package integration_tests

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/deny"
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
	const (
		defaultOffset int = 0
		vatOffset     int = -1
	)

	Context("Cat deny events", func() {
		usrAddress := "0x13141b8a5e4a82ebc6b636849dd6a515185d6236"
		denyIntegrationTest(int64(14764643), test_data.CatAddress(), usrAddress, defaultOffset)
	})

	Context("Flap deny events", func() {
		usrAddress := "0x13141b8a5e4a82ebc6b636849dd6a515185d6236"
		denyIntegrationTest(int64(14764643), test_data.FlapAddress(), usrAddress, defaultOffset)
	})

	Context("Flip deny events", func() {
		usrAddress := "0x13141b8a5e4a82ebc6b636849dd6a515185d6236"
		denyIntegrationTest(int64(14764652), test_data.EthFlipAddress(), usrAddress, defaultOffset)
	})

	Context("Flop deny events", func() {
		usrAddress := "0x13141b8a5e4a82ebc6b636849dd6a515185d6236"
		denyIntegrationTest(int64(14764643), test_data.FlopAddress(), usrAddress, defaultOffset)
	})

	Context("Jug deny events", func() {
		usrAddress := "0x13141b8a5e4a82ebc6b636849dd6a515185d6236"
		denyIntegrationTest(int64(14764643), test_data.JugAddress(), usrAddress, defaultOffset)
	})

	Context("Pot deny events", func() {
		usrAddress := "0x13141b8a5e4a82ebc6b636849dd6a515185d6236"
		denyIntegrationTest(int64(14764643), test_data.PotAddress(), usrAddress, defaultOffset)
	})

	Context("Spot deny events", func() {
		usrAddress := "0x13141b8a5e4a82ebc6b636849dd6a515185d6236"
		denyIntegrationTest(int64(14764643), test_data.SpotAddress(), usrAddress, defaultOffset)
	})

	Context("Vat deny events", func() {
		usrAddress := "0x13141b8a5e4a82ebc6b636849dd6a515185d6236"
		denyIntegrationTest(int64(14764643), test_data.VowAddress(), usrAddress, vatOffset)
	})

	Context("Vow deny events", func() {
		usrAddress := "0x13141b8a5e4a82ebc6b636849dd6a515185d6236"
		denyIntegrationTest(int64(14764643), test_data.VowAddress(), usrAddress, defaultOffset)
	})
})

func denyIntegrationTest(blockNumber int64, contractAddressHex, usrAddressHex string, logNoteArgumentOffset int) {
	It("persists event", func() {
		test_config.CleanTestDB(db)
		logFetcher := fetcher.NewLogFetcher(blockChain)
		denyConfig := transformer.EventTransformerConfig{
			TransformerName: constants.DenyTable,
			Topic:           constants.DenySignature(),
		}
		initializer := event.Transformer{
			Config:    denyConfig,
			Converter: deny.Converter{LogNoteArgumentOffset: logNoteArgumentOffset},
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
		err := db.Select(&dbResult, `SELECT address_id, usr FROM maker.deny`)
		Expect(err).NotTo(HaveOccurred())

		contractAddressID, contractAddressErr := shared.GetOrCreateAddress(contractAddressHex, db)
		Expect(contractAddressErr).NotTo(HaveOccurred())
		usrAddressID, usrAddressErr := shared.GetOrCreateAddress(usrAddressHex, db)
		Expect(usrAddressErr).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].AddressId).To(Equal(contractAddressID))
		Expect(dbResult[0].Usr).To(Equal(usrAddressID))
	})
}

type denyModel struct {
	Usr       int64 `db:"usr"`
	AddressId int64 `db:"address_id"`
}
