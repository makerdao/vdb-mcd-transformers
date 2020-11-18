package integration_tests

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/median_lift"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MedianLift EventTransformer", func() {
	medianLiftConfig := event.TransformerConfig{
		TransformerName:   constants.MedianLiftTable,
		ContractAddresses: test_data.MedianAddresses(),
		ContractAbi:       constants.MedianV100ABI(),
		Topic:             constants.MedianLiftSignature(),
	}

	It("transforms Median Lift log events from Median ETH", func() {
		blockNumber := int64(8934004)
		medianLiftConfig.StartingBlockNumber = blockNumber
		medianLiftConfig.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)

		header, headerErr := persistHeader(db, blockNumber, blockChain)
		Expect(headerErr).NotTo(HaveOccurred())

		medianEthAddress := test_data.MedianEthAddress()
		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, logErr := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(medianEthAddress)},
			[]common.Hash{common.HexToHash(medianLiftConfig.Topic)},
			header)
		Expect(logErr).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

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

		contractAddressID, contractAddressErr := shared.GetOrCreateAddress(test_data.MedianEthAddress(), db)
		Expect(contractAddressErr).NotTo(HaveOccurred())
		Expect(dbResult.AddressID).To(Equal(strconv.FormatInt(contractAddressID, 10)))

		msgSenderAddressID, msgSenderAddressErr := shared.GetOrCreateAddress("0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc", db)
		Expect(msgSenderAddressErr).NotTo(HaveOccurred())
		Expect(dbResult.MsgSender).To(Equal(strconv.FormatInt(msgSenderAddressID, 10)))

		Expect(dbResult.ALength).To(Equal("20"))

		var addresses []string
		addressesError := db.Get(pq.Array(&addresses), `SELECT a FROM maker.median_lift ORDER BY id`)
		Expect(addressesError).NotTo(HaveOccurred())
		Expect(addresses).To(ConsistOf(
			"0xaC8519b3495d8A3E3E44c041521cF7aC3f8F63B3",
			"0x4f95d9B4D842B2E2B1d1AC3f2Cf548B93Fd77c67",
			"0xE6367a7Da2b20ecB94A25Ef06F3b551baB2682e6",
			"0x238A3F4C923B75F3eF8cA3473A503073f0530801",
		))
	})

	It("transforms Median Lift log events from Median WBTC", func() {
		blockNumber := int64(8934019)
		medianLiftConfig.StartingBlockNumber = blockNumber
		medianLiftConfig.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)

		header, headerErr := persistHeader(db, blockNumber, blockChain)
		Expect(headerErr).NotTo(HaveOccurred())

		medianWbtcAddress := test_data.MedianWbtcAddress()
		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, logErr := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(medianWbtcAddress)},
			[]common.Hash{common.HexToHash(medianLiftConfig.Topic)},
			header)
		Expect(logErr).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

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

		contractAddressID, contractAddressErr := shared.GetOrCreateAddress(test_data.MedianWbtcAddress(), db)
		Expect(contractAddressErr).NotTo(HaveOccurred())
		Expect(dbResult.AddressID).To(Equal(strconv.FormatInt(contractAddressID, 10)))

		msgSenderAddressID, msgSenderAddressErr := shared.GetOrCreateAddress("0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc", db)
		Expect(msgSenderAddressErr).NotTo(HaveOccurred())
		Expect(dbResult.MsgSender).To(Equal(strconv.FormatInt(msgSenderAddressID, 10)))

		Expect(dbResult.ALength).To(Equal("20"))

		var addresses []string
		addressesError := db.Get(pq.Array(&addresses), `SELECT a FROM maker.median_lift ORDER BY id`)
		Expect(addressesError).NotTo(HaveOccurred())
		Expect(addresses).To(ConsistOf(
			"0xaC8519b3495d8A3E3E44c041521cF7aC3f8F63B3",
			"0x4f95d9B4D842B2E2B1d1AC3f2Cf548B93Fd77c67",
			"0xE6367a7Da2b20ecB94A25Ef06F3b551baB2682e6",
			"0x238A3F4C923B75F3eF8cA3473A503073f0530801",
		))
	})
})

type MedianLiftModel struct {
	AddressID string `db:"address_id"`
	MsgSender string `db:"msg_sender"`
	ALength   string `db:"a_length"`
}
