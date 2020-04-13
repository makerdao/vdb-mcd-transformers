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

		var dbResult []vatHopeModel
		err := db.Select(&dbResult, `SELECT usr FROM maker.vat_hope`)
		Expect(err).NotTo(HaveOccurred())

		usrAddress := "0x9759A6Ac90977b93B58547b4A71c78317f391A28"
		usrAddressID, usrAddressErr := shared.GetOrCreateAddress(usrAddress, db)
		Expect(usrAddressErr).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Usr).To(Or(Equal(usrAddressID), Equal(usrAddressID)))
	})
})

type vatHopeModel struct {
	Usr int64 `db:"usr"`
}
