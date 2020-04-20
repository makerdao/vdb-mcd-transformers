package integration_tests

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_trade"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogTrade Transformer", func() {
	config := event.TransformerConfig{
		TransformerName:   constants.LogTradeTable,
		ContractAddresses: test_data.OasisAddresses(),
		ContractAbi:       constants.OasisABI(),
		Topic:             constants.LogTradeSignature(),
	}

	It("fetches and transforms a LogTrade event for OASIS_MATCHING_MARKET_ONE contract", func() {
		blockNumber := int64(9439641)
		config.StartingBlockNumber = blockNumber
		config.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      config,
			Transformer: log_trade.Transformer{},
		}
		transformer := initializer.NewTransformer(db)

		oasisOneAddress := constants.GetContractAddress("OASIS_MATCHING_MARKET_ONE")
		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(oasisOneAddress)},
			[]common.Hash{common.HexToHash(config.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResults []logTradeModel
		err = db.Select(&dbResults, `SELECT pay_gem, buy_gem, pay_amt, buy_amt, address_id from maker.log_trade`)
		Expect(err).NotTo(HaveOccurred())

		expectedAddressID, addressErr := shared.GetOrCreateAddress(oasisOneAddress, db)
		Expect(addressErr).NotTo(HaveOccurred())
		expectedPayGemID, payGemErr := shared.GetOrCreateAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", db)
		Expect(payGemErr).NotTo(HaveOccurred())
		expectedBuyGemID, buyGemErr := shared.GetOrCreateAddress("0x6b175474e89094c44da98b954eedeac495271d0f", db)
		Expect(buyGemErr).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		Expect(dbResults[0].PayGem).To(Equal(expectedPayGemID))
		Expect(dbResults[0].BuyGem).To(Equal(expectedBuyGemID))
		Expect(dbResults[0].PayAmt).To(Equal("4590000000000000000"))
		Expect(dbResults[0].BuyAmt).To(Equal("999288900000000000000"))
		Expect(dbResults[0].AddressID).To(Equal(expectedAddressID))
	})

	It("fetches and transforms a LogTrade event for OASIS_MATCHING_MARKET_TWO contract", func() {
		blockNumber := int64(9880684)
		config.StartingBlockNumber = blockNumber
		config.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      config,
			Transformer: log_trade.Transformer{},
		}
		transformer := initializer.NewTransformer(db)

		oasisTwoAddress := constants.GetContractAddress("OASIS_MATCHING_MARKET_TWO")
		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(oasisTwoAddress)},
			[]common.Hash{common.HexToHash(config.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResults []logTradeModel
		err = db.Select(&dbResults, `SELECT pay_gem, buy_gem, pay_amt, buy_amt, address_id from maker.log_trade`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		expectedAddressID, addressErr := shared.GetOrCreateAddress(oasisTwoAddress, db)
		Expect(addressErr).NotTo(HaveOccurred())
		expectedPayGemID, payGemErr := shared.GetOrCreateAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", db)
		Expect(payGemErr).NotTo(HaveOccurred())
		expectedBuyGemID, buyGemErr := shared.GetOrCreateAddress("0x6b175474e89094c44da98b954eedeac495271d0f", db)
		Expect(buyGemErr).NotTo(HaveOccurred())

		Expect(dbResults[0].PayGem).To(Equal(expectedPayGemID))
		Expect(dbResults[0].BuyGem).To(Equal(expectedBuyGemID))
		Expect(dbResults[0].PayAmt).To(Equal("26585039969973563836"))
		Expect(dbResults[0].BuyAmt).To(Equal("3999999999999999999960"))
		Expect(dbResults[0].AddressID).To(Equal(expectedAddressID))
	})
})

type logTradeModel struct {
	PayGem    int64  `db:"pay_gem"`
	BuyGem    int64  `db:"buy_gem"`
	PayAmt    string `db:"pay_amt"`
	BuyAmt    string `db:"buy_amt"`
	AddressID int64  `db:"address_id"`
}
