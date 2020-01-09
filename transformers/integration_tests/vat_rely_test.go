package integration_tests

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_auth"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vat Rely transformer", func() {
	relyConfig := transformer.EventTransformerConfig{
		TransformerName:   constants.VatRelyTable,
		ContractAddresses: []string{test_data.VatAddress()},
		Topic:             constants.RelySignature(),
	}

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("persists event", func() {
		blockNumber := int64(14764552)
		relyConfig.StartingBlockNumber = blockNumber
		relyConfig.EndingBlockNumber = blockNumber

		header, headerErr := persistHeader(db, blockNumber, blockChain)
		Expect(headerErr).NotTo(HaveOccurred())

		address := common.HexToAddress(relyConfig.ContractAddresses[0])
		topics := []common.Hash{common.HexToHash(relyConfig.Topic)}

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, logsErr := logFetcher.FetchLogs([]common.Address{address}, topics, header)
		Expect(logsErr).NotTo(HaveOccurred())

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		relyTransformer := event.ConfiguredTransformer{
			Config:      relyConfig,
			Transformer: vat_auth.Transformer{TableName: constants.VatRelyTable},
		}.NewTransformer(db)

		transformErr := relyTransformer.Execute(headerSyncLogs)
		Expect(transformErr).NotTo(HaveOccurred())

		var dbResult []vatRelyModel
		err := db.Select(&dbResult, `SELECT address_id, usr FROM maker.vat_rely`)
		Expect(err).NotTo(HaveOccurred())

		contractAddressID, contractAddressErr := shared.GetOrCreateAddress(relyConfig.ContractAddresses[0], db)
		Expect(contractAddressErr).NotTo(HaveOccurred())

		usrAddress := "0x0e4725db88bb038bba4c4723e91ba183be11edf3"
		usrAddressID, usrAddressErr := shared.GetOrCreateAddress(usrAddress, db)
		Expect(usrAddressErr).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].AddressID).To(Equal(contractAddressID))
		Expect(dbResult[0].Usr).To(Equal(usrAddressID))
	})
})

type vatRelyModel struct {
	Usr       int64 `db:"usr"`
	AddressID int64 `db:"address_id"`
}
