package integration_tests

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_bump"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogBump Transformer", func() {
	config := event.TransformerConfig{
		TransformerName:   constants.LogBumpTable,
		ContractAddresses: test_data.OasisAddresses(),
		ContractAbi:       constants.OasisABI(),
		Topic:             constants.LogBumpSignature(),
	}

	It("fetches and transforms a LogBump event for OASIS_MATCHING_MARKET_ONE contract", func() {
		blockNumber := int64(8100012)
		config.StartingBlockNumber = blockNumber
		config.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      config,
			Transformer: log_bump.Transformer{},
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

		var dbResults []logBumpModel
		err = db.Select(&dbResults, `SELECT offer_id, pair, maker, pay_gem, buy_gem, pay_amt, buy_amt, timestamp, address_id from maker.log_bump`)
		Expect(err).NotTo(HaveOccurred())

		expectedAddressID, addressErr := shared.GetOrCreateAddress(oasisOneAddress, db)
		Expect(addressErr).NotTo(HaveOccurred())
		expectedMakerID, makerErr := shared.GetOrCreateAddress("0xa4da0f347c6abe0e8bc71b5981fd92b364eda4c2", db)
		Expect(makerErr).NotTo(HaveOccurred())
		expectedPayGemID, payGemErr := shared.GetOrCreateAddress("0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359", db)
		Expect(payGemErr).NotTo(HaveOccurred())
		expectedBuyGemID, buyGemErr := shared.GetOrCreateAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", db)
		Expect(buyGemErr).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		dbResult := dbResults[0]
		Expect(dbResult.OfferID).To(Equal("286499"))
		Expect(dbResult.Pair).To(Equal("0x10aed75aa327f09ef87e5bdfaedf498ca260499a251ae5e049ddbd5e1633cd9c"))
		Expect(dbResult.Maker).To(Equal(expectedMakerID))
		Expect(dbResult.PayGem).To(Equal(expectedPayGemID))
		Expect(dbResult.BuyGem).To(Equal(expectedBuyGemID))
		Expect(dbResult.PayAmt).To(Equal("6153000000000000000"))
		Expect(dbResult.BuyAmt).To(Equal("21000000000000000"))
		Expect(dbResult.Timestamp).To(Equal("1562446107"))
		Expect(dbResult.AddressID).To(Equal(expectedAddressID))
	})
})

type logBumpModel struct {
	OfferID   string `db:"offer_id"`
	Pair      string
	Maker     int64
	PayGem    int64  `db:"pay_gem"`
	BuyGem    int64  `db:"buy_gem"`
	PayAmt    string `db:"pay_amt"`
	BuyAmt    string `db:"buy_amt"`
	Timestamp string
	AddressID int64 `db:"address_id"`
}
