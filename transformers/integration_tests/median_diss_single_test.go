package integration_tests

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/median_diss/single"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// TODO: Update once a diss event has happened
var _ = XDescribe("MedianDissSingle EventTransformer", func() {
	medianDissConfig := event.TransformerConfig{
		TransformerName:   constants.MedianDissSingleTable,
		ContractAddresses: test_data.MedianAddresses(),
		ContractAbi:       constants.MedianV100ABI(),
		Topic:             constants.MedianDissSingleSignature(),
	}

	It("transforms Median diss single log events", func() {
		blockNumber := int64(8936530)
		medianDissConfig.StartingBlockNumber = blockNumber
		medianDissConfig.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, logErr := logFetcher.FetchLogs(
			event.HexStringsToAddresses(medianDissConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(medianDissConfig.Topic)},
			header)
		Expect(logErr).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		transformer := event.ConfiguredTransformer{
			Config:      medianDissConfig,
			Transformer: single.Transformer{},
		}.NewTransformer(db)

		transformErr := transformer.Execute(eventLogs)
		Expect(transformErr).NotTo(HaveOccurred())

		var dbResults []MedianDissSingleModel
		err = db.Select(&dbResults, `SELECT address_id, msg_sender, a FROM maker.median_diss_single ORDER BY address_id`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(2))
		dbResult := dbResults[1]

		contractAddressID, contractAddressErr := shared.GetOrCreateAddress("0x18B4633D6E39870f398597f3c1bA8c4A41294966", db)
		Expect(contractAddressErr).NotTo(HaveOccurred())
		Expect(dbResult.AddressID).To(Equal(strconv.FormatInt(contractAddressID, 10)))

		msgSenderAddressID, msgSenderAddressErr := shared.GetOrCreateAddress("0x000000000000000000000000ddb108893104de4e1c6d0e47c42237db4e617acc", db)
		Expect(msgSenderAddressErr).NotTo(HaveOccurred())
		Expect(dbResult.MsgSender).To(Equal(strconv.FormatInt(msgSenderAddressID, 10)))

		aAddressID, aAddressErr := shared.GetOrCreateAddress("0x000000000000000000000000b4eb54af9cc7882df0121d26c5b97e802915abe6", db)
		Expect(aAddressErr).NotTo(HaveOccurred())
		Expect(dbResult.A).To(Equal(strconv.FormatInt(aAddressID, 10)))
	})
})

type MedianDissSingleModel struct {
	AddressID string `db:"address_id"`
	MsgSender string `db:"msg_sender"`
	A         string
}
