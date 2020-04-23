package integration_tests

import (
	"strconv"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lib/pq"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/median_lift"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MedianLift EventTransformer", func() {
	medianLiftConfig := event.TransformerConfig{
		TransformerName:   constants.MedianLiftTable,
		ContractAddresses: test_data.MedianAddresses(),
		ContractAbi:       constants.MedianABI(),
		Topic:             constants.MedianLiftSignature(),
	}

	XIt("transforms Median Lift log events", func() {
		blockNumber := int64(8936530)
		medianLiftConfig.StartingBlockNumber = blockNumber
		medianLiftConfig.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)

		header, headerErr := persistHeader(db, blockNumber, blockChain)
		Expect(headerErr).NotTo(HaveOccurred())

		// TODO: fetch event from blockchain once one exists
		// logFetcher := fetcher.NewLogFetcher(blockChain)
		// logs, logErr := logFetcher.FetchLogs(
		// 	event.HexStringsToAddresses(medianLiftConfig.ContractAddresses),
		// 	[]common.Hash{common.HexToHash(medianLiftConfig.Topic)},
		// 	header)
		// Expect(logErr).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, []types.Log{test_data.RawMedianLiftLogWithFiveAccounts}, db)

		transformer := event.ConfiguredTransformer{
			Config:      medianLiftConfig,
			Transformer: median_lift.Transformer{},
		}.NewTransformer(db)

		transformErr := transformer.Execute(eventLogs)
		Expect(transformErr).NotTo(HaveOccurred())

		var dbResults []MedianLiftModel
		err := db.Select(&dbResults, `SELECT address_id, msg_sender, a_length FROM maker.median_lift`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		dbResult := dbResults[0]

		contractAddressID, contractAddressErr := shared.GetOrCreateAddress(test_data.EthMedianAddress(), db)
		Expect(contractAddressErr).NotTo(HaveOccurred())
		Expect(dbResult.AddressID).To(Equal(strconv.FormatInt(contractAddressID, 10)))

		msgSenderAddressID, msgSenderAddressErr := shared.GetOrCreateAddress("0xc45E7858EEf1318337A803Ede8C5A9bE12E2B40f", db)
		Expect(msgSenderAddressErr).NotTo(HaveOccurred())
		Expect(dbResult.MsgSender).To(Equal(strconv.FormatInt(msgSenderAddressID, 10)))

		Expect(dbResult.ALength).To(Equal("5"))

		var addresses []string
		addressesError := db.Get(pq.Array(&addresses), `SELECT a FROM maker.median_lift ORDER BY id`)
		Expect(addressesError).NotTo(HaveOccurred())
		Expect(addresses).To(ConsistOf(
			"0x6bDbc0ccC17d72a33Bf72a4657781a37DC2aa94E",
			"0x26c45f7B0E456E36fC85781488A3CD42A57CcbD2",
			"0x20c576F989EE94E571F027b30314aCF709267F7C",
			"0xFCb1fB52E114b364B3Aab63d8a6f65Fe8dcbeF9D",
		))
	})
})

type MedianLiftModel struct {
	AddressID string `db:"address_id"`
	MsgSender string `db:"msg_sender"`
	ALength   string `db:"a_length"`
}
