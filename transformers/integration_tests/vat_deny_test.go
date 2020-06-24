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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vat Deny transformer", func() {
	denyConfig := event.TransformerConfig{
		TransformerName:   constants.VatDenyTable,
		ContractAddresses: []string{test_data.VatAddress()},
		Topic:             constants.DenySignature(),
	}

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("persists event", func() {
		blockNumber := int64(8928152)
		denyConfig.StartingBlockNumber = blockNumber
		denyConfig.EndingBlockNumber = blockNumber

		header, headerErr := persistHeader(db, blockNumber, blockChain)
		Expect(headerErr).NotTo(HaveOccurred())

		address := common.HexToAddress(denyConfig.ContractAddresses[0])
		topics := []common.Hash{common.HexToHash(denyConfig.Topic)}

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, logsErr := logFetcher.FetchLogs([]common.Address{address}, topics, header)
		Expect(logsErr).NotTo(HaveOccurred())

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		denyTransformer := event.ConfiguredTransformer{
			Config:      denyConfig,
			Transformer: vat_auth.Transformer{TableName: constants.VatDenyTable},
		}.NewTransformer(db)

		transformErr := denyTransformer.Execute(headerSyncLogs)
		Expect(transformErr).NotTo(HaveOccurred())

		var usr int64
		err := db.Get(&usr, `SELECT usr FROM maker.vat_deny`)
		Expect(err).NotTo(HaveOccurred())

		usrAddress := "0x403689148fa98a5a6fdcc0b984914ae968d788e5"
		usrAddressID, usrAddressErr := shared.GetOrCreateAddress(usrAddress, db)
		Expect(usrAddressErr).NotTo(HaveOccurred())

		Expect(usr).To(Equal(usrAddressID))
	})
})
