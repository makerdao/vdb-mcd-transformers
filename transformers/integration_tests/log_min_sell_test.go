package integration_tests

import (
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_min_sell"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogMinSell Transformer", func() {
	config := event.TransformerConfig{
		TransformerName:   constants.LogMinSellTable,
		ContractAddresses: test_data.OasisAddresses(),
		ContractAbi:       constants.OasisABI(),
		Topic:             constants.LogMinSellSignature(),
	}

	// OASIS_MATCHING_MARKET_ONE 0x39755357759ce0d7f32dc8dc45414cca409ae24e
	// OASIS_MATCHING_MARKET_TWO 0x794e6e91555438afc3ccf1c5076a74f42133d08d

	It("fetches and transforms a LogMinSell event for OASIS_MATCHING_MARKET_ONE contract", func() {
		blockNumber := int64(8944595)
		config.StartingBlockNumber = blockNumber
		config.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      config,
			Transformer: log_min_sell.Transformer{},
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

		var dbResults []logMinSellModel
		err = db.Select(&dbResults, `SELECT pay_gem, min_amount, address_id from maker.log_min_sell`)
		Expect(err).NotTo(HaveOccurred())

		expectedAddressID, addressErr := shared.GetOrCreateAddress(oasis_one_address, db)
		Expect(addressErr).NotTo(HaveOccurred())
		expectedPayGemID, payGemErr := shared.GetOrCreateAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F", db)
		Expect(payGemErr).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		dbResult := dbResults[0]
		Expect(dbResult.PayGem).To(Equal(expectedPayGemID))
		Expect(dbResult.MinAmount).To(Equal("2000000000000000000"))
		Expect(dbResult.AddressID).To(Equal(expectedAddressID))
	})

	It("fetches and transforms a LogMinSell event for OASIS_MATCHING_MARKET_TWO contract", func() {
		blockNumber := int64(9604711)
		config.StartingBlockNumber = blockNumber
		config.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      config,
			Transformer: log_min_sell.Transformer{},
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

		var dbResults []logMinSellModel
		err = db.Select(&dbResults, `SELECT pay_gem, min_amount, address_id from maker.log_min_sell`)
		Expect(err).NotTo(HaveOccurred())

		expectedAddressID, addressErr := shared.GetOrCreateAddress(oasis_two_address, db)
		Expect(addressErr).NotTo(HaveOccurred())
		expectedPayGemID, payGemErr := shared.GetOrCreateAddress("0x2260fac5e5542a773aa44fbcfedf7c193bc2c599", db)
		Expect(payGemErr).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(2))
		sort.Sort(byMinAmount(dbResults))
		dbResult := dbResults[0]
		Expect(dbResult.PayGem).To(Equal(expectedPayGemID))
		Expect(dbResult.MinAmount).To(Equal("21786"))
		Expect(dbResult.AddressID).To(Equal(expectedAddressID))
	})
})

type logMinSellModel struct {
	PayGem    int64  `db:"pay_gem"`
	MinAmount string `db:"min_amount"`
	AddressID int64  `db:"address_id"`
}

type byMinAmount []logMinSellModel

func (b byMinAmount) Len() int           { return len(b) }
func (b byMinAmount) Less(i, j int) bool { return b[i].MinAmount < b[j].MinAmount }
func (b byMinAmount) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
