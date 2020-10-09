package integration_tests

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_auth"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vat Hope transformer", func() {
	hopeConfig := event.TransformerConfig{
		TransformerName:   constants.VatHopeTable,
		ContractAddresses: []string{test_data.VatAddress()},
		Topic:             constants.VatHopeSignature(),
	}

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("persists event", func() {
		blockNumber := int64(9861643)
		hopeConfig.StartingBlockNumber = blockNumber
		hopeConfig.EndingBlockNumber = blockNumber

		header, headerErr := persistHeader(db, blockNumber, blockChain)
		Expect(headerErr).NotTo(HaveOccurred())

		address := common.HexToAddress(hopeConfig.ContractAddresses[0])
		topics := []common.Hash{common.HexToHash(hopeConfig.Topic)}

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, logsErr := logFetcher.FetchLogs([]common.Address{address}, topics, header)
		Expect(logsErr).NotTo(HaveOccurred())

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		hopeTransformer := event.ConfiguredTransformer{
			Config:      hopeConfig,
			Transformer: vat_auth.Transformer{TableName: constants.VatHopeTable},
		}.NewTransformer(db)

		transformErr := hopeTransformer.Execute(headerSyncLogs)
		Expect(transformErr).NotTo(HaveOccurred())

		var usr int64
		err := db.Get(&usr, `SELECT usr FROM maker.vat_hope`)
		Expect(err).NotTo(HaveOccurred())

		usrAddress := "0x9759A6Ac90977b93B58547b4A71c78317f391A28"
		usrAddressID, usrAddressErr := repository.GetOrCreateAddress(db, usrAddress)
		Expect(usrAddressErr).NotTo(HaveOccurred())

		Expect(usr).To(Equal(usrAddressID))
	})
})
