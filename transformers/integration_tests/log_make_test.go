package integration_tests

import (
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_make"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogMake Transformer", func() {
	config := event.TransformerConfig{
		TransformerName:   constants.LogMakeTable,
		ContractAddresses: test_data.OasisAddresses(),
		ContractAbi:       constants.OasisABI(),
		Topic:             constants.LogMakeSignature(),
	}

	It("fetches and transforms a LogMake event for OASIS_MATCHING_MARKET_ONE contract", func() {
		blockNumber := int64(9440386)
		config.StartingBlockNumber = blockNumber
		config.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      config,
			Transformer: log_make.Transformer{},
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

		var dbResults []logMakeModel
		err = db.Select(&dbResults, `SELECT offer_id, pair, maker, pay_gem, buy_gem, pay_amt, buy_amt, timestamp, address_id from maker.log_make`)
		Expect(err).NotTo(HaveOccurred())

		expectedAddressID, addressErr := shared.GetOrCreateAddress(oasisOneAddress, db)
		Expect(addressErr).NotTo(HaveOccurred())
		expectedMakerID, makerErr := shared.GetOrCreateAddress("0x6Ff7D252627D35B8eb02607c8F27ACDB18032718", db)
		Expect(makerErr).NotTo(HaveOccurred())
		expectedPayGemID, payGemErr := shared.GetOrCreateAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F", db)
		Expect(payGemErr).NotTo(HaveOccurred())
		expectedBuyGemID, buyGemErr := shared.GetOrCreateAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2", db)
		Expect(buyGemErr).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(2))
		sort.Sort(byOfferID(dbResults))
		dbResult := dbResults[0]
		Expect(dbResult.OfferID).To(Equal("811645"))
		Expect(dbResult.Pair).To(Equal("0x7bda8b27e891f9687bd6d3312ab3f4f458e2cc91916429d721d617df7ac29fb8"))
		Expect(dbResult.Maker).To(Equal(expectedMakerID))
		Expect(dbResult.PayGem).To(Equal(expectedPayGemID))
		Expect(dbResult.BuyGem).To(Equal(expectedBuyGemID))
		Expect(dbResult.PayAmt).To(Equal("67707612000000000000000"))
		Expect(dbResult.BuyAmt).To(Equal("307650000000000000000"))
		Expect(dbResult.Timestamp).To(Equal("1581142121"))
		Expect(dbResult.AddressID).To(Equal(expectedAddressID))
	})

	It("fetches and transforms a LogMake event for OASIS_MATCHING_MARKET_TWO contract", func() {
		blockNumber := int64(9866954)
		config.StartingBlockNumber = blockNumber
		config.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      config,
			Transformer: log_make.Transformer{},
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

		var dbResults []logMakeModel
		err = db.Select(&dbResults, `SELECT offer_id, pair, maker, pay_gem, buy_gem, pay_amt, buy_amt, timestamp, address_id from maker.log_make`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		dbResult := dbResults[0]
		expectedAddressID, addressErr := shared.GetOrCreateAddress(oasisTwoAddress, db)
		Expect(addressErr).NotTo(HaveOccurred())
		expectedMakerID, makerErr := shared.GetOrCreateAddress("0xbAEaFc49d8e3a636d61df1F14fd45b97c7018020", db)
		Expect(makerErr).NotTo(HaveOccurred())
		expectedPayGemID, payGemErr := shared.GetOrCreateAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F", db)
		Expect(payGemErr).NotTo(HaveOccurred())
		expectedBuyGemID, buyGemErr := shared.GetOrCreateAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2", db)
		Expect(buyGemErr).NotTo(HaveOccurred())

		Expect(dbResult.OfferID).To(Equal("247773"))
		Expect(dbResult.Pair).To(Equal("0x7bda8b27e891f9687bd6d3312ab3f4f458e2cc91916429d721d617df7ac29fb8"))
		Expect(dbResult.Maker).To(Equal(expectedMakerID))
		Expect(dbResult.PayGem).To(Equal(expectedPayGemID))
		Expect(dbResult.BuyGem).To(Equal(expectedBuyGemID))
		Expect(dbResult.PayAmt).To(Equal("3323372775780000000000"))
		Expect(dbResult.BuyAmt).To(Equal("21889900000000000000"))
		Expect(dbResult.Timestamp).To(Equal("1586819716"))
		Expect(dbResult.AddressID).To(Equal(expectedAddressID))
	})
})

type logMakeModel struct {
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

type byOfferID []logMakeModel

func (b byOfferID) Len() int           { return len(b) }
func (b byOfferID) Less(i, j int) bool { return b[i].OfferID < b[j].OfferID }
func (b byOfferID) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
