package integration_tests

import (
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_kill"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogKill Transformer", func() {
	config := event.TransformerConfig{
		TransformerName:   constants.LogKillTable,
		ContractAddresses: test_data.OasisAddresses(),
		ContractAbi:       constants.OasisABI(),
		Topic:             constants.LogKillSignature(),
	}

	It("fetches and transforms a LogKill event for OASIS_MATCHING_MARKET_ONE contract", func() {
		blockNumber := int64(9613377)
		config.StartingBlockNumber = blockNumber
		config.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      config,
			Transformer: log_kill.Transformer{},
		}
		transformer := initializer.NewTransformer(db)

		oasis_one_address := constants.GetContractAddress("OASIS_MATCHING_MARKET_ONE")
		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(oasis_one_address)},
			[]common.Hash{common.HexToHash(config.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResults []logKillModel
		err = db.Select(&dbResults, `SELECT offer_id, pair, maker, pay_gem, buy_gem, pay_amt, buy_amt, timestamp, address_id from maker.log_kill`)
		Expect(err).NotTo(HaveOccurred())

		expectedAddressID, addressErr := shared.GetOrCreateAddress(oasis_one_address, db)
		Expect(addressErr).NotTo(HaveOccurred())

		expectedMakerID, makerErr := shared.GetOrCreateAddress("0x384d6a80C87D2f185faf095a137888E5E6156e80", db)
		Expect(makerErr).NotTo(HaveOccurred())
		expectedPayGemID, payGemErr := shared.GetOrCreateAddress("0x0D8775F648430679A709E98d2b0Cb6250d2887EF", db)
		Expect(payGemErr).NotTo(HaveOccurred())
		expectedBuyGemID, buyGemErr := shared.GetOrCreateAddress("0x89d24A6b4CcB1B6fAA2625fE562bDD9a23260359", db)
		Expect(buyGemErr).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(15))
		sort.Sort(logKillByOfferID(dbResults))
		dbResult := dbResults[0]
		Expect(dbResult.OfferID).To(Equal("548250"))
		Expect(dbResult.Pair).To(Equal("0x15977f73fa1d2be3f367c4b6006149460588f6b243578cbda0b0cb37e46f4e5c"))
		Expect(dbResult.Maker).To(Equal(expectedMakerID))
		Expect(dbResult.PayGem).To(Equal(expectedPayGemID))
		Expect(dbResult.BuyGem).To(Equal(expectedBuyGemID))
		Expect(dbResult.PayAmt).To(Equal("97221450402889691581"))
		Expect(dbResult.BuyAmt).To(Equal("34990000000000000001"))
		Expect(dbResult.Timestamp).To(Equal("1583440730"))
		Expect(dbResult.AddressID).To(Equal(expectedAddressID))
	})

	FIt("fetches and transforms a LogKill event for OASIS_MATCHING_MARKET_TWO contract", func() {
		blockNumber := int64(9883734)
		config.StartingBlockNumber = blockNumber
		config.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      config,
			Transformer: log_kill.Transformer{},
		}
		transformer := initializer.NewTransformer(db)

		oasis_two_address := constants.GetContractAddress("OASIS_MATCHING_MARKET_TWO")
		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(oasis_two_address)},
			[]common.Hash{common.HexToHash(config.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResults []logKillModel
		err = db.Select(&dbResults, `SELECT offer_id, pair, maker, pay_gem, buy_gem, pay_amt, buy_amt, timestamp, address_id from maker.log_kill`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		dbResult := dbResults[0]
		expectedAddressID, addressErr := shared.GetOrCreateAddress(oasis_two_address, db)
		Expect(addressErr).NotTo(HaveOccurred())
		expectedMakerID, makerErr := shared.GetOrCreateAddress("0x59df5a7df54000fcc09dfb303d24b0d302182540", db)
		Expect(makerErr).NotTo(HaveOccurred())
		expectedPayGemID, payGemErr := shared.GetOrCreateAddress("0x514910771af9ca656af840dff83e8264ecf986ca", db)
		Expect(payGemErr).NotTo(HaveOccurred())
		expectedBuyGemID, buyGemErr := shared.GetOrCreateAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2", db)
		Expect(buyGemErr).NotTo(HaveOccurred())

		Expect(dbResult.OfferID).To(Equal("254470"))
		Expect(dbResult.Pair).To(Equal("0xdf83bd62157b119877a9fd8ef3a81d9c97c1402355414cdecd2102bb699391c7"))
		Expect(dbResult.Maker).To(Equal(expectedMakerID))
		Expect(dbResult.PayGem).To(Equal(expectedPayGemID))
		Expect(dbResult.BuyGem).To(Equal(expectedBuyGemID))
		Expect(dbResult.PayAmt).To(Equal("1613048242893176673399"))
		Expect(dbResult.BuyAmt).To(Equal("32117422686047620972"))
		Expect(dbResult.Timestamp).To(Equal("1587043274"))
		Expect(dbResult.AddressID).To(Equal(expectedAddressID))
	})
})

type logKillModel struct {
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

type logKillByOfferID []logKillModel

func (l logKillByOfferID) Len() int           { return len(l) }
func (l logKillByOfferID) Less(i, j int) bool { return l[i].OfferID < l[j].OfferID }
func (l logKillByOfferID) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
